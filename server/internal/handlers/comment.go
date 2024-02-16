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

func AddComment(commentService services.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		postIdParam := c.Param("post")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var comment models.Comment

		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.TrimSpace(comment.Text) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide comment text"})
			return
		}

		comment.Author = user.ID

		if err := comment.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := commentService.AddComment(comment, user.ID, int64(postId)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetAllPostComments(commentService services.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIdParam := c.Param("post")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		comments, err := commentService.GetAllPostComments(int64(postId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if comments == nil {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": "This post has no comments yet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": comments})
	}
}

func DeleteComment(commentService services.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		postIdParam := c.Param("post")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		commentIdParam := c.Param("comment")
		commentId, err := strconv.Atoi(commentIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := commentService.DeleteComment(int64(commentId), user.ID, int64(postId)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
