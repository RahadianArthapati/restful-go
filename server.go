package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
		//v1.GET("/:id", fetchSingleEmployee)
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
func fetchAllEmployees(c *gin.Context) {
	var e []Employee
	if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.nip_new <> '' and employees.nip_new not like '%-%' and hist_jft.status is not null").Find(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": e})
	}

	// rows, err := db.Table("employees").Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.nip_new <> '' and employees.nip_new not like '%-%' and hist_jft.status is not null").Rows()
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// handleResponse(c, rows, err)

}
func handleResponse(c *gin.Context, rows *sql.Rows, err error) {
	var e []Employee
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		a := Employee{}
		err := rows.Scan()
		if err != nil {
			log.Println(err)
			continue
		}
		e = append(e, a)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if len(e) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No employee found!"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": e})
	}
	fmt.Println(e)
}
