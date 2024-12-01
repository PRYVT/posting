package controller

import (
	"net/http"

	"github.com/PRYVT/posting/pkg/aggregates"
	"github.com/PRYVT/posting/pkg/models/command"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/hash"
	"github.com/gin-gonic/gin"
)

type PostController struct {
}

func NewPostController() *PostController {
	return &PostController{}
}

func (ctrl *PostController) CreatePost(c *gin.Context) {

	token := auth.GetTokenFromHeader(c)
	userUuid, err := auth.GetUserUuidFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var m command.CreatePost
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postUuid := hash.GenerateGUID(m.Id)
	ua, err := aggregates.NewPostAggregate(postUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	m.UserId = userUuid
	err = ua.CreatePost(m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
