package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		id           uuid.UUID
		name         string
		email        string
		description  string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Cartomante Store", "email": "cartomante.store@gmail.com", "description":"teste kdsflksdfsdfdsfde 0000sdfsdf", "password": "password"}`,
			statusCode:   201,
			name:         "Cartomante Store",
			email:        "cartomante.store@gmail.com",
			description:  "teste kdsflksdfsdfdsfde 0000sdfsdf",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"So Sneakers", "email": "cartomante.store@gmail.com", "description":"teste de ljdalsf",  "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"name":"So Sneakers", "email": "cartomante.store.com", "description":"teste de 1sdfsdf234",  "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"name":"", "email": "cartomante.store.com", "description":"teste de sdfsdsd",  "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name":"So Sneakers", "email": "", "description":"teste de sdfsdfs",  "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"name":"So Sneakers", "email": "sosneakers@gmail.com", "description":"teste de sdfsdfsdfsdfsfsdf", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/shops", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateShop)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["email"], v.email)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetShops(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedShops()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/shops", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetShops)
	handler.ServeHTTP(rr, req)

	var shops []models.Shop
	err = json.Unmarshal([]byte(rr.Body.String()), &shops)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(shops), 2)
}

func TestGetShopByID(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}
	shop, err := seedOneShop()
	if err != nil {
		log.Fatal(err)
	}
	shopSample := []struct {
		id           uuid.UUID
		statusCode   int
		name         string
		description  string
		email        string
		errorMessage string
	}{
		{
			id:         shop.ID,
			statusCode: 200,
			name:       shop.Name,
			email:      shop.Email,
		},
		{
			id:         uuid.Nil,
			statusCode: 400,
		},
	}
	for _, v := range shopSample {

		req, err := http.NewRequest("GET", "/shops", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		convertID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": convertID.String()})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetShop)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, shop.Name, responseMap["name"])
			assert.Equal(t, shop.Email, responseMap["email"])
			assert.Equal(t, shop.Description, responseMap["description"])
		}
	}
}

func TestUpdateShop(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID, otherID uuid.UUID

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}
	Shops, err := seedShops() //we need atleast two Shops to properly check the update
	if err != nil {
		log.Fatalf("Error seeding shop: %v\n", err)
	}
	// Get only the first Shop
	for _, Shop := range Shops {
		if Shop.Email == "farmacia.central@gmail.com" {
			otherID = Shop.ID
			continue
		}
		AuthID = Shop.ID
		AuthEmail = Shop.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// fmt.Println(tokenString)

	samples := []struct {
		id                uuid.UUID
		updateJSON        string
		statusCode        int
		updateName        string
		updateEmail       string
		updateDescription string
		tokenGiven        string
		errorMessage      string
	}{
		{
			// Convert int32 to int first before converting to string
			id:                AuthID,
			updateJSON:        `{"name":"Subway Angola", "email":"subway.angola@gmail.com", "description":"Melhor Sandes do mundo", "password":"password"}`,
			statusCode:        200,
			updateName:        "Subway Angola",
			updateEmail:       "subway.angola@gmail.com",
			updateDescription: "Melhor Sandes do mundo",
			tokenGiven:        tokenString,
			errorMessage:      "",
		},
		{
			// When password field is empty
			id:           AuthID,
			updateJSON:   `{"name":"Farmacia Central", "email":"farmacia.central@gmail.com", "description":"Testando essa cena", "password":""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		{
			// When no token was passed
			id:           AuthID,
			updateJSON:   `{"name":"Man Mozart", "email":"man.mozart@gmail.com", "description":"Testando essa cena1", "password":"password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           AuthID,
			updateJSON:   `{"name":"Woman Mozart", "email":"woman.mozart@gmail.com", "description":"Testando essa cena", "password":"password"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "farmacia.central@gmail.com" belongs to Shop 2
			id:           AuthID,
			updateJSON:   `{"name":"Frank Lopes", "email":"farmacia.central@gmail.com", "description":"Testando essa cena12", "password":"password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			id:           AuthID,
			updateJSON:   `{"name":"Kan", "email":"kan.mozartgmail.com", "description":"Testando essa cena21", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id:           AuthID,
			updateJSON:   `{"name":"", "email":"sir.mozart@gmail.com", "description":"Testando essa cena", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			id:           AuthID,
			updateJSON:   `{"name":"Rumi Mozart", "email":"", "description":"Testando essa cena", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Email",
		},
		{
			// When Shop 2 is using Shop 1 token
			id:           otherID,
			updateJSON:   `{"name":"Mike", "email": "mike.mozart@gmail.com", "description":"Testando essa cenasddds", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("PUT", "/shops", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		updateConvertID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": updateConvertID.String()})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateShop)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
			assert.Equal(t, responseMap["email"], v.updateEmail)
			assert.Equal(t, responseMap["description"], v.updateDescription)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteShop(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID, otherID uuid.UUID

	err := refreshShopTable()
	if err != nil {
		log.Fatal(err)
	}

	Shops, err := seedShops() //we need atleast two Shops to properly check the update
	if err != nil {
		log.Fatalf("Error seeding Shop: %v\n", err)
	}
	// Get only the first and log him in
	for _, Shop := range Shops {
		if Shop.Email == "farmacia.central@gmail.com" {
			otherID = Shop.ID
			continue
		}
		AuthID = Shop.ID
		AuthEmail = Shop.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}
	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	ShopSample := []struct {
		id           uuid.UUID
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           AuthID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           AuthID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           AuthID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		// {
		// 	id:         uuid.Nil,
		// 	tokenGiven: tokenString,
		// 	statusCode: 400,
		// },
		{
			// Shop 2 trying to use Shop 1 token
			id:           otherID,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range ShopSample {

		req, err := http.NewRequest("GET", "/shops", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		convertID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": convertID.String()})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteShop)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
