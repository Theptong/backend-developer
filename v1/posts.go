package v1

import (
	"backend-test/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "strconv"

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
	filter := make(map[string]interface{})
	if c.Request.URL.Query().Get("limit") != "" {
		filter["limit"] = c.Request.URL.Query().Get("limit")
	}
	if c.Request.URL.Query().Get("page") != "" {
		filter["page"] = c.Request.URL.Query().Get("page")
	}
	fmt.Println("filter[limit]", filter["limit"])
	fmt.Println("filter[page]", filter["page"])
	if filter["limit"] != nil && filter["page"] != nil {
		Limit, _ := strconv.Atoi(fmt.Sprint(filter["limit"]))
		Page, _ := strconv.Atoi(fmt.Sprint(filter["page"]))
		Offset := 0
		if Page >= 0 {
			Offset = (Page - 1) * Limit
		} else {
			Offset = 0
		}
		db.Database.Limit(Limit).Offset(Offset).Find(&posts)
	} else {
		db.Database.Where(_where).Find(&posts)
	}

	// fmt.Println("posts:", posts)

	c.JSON(http.StatusOK, gin.H{"results": &posts})
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
			CreatedAt := strings.Split(fmt.Sprint(obj.CreatedAt), "T")
			// // fmt.Println("CreatedAt::",CreatedAt[0])
			Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ")
			if Date[0] == _type { //ถ้ามีcreated_atตรงกันให้แสดง
				_where["created_at"] = _type
			}

			if _type == "true" || _type == "false" {
				boolValue, _ := strconv.ParseBool(_type)
				if obj.Published == boolValue { //ถ้ามีpublishedตรงกันให้แสดง
					_where["published"] = boolValue
				}
			}
			if fmt.Sprint(obj.Id) == _type {
				_where["id"] = _type
			}
		}
		fmt.Println("_where[title]::", _where["title"])
		fmt.Println("_where[content]::", _where["content"])
		fmt.Println("_where[created_at]::", _where["created_at"])
		fmt.Println("_where[published]::", _where["published"])
		fmt.Println("_where[id]::", _where["id"])
	}
	// if _where != nil {
	// 	db.UpdateCollectionByViewCount(c)
	// }
	db.UpdateCollectionByViewCount(c)
	if _where["id"] != nil {
		db.GetCollectionByUUID(c)
	} else if _where["created_at"] != nil {
		db.GetCollectionByTime(c)
	} else if _where["published"] != nil {
		db.GetCollectionByPublished(c)

	} else {
		db.Database.Where(_where).Find(&posts)
		var dataPosts models.Posts
		if posts[0].Id != uuid.Nil {
			dataPosts.Id = posts[0].Id
		}
		if posts[0].Title != "" {
			dataPosts.Title = posts[0].Title
		}
		if posts[0].Content != "" {
			dataPosts.Content = posts[0].Content
		}
		if &posts[0].Published != nil {
			dataPosts.Published = posts[0].Published
		}
		if &posts[0].ViewCount != nil {
			ViewCount := 0
			ViewCount = posts[0].ViewCount + 1
			// db.UpdateCollectionByViewCount(c)
			dataPosts.ViewCount = ViewCount
		}
		if &posts[0].CreatedAt != nil {
			dataPosts.CreatedAt = posts[0].CreatedAt
		}
		if &posts[0].UpdatedAt != nil {
			dataPosts.UpdatedAt = posts[0].UpdatedAt
		}

		fmt.Println("_where", _where)
		if _where["title"] != nil || _where["content"] != nil || _where["id"] != nil {
			c.JSON(http.StatusOK, &dataPosts)
			fmt.Println("dataPosts.Published:", dataPosts.Published)
			if dataPosts.Published == true {
				db.UpdateCollectionByViewCount(c)
			}
			// db.UpdateCollectionByViewCount(c)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		}
	}
}

// // GET BY Limit
// func (db *DBController) GetCollectionByLimit(c *gin.Context) {
// 	_type := c.Param("id")

// 	var posts []models.Posts

// 	Limit, _ := strconv.Atoi(fmt.Sprint(_type))
// 	db.Database.Limit(Limit).Find(&posts)
// 	fmt.Println("db.Database.Limit(Limit).Find(&posts)", db.Database.Limit(Limit).Find(&posts))
// 	if len(posts) > 0 {
// 		c.JSON(http.StatusOK, &posts)
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
// 	}

// }

