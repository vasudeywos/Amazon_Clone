package controllers

import (
	"encoding/json"
	"errors"
	"mymod/models"
	"fmt"
	"github.com/gin-contrib/sessions"
)

type Cart struct {
	UserID      uint
	Items []CartItem `json:"items"`  //NO spaces all else reads error
}

type CartItem struct {
    ProductName string `json:"product_name"`
    Quantity    int    `json:"quantity"`
}

func(crt *Cart)AddItem(addcrt CartItem){
	//var addcrt CartItem
	//If already in Cart
	for i,itm :=range crt.Items{
		if itm.ProductName==addcrt.ProductName{
			crt.Items[i].Quantity+=addcrt.Quantity
		}
	}
	//If not in Cart
	// user,err:=utils.CurrentUser(c)
	// if err!=nil{
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// }
	crt.Items = append(crt.Items,addcrt )
}

func(crt *Cart)DeleteItem(addcrt CartItem){
	for i,itm :=range crt.Items{
		if itm.ProductName==addcrt.ProductName{
			crt.Items=append(crt.Items[:i],crt.Items[i+1:]...)
		}
	}
}

func SaveCart(session sessions.Session,crt *Cart) error{
	userno:=crt.UserID
	cartserialize,_:=json.Marshal(crt) //As not HTTP have to use Marshall
	session.Set(fmt.Sprintf("cart_user%d",userno),string(cartserialize)) //create personalize cart,becoz multiple users
	session.Save()
	return nil
}

func Deserialize(receiv string) (*Cart,error){
	var cart Cart
	err:=json.Unmarshal([]byte(receiv),&cart)
	if err!=nil{
		return nil,err
	}
	return &cart,nil
}

func GetCart(session sessions.Session,user models.User) (*Cart,error){
	name:=fmt.Sprintf("cart_user%d",user.ID)
	data:=session.Get(name)
	if data == nil {
		return &Cart{}, nil
	}

	cartser, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid cart data in session")
	}

	return Deserialize(cartser)
}

func GetCartForUser(session sessions.Session,user models.User) (*Cart, error) {
	name:=fmt.Sprintf("cart_user%d",user.ID) //create personalize cart,becoz multiple users
    cartData := session.Get(name)
    if cartData == nil {
        // If the cart doesn't exist then new cart.
        return &Cart{}, nil
    }

    cartser, ok := cartData.(string)
    if !ok {
        return nil, errors.New("invalid cart data in session")
    }
    return Deserialize(cartser)
}