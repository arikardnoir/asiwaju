package modeltests

import (
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllCollectionPoints(t *testing.T) {

	err := refreshCollectionPointTable()
	if err != nil {
		log.Fatalf("Error refreshing collection Point table %v\n", err)
	}
	err = seedCollectionPoints()
	if err != nil {
		log.Fatalf("Error seeding collection points  table %v\n", err)
	}
	collectionPoints, err := collectionPointInstance.FindAllCollectionPoint(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the collection points: %v\n", err)
		return
	}
	assert.Equal(t, len(*collectionPoints), 2)
}

func TestGetCollectionPointByID(t *testing.T) {
	err := refreshCollectionPointTable()
	if err != nil {
		log.Fatalf("Error refreshing collection Point table %v\n", err)
	}
	collectionPoint, err := seedOneCollectionPoint()
	if err != nil {
		log.Fatalf("Error seeding collection Point table %v\n", err)
	}
	foundCollectionPoints, err := collectionPointInstance.FindCollectionPointByID(server.DB, collectionPoint.ID)
	if err != nil {
		log.Fatalf("This is the error getting one Collection Point %v\n", err)
	}

	assert.Equal(t, foundCollectionPoints.ID, collectionPoint.ID)
	assert.Equal(t, foundCollectionPoints.Name, collectionPoint.Name)
	assert.Equal(t, foundCollectionPoints.Instrutions, collectionPoint.Instrutions)
	assert.Equal(t, foundCollectionPoints.PhoneNumber, collectionPoint.PhoneNumber)
	assert.Equal(t, foundCollectionPoints.PhotoName, collectionPoint.PhotoName)
	assert.Equal(t, foundCollectionPoints.City, collectionPoint.City)
	assert.Equal(t, foundCollectionPoints.Neighborhood, collectionPoint.Neighborhood)
	assert.Equal(t, foundCollectionPoints.Street, collectionPoint.Street)
	assert.Equal(t, foundCollectionPoints.Lat, collectionPoint.Lat)
	assert.Equal(t, foundCollectionPoints.Log, collectionPoint.Log)
}
