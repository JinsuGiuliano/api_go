package controller

import (
	ml "example/api_go/model"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (cont *Controller) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, cont.data.GetUsers())
}

func (cont *Controller) GetUserByID(c *gin.Context) {
	var user ml.User
	id := c.Param("id")
	user, err := cont.data.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (cont *Controller) InsertUser(c *gin.Context) {
	var user ml.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, cont.data.InsertUser(user))
}

func (cont *Controller) DeleteUSer(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusAccepted, cont.data.DeleteUser(id))
}
