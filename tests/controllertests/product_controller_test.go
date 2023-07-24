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

func TestCreateProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, _, err := server.SignIn(users[0].Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	choosedID := users[0].ID
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		brand        string
		price        float64
		image        string
		owner_id     uuid.UUID
		description  string
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			name:         "BUFFET ALMOÇO NA MESA",
			brand:        "Pizza Hut",
			price:        40,
			image:        "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg",
			owner_id:     choosedID,
			description:  "Serviço",
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"name":"", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Brand",
		},
		{
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 0.00, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Price",
		},
		{
			inputJSON:    `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Image",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("Product", "/products", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateProduct)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		// trunk-ignore(golangci-lint/gosimple)
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		fmt.Printf("Errors: %v", v.errorMessage)
		fmt.Printf("Errors Other: %v", responseMap["error"])
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["brand"], v.brand)
			assert.Equal(t, responseMap["price"], v.price)
			assert.Equal(t, responseMap["image"], v.image)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetProducts(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, err := seedUsersAndProducts()
	if err != nil {
		log.Fatal(err)
	}

	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(users[0].Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetProducts)

	req.Header.Set("Authorization", tokenString)

	handler.ServeHTTP(rr, req)

	var products []models.Product
	// trunk-ignore(golangci-lint/gosimple)
	err = json.Unmarshal([]byte(rr.Body.String()), &products)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(products), 1)
}
func TestGetProductByID(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatal(err)
	}
	user, product, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	//Login the Shop and get the authentication token
	token, _, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	productSample := []struct {
		id           uuid.UUID
		statusCode   int
		name         string
		brand        string
		price        float64
		image        string
		owner_id     uuid.UUID
		description  string
		errorMessage string
	}{
		{
			id:          product.ID,
			statusCode:  200,
			name:        product.Name,
			brand:       product.Brand,
			price:       product.Price,
			image:       product.Image,
			description: product.Description,
			owner_id:    product.OwnerID,
		},
		{
			id:         uuid.Nil,
			statusCode: 404,
		},
	}
	for _, v := range productSample {

		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		convertID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": convertID.String()})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetProduct)

		req.Header.Set("Authorization", tokenString)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		fmt.Printf("Error to: %s", responseMap["error"])

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, product.Name, responseMap["name"])
			assert.Equal(t, product.Brand, responseMap["brand"])
			assert.Equal(t, product.Price, responseMap["price"])
			assert.Equal(t, product.Image, responseMap["image"])
			assert.Equal(t, product.Description, responseMap["description"])
		}

	}
}

func TestUpdateProduct(t *testing.T) {

	var ProductUserEmail, ProductUserPassword string
	var AuthProductOwnerID uuid.UUID
	var AuthProductID uuid.UUID

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatal(err)
	}
	users, products, err := seedUsersAndProducts()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	ProductUserEmail = users[0].Email
	ProductUserPassword = "password" //Note the password in the database is already hashed, we want unhashed

	//Login the user and get the authentication token
	token, _, err := server.SignIn(ProductUserEmail, ProductUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first product
	AuthProductID = products[0].ID
	AuthProductOwnerID = products[0].OwnerID

	samples := []struct {
		id           uuid.UUID
		updateJSON   string
		statusCode   int
		name         string
		brand        string
		price        float64
		image        string
		owner_id     uuid.UUID
		description  string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   200,
			name:         "BUFFET ALMOÇO NA MESA",
			brand:        "Pizza Hut",
			price:        40,
			image:        "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg",
			description:  "Serviço",
			owner_id:     AuthProductOwnerID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           AuthProductID,
			updateJSON:   `{"name":"", "brand":"Pizza Hut", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"", "price": 40, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Brand",
		},
		{
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 0.00, "image": "https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Price",
		},
		{
			id:           AuthProductID,
			updateJSON:   `{"name":"BUFFET ALMOÇO NA MESA", "brand":"Pizza Hut", "price": 40, "image": "", "description": "Serviço"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Image",
		},
		{
			id:           uuid.Nil,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/products", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		convertedID := v.id
		req = mux.SetURLVars(req, map[string]string{"id": convertedID.String()})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateProduct)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		fmt.Printf("Error: %s", responseMap["error"])
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["brand"], v.brand)
			assert.Equal(t, responseMap["price"], v.price)
			assert.Equal(t, responseMap["image"], v.image)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteProduct(t *testing.T) {

	var ProductUserEmail, ProductUserPassword string
	var ProductUserID uuid.UUID
	var AuthProductID uuid.UUID

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatal(err)
	}
	users, products, err := seedUsersAndProducts()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	ProductUserEmail = users[1].Email
	ProductUserPassword = "password" //Note the password in the database is already hashed, we want unhashed

	//Login the user and get the authentication token
	token, _, err := server.SignIn(ProductUserEmail, ProductUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	//Get only the first product
	firstProductID := products[0].ID
	firstUserID := products[0].OwnerID

	// Get only the second product
	AuthProductID = products[1].ID
	ProductUserID = products[1].OwnerID

	productSample := []struct {
		id           uuid.UUID
		owner_id     uuid.UUID
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           AuthProductID,
			owner_id:     ProductUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           AuthProductID,
			owner_id:     ProductUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           AuthProductID,
			owner_id:     ProductUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         uuid.Nil,
			tokenGiven: tokenString,
			statusCode: 404,
		},
		{
			id:           firstProductID,
			owner_id:     firstUserID,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range productSample {

		convertID := v.id
		req, _ := http.NewRequest("GET", "/products", nil)
		req = mux.SetURLVars(req, map[string]string{"id": convertID.String()})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteProduct)

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
