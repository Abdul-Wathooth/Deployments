package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type ListItem struct {
	Id        string    `json:"id"`
	Item      string    `json:"item"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB
var err error

func SetupPostgres() {
	// For Docker
	db, err = sql.Open("postgres", "postgres://postgres:password@postgres/todo?sslmode=disable")

	if err != nil {
		fmt.Println(err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err.Error())
	}

	log.Println("connected to postgres")
}

func TodoItems(c *gin.Context) {
	rows, err := db.Query("SELECT id, item, done, created_at, updated_at FROM list")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	items := make([]ListItem, 0)
	
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			item := ListItem{}
			if err := rows.Scan(&item.Id, &item.Item, &item.Done, &item.CreatedAt, &item.UpdatedAt); err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
				return
			}
			item.Item = strings.TrimSpace(item.Item)
			items = append(items, item)
		}
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func CreateTodoItem(c *gin.Context) {
	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if len(strings.TrimSpace(req.Item)) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter an item"})
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO list(item, done) VALUES($1, $2) RETURNING id;", req.Item, false).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	var newItem ListItem
	err = db.QueryRow("SELECT id, item, done, created_at, updated_at FROM list WHERE id = $1", id).Scan(
		&newItem.Id, &newItem.Item, &newItem.Done, &newItem.CreatedAt, &newItem.UpdatedAt,
	)
	newItem.Item = strings.TrimSpace(newItem.Item)

	log.Println("created todo item", id)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusCreated, gin.H{"item": newItem})
}

func UpdateTodoItem(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Item *string `json:"item"`
		Done *bool   `json:"done"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM list WHERE id=$1);", id).Scan(&exists)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	if req.Item != nil {
		_, err := db.Exec("UPDATE list SET item=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2;", *req.Item, id)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
			return
		}
	}

	if req.Done != nil {
		_, err := db.Exec("UPDATE list SET done=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2;", *req.Done, id)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
			return
		}
	}

	var updatedItem ListItem
	err = db.QueryRow("SELECT id, item, done, created_at, updated_at FROM list WHERE id = $1", id).Scan(
		&updatedItem.Id, &updatedItem.Item, &updatedItem.Done, &updatedItem.CreatedAt, &updatedItem.UpdatedAt,
	)
	updatedItem.Item = strings.TrimSpace(updatedItem.Item)

	log.Println("updated todo item", id)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, gin.H{"item": updatedItem})
}

func DeleteTodoItem(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM list WHERE id=$1);", id).Scan(&exists)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	_, err = db.Exec("DELETE FROM list WHERE id=$1;", id)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	log.Println("deleted todo item", id)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted todo item", "id": id})
}
