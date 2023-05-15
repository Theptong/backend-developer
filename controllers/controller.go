package controllers

import "backend-test/models"

// "backend-test/models"
// "fmt"
// "net/http"
// "strings"
// // "time"

// "github.com/gin-gonic/gin"
// "github.com/google/uuid"
// "gorm.io/gorm"

type PostsCtrl struct{}

func (pc PostsCtrl) GetListPosts() models.Posts {
	var result models.Posts // ตัวแปล return

	// result = dataList

	return result
}

// type DBController struct {
// 	Database *gorm.DB
// }

// // GET
// func (db *DBController) GetCollection(c *gin.Context) {
// 	_type := c.Query("type")
// 	_where := map[string]interface{}{}

// 	if _type != "" {
// 		_where["type"] = _type
// 	}

// 	var posts []models.Posts
// 	db.Database.Where(_where).Find(&posts)
// 	fmt.Println("posts:", posts)
// 	// for i, _ := range posts {
// 	// 	db.Database.Model(posts[i]).Association("Groups").Find(&posts[i].Groups)
// 	// }

// 	c.JSON(http.StatusOK, gin.H{"results": &posts})
// }

// // GET BY ID
// func (db *DBController) GetCollectionById(c *gin.Context) {
// 	id := c.Param("id")
// 	var posts models.Posts

// 	db.Database.First(&posts, id)
// 	fmt.Println("db.Database.First(&posts, id)::",db.Database.First(&posts, id))
// 	// db.Database.Model(&posts).Association("Groups").Find(&posts.Groups)
// 	if strings.TrimSpace(id) != "" {
// 		c.JSON(http.StatusOK, gin.H{"results": &posts})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"result": "title is required"})
// 	}
// }

// // POST
// func (db *DBController) CreateCollection(c *gin.Context) {
// 	var data models.Posts
// 	// err := c.ShouldBind(&data)
// 	//random Created id
// 	myUUID := uuid.New()
// 	fmt.Println("myUUID::", myUUID.String())
// 	data.Uuid = myUUID.String()
// 	// data.CreatedAt = time.Now()
// 	// time := time.Now().Format("2006-01-02T15:04:05")
// 	// data.UpdatedAt = time.Now()
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		fmt.Println("err::", err)
// 		// logging.Logger(setting.LogLevelSetting.Error, err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
// 		return
// 	}

// 	result := db.Database.Create(&data)
// 	fmt.Println("result::", result)

// 	c.JSON(http.StatusCreated, gin.H{"results": &data})
// }

// // // PATCH
// // func (db *DBController) UpdateCollection(c *gin.Context) {

// // 	var collection models.Collections
// // 	err := c.ShouldBind(&collection)

// // 	result := db.Database.Updates(collection)

// // 	if result.Error != nil || err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"meassage": "Bad request."})
// // 	} else {
// // 		c.JSON(http.StatusOK, gin.H{"results": &collection})
// // 	}
// // }

// // // DELETE
// // func (db *DBController) DeleteCollection(c *gin.Context) {
// // 	id := c.Param("id")
// // 	var collections models.Collections
// // 	db.Database.Delete(&collections, id)

// // 	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK})
// // }
