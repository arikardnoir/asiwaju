package seed

import (
	"log"
	"github.com/google/uuid"
	"github.com/arikardnoir/asiwaju/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		ID:          uuid.Must(uuid.NewRandom()),
		Fullname: "Eloisa da Silva",
		Nickname: "elo.silva",
		Email:    "eloisa@gmail.com",
		Password: "password",
	},
	models.User{
		ID:          uuid.Must(uuid.NewRandom()),
		Fullname: "Adriel Van-dunem",
		Nickname: "adri.van",
		Email:    "adriel@gmail.com",
		Password: "password",
	},
}

var products = []models.Product{
	models.Product{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:   "AIR JORDAN 1 RETRO LOW OG EX",
		Brand: "Nike",
		Image:   "https://www.nike.com.br/air-jordan-1-retro-low-og-ex-023577.html?cor=ID#pid1",
		Size:   "38;39;40;41;42",
		Model:   "Air Jordan 1",
		Price:   90.000,
		Description: "Chame-o de obra-prima inacabada. Esta versão trabalhada do AJ1 Low tem tudo a ver com bordas expostas e desgastadas, trazendo uma estética desconstruída para seu têni favorito.",
	},
	models.Product{
		ID:          uuid.Must(uuid.NewRandom()),
		Name:   "Gomes Da Costa Atum Sólido em Óleo Delivery",
		Brand:  "Gomes Da Costa",
		Image:  "https://images.rappi.com.br/products/630151d9-9e33-460a-bed0-4cc527424a74.jpg?d=128x128&e=webp&q=70",
		Size:   "",
		Model:  "",
		Price:  5.000,
		Description: "Produzido com o lombo do atum, a parte mais nobre do peixe, e por isso é muito valorizado pela sua qualidade e sabor diferenciado.",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Product{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Product{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		products[i].OwnerID = users[i].ID

		err = db.Debug().Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed products table: %v", err)
		}
	}
}