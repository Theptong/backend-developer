package main

import (
	"backend-test/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

// postgres: //postgres:123456@localhost:54321/blog?sslmode=disable

// 	DB_HOST = localhost
// 	DB_NAME = blog
// 	DB_USER = postgres
// 	DB_PORT = 54321
// 	DB_PASSWORD = 123456

	// DB_HOST := "@dev.opensource-technology.com"
	// DB_NAME := "blog"
	// DB_USER := "posts"
	// DB_PORT := "5523"
	// DB_PASSWORD := "38S2GPNZut4Tmvan"
	PORT := "8080"

	// fmt.Println(":PORT:", PORT)
	// psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", DB_HOST, DB_USER, DB_NAME, DB_PORT, DB_PASSWORD)
	// fmt.Println(":host:", "host=%s",DB_HOST)
	// fmt.Println(":user:", "user=%s",DB_USER)
	// fmt.Println(":dbname:", "dbname=%s",DB_NAME)
	// fmt.Println(":port:", "port=%s",DB_PORT)
	// fmt.Println(":password:", "password=%s",DB_PASSWORD)
	psqlInfo := "postgres://posts:38S2GPNZut4Tmvan@dev.opensource-technology.com:5523/posts?sslmode=disable"
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	router := gin.New()
	api := router.Group("/api")
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "HELLO GOLANG RESTFUL API.",
		})
	})
	//เข้า Routes
	routers.SetCollectionRoutes(api, db)

	port := fmt.Sprintf(":%v", PORT)
	fmt.Println("Server Running on Port", port)
	http.ListenAndServe(port, router)
}
