package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type (
	Employee struct {
		ID       uint   `json:"id"`
		Nip_New  string `json:"nip"`
		Golongan string `json:"golongan"`
		Unit     string `json:"unit"`
		Fullname string `json:"name"`
		Status   string `json:"status"`
	}
	EmployeeDetail struct {
		ID         uint   `json:"id"`
		Status     string `json:"status"`
		Jabatan    string `json:"jabatan"`
		No_SK      string `json:"no_sk"`
		Date_SK    string `json:"date_sk"`
		Date_Start string `json:"date_start"`
		Date_End   string `json:"date_end"`
	}
)

func main() {
	startDB()
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(CORSMiddleware())
	//r.Static("/", "./public")
	e := r.Group("/api/v1/employees")
	{
		e.GET("/", fetchEmployees)
		e.GET("/user/:id", fetchSingleEmployee)
		e.POST("/upload", uploadEmployeeData)
		e.POST("/print", generate_report)
		//v1.PUT("/:id", updateSingleEmployee)
	}
	// Listen and serve on 0.0.0.0:8080
	//r.Run(":3000")
	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Fatal("Server close under request", err)
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
func startDB() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
}
