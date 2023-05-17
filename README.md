# backend-developer
ขั้นตอนแรก ลงพวก mod ที่จำเป็นจ้องใช้งาน
```
go mod init ("name-project")
github.com/gin-gonic/gin
github.com/joho/godotenv
gorm.io/driver/postgres
gorm.io/gorm
gorm.io/gorm เพื่อเชื่อมต่อไปยัง ฐานข้อมูล postgresSQL
go get tidy เพื่อ ดาวน์โหลด ส่วนที่ยังไม่ได้ลง 
```
สร้างโฟลเดอร์ที่จำเป็น

```
main.go 
routers
models
v1
```


```
PORT := "8080"
ยิงไปยังPORT ไหน

psqlInfo := "postgres://posts:38S2GPNZut4Tmvan@dev.opensource-technology.com:5523/posts?sslmode=disable"
db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
```
เช็ค database ว่าเชื่อมต่อหรือไม่ ถ้าเชื่อ จะเข้า db ถ้าไม่เชื่อมจะเข้า err 
```
    router := gin.New()
	api := router.Group("/api")

	//เข้า Routes
	routers.SetCollectionRoutes(api, db)

	port := fmt.Sprintf(":%v", PORT)
	fmt.Println("Server Running on Port", port)
	http.ListenAndServe(port, router)
    ยิงมา localhost:8080 ตาม Port ที่ตั้งไว้
```

routers.SetCollectionRoutes(api, db)
```
ctrls := v1.DBController{Database: db}
ctrls ระบุยังที่จะไป
SetCollectionRoutes เราจะ ระบุ พาร์ท ที่ต้องการจะยิงไป
POST
PUT
GET
GET ID
DELETE
router.GET("collections", ctrls.GetCollection)           // GET
router.GET("collections/:id", ctrls.GetCollectionById)   // GET BY ID
router.POST("collections", ctrls.CreateCollection)       // POST
router.PATCH("collections/:id", ctrls.UpdateCollection)  // PATCH
router.DELETE("collections/:id", ctrls.DeleteCollection) // DELETE
```

v1 หรือ Controller

```
หลักๆจะสร้าง medtod 
GetCollection
GetCollectionById
CreateCollection
UpdateCollection
DeleteCollection
```
structs
```
import (
	"time" เพื่อใช้  time.Time

	"github.com/google/uuid" เพื่อใช้  uuid.UUID
)
ก้อนเรียกค่า list
type ListPosts struct {
	Posts     []Posts `json:"posts"`
	Count     int   `json:"count"` 
	Limit     int   `json:"limit"` 
	Page      int   `json:"page"`
	TotalPage int   `json:"total_page"`
}
ก้อนที่จะยิงค่าไป
type Posts struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Published bool      `json:"published,omitempty"`
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
```


```
func (db *DBController) GetCollection(c *gin.Context) {
    ```
	_where := map[string]interface{}{}
	var posts []models.Posts
	var dataList models.ListPosts
	db.Database.Where(_where).Find(&posts)

    db.Database.Where(_where).Find(&posts) จะไปเรียกยังฐานข้อมูลที่ เรายิงไป
    ```
	if len(posts) >0 { เช็คว่า มีค่าว่างมั้ย
	filter := make(map[string]interface{})
    ฟิลเตอร์ limit กับ page ที่ใช้
	if c.Request.URL.Query().Get("limit") != "" {
		filter["limit"] = c.Request.URL.Query().Get("limit")
	}

	if c.Request.URL.Query().Get("page") != "" {
		filter["page"] = c.Request.URL.Query().Get("page")
	}
    ถ้ามีการ limit และ page จะเข้าเงื่อนไข if 
	if filter["limit"] != nil && filter["page"] != nil {
        
		Limit, _ := strconv.Atoi(fmt.Sprint(filter["limit"]))
		Page, _ := strconv.Atoi(fmt.Sprint(filter["page"]))
		Offset := 0 
		if Page >= 0 {
            คำนวนหา offset ที่จะต้องยิงไปยัง ฐาน
			Offset = (Page - 1) * Limit
		} else {
			Offset = 0
		}
		db.Database.Where(_where).Find(&posts)
        เรียกฐานข้อมูล posts เพื่อไปคิด จำนวนทั้งหมดที่ Count
		dataList.Count = len(posts) รับ ค่า len จากฐานทั้งหมด เพื่อเอาไปใช้ ยิงสูตรคำนวน ต่อ
		db.Database.Limit(Limit).Offset(Offset).Find(&posts)
        ใช้ limit กับ offset ยิงไป ยังฐาน เพื่อเอา limit กับ page 
		dataList.Posts = append(dataList.Posts, posts...)
        แมพ ค่า post เข้ากับ dataList 
		dataList.Limit = Limit
		dataList.Page = Page
        สูตรคำนวณหาค่า totalPage 

		total := (dataList.Count / dataList.Limit)

		remainder  := (dataList.Count % dataList.Limit)
        เศษ
		if remainder  == 0 {
			dataList.TotalPage = total
		} else {
			dataList.TotalPage = total + 1
		}
	} else {
        ถ้าไม่มีการ limit และ page จะเข้า else
		db.Database.Where(_where).Find(&posts)
        เรียก ข้อมูลที่ฐานข้อมูล 
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
        ใส่ค่าตัวแปรทั้งหมดที่ได้ในก้อนข้อมูล 
	}
    ประกาศออก
		c.JSON(http.StatusOK,&dataList)
	} else {
        ถ้าเป็นค่าว่าง ให้ส่ง arr ว่างออกไป
		c.JSON(http.StatusOK, make([]models.ListPosts, 0))
	}
}
```