// GET BY UUID
func (db *DBController) GetCollectionByUUID(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	var dataPosts models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {

			Id := fmt.Sprint(obj.Id)
			if Id == _type { //ถ้ามีtitleตรงกันให้แสดง
				_where["id"] = _type
				if &obj.Id != nil {
					dataPosts.Id = obj.Id
				}
				if &obj.Title != nil {
					dataPosts.Title = obj.Title
				}
				if &obj.Content != nil {
					dataPosts.Content = obj.Content
				}
				if &obj.Published != nil {
					dataPosts.Published = obj.Published
				}
				if &obj.ViewCount != nil {
					dataPosts.ViewCount = obj.ViewCount
				}
				if &obj.CreatedAt != nil {
					dataPosts.CreatedAt = obj.CreatedAt
				}
				if &obj.CreatedAt != nil {
					dataPosts.UpdatedAt = obj.UpdatedAt
				}

			}

		}
	}

	if _where["id"] != nil {
		c.JSON(http.StatusOK, &dataPosts)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}

// GET BY Time
func (db *DBController) GetCollectionByTime(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {

			CreatedAt := strings.Split(fmt.Sprint(obj.CreatedAt), "T")
			// // fmt.Println("CreatedAt::",CreatedAt[0])
			Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ")
			if Date[0] == _type { //ถ้ามีtitleตรงกันให้แสดง
				_where["created_at"] = _type
			}

		}
	}
	db.Database.Where("created_at > ?", _where["created_at"]).Find(&posts)

	if _where["created_at"] != nil {
		c.JSON(http.StatusOK, &posts)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}

// GET BY published
func (db *DBController) GetCollectionByPublished(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {
			boolValue, err := strconv.ParseBool(_type)
			if err != nil {
				log.Fatal(err)
			}
			if obj.Published == boolValue { //ถ้ามีtitleตรงกันให้แสดง
				_where["published"] = boolValue
			}
			// CreatedAt := strings.Split(fmt.Sprint(obj.CreatedAt), "T")
			// // // fmt.Println("CreatedAt::",CreatedAt[0])
			// Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ")
			// if Date[0] == _type { //ถ้ามีtitleตรงกันให้แสดง
			// 	_where["created_at"] = _type
			// }

		}
	}
	db.Database.Where(_where).Find(&posts)

	if _where["published"] != nil {
		c.JSON(http.StatusOK, &posts)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}

// Update view_count
func (db *DBController) UpdateCollectionByViewCount(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	var dataPosts models.Posts
	ViewCount := 0
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {
			ViewCount = obj.ViewCount + 1
			Id := fmt.Sprint(obj.Id)
			if Id == _type { //ถ้ามีtitleตรงกันให้แสดง
				if fmt.Sprint(obj.Id) == _type {
					_where["id"] = _type
					if obj.Published == true {
						obj.ViewCount = ViewCount
						if &obj.Id != nil {
							dataPosts.Id = obj.Id
						}
						if &obj.Title != nil {
							dataPosts.Title = obj.Title
						}
						if &obj.Content != nil {
							dataPosts.Content = obj.Content
						}
						if &obj.Published != nil {
							dataPosts.Published = obj.Published
						}
						if &obj.ViewCount != nil {
							dataPosts.ViewCount = obj.ViewCount
						}
						if &obj.CreatedAt != nil {
							dataPosts.CreatedAt = obj.CreatedAt
						}
						if &obj.CreatedAt != nil {
							dataPosts.UpdatedAt = obj.UpdatedAt
						}

						db.Database.Model(&dataPosts).Update("view_count", ViewCount).Where(_where)
					}
				}
			}

		}
	}
}

// POST
func (db *DBController) CreateCollection(c *gin.Context) {
	var data models.Posts
	// err := c.ShouldBind(&data)
	//random Created id
	// myUUID := uuid.New()
	// // fmt.Println("myUUID::", myUUID.String())
	// data.Id = myUUID

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("err::", err)
		// logging.Logger(setting.LogLevelSetting.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if data.Title != ""{
		db.Database.Select("title", "content", "published").Create(&data)
		// db.Database.Create(&data)
	var dataPost models.ListPosts
	var Posts models.Posts
	
	if &data.Id != nil {
		Posts.Id = data.Id
	}
	if &data.Title != nil {
		Posts.Title = data.Title
	}
	if &data.Content != nil {
		Posts.Content = data.Content
	}
	if &data.Published != nil {
		Posts.Published = data.Published
	}
	if &data.ViewCount != nil {
		Posts.ViewCount = data.ViewCount
	}
	if &data.CreatedAt != nil {
		Posts.CreatedAt = data.CreatedAt
	}
	if &data.CreatedAt != nil {
		Posts.UpdatedAt = data.UpdatedAt
	}
	dataPost.Posts = append(dataPost.Posts, Posts)

	if len(dataPost.Posts) > 0 {
		c.JSON(http.StatusCreated, &dataPost)
	}
	}else{
		c.JSON(http.StatusCreated, make([]models.ListPosts, 0))
	}

}
// // POST
// func (db *DBController) CreateCollection(c *gin.Context) {
// 	var data models.Posts
// 	// err := c.ShouldBind(&data)
// 	//random Created id
// 	myUUID := uuid.New()
// 	// fmt.Println("myUUID::", myUUID.String())
// 	data.Id = myUUID

// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		fmt.Println("err::", err)
// 		// logging.Logger(setting.LogLevelSetting.Error, err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
// 		return
// 	}

// 	db.Database.Create(&data)
// 	var dataPost models.ListPosts

// 	dataPost.Posts = append(dataPost.Posts, data)

// 	if len(dataPost.Posts) > 0 {
// 		c.JSON(http.StatusCreated, &dataPost)
// 	} else {
// 		c.JSON(http.StatusCreated, make([]models.ListPosts, 0))
// 	}

// }

// func (db *DBController) UpdateCollection(c *gin.Context) {
// 	_type := c.Param("id")
// 	_where := map[string]interface{}{}
// 	var posts []models.Posts
// 	var dataPosts models.Posts
// 	if _type != "" {
// 		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
// 		//SELECT * from posts			^^^
// 		for _, obj := range posts {

// 			Id := fmt.Sprint(obj.Id)
// 			if Id == _type { //ถ้ามีtitleตรงกันให้แสดง
// 				if obj.Title == _type { //ถ้ามีtitleตรงกันให้แสดง
// 					_where["title"] = _type
// 				} else if obj.Content == _type { //ถ้ามีcontentตรงกันให้แสดง
// 					_where["content"] = _type
// 				} else if fmt.Sprint(obj.Id) == _type {
// 					_where["id"] = _type
// 				}
// 				if &obj.Id != nil {
// 					dataPosts.Id = obj.Id
// 				}
// 				if &obj.Title != nil {
// 					dataPosts.Title = obj.Title
// 				}
// 				if &obj.Content != nil {
// 					dataPosts.Content = obj.Content
// 				}
// 				if &obj.Published != nil {
// 					dataPosts.Published = obj.Published
// 				}
// 				if &obj.ViewCount != nil {
// 					dataPosts.ViewCount = obj.ViewCount
// 				}
// 				if &obj.CreatedAt != nil {
// 					dataPosts.CreatedAt = obj.CreatedAt
// 				}
// 				if &obj.CreatedAt != nil {
// 					dataPosts.CreatedAt = obj.CreatedAt
// 				}

// 			}

// 		}
// 	}
// 	// if _where["id"] != nil {
// 	// 	db.GetCollectionByUUID(c)
// 	// }
// 	db.Database.Where(_where).Find(&posts).Update(_type,dataPosts.Published)
// 	fmt.Println("db.Database::",db.Database)
// 	// db.Database.Where(_where).Update(&dataPosts)
// 	if _where["id"] != nil {
// 		c.JSON(http.StatusOK,dataPosts)
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
// 	}
// }

// DELETE BY UUID
func (db *DBController) DeleteCollection(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var posts []models.Posts
	var dataPosts models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
		for _, obj := range posts {

			Id := fmt.Sprint(obj.Id)
			if Id == _type { //ถ้ามีtitleตรงกันให้แสดง
				if fmt.Sprint(obj.Id) == _type {
					_where["id"] = _type
				}
				if &obj.Id != nil {
					dataPosts.Id = obj.Id
				}
				if &obj.Title != nil {
					dataPosts.Title = obj.Title
				}
				if &obj.Content != nil {
					dataPosts.Content = obj.Content
				}
				if &obj.Published != nil {
					dataPosts.Published = obj.Published
				}
				if &obj.ViewCount != nil {
					dataPosts.ViewCount = obj.ViewCount
				}
				if &obj.CreatedAt != nil {
					dataPosts.CreatedAt = obj.CreatedAt
				}
				if &obj.CreatedAt != nil {
					dataPosts.UpdatedAt = obj.UpdatedAt
				}

			}

		}
	}

	if _where["id"] != nil {
		db.Database.Where(_where).Delete(&posts)
		c.JSON(http.StatusOK, gin.H{"Delete": "UUID : " + fmt.Sprint(dataPosts.Id)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}
