package controllers

import(

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

func AddProductToCart(){

}
func RemoveCartItem(){

}
func BuyItemFromCart(){

}
func InstantBuyer(){

}