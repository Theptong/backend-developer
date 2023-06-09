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
เปลี่ยนจาก gorm เป็น sql
type DBController struct {
	Database *sql.DB
}
```

func checkErr(err error) {
if err != nil {
panic(err)
}
}

```
// GET
เป็นตัว Query ที่จะไปดึงจากฐานข้อมูล
func (db *DBController) QueryCollection(c *gin.Context) []models.Posts {
	// _where := map[string]interface{}{}
	var dataList []models.Posts
	// var dataList models.ListPosts

	rows, err := db.Database.Query(`SELECT * FROM posts`)
	SELECT * FROM posts คือการเลือกทุกก้อนจากฐานข้อมูล
	if err != nil {
		panic(err)
	}
	// var list = []models.Posts{}
	for rows.Next() {
		var Id uuid.UUID
		var Title string
		var Content string
		var Published bool
		var ViewCount int
		var CreatedAt time.Time
		var UpdatedAt time.Time
		err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
		checkErr(err)
	rows.Next ค่าแต่ละก้อนในฐาน
		posts := models.Posts{
			Id:        Id,
			Title:     Title,
			Content:   Content,
			Published: Published,
			ViewCount: ViewCount,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		}
		มายัดค่าตัวแปรเข้า posts และ ส่งออกเป็น arr dataList
		dataList = append(dataList, posts)
	}
	return dataList
}
```

```
การฟิลเตอร์ limit page
func (db *DBController) LimitCollection(Offset, Limit int) []models.Posts {
	// _where := map[string]interface{}{}
	var dataList []models.Posts
	// var dataList models.ListPosts
	sql := `SELECT * FROM posts OFFSET $1 LIMIT $2;`

SELECT * FROM posts OFFSET $1 LIMIT $2 การระบุ limit และ offset
	rows, err := db.Database.Query(sql, Offset, Limit)
	// db.Database.Query(.Limit)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var Id uuid.UUID
		var Title string
		var Content string
		var Published bool
		var ViewCount int
		var CreatedAt time.Time
		var UpdatedAt time.Time
		err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
		checkErr(err)

		posts := models.Posts{
			Id:        Id,
			Title:     Title,
			Content:   Content,
			Published: Published,
			ViewCount: ViewCount,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		}
		ยัดค่าเข้าตัวแปร ส่งออกเป็น arr list
		dataList = append(dataList, posts)
	}
	return dataList
}
```

```
// GET
func (db *DBController) GetCollection(c *gin.Context) {
// _where := map[string]interface{}{}
// var posts []models.Posts
var dataList models.ListPosts

    rows := db.QueryCollection(c)

    if len(rows) > 0 {
    	filter := make(map[string]interface{})
    	if c.Request.URL.Query().Get("limit") != "" {
    		filter["limit"] = c.Request.URL.Query().Get("limit")
    	}
    	if c.Request.URL.Query().Get("page") != "" {
    		filter["page"] = c.Request.URL.Query().Get("page")
    	}
	filter limit และ page 
    	if filter["limit"] != nil && filter["page"] != nil {
			ถ้า filter limit และ page  ไม่ใช่ ค่าว่างให้เข้าเงื่อนไข
    		Limit, _ := strconv.Atoi(fmt.Sprint(filter["limit"]))
			แปลง filter Limit เป็น int 
    		Page, _ := strconv.Atoi(fmt.Sprint(filter["page"]))
			แปลง filter Page เป็น int 
    		Offset := 0
    		if Page >= 0 {
    			Offset = (Page - 1) * Limit
    		} else {
    			Offset = 0
    		}

    		if Limit > 0 {
				ถ้า limit มากกว่า 0 เข้า เงื่อนไข 
				ต้องกระกาศก่อน dataList.Count เพราะต้องไปดึงข้อมูลจากฐานข้อมูลว่าจำนวนก้อนมีมากสุดเท่าไร ก่อนหาไป ฟิลเตอร์
    			dataList.Count = len(rows)
				rows ที่เข้า ฟังชั่น LimitCollection
    			rows := db.LimitCollection(Offset, Limit)
    			dataList.Posts = append(dataList.Posts, rows...)
    			dataList.Limit = Limit
    			dataList.Page = Page
    			total := (dataList.Count / dataList.Limit)

    			remainder := (dataList.Count % dataList.Limit)
    			if remainder == 0 {
    				dataList.TotalPage = total
    			} else {
    				dataList.TotalPage = total + 1
    			}
    		} else {
    			c.JSON(http.StatusOK, make([]models.ListPosts, 0))
    		}

    	} else {
    		dataList.Posts = rows
    		dataList.Count = len(rows)
    		dataList.Limit = len(rows)
    		total := (dataList.Count / dataList.Limit)
    		dataList.Page = total
    		remainder := (dataList.Count % dataList.Limit)
    		if remainder == 0 {
    			dataList.TotalPage = total
    		} else {
    			dataList.TotalPage = total + 1
    		}
    	}
    	c.JSON(http.StatusOK, &dataList)
    } else {
    	c.JSON(http.StatusOK, make([]models.ListPosts, 0))
    }

}
```
```
// // GET BY ID
func (db *DBController) GetCollectionById(c *gin.Context) {
_type := c.Param("id") รับ param
filter := map[string]interface{}{}

    var dataPost models.Posts
    posts := db.QueryCollection(c) // เรียก list จากฐานข้อมูล 
    for _, obj := range posts { วนฟอลเพื่อหาข้อมูล param ตรงกับ list ในก้อนมั้ย

    	if fmt.Sprint(obj.Id) == _type { เช็ค id ตรงกันมั้ย
    		if &obj.Id != nil { 
    			filter["id"] = obj.Id ถ้าข้อมูลตรง ให้ยัดเข้า filter["id"]
    		}
    	} else if obj.Title == _type { เช็คชื่อ ตรงกันมั้ย
    		if &obj.Title != nil {
    			filter["title"] = obj.Title ถ้าข้อมูลตรง ให้ยัดเข้า filter["title"]
    		}
    	} else if obj.Content == _type { เช็คชื่อ ตรงกันมั้ย
    		if &obj.Content != nil {
    			filter["content"] = obj.Content obj.Title ถ้าข้อมูลตรง ให้ยัดเข้า filter["content"]
    		}
    	}
    	if _type == "true" || _type == "false" {  เช็คว่าparamส่ง true และ false มามั้ย
    		boolValue, _ := strconv.ParseBool(_type) ถ้าข้อมูลตรง ให้แปลง เป็น bool
    		if obj.Published == boolValue { //ถ้ามีpublishedตรงกันให้แสดง
    			filter["published"] = boolValue ถ้าข้อมูลตรง ให้ยัดเข้า filter["published"]
    		}
    	}
    	CreatedAt := strings.Split(fmt.Sprint(obj.CreatedAt), "T") แยกวันและเวลาออกระหว่างตัว T
    	// // fmt.Println("CreatedAt::",CreatedAt[0])
    	Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ") //2023-05-23  แยกวันและเวลาออก
    	if Date[0] == _type { //ถ้ามีcreated_atตรงกันให้แสดง
    		filter["created_at"] = _type ถ้าข้อมูลตรง ให้ยัดเข้า filter["created_at"]
    	}
    }

    db.UpdateCollectionByViewCount(c) //นับจำนวนเข้าดู

    if filter != nil {
    	if filter["id"] != nil {  ถ้า filter["id"] ไม่ใช่ค่าว่าง เข้าเงื่อนไข
    		QueryCollection := db.QueryCollectionById(_type) ฟังชั่น ค้นหา UUID
    		dataPost = QueryCollection
    		c.JSON(http.StatusOK, dataPost)
    	} else if filter["title"] != nil { ถ้า filter["title"] ไม่ใช่ค่าว่าง เข้าเงื่อนไข
    		QueryCollection := db.QueryCollectionByTitle(_type) ฟังชั่น ค้นหา Title
    		dataPost = QueryCollection
    		c.JSON(http.StatusOK, dataPost)
    	} else if filter["content"] != nil { ถ้า filter["content"] ไม่ใช่ค่าว่าง เข้าเงื่อนไข
    		QueryCollection := db.QueryCollectionByContent(_type) ฟังชั่น ค้นหา Content
    		dataPost = QueryCollection
    		c.JSON(http.StatusOK, dataPost)
    	} else if filter["published"] != nil { ถ้า filter["published"] ไม่ใช่ค่าว่าง เข้าเงื่อนไข
    		var dataPost []models.Posts
    		QueryCollection := db.QueryCollectionByPublished(_type) ฟังชั่น ค้นหา Published
    		dataPost = QueryCollection
    		c.JSON(http.StatusOK, dataPost)
    	} else if filter["created_at"] != nil { ถ้า filter["created_at"] ไม่ใช่ค่าว่าง เข้าเงื่อนไข
    		var dataPost []models.Posts
    		QueryCollection := db.QueryCollectionByDate(_type) ฟังชั่น ค้นหา created_at
    		dataPost = QueryCollection
    		c.JSON(http.StatusOK, dataPost)
    	}
    } else {
    	if posts != nil {
    		c.JSON(http.StatusOK, posts)
    	} else {
    		c.JSON(http.StatusOK, make([]models.Posts, 0))
    	}
    }

}
```
```
// GET BY UUID
ฟังชั่นค้นหาจาก UUID
func (db _DBController) QueryCollectionById(id string) models.Posts {
var dataList models.Posts
sql := `SELECT _ FROM posts where id = $1 `

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}
```

```
// GET BY Time
ฟังชั่นค้นหาจาก Date
func (db \*DBController) QueryCollectionByDate(id string) []models.Posts {

    var dataList []models.Posts

    sql := `SELECT * FROM posts where created_at > $1`

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = append(dataList, posts)
    }
    return dataList

}
```

```
ฟังชั่นค้นหาจาก title
// GET BY Title
func (db \*DBController) QueryCollectionByTitle(id string) models.Posts {

    var dataList models.Posts

    sql := `SELECT * FROM posts where title = $1`

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}
```
```
// GET BY Content
ฟังชั่นค้นหาจาก Content
func (db \*DBController) QueryCollectionByContent(id string) models.Posts {

    var dataList models.Posts

    sql := `SELECT * FROM posts where content = $1`

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}
```

```
// GET BY Published
ฟังชั่นค้นหาจาก Published
func (db \*DBController) QueryCollectionByPublished(id string) []models.Posts {

    var dataList []models.Posts

    sql := `SELECT * FROM posts where published = $1`

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = append(dataList, posts)
    }
    return dataList

}
```

```
// Update view_count
ฟังชั่นนับจำนวน view_count
func (db *DBController) UpdateCollectionByViewCount(c *gin.Context) {
\_type := c.Param("id")

    filter := map[string]interface{}{}
    var dataPosts models.Posts
    ViewCount := 0
    if _type != "" {
    	posts := db.QueryCollection(c)

    	for _, obj := range posts {
    		if fmt.Sprint(obj.Id) == _type { ถ้าไอดี ตรงกัน 
    			ViewCount = obj.ViewCount + 1    //obj.ViewCount + 1 ค่าเดิม +1 ยัดเข้า ViewCount
    			if &obj.Id != nil {
    				filter["id"] = obj.Id
    			}
    			if obj.Published == true {
    				filter["published"] = obj.Published
    			}
    		}
    		if filter["id"] != nil {
    			if filter["published"] == true {
    				QueryCollection := db.QueryCollectionById(_type)
    				if &QueryCollection.Id != nil {
    					dataPosts.Id = QueryCollection.Id
    				}
    				if &QueryCollection.Title != nil {
    					dataPosts.Title = QueryCollection.Title
    				}
    				if &QueryCollection.Content != nil {
    					dataPosts.Content = QueryCollection.Content
    				}
    				if &QueryCollection.Published != nil {
    					dataPosts.Published = QueryCollection.Published
    				}
    				if &QueryCollection.ViewCount != nil {
    					dataPosts.ViewCount = ViewCount
    				}
    				if &QueryCollection.CreatedAt != nil {
    					dataPosts.CreatedAt = QueryCollection.CreatedAt
    				}
    				if &QueryCollection.CreatedAt != nil {
    					dataPosts.UpdatedAt = QueryCollection.UpdatedAt
    				}
    				db.QueryCollectionViewCount(ViewCount, _type)
    			}
    		}
    	}
    }

}
```

```
เพิ่มในฐานข้อมูล 
func (db \*DBController) QueryCollectionViewCount(ViewCount int, id string) models.Posts {

    var dataList models.Posts

    sql := `update posts
    set view_count = $1
    where id = $2`

    rows, err := db.Database.Query(sql, ViewCount, id)

    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

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

    if data.Title != "" {
		เข้าฟั่งชั่น เพิ่มในฐานข้อมูล
    	db.CreateNewCollection(data.Title, data.Content, data.Published)

    	database := db.QueryCollection(c)
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
    		Posts.CreatedAt = ObJ.CreatedAt
    	}
    	if &data.CreatedAt != nil {
    		Posts.UpdatedAt = ObJ.UpdatedAt
    	}

    	c.JSON(http.StatusCreated, &Posts)

    } else {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
    }

}
```
```
func (db \*DBController) CreateNewCollection(title, content string, published bool) models.Posts {
var dataList models.Posts

    sql := `INSERT INTO posts (title, content, published)
    VALUES ($1, $2, $3)`

    rows, err := db.Database.Query(sql, title, content, published)
    fmt.Println("CreateNewCollection::", rows)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}
