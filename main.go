package main

import (
	"backend-test/routers"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	PORT := "8080"

	db := Connect
	router := gin.New()
	api := router.Group("/api")
	// fmt.Println("rows::", rows)
	routers.SetCollectionRoutes(api, db())
	port := fmt.Sprintf(":%v", PORT)
	fmt.Println("Server Running on Port", port)
	http.ListenAndServe(port, router)
}

func Connect() *sql.DB {
	psqlInfo := "postgres://posts:38S2GPNZut4Tmvan@dev.opensource-technology.com:5523/posts?sslmode=disable"
	// connectionStr := "postgres://postgres:123456@localhost:54321/blog?sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}
