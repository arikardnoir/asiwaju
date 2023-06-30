
package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/arikardnoir/asiwaju/api/controllers"
	"github.com/arikardnoir/asiwaju/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var productInstance = models.Product{}

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

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		ID:       uuid.Must(uuid.NewRandom()),
		Fullname: "Kayla Maziano",
		Nickname: "kayla.maziano",
		Email:    "kay.maziano@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			ID:       uuid.Must(uuid.NewRandom()),
			Fullname: "Benoniel Correia",
			Nickname: "benoniel.correia",
			Email:    "ben.correia@gmail.com",
			Password: "password",
		},
		models.User{
			ID:       uuid.Must(uuid.NewRandom()),
			Fullname: "Aniel Lopes",
			Nickname: "aniel.lopes",
			Email:    "aniel.lopes@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndProductTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Product{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Product{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneProduct() (models.User,models.Product, error) {

	err := refreshUserAndProductTable()
	if err != nil {
		return models.User{}, models.Product{}, err
	}
	user := models.User{
		ID:       uuid.Must(uuid.NewRandom()),
		Fullname: "Aaron Lopes",
		Nickname: "aaron.lopes",
		Email:    "aaron.lopes@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, models.Product{}, err
	}
	product := models.Product{
		ID:       uuid.Must(uuid.NewRandom()),
		Name:   			"BUFFET ALMOÇO NA MESA",
		Brand:  			"Pizza Hut",
		Price:  			40,
		Image:  			"https://www.pizzahut.pt/wp-content/uploads/BUFFET_ALMOCO_na_mesa_8_95_30_junho-scaled.jpg",
		OwnerID: 			user.ID,
		Description:  "Serviço exclusivo para consumo no serviço à mesa dos restaurantes aderentes. Válido de 2.ª a 6.ª feira, das 12:00h às 16:00h, exceto feriados. Imagens ilustrativas. IVA incluído à taxa legal em vigor.",
	}
	err = server.DB.Model(&models.Product{}).Create(&product).Error
	if err != nil {
		return models.User{}, models.Product{}, err
	}
	return user,product, nil
}

func seedUsersAndProducts() ([]models.User, []models.Product, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Product{}, err
	}
	var users = []models.User{
		models.User{
			ID:       uuid.Must(uuid.NewRandom()),
			Fullname: "Eloisa Lopes",
			Nickname: "eloisa.lopes",
			Email:    "eloisa.lopes@gmail.com",
			Password: "password",
		},
		models.User{
			ID:       uuid.Must(uuid.NewRandom()),
			Fullname: "Zoe Maziano",
			Nickname: "zoe.maziano",
			Email:    "zozo.maziano@gmail.com",
			Password: "password",
		},
	}
	var products = []models.Product{
		models.Product{
			ID:       uuid.Must(uuid.NewRandom()),
			Name:    			"Shawarma de Frango",
			Brand:  			"Alchaer Restaurante",
			Price:				31.12,
			Image:  			"https://images.rappi.com.br/products/06c0a5c9-9db5-4af9-b86b-da49927fb673-1673533770540.png?e=webp&d=511x511&q=85",
			OwnerID: 			users[0].ID,
			Description:  "Pão sírio assado na hora com peito de frango, picles, batata frita, pasta de alho e molho de romã.",
		},
		models.Product{
			ID:       uuid.Must(uuid.NewRandom()),
			Name:    "MacBook Pro 13”",
			Brand:  "Apple Inc.",
			Price: 1299,
			Image:  "https://www.apple.com/v/macbook-pro/ah/images/overview/hero_13__d1tfa5zby7e6_large.jpg",
			OwnerID: users[1].ID,
			Description:  "The new M2 chip makes the 13‑inch MacBook Pro more capable than ever. The same compact design supports up to 20 hours of battery life1 and an active cooling system to sustain enhanced performance. Featuring a brilliant Retina display, a FaceTime HD camera, and studio‑quality mics, it’s our most portable pro laptop.",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		products[i].OwnerID = users[i].ID

		err = server.DB.Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed products table: %v", err)
		}
	}
	return users, products, nil
}
