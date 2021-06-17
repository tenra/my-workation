package main

import (
  "log"
  "net/http"
  //"fmt"
  //"os"
  "strconv"
  //"time"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
  "github.com/joho/godotenv"
)

type User struct {
  gorm.Model
  Name string
  Email string
}

type Place struct {
	gorm.Model
	Content string
}

func gormConnect() *gorm.DB {
  err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
  DBMS := "mysql"
  USER := "go_test"
  PASS := "password"
  PROTOCOL := "tcp(db:3306)"
  DBNAME := "go_database"
  CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
  db, err := gorm.Open(DBMS, CONNECT)
  if err != nil {
    panic(err.Error())
  }
  return db
}


// DBの初期化
func dbInit() {
	db := gormConnect()
	// コネクション解放
	defer db.Close()
	db.AutoMigrate(&Place{})
	db.AutoMigrate(&User{})
}

// データインサート処理
func dbInsert(content string) {
	db := gormConnect()
	defer db.Close()
	// Insert処理
	db.Create(&Place{Content: content})
}

// db更新
func dbUpdate(id int, placeText string) {
	db := gormConnect()
	var place Place
	db.First(&place, id)
	place.Content = placeText
	db.Save(&place)
	db.Close()
}

// Place全件取得
func getAllPlaces() []Place {
	db := gormConnect()

	defer db.Close()
	var places []Place
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&places)
	return places
}

// Place一件取得
func getPlace(id int) Place {
	db := gormConnect()
	var place Place
	db.First(&place, id)
	db.Close()
	return place
}

// Place削除
func deletePlace(id int) {
	db := gormConnect()
	var place Place
	db.First(&place, id)
	db.Delete(&place)
	db.Close()
}


func main() {
  db := gormConnect()
  db.AutoMigrate(&User{}, &Place{})
  defer db.Close()

  router := gin.Default()
  router.LoadHTMLGlob("templates/*.html")

  // 一覧
  router.GET("/", func(c *gin.Context) {
    places := getAllPlaces()
    c.HTML(200, "index.html", gin.H{"places": places})
  })

  // Place登録
	router.POST("/new", func(c *gin.Context) {
		var form Place
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			places := getAllPlaces()
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"places": places, "err": err})
			c.Abort()
		} else {
			content := c.PostForm("content")
			dbInsert(content)
			c.Redirect(302, "/")
		}
	})

	// Place詳細
	router.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		// パラメータから受け取った値をint化
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		place := getPlace(id)
		c.HTML(200, "edit.html", gin.H{"place": place})
	})

	// Place更新
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		place := c.PostForm("place")
		dbUpdate(id, place)
		c.Redirect(302, "/")
	})

	// Place削除確認
	router.GET("/delete_confirm/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		place := getPlace(id)
		c.HTML(200, "delete.html", gin.H{"place": place})
	})

	// Place削除
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		deletePlace(id)
		c.Redirect(302, "/")

	})
  /*
  router.POST("/new", func(ctx *gin.Context) {
    db := sqlConnect()
    name := ctx.PostForm("name")
    email := ctx.PostForm("email")
    fmt.Println("create user " + name + " with email " + email)
    db.Create(&User{Name: name, Email: email})
    defer db.Close()

    ctx.Redirect(302, "/")
  })

  router.POST("/delete/:id", func(ctx *gin.Context) {
    db := sqlConnect()
    n := ctx.Param("id")
    id, err := strconv.Atoi(n)
    if err != nil {
      panic("id is not a number")
    }
    var user User
    db.First(&user, id)
    db.Delete(&user)
    defer db.Close()

    ctx.Redirect(302, "/")
  })
  */

  router.Run()
}// main
