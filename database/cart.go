package database

import(
	"errors"
	"log"
	"context"
)

var(
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodeProduct = errors.New("can't find the product")
	ErrorIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser= errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from cart")
	ErrCantGetItem= errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem= errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context ,prodCollection *mongo.Collection , userCollection *mongo.Collection , productID primitive.ObjectID ,userID string)error {
	searchfromdb,err := prodCollection.Find(ctx,bson.M{"_id":productID})
	if err != nil{
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart [] models.ProductUser
	err = searchfromdb.All(ctx,&productCart)
	if err != nil{
		log.Println(err)
		return ErrCantDecodeProduct
	}
	id ,err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	filter :=bson.D{primive.E{Key:"_id",Value: id}}
	update :=bson.D{{Key:"$push",Value:bson.D{primive.E{Key:"usercart",Value:bson.D{Key:"$each",Value:productCart}}}}}
	_,err = userCollection.UpdateOne(ctx,filter,update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}
func RemoveCartItem(){

}
func BuyItemFromCart(){

}
func InstantBuyer(){

}