```
// GET BY ID
func (db *DBController) GetCollectionById(c *gin.Context) {
	_type := c.Param("id") รับค่าparam 
	_where := map[string]interface{}{} ประกาศ _where map[string]interface เพื่อรับค่าไปเช็ค
	var posts []models.Posts
	if _type != "" {
		db.Database.Where(_where).Find(&posts) //เรียก ฐานข้อมูล
		//SELECT * from posts			^^^
        ```
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
    ```
    ฟังชั่น นับจำนวนเข้าไปดู 
	db.UpdateCollectionByViewCount(c)
    ```
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
			ViewCount = obj.ViewCount + 1 การวนหา ถ้าเจอ  ViewCount จะ +1
			Id := fmt.Sprint(obj.Id)
			if Id == _type { //ถ้ามีIdตรงกันให้แสดง
				if fmt.Sprint(obj.Id) == _type {
					_where["id"] = _type ยัดค่าตัวแปร  _type ไปยัง _where["id"]
					if obj.Published == true { ถ้า obj.Published = true
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
    ```
    ถ้ามี ไอดี ตรงกัน จะเข้า เงื่อนไข
	if _where["id"] != nil {
		db.GetCollectionByUUID(c)
        จะเข้า ฟั่งชั่น  GetCollectionByUUID
	} 
    ```
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
			if Id == _type { //ถ้ามีIdตรงกันให้แสดง
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
    ```
     ถ้ามี created_at ตรงกัน จะเข้า เงื่อนไข 
    else if _where["created_at"] != nil {
		db.GetCollectionByTime(c)
	}
    ฟังชั่นค้นหาจากเวลา 
    ```
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
			if Date[0] == _type { //ถ้ามีtimeตรงกันให้แสดง
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
    ```
     ถ้ามี published ตรงกัน จะเข้า เงื่อนไข 
     else if _where["published"] != nil {
		db.GetCollectionByPublished(c)
    ฟังชั่นค้นหาจากที่ตีพิมพ์ true & false 

    ```
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
			if obj.Published == boolValue { //ถ้ามีboolตรงกันให้แสดง
				_where["published"] = boolValue
			}

		}
	}
	db.Database.Where(_where).Find(&posts)

	if _where["published"] != nil {
		c.JSON(http.StatusOK, &posts)
	} else {
        ถ้าไม่มีก้อน ข้อมูลจาก ฐาน ให้แสดงเป็น []  กรณีที่จะไม่มีได้ก็ต่อเมื่อ ฐานข้อมูลไม่มีข้อมูล
		c.JSON(http.StatusOK, make([]models.Posts, 0))
	}
}
    ```
	} else {
        ถ้าไม่เข้าเงื่อนไขอะไรเลย
		db.Database.Where(_where).Find(&posts)
        เรียกฐานข้อมูล
		var dataPosts models.Posts ประกาศค่าตัวแปร
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
        นำข้อมูลที่เรียก ยัดเข้าค่าตัวแปรทั้งหมด 
        
		if _where["title"] != nil || _where["content"] != nil || _where["id"] != nil {
            เช็ค title หรือ content หรือ id ค่าใดก็ได้ ไม่ใช่ค่าว่าง นำมาแสดง
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
```

```
// POST
func (db *DBController) CreateCollection(c *gin.Context) {
	var data models.Posts

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("err::", err)
		// logging.Logger(setting.LogLevelSetting.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
    ถ้า Title ไม่ใช่ค่าว่างจะเข้าเงื่อนไข
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
        ถ้า Title เท่ากับค่าว่างจะเข้าเงื่อนไข
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}

}
```

```
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
        ถ้าไอดีไม่ใช่ค่าว่าง จะเข้าเงื่อนไข และยิงไปยิงฐานข้อมูลก้อนนั้น ลบได้เฉพาะ ไอดี เท่านั้น
		db.Database.Where(_where).Delete(&posts)
		c.JSON(http.StatusOK, gin.H{"Delete": "UUID : " + fmt.Sprint(dataPosts.Id)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}
}
```

```
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
    ถ้าTitle ไม่ใช่ค่าว่างจะเข้าเงื่อนไข
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
            เก็บค่า Id  CreatedAt UpdatedAt จากฐานข้อมูลเท่านั้น 
			}
		}
        ถ้าไอดีเท่ากับ param
		if _where["id"] == _type {
			db.Database.Model(&database).Where(_where).Updates(map[string]interface{}{"title": data.Title, "content": data.Content, "published": data.Published})
            จะเข้าไปยังข้อมูลที่ต้องการแก้ไข สามารถแก้ไข title และ  content  published ได้เท่านั้น จะเปลี่ยนค่าตัวไหนก็ได้
			c.JSON(http.StatusCreated, &data)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}

}
```