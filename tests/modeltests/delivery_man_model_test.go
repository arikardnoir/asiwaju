package modeltests

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllDeliveryMen(t *testing.T) {

	err := refreshDeliveryMan()
	if err != nil {
		log.Fatalf("Error refreshing delivery man table %v\n", err)
	}
	err = seedDeliveryMen()
	if err != nil {
		log.Fatalf("Error seeding delivery man  table %v\n", err)
	}
	deliveryMen, err := deliveryManInstance.FindAllDeliveryMan(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the delivery men: %v\n", err)
		return
	}
	assert.Equal(t, len(*deliveryMen), 2)
}

func TestGetDeliveryManByID(t *testing.T) {
	err := refreshDeliveryMan()
	if err != nil {
		log.Fatalf("Error refreshing delivery man table %v\n", err)
	}
	deliveryMan, err := seedOneDeliveryMan()
	if err != nil {
		log.Fatalf("Error seeding delivery man table %v\n", err)
	}
	foundDeliveryMan, err := deliveryManInstance.FindDeliveryManByID(server.DB, deliveryMan.ID)
	if err != nil {
		log.Fatalf("This is the error getting one delivery man %v\n", err)
	}

	assert.Equal(t, foundDeliveryMan.ID, deliveryMan.ID)
	assert.Equal(t, foundDeliveryMan.FirstName, deliveryMan.FirstName)
	assert.Equal(t, foundDeliveryMan.LastName, deliveryMan.LastName)
	assert.Equal(t, foundDeliveryMan.PhoneNumber, deliveryMan.PhoneNumber)
	assert.Equal(t, foundDeliveryMan.Email, deliveryMan.Email)
	assert.Equal(t, foundDeliveryMan.Status, deliveryMan.Status)
}
