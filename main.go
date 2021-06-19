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

type Spot struct {
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
	db.AutoMigrate(&Spot{})
	db.AutoMigrate(&User{})
}

// データインサート処理
func dbInsert(content string) {
	db := gormConnect()
	defer db.Close()
	// Insert処理
	db.Create(&Spot{Content: content})
}

// db更新
func dbUpdate(id int, spotText string) {
	db := gormConnect()
	var spot Spot
	db.First(&spot, id)
	spot.Content = spotText
	db.Save(&spot)
	db.Close()
}

// Spot全件取得
func getAllSpots() []Spot {
	db := gormConnect()

	defer db.Close()
	var spots []Spot
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&spots)
	return spots
}

// Spot一件取得
func getSpot(id int) Spot {
	db := gormConnect()
	var spot Spot
	db.First(&spot, id)
	db.Close()
	return spot
}

// Spot削除
func deleteSpot(id int) {
	db := gormConnect()
	var spot Spot
	db.First(&spot, id)
	db.Delete(&spot)
	db.Close()
}


func main() {
  /*db := gormConnect()
  db.AutoMigrate(&User{}, &Spot{})
  defer db.Close()*/

  router := gin.Default()
  router.LoadHTMLGlob("templates/*.html")

  dbInit()

  // 一覧
  router.GET("/", func(c *gin.Context) {
    spots := getAllSpots()
    c.HTML(200, "index.html", gin.H{"spots": spots})
  })

  // Spot登録
	router.POST("/new", func(c *gin.Context) {
		var form Spot
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			spots := getAllSpots()
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"spots": spots, "err": err})
			c.Abort()
		} else {
			content := c.PostForm("content")
			dbInsert(content)
			c.Redirect(302, "/")
		}
	})

	// Spot詳細
	router.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		// パラメータから受け取った値をint化
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		spot := getSpot(id)
		c.HTML(200, "edit.html", gin.H{"spot": spot})
	})

	// Spot更新
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		spot := c.PostForm("spot")
		dbUpdate(id, spot)
		c.Redirect(302, "/")
	})

	// Spot削除確認
	router.GET("/delete_confirm/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		spot := getSpot(id)
		c.HTML(200, "delete.html", gin.H{"spot": spot})
	})

	// Spot削除
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		deleteSpot(id)
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
