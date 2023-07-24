package api

import (
	"fmt"
	"log"
	"os"

	"github.com/arikardnoir/asiwaju/api/controllers"
	"github.com/arikardnoir/asiwaju/api/seed"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

//Run the server
func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	//server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	server.Initialize("postgres", "mjanseeypnboap", "57043a3e1e4a5c356ec48186db207e14255c14920766e2291c5d8098f3cb86e6", "5432", "ec2-54-84-182-168.compute-1.amazonaws.com", "d2nnv4kopse6e")

	seed.Load(server.DB)

	server.Run(GetPort())
}

//GetPort is getting the port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "5000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
