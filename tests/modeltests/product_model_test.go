
package modeltests

import (
	"log"
	"testing"
	"github.com/google/uuid"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/arikardnoir/asiwaju/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllProducts(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and product table %v\n", err)
	}
	_, _, err = seedUsersAndProducts()
	if err != nil {
		log.Fatalf("Error seeding user and product  table %v\n", err)
	}
	products, err := productInstance.FindAllOpenProducts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the products: %v\n", err)
		return
	}
	assert.Equal(t, len(*products), 2)
}

func TestSaveProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error user and product refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newProduct := models.Product{
		ID:       uuid.Must(uuid.NewRandom()),
		Name:   			"BUFFET ALMOÇO NA MESA",
		Brand:  			"Pizza Hut",
		Price:  			40,
		Image:  			"https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg",
		OwnerID: 			user.ID,
		Description:  "Serviço exclusivo para consumo no serviço à mesa dos restaurantes aderentes. Válido de 2.ª a 6.ª feira, das 12:00h às 16:00h, exceto feriados. Imagens ilustrativas. IVA incluído à taxa legal em vigor.",
	}
	savedProduct, err := newProduct.SaveProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the product: %v\n", err)
		return
	}
	assert.Equal(t, newProduct.ID, savedProduct.ID)
	assert.Equal(t, newProduct.Name, savedProduct.Name)
	assert.Equal(t, newProduct.Brand, savedProduct.Brand)
	assert.Equal(t, newProduct.Price, savedProduct.Price)
	assert.Equal(t, newProduct.Description, savedProduct.Description)
	assert.Equal(t, newProduct.OwnerID, savedProduct.OwnerID)

}

func TestGetProductByID(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and product table: %v\n", err)
	}
	product, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundProduct, err := productInstance.FindProductByID(server.DB, product.ID, product.OwnerID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundProduct.ID, product.ID)
	assert.Equal(t, foundProduct.Name, product.Name)
	assert.Equal(t, foundProduct.Brand, product.Brand)
	assert.Equal(t, foundProduct.Price, product.Price)
	assert.Equal(t, foundProduct.Description, product.Description)
}

func TestUpdateAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and product table: %v\n", err)
	}
	product, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	productUpdate := models.Product{
		ID:       		product.ID,
		Name:    			"Shawarma de Frango",
		Brand:  			"Alchaer Restaurante",
		Price:				31.12,
		Image:  			"https://images.rappi.com.br/products/06c0a5c9-9db5-4af9-b86b-da49927fb673-1673533770540.png?e=webp&d=511x511&q=85",
		OwnerID: 			product.OwnerID,
		Description:  "Pão sírio assado na hora com peito de frango, picles, batata frita, pasta de alho e molho de romã.",
	
	}
	updatedProduct, err := productUpdate.UpdateAProduct(server.DB, product.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, updatedProduct.ID, productUpdate.ID)
	assert.Equal(t, updatedProduct.Name, productUpdate.Name)
	assert.Equal(t, updatedProduct.Brand, productUpdate.Brand)
	assert.Equal(t, updatedProduct.Price, productUpdate.Price)
	assert.Equal(t, updatedProduct.Description, productUpdate.Description)
	assert.Equal(t, updatedProduct.OwnerID, productUpdate.OwnerID)
}

func TestDeleteAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and product table: %v\n", err)
	}
	product, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := productInstance.DeleteAProduct(server.DB, product.ID, product.OwnerID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
