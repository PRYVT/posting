package controller

import (
	"net/http"

	"github.com/PRYVT/posting/pkg/query/store/repository"
	"github.com/PRYVT/posting/pkg/query/utils"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	postRepo     *repository.PostRepository
	tokenManager *auth.TokenManager
}

func NewPostController(userRepo *repository.PostRepository, tokenManager *auth.TokenManager) *PostController {
	return &PostController{postRepo: userRepo, tokenManager: tokenManager}
}

func (ctrl *PostController) GetPost(c *gin.Context) {

	postUuid, err := utils.GetPostIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := ctrl.postRepo.GetPostById(postUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) GetPosts(c *gin.Context) {

	limit := utils.GetLimit(c)
	offset := utils.GetOffset(c)

	users, err := ctrl.postRepo.GetAllPosts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)

}
