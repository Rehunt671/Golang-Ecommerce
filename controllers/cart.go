package controllers

import(
	"time"
	"context"
	"errors"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/kingsglaive/database"
)

type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}
func NewApplication(prodCollection,userCollection *mongo.Collection) *Application{
	return &Application{
			prodCollection:prodCollection,
			userCollection:userCollection,
	}
}
func (app *Application) AddToCart() gin.HandlerFunc{
	return func (c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("product id is empty"))
			return
		}
		userQueryID  := c.Query("userID")
		if userQueryID ==""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("user is is empty"))
			return 
		}
		productID,err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError ,errors.New("can't convert productID to primitive type"))
			return 
		}
		
		var ctx,cancel = context.WithTimeout(context.Background(),5 *time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx,app.prodCollection,app.userCollection,productID,userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON(200,"Successfully added to the cart")
	}
}

func (app *Application)RemoveItem() gin.HandlerFunc{
		return func(c *gin.Context) {
			productQueryID := c.Query("id")
			if productQueryID == ""{
				log.Println("product id is empty")
				_ = c.AbortWithError(http.StatusBadRequest,errors.New("product id is empty"))
				return
			}
			userQueryID  := c.Query("userID")
			if userQueryID ==""{
				log.Println("user id is empty")
				_ = c.AbortWithError(http.StatusBadRequest,errors.New("user is is empty"))
				return 
			}
			productID,err := primitive.ObjectIDFromHex(productQueryID)
			if err != nil{
				log.Println(err)
				c.AbortWithError(http.StatusInternalServerError ,errors.New("can't convert productID to primitive type"))
				return 
			}
			
			var ctx,cancel = context.WithTimeout(context.Background(),5 *time.Second)
			defer cancel()

			err  = database.RemoveCartItem(ctx,app.prodCollection,app.userCollection,productID	,userQueryID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError,err)
				return
			}
			c.IndentedJSON(200,"Successfully remove item from cart")
	
		}
}

func (app *Application)GetItemFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id ==""{
			c.Header("Context-type","application/json")
			c.JSON(http.StatusNotFound,gin.H{"error":"invalid id"})
			c.Abort()
			return
		}
		usert_id := primitive.ObjectIDFromHex(user_id)
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx,bson.D{primitive.E{Key:"_id",Value:usert_id}}).Decode(&filledcart)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(500,"not found")
			return 
		}
		filter_match := bson.D{{Key:"$match"},Value:bson.D{primitive.E{Key:"_id",usert_id}}}
		unwind := bson.D{{Key:"$unwind" , Value:bson.D{primitive.E{Key:"path",Value:"$usercart"}}}}
		grouping := bson.D{{Key:"$group"},Value:bson.D{primitive.E{Key:"_id",Value:"$_id"},{Key:"$sum",Value:"$usercart.price"}}}
		poincursor,err := UserCollection.Aggregate(ctx,mogo.Pipeline(filter_match,unwind,grouping))
		if err != nil{
			log.Println(err)
		}
		var listing []bson.M
		pointcursor.All(ctx,&listing);err !=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _,json := range listing{
			c.IndentedJSON(200,listing["total"])
			c.IndentedJSON(200,filledcart.UserCart)

		}
		ctx.Done()
	}
}
func (app *Application)BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("UserID is empty"))
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx,app.userCollection,userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError , err)
		}
		c.IndentedJSON("successfully placed the order")
	}

}

func (app *Application)InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
			if productQueryID == ""{
				log.Println("product id is empty")
				_ = c.AbortWithError(http.StatusBadRequest,errors.New("product id is empty"))
				return
			}
			userQueryID  := c.Query("userID")
			if userQueryID ==""{
				log.Println("user id is empty")
				_ = c.AbortWithError(http.StatusBadRequest,errors.New("user is is empty"))
				return 
			}
			productID,err := primitive.ObjectIDFromHex(productQueryID)
			if err != nil{
				log.Println(err)
				c.AbortWithError(http.StatusInternalServerError ,errors.New("can't convert productID to primitive type"))
				return 
			}
			
			var ctx,cancel = context.WithTimeout(context.Background(),5 *time.Second)
			defer cancel()
			err = database.InstantBuyer(ctx,app.prodCollection,app.userCollection,productID,userQueryID)
			if err != nil{
				c.IndentedJSON(http.StatusInternalServerError,err)
			}
			c.IndentedJSON(200,"sccessfully placed the order")
	}
}

