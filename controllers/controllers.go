package controllers


import(
	"github.com/gin-gonic/gin"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/go-playground/validator/v10"
	"github.com/kingsglaive/database"
	"golang.org/x/crypto/bcrypt"
)
var UserCollection  *mongo.Collection = database.UserData(database.Client,"Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client,"Products")
var Validate = validator.New()
func HashPassWord(password string) string{
	bytes,err := bcypt.GenerateFromPassword([]byte(password),14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string,givenPassword string)(bool,string){
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	valid := true
	mas := ""
	if err != nil {
		 msg  = "Login or Password is incorrect"
		 valid = false
	}
	return valid,msg
}


func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeOut(context.Background(),100*time.Second)
		defer cancel()
		
		var user models.User 
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		validateErr := Validate.Struct(user)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr})
			return
		}
		count, err := UserCollection.CountDocuments(ctx,bson.M{"email": user.Email})
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError , gin.H{"error" : err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest,gin.H{"error": "user already exists"})
		}
		count,err = UserCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return
		}
		if count > 0{
			c.JSON(http.StatusBadRequest,gin.H{"error":"this phone no,is already in use"})
			return
		}
		password := HashPassWord(*user.Password)
		user.Password := &password
		user.Created_At , _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_At,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID	= primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token,refresh_token,_ := generate.TokenGenerator(*user.Email,*user.First_Name,*user.Last_Name,user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.UserCart = make([]modles.ProductUser,0)
		user.Address_Details = make([]models.Address,0)
		user.Order_Status = make([] models.Order,0)
		_,inserterr := UserCollection.InsertOne(ctx,user)
		if inserterr != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"the user did not get created"})
			return 
		}
		defer cancel()
		c.JSON(http.StatusCreated , "Successfully signed in")
	}
}
func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		context.WithTimeOut(context.Background() , 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error",err})
		}
		err := UserCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&founduser)
		defer cancel()
		
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}
		PasswordIsValid,msg := VerifyPassword(*user.Password,*founduser.Password)
		defer cancel()
		if !PasswordIsValid{
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			fmt.Println(msg)
			return 
		}
		token, refresh_token := generate.TokenGenerator(*founduser.Email,*founduser.First_Name,*founduser.Last_Name,founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(tokens,refresh_token,founduser.User_ID)

		c.JSON(http.StatusFound,founduser)

	}
}

func ProductViewerAdmin() gin.HandlerFunc{

}



func SearchProduct() gin.HandlerFunc{

	return func(c *gin.Context) {
		var productList []models.Product
		var ctx,cancel = context.WithTimeOut(context.Background() , 100 *time.Second)
		defer cancel()

		cursor,err := ProductCollection.Find(ctx,bson.D{{}})
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,"soemthing went wrong, please try after some time")
			
			err = cursor.All(ctx,&productlist)

			if err != nil{
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return 
			}
			defer cursor.Close()
			if err := cursor.err() ; err != nil{
				log.Println(err)
				c.IndentedJSON(400,"invalid")
				return
			}
			defer cancel()
			c.IndentedJSON(200,productlist)
		}
	}
}


func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context){
		var searchproducts []models.Product
		queryParam := c.Query("name")

		if queryParam == ""{
			log.Println("query is empty")
			c.Header("Content-Type","application/json")
			c.JSON(HTTP.StatusNotFound,gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}
		var ctx,cancel = context.WithTimeOut(context.Background(),100*time.Second)
		defer cancel()
		
		searchquerydb,err := ProductCollection.Find(ctx,bson.M{"product_name":bson.M{"$regex":queryParam}})
		if err != nil{
			c.IndentedJSON(404,"something went wrong while fetching the data")
			return 
		}
		searchquerydb.All(ctx,&searchproducts)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(400,"invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err();err != nil{
			log.Println(err)
			c.IndentedJSON(400,"invalid request")
			return
		}
		c.IndentedJSON(200,searchproducts)
	}
}