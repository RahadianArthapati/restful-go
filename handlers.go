package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchEmployees(c *gin.Context) {
	var e []Employee
	id := c.DefaultQuery("id", "all")
	keyword := c.Query("keyword")
	field := c.Query("field")
	if id == "all" && keyword == "" {
		if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.nip_new <> '' and employees.nip_new not like '%-%' and hist_jft.status is not null").Find(&e).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": true,
				"list":    e,
				"page":    "{total: 100, current: 1, pageSize: 10}",
			})
		}
	} else if keyword != "" && field != "" {
		if field == "name" {
			if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.nip_new <> '' and employees.nip_new not like '%-%' and hist_jft.status is not null and employees.fullname LIKE ?", "%"+keyword+"%").Find(&e).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"success": true,
					"list":    e,
					"page":    "{total: 1, current: 1, pageSize: 10}",
				})
			}
		} else if field == "nip" {
			if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.nip_new <> '' and employees.nip_new not like '%-%' and hist_jft.status is not null and employees.nip_new LIKE ?", "%"+keyword+"%").Find(&e).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"success": true,
					"list":    e,
					"page":    "{total: 1, current: 1, pageSize: 10}",
				})
			}
		}

	} else {
		if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.id = ?", id).Find(&e).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": true,
				"data":    e,
			})
		}

	}
}
func fetchSingleEmployee(c *gin.Context) {
	var e Employee
	id := c.Param("id")
	log.Println("ID : ", id)
	if err := db.Select("employees.id, employees.nip_new, employees.fullname, hist_jft.status").Joins("JOIN hist_jft ON employees.id = hist_jft.id").Where("employees.id = ?", id).Find(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": e})
	}
}
func uploadEmployeeData(c *gin.Context) {
	id := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	if err := c.SaveUploadedFile(file, "./employee_data/"+file.Filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with id ", file.Filename, id))
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
