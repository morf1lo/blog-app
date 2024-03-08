package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils"
)

func (h *Handler) createPost(c *gin.Context) {
	user := utils.GetUserFromRequest(c)

	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorID = user.ID

	if err := post.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Post.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) getPostById(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.services.Post.FindPostById(int64(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": post})
}

func (h *Handler) getAuthorPosts(c *gin.Context) {
	author := c.Param("id")

	authorInt, err := strconv.Atoi(author)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, err := h.services.Post.FindAuthorPosts(int64(authorInt))
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

func (h *Handler) updatePost(c *gin.Context) {
	user := utils.GetUserFromRequest(c)

	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateOptions models.PostUpdateOptions

	if err := c.ShouldBindJSON(&updateOptions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Post.UpdatePost(updateOptions, int64(postID), user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) likePost(c *gin.Context) {
	user := utils.GetUserFromRequest(c)

	postIdParam := c.Param("id")
	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
	}

	if err := h.services.Post.LikePost(int64(postId), user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) deletePost(c *gin.Context) {
	user := utils.GetUserFromRequest(c)

	postIdParam := c.Param("id")
	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Post.DeletePost(int64(postId), user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) getUserLikes(c *gin.Context) {
	user := utils.GetUserFromRequest(c)

	likes, err := h.services.Post.FindUserLikes(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": likes})
}

func (h *Handler) searchPosts(c *gin.Context) {
	q := c.Query("q")

	posts, err := h.services.Post.SearchPosts(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": posts})
}
