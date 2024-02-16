package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/services"
	"github.com/morf1lo/blog-app/internal/utils"
)

type UpdateOptions struct {
	Title		string	`json:"title"`
	Text		string	`json:"text"`
}

func (u *UpdateOptions) filterUpdateOptions() (string, []interface{}) {
	query := "UPDATE posts SET"
	var values []interface{}

	if strings.TrimSpace(u.Title) == "" && strings.TrimSpace(u.Text) == "" {
		return "", nil
	}

	if strings.TrimSpace(u.Title) != "" {
		query += " title = ?,"
		values = append(values, strings.TrimSpace(u.Title))
	}

	if strings.TrimSpace(u.Text) != "" {
		query += " text = ?,"
		values = append(values, strings.TrimSpace(u.Text))
	}

	query = strings.TrimSuffix(query, ",")

	return query, values
}

func CreatePost(postService services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		var post models.Post

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post.Author = user.ID

		if err := post.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := postService.CreatePost(post); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetAuthorPosts(postService services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		author := c.Param("id")

		authorInt, err := strconv.Atoi(author)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		posts, err := postService.GetAuthorPosts(int64(authorInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if posts == nil {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": "This user has no posts yet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": posts})
	}
}

func UpdatePost(postService services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		postIdParam := c.Param("id")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateOptions UpdateOptions

		if err :=  c.ShouldBindJSON(&updateOptions); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updQuery, values := updateOptions.filterUpdateOptions()
		if updQuery == "" {
			c.JSON(http.StatusOK, gin.H{"success": true})
			return
		}

		updQuery += " WHERE id = ? AND author = ?"
		values = append(values, uint64(postId), user.ID)

		if err := postService.UpdatePost(updQuery, values); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func LikePost(postService services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		postIdParam := c.Param("id")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
				return
		}

		if err := postService.LikePost(int64(postId), user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func DeletePost(postService services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		postIdParam := c.Param("id")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := postService.DeletePost(int64(postId), user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
