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
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	//r.Static("/", "./public")
	e := r.Group("/api/v1/employees")
	{
		e.GET("/", fetchAllEmployees)
		e.GET("/:id", fetchSingleEmployee)
		e.POST("/:id/upload", uploadEmployeeData)
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
