package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func queryDB(c *gin.Context) {
	var requestData struct {
		Query string `form:"query" json:"query"`
	}

	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	query := requestData.Query
	fmt.Println(query)
	fmt.Println(os.Getenv("DBUSER"))
	fmt.Println(os.Getenv("DBPASS"))

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "test_db",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Query"})
	}
	defer rows.Close()
	var queryResult []string
	for rows.Next() {
		if err := rows.Scan(&queryResult); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Query"})
		}
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Query"})
	}

	c.JSON(http.StatusOK, gin.H{"queryResult": rows})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/query", queryDB)

	r.Run(":3334")
}
