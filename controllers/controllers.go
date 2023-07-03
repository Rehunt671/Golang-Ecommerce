package controllers


import(
	"github.com/gin-gonic/gin"
)

func HashPassWord(password string) string{
 
}

func VerifyPassword(userpassword string,givenPassword string)(bool,string){

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

}


func SearchProductByQuery() gin.HandlerFunc{
	
}