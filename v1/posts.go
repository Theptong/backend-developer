package v1

import (
	"backend-test/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DBController struct {
	Database *sql.DB
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GET
func (db *DBController) QueryCollection(c *gin.Context) []models.Posts {
	// _where := map[string]interface{}{}
	var dataList []models.Posts
	// var dataList models.ListPosts

	rows, err := db.Database.Query(`SELECT * FROM posts`)
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

func (db *DBController) LimitCollection(Offset, Limit int) []models.Posts {
	// _where := map[string]interface{}{}
	var dataList []models.Posts
	// var dataList models.ListPosts
	sql := `SELECT * FROM posts OFFSET $1 LIMIT $2;`

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
		dataList = append(dataList, posts)
	}
	return dataList
}

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

		if filter["limit"] != nil && filter["page"] != nil {
			Limit, _ := strconv.Atoi(fmt.Sprint(filter["limit"]))
			Page, _ := strconv.Atoi(fmt.Sprint(filter["page"]))
			Offset := 0
			if Page >= 0 {
				Offset = (Page - 1) * Limit
			} else {
				Offset = 0
			}

			if Limit > 0 {
				dataList.Count = len(rows)

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

// // GET BY ID
func (db *DBController) GetCollectionById(c *gin.Context) {
	_type := c.Param("id")
	filter := map[string]interface{}{}

	var dataPost models.Posts
	posts := db.QueryCollection(c)
	for _, obj := range posts {

		if fmt.Sprint(obj.Id) == _type {
			if &obj.Id != nil {
				filter["id"] = obj.Id
			}
		} else if obj.Title == _type {
			if &obj.Title != nil {
				filter["title"] = obj.Title
			}
		} else if obj.Content == _type {
			if &obj.Content != nil {
				filter["content"] = obj.Content
			}
		}
		if _type == "true" || _type == "false" {
			boolValue, _ := strconv.ParseBool(_type)
			if obj.Published == boolValue { //ถ้ามีpublishedตรงกันให้แสดง
				filter["published"] = boolValue
			}
		}
		CreatedAt := strings.Split(fmt.Sprint(obj.CreatedAt), "T")
		// // fmt.Println("CreatedAt::",CreatedAt[0])
		Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ")
		if Date[0] == _type { //ถ้ามีcreated_atตรงกันให้แสดง
			filter["created_at"] = _type
		}else{
			filter["created_at"] = _type
		}
	}

	db.UpdateCollectionByViewCount(c) //นับจำนวนเข้าดู
	if len(posts) > 0 {
		if filter["id"] != nil {
			QueryCollection := db.QueryCollectionById(_type)
			dataPost = QueryCollection
			c.JSON(http.StatusOK, dataPost)
		} else if filter["title"] != nil {
			QueryCollection := db.QueryCollectionByTitle(_type)
			dataPost = QueryCollection
			c.JSON(http.StatusOK, dataPost)
		} else if filter["content"] != nil {
			QueryCollection := db.QueryCollectionByContent(_type)
			dataPost = QueryCollection
			c.JSON(http.StatusOK, dataPost)
		} else if filter["published"] != nil {
			var dataPost []models.Posts
			QueryCollection := db.QueryCollectionByPublished(_type)
			dataPost = QueryCollection
			c.JSON(http.StatusOK, dataPost)
		} else if filter["created_at"] != nil {
			var dataPost []models.Posts
			timeNow := time.Now().AddDate(0,0,1)
			CreatedAt := strings.Split(fmt.Sprint(timeNow), "T")
			Date := strings.Split(fmt.Sprint(CreatedAt[0]), " ")
			QueryCollection := db.QueryCollectionByDate(_type,fmt.Sprint(Date[0]))
			dataPost = QueryCollection
			c.JSON(http.StatusOK, dataPost)
		}
	} else {
		c.JSON(http.StatusOK, make([]models.Posts, 0))
	}
}

// GET BY UUID
func (db *DBController) QueryCollectionById(id string) models.Posts {
	var dataList models.Posts
	sql := `SELECT * FROM posts where id = $1 `

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

// GET BY Time
func (db *DBController) QueryCollectionByDate(id string,today string) []models.Posts {

	var dataList []models.Posts
	// sql := `SELECT * FROM posts where created_at > $1
	sql := `SELECT * FROM posts where created_at between $1 and $2`

	rows, err := db.Database.Query(sql, id, today)
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

// GET BY Title
func (db *DBController) QueryCollectionByTitle(id string) models.Posts {

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

// GET BY Content
func (db *DBController) QueryCollectionByContent(id string) models.Posts {

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

// GET BY Published
func (db *DBController) QueryCollectionByPublished(id string) []models.Posts {

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

// Update view_count
func (db *DBController) UpdateCollectionByViewCount(c *gin.Context) {
	_type := c.Param("id")

	filter := map[string]interface{}{}
	var dataPosts models.Posts
	ViewCount := 0
	if _type != "" {
		posts := db.QueryCollection(c)

		for _, obj := range posts {
			if fmt.Sprint(obj.Id) == _type {
				ViewCount = obj.ViewCount + 1
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

func (db *DBController) QueryCollectionViewCount(ViewCount int, id string) models.Posts {

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
func (db *DBController) CreateNewCollection(title, content string, published bool) models.Posts {
	var dataList models.Posts

	sql := `INSERT INTO posts (title, content, published)
	VALUES ($1, $2, $3)`

	rows, err := db.Database.Query(sql, title, content, published)
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

// DELETE BY UUID
func (db *DBController) DeleteCollection(c *gin.Context) {
	_type := c.Param("id")
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

func (db *DBController) DeleteCollectionById(id string) models.Posts {

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

// Update
func (db *DBController) UpdateCollection(c *gin.Context) {
	_type := c.Param("id")
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

func (db *DBController) QueryUpdateCollection(id, title, content string, published bool) models.Posts {

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
