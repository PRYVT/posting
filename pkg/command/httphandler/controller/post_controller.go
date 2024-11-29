package controller

import (
	"net/http"

	"github.com/PRYVT/posting/pkg/aggregates"
	"github.com/PRYVT/posting/pkg/models/command"
	"github.com/PRYVT/utils/pkg/hash"
	"github.com/gin-gonic/gin"
)

type PostController struct {
}

func NewPostController() *PostController {
	return &PostController{}
}

func (ctrl *PostController) CreatePost(c *gin.Context) {

	var m command.CreatePost
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userUuid := hash.GenerateGUID(m.Id)
	ua, err := aggregates.NewPostAggregate(userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = ua.CreatePost(m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
