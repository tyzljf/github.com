package main

import (
	"os"
	"log"
	"net/http"
	"github.com/golang/standard-rest-api/utils/database"
	"github.com/golang/standard-rest-api/utils/caching"
	"github.com/golang/standard-rest-api/controllers"
	"github.com/golang/standard-rest-api/routers"
)

func main() {
	db, err := database.Connect(os.Getenv("PGUSER"), os.Getenv("PGPASS"), os.Getenv("PGDB"), os.Getenv("PGHOST"), os.Getenv("PGPORT"))
	if err != nil {
		log.Fatal(err)
	}
	cache := &caching.Redis{
		Client: caching.Connect(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0),
	}

	userController := controllers.NewUserController(db, cache)
	jobController := controllers.NewJobController(db, cache)

	mux := http.NewServeMux()
	routers.CreateRouters(mux, userController, jobController)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
