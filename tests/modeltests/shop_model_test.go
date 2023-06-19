package modeltests

import (
	"log"
	"testing"

	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllShops(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatalf("Error refreshing shop table %v\n", err)
	}
	err = seedShops()
	if err != nil {
		log.Fatalf("Error seeding shops table %v\n", err)
	}
	shops, err := shopInstance.FindAllShops(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the shops: %v\n", err)
		return
	}
	assert.Equal(t, len(*shops), 2)
}

func TestSaveShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatalf("Error shops refreshing table %v\n", err)
	}

	newShop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "Shopping Teste",
		Email:       "shopping.teste@gmail.com",
		Description: "Shopping de teste",
	}
	savedShop, err := newShop.SaveShop(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the shop: %v\n", err)
		return
	}
	assert.Equal(t, newShop.ID, savedShop.ID)
	assert.Equal(t, newShop.Name, savedShop.Name)
	assert.Equal(t, newShop.Email, savedShop.Email)
	assert.Equal(t, newShop.Description, savedShop.Description)

}

func TestGetShopByID(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatalf("Error refreshing shop table: %v\n", err)
	}
	shop, err := seedOneShop()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundShop, err := shopInstance.FindShopByID(server.DB, shop.ID)
	if err != nil {
		t.Errorf("this is the error getting one shop: %v\n", err)
		return
	}
	assert.Equal(t, foundShop.ID, shop.ID)
	assert.Equal(t, foundShop.Name, shop.Name)
	assert.Equal(t, foundShop.Email, shop.Email)
	assert.Equal(t, foundShop.Description, shop.Description)
}

func TestUpdateAShop(t *testing.T) {

	err := refreshShopTable()
	if err != nil {
		log.Fatalf("Error refreshing shop table: %v\n", err)
	}
	shop, err := seedOneShop()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	shopUpdate := models.Shop{
		ID:          shop.ID,
		Name:        "Shopping Luanda",
		Email:       "shopping.luanda@gmail.com",
		Description: "Estamos testando esta opcao",
	}
	updatedShop, err := shopUpdate.UpdateAShop(server.DB, shop.ID)
	if err != nil {
		t.Errorf("this is the error updating the shop: %v\n", err)
		return
	}
	assert.Equal(t, updatedShop.ID, shopUpdate.ID)
	assert.Equal(t, updatedShop.Name, shopUpdate.Name)
	assert.Equal(t, updatedShop.Email, shopUpdate.Email)
	assert.Equal(t, updatedShop.Description, shopUpdate.Description)
}

func TestDeleteAShop(t *testing.T) {

	err := refreshShopAndAddressTable()
	if err != nil {
		log.Fatalf("Error refreshing shop table: %v\n", err)
	}
	shop, err := seedOneShop()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := shopInstance.DeleteAShop(server.DB, shop.ID)
	if err != nil {
		t.Errorf("this is the error deleting the shop: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
