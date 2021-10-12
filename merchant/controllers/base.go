package controllers

import (
	"GOMISPLUS/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	//mysql database driver
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres" //postgres database driver
)

type Server struct {
	PG     *gorm.DB
	Router *mux.Router
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "majoo123"
	dbname   = "merchant"
)

func (server *Server) Initialize() {
	var err error

	psqlconn := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=disable TimeZone=Asia/Jakarta", host, port, user, password, dbname)

	server.PG, err = gorm.Open(postgres.Open(psqlconn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Cannot connect to %s database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database")
	}

	server.PG.AutoMigrate(&models.User{})
	server.PG.AutoMigrate(&models.Product{})

	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(server.Router)))
}
