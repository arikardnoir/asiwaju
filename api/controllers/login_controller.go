package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arikardnoir/asiwaju/api/auth"
	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/arikardnoir/asiwaju/api/responses"
	"github.com/arikardnoir/asiwaju/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

//LoginResponse is the response after login
type LoginResponse struct {
	token string
	shop  models.Shop
}

//Login that make login
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	shop := models.Shop{}
	err = json.Unmarshal(body, &shop)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	shop.Prepare()
	err = shop.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, shops, err := server.SignIn(shop.Email, shop.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		fmt.Println(formattedError)
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responseShop := models.SanitizeShop(shops)

	response := map[string]interface{}{
		"data":  responseShop,
		"token": token,
	}

	responses.JSON(w, http.StatusOK, response)
}

//SignIn that make sign in
func (server *Server) SignIn(email, password string) (string, models.Shop, error) {

	var err error

	shop := models.Shop{}

	err = server.DB.Debug().Model(models.Shop{}).Where("email = ?", email).Take(&shop).Error
	if err != nil {
		return "", shop, err
	}

	err = models.VerifyPassword(shop.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", shop, err
	}

	token, err := auth.CreateToken(shop.ID)
	if err != nil {
		return "", shop, err
	}

	return token, shop, nil
}
