package controller

import (
	ml "example/api_go/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

// getUser get the user from the session key
func (cont *Controller) getUser(c *gin.Context) (user ml.User, err error) {
	// This is set on the AuthRequired middleware
	response, ok := c.Get(responseKey)
	if !ok {
		return ml.User{}, fmt.Errorf("there is no user associated with the given key")
	}

	user, err = cont.data.UserByEmail(string(response.([]uint8)))
	if err != nil {
		return ml.User{}, fmt.Errorf("unable to get user from email")
	}

	return user, nil
}