```
```
// DELETE BY UUID
func (db *DBController) DeleteCollection(c *gin.Context) {
\_type := c.Param("id")
filter := map[string]interface{}{}

    var dataPosts models.Posts
    posts := db.QueryCollection(c)
    if _type != "" {
    	for _, obj := range posts {

    		Id := fmt.Sprint(obj.Id)
    		if Id == _type { //ถ้ามีidตรงกันให้แสดง
    			if fmt.Sprint(obj.Id) == _type {
    				filter["id"] = _type
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

    if filter["id"] != nil {
    	db.DeleteCollectionById(_type)
    	c.JSON(http.StatusOK, gin.H{"Delete": "UUID : " + fmt.Sprint(dataPosts.Id)})
    } else {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
    }

}
```
```
ฟังชั่น ลบในฐานข้อมูล
func (db \*DBController) DeleteCollectionById(id string) models.Posts {

    var dataList models.Posts

    sql := `delete from posts where id = $1`

    rows, err := db.Database.Query(sql, id)
    // db.Database.Query(.Limit)
    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}
```
```
// Update
func (db *DBController) UpdateCollection(c *gin.Context) {
\_type := c.Param("id")
filter := map[string]interface{}{}
var data models.Posts

    if err := c.ShouldBindJSON(&data); err != nil {
    	fmt.Println("err::", err)
    	// logging.Logger(setting.LogLevelSetting.Error, err)
    	c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
    	return
    }

    if data.Title != "" {
    	// var database []models.Posts

    	database := db.QueryCollection(c)
    	for _, obj := range database {
    		Id := fmt.Sprint(obj.Id)
    		if Id == _type { //
    			filter["id"] = _type

    			if &data.Id != nil {
    				data.Id = obj.Id
    			}
    			if &data.ViewCount != nil {
    				data.ViewCount = obj.ViewCount
    			}
    			if &data.CreatedAt != nil {
    				data.CreatedAt = obj.CreatedAt
    			}
    			if &data.CreatedAt != nil {
    				data.UpdatedAt = obj.UpdatedAt
    			}

    		}
    	}
    	if filter["id"] == _type {
    		db.QueryUpdateCollection(fmt.Sprint(filter["id"]), data.Title, data.Content, data.Published)
    		c.JSON(http.StatusCreated, &data)
    	}
    } else {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
    }

}
```
```
func (db \*DBController) QueryUpdateCollection(id, title, content string, published bool) models.Posts {

    var dataList models.Posts

    sql := `update posts
    set title = $2 , content = $3 , published = $4
    where id = $1`

    rows, err := db.Database.Query(sql, id, title, content, published)

    if err != nil {
    	panic(err)
    }
    for rows.Next() {
    	var Id uuid.UUID
    	var Title string
    	var Content string
    	var Published bool
    	var ViewCount int
    	var CreatedAt time.Time
    	var UpdatedAt time.Time
    	err = rows.Scan(&Id, &Title, &Content, &Published, &ViewCount, &CreatedAt, &UpdatedAt)
    	checkErr(err)

    	posts := models.Posts{
    		Id:        Id,
    		Title:     Title,
    		Content:   Content,
    		Published: Published,
    		ViewCount: ViewCount,
    		CreatedAt: CreatedAt,
    		UpdatedAt: UpdatedAt,
    	}
    	dataList = posts
    }
    return dataList

}

```
