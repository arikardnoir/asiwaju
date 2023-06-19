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

// Login that make login
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, users, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		fmt.Println(formattedError)
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responseuser := models.SanitizeUser(users)

	response := map[string]interface{}{
		"data":  responseuser,
		"token": token,
	}

	responses.JSON(w, http.StatusOK, response)
}

// SignIn that make sign in
func (server *Server) SignIn(email, password string) (string, models.User, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", user, err
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", user, err
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", user, err
	}

	return token, user, nil
}
