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
	_where := map[string]interface{}{}
	var posts []models.Posts
	var dataList models.ListPosts
	db.Database.Where(_where).Find(&posts)
	if len(posts) >0 {
	filter := make(map[string]interface{})
	if c.Request.URL.Query().Get("limit") != "" {
		filter["limit"] = c.Request.URL.Query().Get("limit")
	}
	if c.Request.URL.Query().Get("page") != "" {
		filter["page"] = c.Request.URL.Query().Get("page")
	}

	if filter["limit"] != nil && filter["page"] != nil {
		Limit, _ := strconv.Atoi(fmt.Sprint(filter["limit"]))
		Page, _ := strconv.Atoi(fmt.Sprint(filter["page"]))
		Offset := 0
		if Page >= 0 {
			Offset = (Page - 1) * Limit
		} else {
			Offset = 0
		}
		db.Database.Where(_where).Find(&posts)
		dataList.Count = len(posts)
		db.Database.Limit(Limit).Offset(Offset).Find(&posts)

		dataList.Posts = append(dataList.Posts, posts...)
		dataList.Limit = Limit
		dataList.Page = Page
		total := (dataList.Count / dataList.Limit)

		remainder  := (dataList.Count % dataList.Limit)
		if remainder  == 0 {
			dataList.TotalPage = total
		} else {
			dataList.TotalPage = total + 1
		}
	} else {
		db.Database.Where(_where).Find(&posts)
		
		dataList.Posts = append(dataList.Posts, posts...)
		dataList.Count = len(posts)
		dataList.Limit = len(posts)
		dataList.Page = 1
		total := (dataList.Count / dataList.Limit)
		// fmt.Println("total::",total)
		remainder  := (dataList.Count % dataList.Limit)
		// fmt.Println("totalpersen::",remainder)
		if remainder  == 0 {
			dataList.TotalPage = total
		} else {
			dataList.TotalPage = total + 1
		}
	}
		c.JSON(http.StatusOK,&dataList)
	} else {
		c.JSON(http.StatusOK, make([]models.ListPosts, 0))
	}
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
			if fmt.Sprint(obj.Id) == _type { //ถ้ามีidตรงกันให้แสดง
				_where["id"] = _type
			}
		}
	}

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
			dataPosts.ViewCount = posts[0].ViewCount
		}
		if &posts[0].CreatedAt != nil {
			dataPosts.CreatedAt = posts[0].CreatedAt
		}
		if &posts[0].UpdatedAt != nil {
			dataPosts.UpdatedAt = posts[0].UpdatedAt
		}

		if _where["title"] != nil || _where["content"] != nil || _where["id"] != nil {
			c.JSON(http.StatusOK, &dataPosts)
			if dataPosts.Published == true {
				//ถ้าPublished เป็น true  ViewCount +1 ทุกครั้งที่กดดู
				db.UpdateCollectionByViewCount(c)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		}
	}
}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
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

		}
	}
	db.Database.Where(_where).Find(&posts)

	if _where["published"] != nil {
		c.JSON(http.StatusOK, &posts)
	} else {
		c.JSON(http.StatusOK, make([]models.Posts, 0))
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

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("err::", err)
		// logging.Logger(setting.LogLevelSetting.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if data.Title != "" {
		db.Database.Select("title", "content", "published").Create(&data)
		// db.Database.Create(&data)
		////////////////////////////////////
		_where := map[string]interface{}{}
		var database []models.Posts
		db.Database.Where(_where).Find(&database) //เรียก ฐานข้อมูล
		ObJ := database[len(database)-1]
		//ไปดึงก้อนข้อมูลจาก ฐาน เอา ลิชล่างสุดที่พึ่งสร้าง เอาเฉพาะ ค่า ID
		//กรณีถ้าไม่ไปดึง ค่า ID จะเป็น 00000000-0000-0000-0000-000000000000
		//////////////////////////////////
		var Posts models.Posts

		if &data.Id != nil {
			Posts.Id = ObJ.Id
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

		c.JSON(http.StatusCreated, &Posts)
		
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}

}

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

// Update
func (db *DBController) UpdateCollection(c *gin.Context) {
	_type := c.Param("id")
	_where := map[string]interface{}{}
	var data models.Posts

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("err::", err)
		// logging.Logger(setting.LogLevelSetting.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if data.Title != "" {
		var database []models.Posts

		db.Database.Where(_where).Find(&database) //เรียก ฐานข้อมูล
		for _, obj := range database {
			Id := fmt.Sprint(obj.Id)
			if Id == _type { //
				_where["id"] = _type

				if &data.Id != nil {
					data.Id = obj.Id
				}
				if &data.CreatedAt != nil {
					data.CreatedAt = obj.CreatedAt
				}
				if &data.CreatedAt != nil {
					data.UpdatedAt = obj.UpdatedAt
				}

			}
		}
		if _where["id"] == _type {
			db.Database.Model(&database).Where(_where).Updates(map[string]interface{}{"title": data.Title, "content": data.Content, "published": data.Published})
			c.JSON(http.StatusCreated, &data)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}

}
