package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/arikardnoir/asiwaju/api/controllers"
	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var shopInstance = models.Shop{}
var addressInstance = models.Address{}
var collectionPointInstance = models.CollectionPoint{}
var pickupManInstance = models.PickupMan{}
var pickupPointInstance = models.PickupPoint{}
var deliveryManInstance = models.DeliveryMan{}
var orderInstance = models.Order{}
var packageInstance = models.Package{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error
	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

/* Start of Model Test on Shop */

func refreshShopTable() error {
	err := server.DB.DropTableIfExists(&models.Shop{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Shop{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneShop() (models.Shop, error) {

	refreshShopTable()

	shop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "B Eventos",
		Email:       "beventos@gmail.com",
		Password:    "password",
		Description: "Estamos testando",
	}

	err := server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		log.Fatalf("cannot seed shops table: %v", err)
	}
	return shop, nil
}

func seedShops() error {

	shops := []models.Shop{
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Euclides Shop",
			Email:       "euclides.shop@gmail.com",
			Password:    "password",
			Description: "Loja do Euclides",
		},
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Lopes Shopping",
			Email:       "lopes.shopping@gmail.com",
			Password:    "password",
			Description: "Loja dos Lopes",
		},
	}

	for i := range shops {
		err := server.DB.Model(&models.Shop{}).Create(&shops[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/* End of Model Test on Shop */

/* Start of Model Test on Address */

func refreshShopAndAddressTable() error {

	err := server.DB.DropTableIfExists(&models.Shop{}, &models.Address{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Shop{}, &models.Address{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneShopAndOneAddress() (models.Address, error) {

	err := refreshShopAndAddressTable()
	if err != nil {
		return models.Address{}, err
	}
	shop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "Malu Temeperos",
		Email:       "malu.temperos@gmail.com",
		Password:    "password",
		Description: "Negocio da Malu",
	}
	err = server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.Address{}, err
	}
	address := models.Address{
		ID:           uuid.Must(uuid.NewRandom()),
		Country:      "Angola",
		State:        "Luanda",
		City:         "Belas",
		Neighborhood: "Golfe 2",
		Street:       "Santa Teresinha",
		Number:       43433,
		Description:  "Junto a igreja mundial",
		ShopID:       shop.ID,
	}
	err = server.DB.Model(&models.Address{}).Create(&address).Error
	if err != nil {
		return models.Address{}, err
	}
	return address, nil
}

func seedShopsAndAddresses() ([]models.Shop, []models.Address, error) {

	var err error

	if err != nil {
		return []models.Shop{}, []models.Address{}, err
	}
	var shops = []models.Shop{
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Shoezone",
			Email:       "shoezone@gmail.com",
			Password:    "password",
			Description: "A showzone é uma loja de sneakers",
		},
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Shoemania",
			Email:       "shoemania@gmail.com",
			Password:    "password",
			Description: "Showmania é um paraiso de calçados",
		},
	}
	var address = []models.Address{
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Country:      "Angola",
			State:        "Luanda",
			City:         "Belas",
			Neighborhood: "Golfe 2",
			Street:       "Santa Teresinha",
			Number:       34454555,
			Description:  "Testando o Endereço",
		},
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Country:      "Angola",
			State:        "Luanda",
			City:         "Viana",
			Neighborhood: "Condominio Veredea das Flores",
			Street:       "Rua das Bromelias",
			Number:       3487348,
			Description:  "Testando o Endereço",
		},
	}

	for i := range shops {
		err = server.DB.Model(&models.Shop{}).Create(&shops[i]).Error
		if err != nil {
			log.Fatalf("cannot seed shops table: %v", err)
		}
		address[i].ShopID = shops[i].ID

		err = server.DB.Model(&models.Address{}).Create(&address[i]).Error
		if err != nil {
			log.Fatalf("cannot seed address table: %v", err)
		}
	}
	return shops, address, nil
}

/* End of Model Test on Address */

/* Start of Model Test on CollectionPoint */

func refreshCollectionPointTable() error {
	err := server.DB.DropTableIfExists(&models.CollectionPoint{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.CollectionPoint{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneCollectionPoint() (models.CollectionPoint, error) {

	refreshCollectionPointTable()

	collectionPoint := models.CollectionPoint{
		ID:           uuid.Must(uuid.NewRandom()),
		Name:         "Mamadou Golfe 2",
		Instrutions:  "Lado Oposto a Igreja Josafat",
		PhoneNumber:  "+244934445569",
		PhotoName:    "hjdhsfsjndsfjds.jpg",
		City:         "Luanda",
		Neighborhood: "Golfe 2",
		Street:       "Rua Santo Ambrosio",
		Log:          24238573349,
		Lat:          -1234934455,
	}

	err := server.DB.Model(&models.CollectionPoint{}).Create(&collectionPoint).Error
	if err != nil {
		log.Fatalf("cannot seed shops table: %v", err)
	}
	return collectionPoint, nil
}

func seedCollectionPoints() error {

	collectionPoints := []models.CollectionPoint{
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Name:         "Farmacia Luandense",
			Instrutions:  "Junto ao posto de combustivel da Pumangol, Talatona",
			PhoneNumber:  "+244934001569",
			PhotoName:    "hjdhsfsjndsfjds.jpg",
			City:         "Luanda",
			Neighborhood: "Talatona",
			Street:       "Rua Samora Machel",
			Log:          24238573349,
			Lat:          -1234934455,
		},
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Name:         "Salão de Beleza Eloisa",
			Instrutions:  "Junto ao Colégio Ulumbu",
			PhoneNumber:  "+244934445000",
			PhotoName:    "dfgdfdfdfgawr.jpg",
			City:         "Luanda",
			Neighborhood: "Golfe 2",
			Street:       "Rua do Ulumbu",
			Log:          342545435334,
			Lat:          -1234456334,
		},
	}

	for i := range collectionPoints {
		err := server.DB.Model(&models.CollectionPoint{}).Create(&collectionPoints[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/* End of Model Test on CollectionPoint */

/* Start of Model Test on PickupMan */

func refreshPickupMan() error {
	err := server.DB.DropTableIfExists(&models.PickupMan{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.PickupMan{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOnePickupMan() (models.PickupMan, error) {

	refreshPickupMan()

	pickupMan := models.PickupMan{
		ID:          uuid.Must(uuid.NewRandom()),
		FirstName:   "Eloisa",
		LastName:    "Lopes",
		Email:       "eloisa.lopes@gmail.com",
		PhoneNumber: "+244934959569",
		Status:      "available",
	}

	err := server.DB.Model(&models.PickupMan{}).Create(&pickupMan).Error
	if err != nil {
		log.Fatalf("cannot seed shops table: %v", err)
	}
	return pickupMan, nil
}

func seedPickupMen() error {

	pickupMen := []models.PickupMan{
		{
			ID:          uuid.Must(uuid.NewRandom()),
			FirstName:   "Eloisa",
			LastName:    "Lopes",
			Email:       "eloisa.lopes@gmail.com",
			PhoneNumber: "+244934959569",
			Status:      "available",
		},
		{
			ID:          uuid.Must(uuid.NewRandom()),
			FirstName:   "Tiete",
			LastName:    "Madalena",
			Email:       "tiete.madalena@gmail.com",
			PhoneNumber: "+244944919069",
			Status:      "running",
		},
	}

	for i := range pickupMen {
		err := server.DB.Model(&models.PickupMan{}).Create(&pickupMen[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/* End of Model Test on PickupMan */

/* Start of Model Test on PickupPoint */

func refreshShopAndPickupPointTable() error {

	err := server.DB.DropTableIfExists(&models.Shop{}, &models.PickupPoint{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Shop{}, &models.PickupPoint{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneShopAndOnePickupPoint() (models.PickupPoint, error) {

	err := refreshShopAndPickupPointTable()
	if err != nil {
		return models.PickupPoint{}, err
	}
	shop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "Malu Temeperos",
		Email:       "malu.temperos@gmail.com",
		Password:    "password",
		Description: "Negocio da Malu",
	}
	err = server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.PickupPoint{}, err
	}
	pickupPoint := models.PickupPoint{
		ID:                uuid.Must(uuid.NewRandom()),
		Name:              "Mamadou Dialou",
		PickupInstrutions: "Proximo do Hotel Vergas",
		PhoneNumber:       "+244934959569",
		City:              "Luanda",
		Neighborhood:      "Golfe 2",
		Street:            "Rua Santo Antonio",
		ShopID:            shop.ID,
	}
	err = server.DB.Model(&models.PickupPoint{}).Create(&pickupPoint).Error
	if err != nil {
		return models.PickupPoint{}, err
	}
	return pickupPoint, nil
}

func seedShopsAndPickupPoints() ([]models.Shop, []models.PickupPoint, error) {

	var err error

	if err != nil {
		return []models.Shop{}, []models.PickupPoint{}, err
	}
	var shops = []models.Shop{
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Shoezone",
			Email:       "shoezone@gmail.com",
			Password:    "password",
			Description: "A showzone é uma loja de sneakers",
		},
		{
			ID:          uuid.Must(uuid.NewRandom()),
			Name:        "Shoemania",
			Email:       "shoemania@gmail.com",
			Password:    "password",
			Description: "Showmania é um paraiso de calçados",
		},
	}
	var pickupPoints = []models.PickupPoint{
		{
			ID:                uuid.Must(uuid.NewRandom()),
			Name:              "Roberta Narciso",
			PickupInstrutions: "Condominio do BNA",
			PhoneNumber:       "+244934991569",
			City:              "Luanda",
			Neighborhood:      "Bairro dos Militantes",
			Street:            "Rua Santo Antonio",
		},
		{
			ID:                uuid.Must(uuid.NewRandom()),
			Name:              "Nuria Suende",
			PickupInstrutions: "Junto ao mercado Tropical",
			PhoneNumber:       "+244922929283",
			City:              "Luanda",
			Neighborhood:      "Camama",
			Street:            "Rua Centro Tropical",
		},
	}

	for i := range shops {
		err = server.DB.Model(&models.Shop{}).Create(&shops[i]).Error
		if err != nil {
			log.Fatalf("cannot seed shops table: %v", err)
		}
		pickupPoints[i].ShopID = shops[i].ID

		err = server.DB.Model(&models.PickupPoint{}).Create(&pickupPoints[i]).Error
		if err != nil {
			log.Fatalf("cannot seed collection point table: %v", err)
		}
	}
	return shops, pickupPoints, nil
}

/* End of Model Test on PickupPoint */

/* Start of Model Test on DeliveryMan */

func refreshDeliveryMan() error {
	err := server.DB.DropTableIfExists(&models.DeliveryMan{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.DeliveryMan{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneDeliveryMan() (models.DeliveryMan, error) {

	refreshDeliveryMan()

	deliveryMan := models.DeliveryMan{
		ID:          uuid.Must(uuid.NewRandom()),
		FirstName:   "Davidson",
		LastName:    "Bengui",
		Email:       "davidson.bengui@gmail.com",
		PhoneNumber: "+244934595569",
		Status:      "available",
	}

	err := server.DB.Model(&models.DeliveryMan{}).Create(&deliveryMan).Error
	if err != nil {
		log.Fatalf("cannot seed Deliveryman table: %v", err)
	}
	return deliveryMan, nil
}

func seedDeliveryMen() error {

	deliveryMen := []models.DeliveryMan{
		{
			ID:          uuid.Must(uuid.NewRandom()),
			FirstName:   "Pedro",
			LastName:    "Nunes",
			Email:       "pedro.nunes@gmail.com",
			PhoneNumber: "+244934001569",
			Status:      "available",
		},
		{
			ID:          uuid.Must(uuid.NewRandom()),
			FirstName:   "Claudio",
			LastName:    "Lopes",
			Email:       "killah@gmail.com",
			PhoneNumber: "+244922019069",
			Status:      "running",
		},
	}

	for i := range deliveryMen {
		err := server.DB.Model(&models.DeliveryMan{}).Create(&deliveryMen[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/* End of Model Test on DeliveryMan */

/* Start of Model Test on Order */

func refreshOrder() error {
	err := server.DB.DropTableIfExists(&models.Order{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Order{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneOrder() (models.Order, error) {

	refreshOrder()

	order := models.Order{
		ID: uuid.Must(uuid.NewRandom()),
	}

	err := server.DB.Model(&models.Order{}).Create(&order).Error
	if err != nil {
		log.Fatalf("cannot seed Order table: %v", err)
	}
	return order, nil
}

func seedOrders() error {

	orders := []models.Order{
		{
			ID: uuid.Must(uuid.NewRandom()),
		},
		{
			ID: uuid.Must(uuid.NewRandom()),
		},
	}

	for i := range orders {
		err := server.DB.Model(&models.Order{}).Create(&orders[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/* End of Model Test on Order */

/* Start of Model Test on Packages */

func refreshShopCollectionPointPickupPointPackage() error {
	err := server.DB.DropTableIfExists(&models.Shop{}, &models.CollectionPoint{}, &models.PickupPoint{}, &models.Package{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Shop{}, &models.CollectionPoint{}, &models.PickupPoint{}, &models.Package{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneCollectionPointOnePickupPointAndOnePackage() (models.Package, models.CollectionPoint, models.PickupPoint, error) {

	err := refreshShopCollectionPointPickupPointPackage()
	if err != nil {
		return models.Package{}, models.CollectionPoint{}, models.PickupPoint{}, err
	}

	shop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "Pedro Intermediações",
		Email:       "p.intermediações@gmail.com",
		Password:    "password",
		Description: "Intermediações de Automoveis",
	}
	err = server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.Package{}, models.CollectionPoint{}, models.PickupPoint{}, err
	}
	collectionPoint := models.CollectionPoint{
		ID:           uuid.Must(uuid.NewRandom()),
		Name:         "Mamadou Dialou",
		Instrutions:  "Proximo do Hotel Vergas",
		PhoneNumber:  "+244934959569",
		PhotoName:    "hjdhsfsjndsfjds.jpg",
		City:         "Luanda",
		Neighborhood: "Golfe 2",
		Street:       "Rua Santo Antonio",
		Log:          24230857349,
		Lat:          -1234934455,
	}
	err = server.DB.Model(&models.CollectionPoint{}).Create(&collectionPoint).Error
	if err != nil {
		return models.Package{}, models.CollectionPoint{}, models.PickupPoint{}, err
	}
	pickupPoint := models.PickupPoint{
		ID:                uuid.Must(uuid.NewRandom()),
		Name:              "Mamadou Dialou",
		PickupInstrutions: "Proximo do Hotel Vergas",
		PhoneNumber:       "+244934959569",
		City:              "Luanda",
		Neighborhood:      "Golfe 2",
		Street:            "Rua Santo Antonio",
		ShopID:            shop.ID,
	}
	err = server.DB.Model(&models.PickupPoint{}).Create(&pickupPoint).Error
	if err != nil {
		return models.Package{}, models.CollectionPoint{}, models.PickupPoint{}, err
	}
	packageSingle := models.Package{
		ID:                   uuid.Must(uuid.NewRandom()),
		RecipientName:        "Pedro Nunes",
		RecipientPhoneNumber: "+244924988823",
		IdentityCard:         "28780048LA839",
		ChargeValue:          500,
		ChargeMethod:         1, //Method: Money
		DimensionWidth:       10,
		DimensionHeight:      10,
		DimensionLength:      10,
		CollectionPointID:    collectionPoint.ID,
		PickupPointID:        pickupPoint.ID,
		ShopID:               shop.ID,
	}
	err = server.DB.Model(&models.Package{}).Create(&packageSingle).Error
	if err != nil {
		return models.Package{}, models.CollectionPoint{}, models.PickupPoint{}, err
	}
	return packageSingle, collectionPoint, pickupPoint, nil
}

func seedCollectionPointPickupPointPackages() (models.Shop, error) {

	shop := models.Shop{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:        "BOX Office",
		Email:       "box.office@gmail.com",
		Description: "Testando a Loja do BOX",
		Password:    "password",
	}

	collectionPoints := []models.CollectionPoint{
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Name:         "Mamadou Dialou",
			Instrutions:  "Proximo do Hotel Vergas",
			PhoneNumber:  "+244934959569",
			PhotoName:    "hjdhsfsjndsfjds.jpg",
			City:         "Luanda",
			Neighborhood: "Golfe 2",
			Street:       "Rua Santo Antonio",
			Log:          24230857349,
			Lat:          -1234934455,
		},
		{
			ID:           uuid.Must(uuid.NewRandom()),
			Name:         "Farmacia Jaquelina",
			Instrutions:  "De frente a ENDE do Golfe 2",
			PhoneNumber:  "+244985929283",
			PhotoName:    "ijdfhdsfjsdf.jpg",
			City:         "Luanda",
			Neighborhood: "Golfe 2",
			Street:       "Soba Capassa",
			Log:          8958484543,
			Lat:          -348843453,
		},
	}

	pickupPoints := []models.PickupPoint{
		{
			ID:                uuid.Must(uuid.NewRandom()),
			Name:              "Mamadou Dialou",
			PickupInstrutions: "Proximo do Hotel Vergas",
			PhoneNumber:       "+244934959569",
			City:              "Luanda",
			Neighborhood:      "Golfe 2",
			Street:            "Rua Santo Antonio",
			ShopID:            shop.ID,
		},
		{
			ID:                uuid.Must(uuid.NewRandom()),
			Name:              "Farmacia Jaquelina",
			PickupInstrutions: "De frente a ENDE do Golfe 2",
			PhoneNumber:       "+244985929283",
			City:              "Luanda",
			Neighborhood:      "Golfe 2",
			Street:            "Soba Capassa",
			ShopID:            shop.ID,
		},
	}

	packages := []models.Package{
		{
			ID:                   uuid.Must(uuid.NewRandom()),
			RecipientName:        "Pedro Nunes",
			RecipientPhoneNumber: "+244924988823",
			IdentityCard:         "28780048LA839",
			ChargeValue:          500,
			ChargeMethod:         1, //Method: Money
			DimensionWidth:       10,
			DimensionHeight:      10,
			DimensionLength:      10,
			CollectionPointID:    collectionPoints[0].ID,
			PickupPointID:        pickupPoints[0].ID,
			ShopID:               shop.ID,
		},
		{
			ID:                   uuid.Must(uuid.NewRandom()),
			RecipientName:        "Danilson Bengui",
			RecipientPhoneNumber: "+244942343536",
			IdentityCard:         "000783948LA839",
			ChargeValue:          750,
			ChargeMethod:         3, //Method: Debit Card
			DimensionWidth:       20,
			DimensionHeight:      12,
			DimensionLength:      22,
			CollectionPointID:    collectionPoints[1].ID,
			PickupPointID:        pickupPoints[1].ID,
			ShopID:               shop.ID,
		},
	}

	err := server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		log.Fatalf("cannot seed shop table: %v", err)
	}

	for i := range collectionPoints {

		err = server.DB.Model(&models.CollectionPoint{}).Create(&collectionPoints[i]).Error
		if err != nil {
			log.Fatalf("cannot seed collection point table: %v", err)
		}

		pickupPoints[i].ShopID = shop.ID

		err = server.DB.Model(&models.PickupPoint{}).Create(&pickupPoints[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Pickup point table: %v", err)
		}

		packages[i].CollectionPointID = collectionPoints[i].ID
		packages[i].PickupPointID = pickupPoints[i].ID
		packages[i].ShopID = shop.ID

		err = server.DB.Model(&models.Package{}).Create(&packages[i]).Error
		if err != nil {
			log.Fatalf("cannot seed package table: %v", err)
		}
	}
	return shop, nil
}

/* End of Model Test on Packages */
