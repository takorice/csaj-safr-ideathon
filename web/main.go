package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ReactionSummary struct {
	Id             int64     `json:"id" gorm:"AUTO_INCREMENT;column:id;primary_key"`
	UserId         int64     `json:"user_id"`
	LessonId       int64     `json:"lesson_id"`
	EmotionalValue float32   `json:"emotional_value" sql:"type:decimal;"`
	ReactedAt      time.Time `json:"reacted_at" time_format:"2006-01-02"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at" sql:"not null;"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at" sql:"not null;"`
}

func gormConnect() *gorm.DB {
	connection := os.Getenv("DATABASE_URL")
	val, ret := os.LookupEnv("SSL_MODE")
	if ret == false {
		val = "false"
	}

	sslmode, _ := strconv.ParseBool(val)
	if sslmode {
		connection += "?sslmode=require"
	}
	fmt.Println(connection)
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	db := gormConnect()

	r := gin.Default()
	r.GET("/api/reaction_summaries", func(c *gin.Context) {
		summaries := []ReactionSummary{}

		// クエリパラメータを取得
		reactedAtFrom := c.Query("reacted_at_from")
		reactedAtTo := c.Query("reacted_at_to")

		// 新規クエリを作成
		tx := db
		// クエリパラメータ設定
		if reactedAtFrom != "" {
			tx = tx.Where("reacted_at >= ?", reactedAtFrom)
		}
		if reactedAtTo != "" {
			tx = tx.Where("reacted_at <= ?", reactedAtTo)
		}

		// SQL クエリのデバッグ出力
		//tx.Debug().Find(&summaries)
		// SQL 実行
		tx.Find(&summaries)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, gin.H{
			"reaction_summaries": summaries,
		})
	})

	defer db.Close()

	db.AutoMigrate(&ReactionSummary{})

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Printf("Starting server at Port %d", port)
	r.Run(fmt.Sprintf(":%d", port))
}
