package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
func fetchSingleEmployee(c *gin.Context) {
	var e Employee
	id := c.Param("id")
	if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.id = ?", id).Find(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": e})
	}
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
