package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type (
	Employee struct {
		ID       uint   `json:"id"`
		Nip_New  string `json:"nip"`
		Fullname string `json:"name"`
		Status   string `json:"status"`
	}
)

func main() {
	startDB()
	r := gin.New()
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1/employees")
	{
		v1.GET("/", fetchAllEmployees)
		v1.GET("/:id", fetchSingleEmployee)
		//v1.PUT("/:id", updateSingleEmployee)
	}
	// Listen and serve on 0.0.0.0:8080
	r.Run(":3000")
}
func startDB() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
}
