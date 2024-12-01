package controller

import (
	"net/http"

	"github.com/PRYVT/posting/pkg/models/query"
	"github.com/PRYVT/posting/pkg/query/store/repository"
	"github.com/PRYVT/posting/pkg/query/utils"
	"github.com/PRYVT/utils/pkg/eventpolling"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	postRepo   *repository.PostRepository
	userEventH eventpolling.EventHanlder
}

func NewPostController(userRepo *repository.PostRepository, userEventH eventpolling.EventHanlder) *PostController {
	return &PostController{postRepo: userRepo, userEventH: userEventH}
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

	posts, err := ctrl.postRepo.GetAllPosts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if posts == nil {
		posts = []query.Post{}
	}
	c.JSON(http.StatusOK, posts)

}
