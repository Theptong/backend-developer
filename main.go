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

	


	DB_HOST := "localhost"
	DB_NAME := "blog"
	DB_USER := "postgres"
	DB_PORT := "54321"
	DB_PASSWORD := "123456"
	PORT := "8080"

	fmt.Println(":PORT:",PORT)
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", DB_HOST, DB_USER, DB_NAME, DB_PORT, DB_PASSWORD)
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
