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
func (app *Application) AddToCart() gin.Handler{
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

}

func (app *Application)GetItemFromCart() gin.HandlerFunc{

}
func (app *Application)BuyFromCart() gin.HandlerFunc{


}

func (app *Application)InstantBuy() gin.HandlerFunc{

}

