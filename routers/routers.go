package routers

import (
	v1 "backend-test/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetCollectionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	ctrls := v1.DBController{Database: db}

	router.GET("collections", ctrls.GetCollection)         // GET
	router.GET("collections/:id", ctrls.GetCollectionById) // GET BY ID
	// router.POST("collections", ctrls.CreateCollection)     // POST
	// router.PATCH("collections", ctrls.UpdateCollection)      // PATCH
	// router.DELETE("collections/:id", ctrls.DeleteCollection) // DELETE
}
