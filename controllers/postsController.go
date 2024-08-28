package controllers

import (
	"net/http"
	"rental-mobil/initializers"
	"rental-mobil/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {
	// get data
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// create post
	post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// return post
	c.JSON(200, gin.H{
		"posts": post,
	})
}

func PostsCreateById(c *gin.Context) {
	id := c.Param("id")

	// get data
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// check if id already exists
	var count int64
	initializers.DB.Model(&models.Post{}).Where("id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"message": "ID already exists",
		})
		return
	}

	// parse id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	idInt64 := int64(idInt)

	// create post with specific id
	post := models.Post{ID: idInt64, Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to create post",
		})
		return
	}

	// return post
	c.JSON(200, gin.H{
		"posts": post,
	})
}

func PostsIndex(c *gin.Context) {
	// get post
	var posts []models.Post
	initializers.DB.Order("id asc").Find(&posts)
	// response

	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostsShow(c *gin.Context) {
	// get id url
	id := c.Param("id")
	// get post
	var post models.Post
	initializers.DB.First(&post, id)

	// response
	c.JSON(200, gin.H{
		"posts": post,
	})
}

func PostsUpdate(c *gin.Context) {
	// get id url
	id := c.Param("id")

	// get data req body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// find updating post
	var post models.Post
	initializers.DB.First(&post, id)

	// updated
	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})
	// response
	c.JSON(200, gin.H{
		"posts": post,
	})
}

func PostsDelete(c *gin.Context) {
	// get id url
	id := c.Param("id")

	// find post delete
	var post models.Post
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Error delete the data",
		})
		return
	}

	// delete
	initializers.DB.Delete(&models.Post{}, id)

	// response
	c.JSON(200, gin.H{
		"message": "Success Delete the Data",
	})
	return
}
