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
	// router.GET("collections/limit/:id", ctrls.GetCollectionByLimit) // GET BY limit
	router.GET("collections/time/:id", ctrls.GetCollectionByTime) // GET BY time
	router.GET("collections/published/:id", ctrls.GetCollectionByPublished) // GET BY time
	// router.GET("collections/between/:id/:subid", ctrls.GetCollectionByBETWEEN) // GET BY time
	router.POST("collections", ctrls.CreateCollection)     // POST
	// router.GET("collections/:id", ctrls.UpdateCollectionByViewCount)      // PATCH
	router.DELETE("collections/:id", ctrls.DeleteCollection) // DELETE
}
