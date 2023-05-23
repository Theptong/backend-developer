package routers

import (
	v1 "backend-test/v1"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetCollectionRoutes(router *gin.RouterGroup, db *sql.DB) {
	ctrls := v1.DBController{Database: db}

	router.GET("collections", ctrls.GetCollection) // GET
	router.GET("collections/:id", ctrls.GetCollectionById)   // GET BY ID
	router.POST("collections", ctrls.CreateCollection)       // POST
	router.PATCH("collections/:id", ctrls.UpdateCollection)  // PATCH
	router.DELETE("collections/:id", ctrls.DeleteCollection) // DELETE
}
