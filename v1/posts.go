package v1

import (
	"backend-test/models"
	"fmt"
	"net/http"

	// "strings"

	// "strings"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBController struct {
	Database *gorm.DB
}

// GET
func (db *DBController) GetCollection(c *gin.Context) {
	// _type := c.Query("type")
	_where := map[string]interface{}{}

	// if _type != "" {
	// 	_where["type"] = _type
	// }

	var posts []models.Posts
	db.Database.Where(_where).Find(&posts)
	// fmt.Println("posts:", posts)

	c.JSON(http.StatusOK, gin.H{"results": &posts})
}

var filterpostrs = map[string]string{
	"title":     "title",
	"content":   "content",
	"published": "published",
}

// GET BY ID
func (db *DBController) GetCollectionById(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {
			if obj.Title == _type { //ถ้ามีtitleตรงกันให้แสดง
				_where["title"] = _type
			} else if obj.Content == _type { //ถ้ามีcontentตรงกันให้แสดง
				_where["content"] = _type
			}
		}
		// 	if _where["title"] == nil {
		// 		_where["content"] = _type
		// 	}else if _where["content"] == nil {
		// 		_where["published"] = _type
		// 	}else{
		// 		_where["title"] = _type
		// 	}
		fmt.Println("_where[]::", _where["title"])
		fmt.Println("_where[]::", _where["content"])
	}
	db.Database.Where(_where).Find(&posts)
	//SELECT * from posts where ?? =
	var dataPosts models.Posts
	if posts[0].Id != uuid.Nil {
		dataPosts.Id = posts[0].Id
	}
	if posts[0].Title != "" {
		dataPosts.Title = posts[0].Title
	}
	if posts[0].Content != ""{
		dataPosts.Content = posts[0].Content
	}
	
	// if posts[0].Published != {}
	dataPosts.Published = posts[0].Published
	dataPosts.ViewCount = posts[0].ViewCount
	dataPosts.CreatedAt = posts[0].CreatedAt
	dataPosts.UpdatedAt = posts[0].UpdatedAt
	fmt.Println("_where",_where)
	if _where["title"] != nil || _where["content"] != nil {
		c.JSON(http.StatusOK, &dataPosts)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}

// POST
func (db *DBController) CreateCollection(c *gin.Context) {
	var data models.Posts
	// err := c.ShouldBind(&data)
	//random Created id
	myUUID := uuid.New()
	fmt.Println("myUUID::", myUUID.String())
	data.Id = myUUID
	// data.Uuid = myUUID.String()
	// data.CreatedAt = time.Now()
	// time := time.Now().Format("2006-01-02T15:04:05")
	// data.UpdatedAt = time.Now()
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("err::", err)
		// logging.Logger(setting.LogLevelSetting.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	result := db.Database.Create(&data)
	fmt.Println("result::", result)

	c.JSON(http.StatusCreated, gin.H{"results": &data})
}
