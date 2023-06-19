package controllertests

import (
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

func TestGetCollectionPoints(t *testing.T) {
	var AuthEmail, AuthPassword string

	err := refreshCollectionPointTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedCollectionPoints()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshShopTable()
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
			continue
		}
		AuthEmail = Shop.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}

	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest("GET", "/collection-points", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetCollectionPoints)
	//Set the Token to request header
	req.Header.Set("Authorization", tokenString)
	//Sending request
	handler.ServeHTTP(rr, req)

	var collectionPoints []models.CollectionPoint
	err = json.Unmarshal([]byte(rr.Body.String()), &collectionPoints)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(collectionPoints), 2)
}

func TestGetCollectionPointByID(t *testing.T) {
	var AuthEmail, AuthPassword string

	err := refreshCollectionPointTable()
	if err != nil {
		log.Fatal(err)
	}

	collectionPoint, err := seedOneCollectionPoint()
	if err != nil {
		log.Fatal(err)
	}

	err = refreshShopTable()
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
			continue
		}
		AuthEmail = Shop.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}

	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	collectionPointSample := []struct {
		id           uuid.UUID
		statusCode   int
		tokenGiven   string
		name         string
		instructions string
		phoneNumber  string
		photoName    string
		city         string
		neighborhood string
		street       string
		log          float64
		lat          float64
		errorMessage string
	}{
		{

			id:           collectionPoint.ID,
			statusCode:   200,
			name:         collectionPoint.Name,
			instructions: collectionPoint.Instrutions,
			phoneNumber:  collectionPoint.PhoneNumber,
			photoName:    collectionPoint.PhotoName,
			city:         collectionPoint.City,
			neighborhood: collectionPoint.Neighborhood,
			street:       collectionPoint.Street,
			log:          collectionPoint.Log,
			lat:          collectionPoint.Lat,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:         uuid.Nil,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			id:           collectionPoint.ID,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			id:           collectionPoint.ID,
			statusCode:   401,
			tokenGiven:   "Incorrect Token",
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range collectionPointSample {

		req, err := http.NewRequest("GET", "/collection-points", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		convertID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": convertID.String()})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetCollectionPoint)
		//Set the Token to request header
		req.Header.Set("Authorization", v.tokenGiven)
		//Sending request
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, collectionPoint.Name, responseMap["name"])
			assert.Equal(t, collectionPoint.Instrutions, responseMap["pickup_instructions"])
			assert.Equal(t, collectionPoint.PhoneNumber, responseMap["phone_number"])
			assert.Equal(t, collectionPoint.PhotoName, responseMap["photo_name"])
			assert.Equal(t, collectionPoint.City, responseMap["city"])
			assert.Equal(t, collectionPoint.Neighborhood, responseMap["neighborhood"])
			assert.Equal(t, collectionPoint.Street, responseMap["street"])
			assert.Equal(t, collectionPoint.Log, responseMap["log"])
			assert.Equal(t, collectionPoint.Lat, responseMap["lat"])
		}
		if v.statusCode == 401 {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
