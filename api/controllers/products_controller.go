package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arikardnoir/asiwaju/api/auth"
	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/arikardnoir/asiwaju/api/responses"
	"github.com/arikardnoir/asiwaju/api/utils/formaterror"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (server *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	product.Prepare()
	err = product.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	oid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if oid != product.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	productCreated, err := product.SaveProduct(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, productCreated.ID))
	responses.JSON(w, http.StatusCreated, productCreated)
}

func (server *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	product := models.Product{}

	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, products)
}

func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	product := models.Product{}

	productReceived, err := product.FindProductByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productReceived)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Product id is valid
	pid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	oid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Product exist
	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Product not found"))
		return
	}

	// If a user attempt to update a Product not belonging to him
	if oid != product.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data Producted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	productUpdate := models.Product{}
	err = json.Unmarshal(body, &productUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if oid != productUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	productUpdate.Prepare()
	err = productUpdate.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate.ID = product.ID //this is important to tell the model the Product id to update, the other update field are set above

	productUpdated, err := productUpdate.UpdateAProduct(server.DB, productUpdate.ID)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, productUpdated)
}

func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Product id given to us?
	pid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	oid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Product exist
	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Product?
	if oid != product.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = product.DeleteAProduct(server.DB, pid, oid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
