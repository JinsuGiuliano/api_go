package controller

import (
	"errors"
	"example/api_go/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (cont *Controller) TokenOrCookieAuth(c *gin.Context) {
	// check header for Token Auth
	reqToken := c.GetHeader("Authorization")
	if reqToken != "" {
		cont.TokenAuth(c)
		return
	}
	cont.CookieAuth(c)
}

func (cont *Controller) TokenAuth(c *gin.Context) {
	// setHeaders(c)
	reqToken := c.GetHeader("Authorization")
	if reqToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "authorization header not found"})
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Token Auth Unauthorized 1"})
		return
	}
	if splitToken[1] == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Token Auth Unauthorized 2"})
		return
	}
	tokenInfo := strings.Split(splitToken[1], ":")
	if len(tokenInfo) == 1 {
		cont.AccessTokenAuth(c, tokenInfo)
		return
	}
	cont.UserTokenAuth(c, tokenInfo)
}

// CookieAuth is used by the API endpoints to check if the user is logged in.
// If so, it goes next and return the endpoint response.
// If not, it returns an Unauthorized or Error HTTP status code along with an error description.
func (cont *Controller) CookieAuth(c *gin.Context) {
	// setHeaders(c)
	// We can obtain the session token from the requests cookies, which come with every handler
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// If the cookie is not set, return an unauthorized status
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Unauthorized - AuthRequired 1"})
			return
		}
		// For any other type of error, return a bad handler status
		c.AbortWithStatusJSON(http.StatusBadRequest, model.HTTPError{Error: "Bad handler - AuthRequired 2"})
		return
	}

	// We then get the name of the user from our cache, where we set the session token
	userID, err := cont.data.GetUserByToken(sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.HTTPError{Error: err.Error() + " - AuthRequired 3"})
		return
	}
	if userID.Token == "" {
		// If the session token is not present in cache, return an unauthorized error
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Unauthorized  - AuthRequired 4"})
		return
	}

	// Set this info to be used on the child handlers
	c.Set(responseKey, userID)
	c.Set(sessionTokenKey, sessionToken)
	c.Next()
}

func (cont *Controller) AccessTokenAuth(c *gin.Context, tokenInfo []string) {
	qToken := cont.data.GetValidAccessToken(tokenInfo[0], false)
	if qToken != nil || qToken.Token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "invalid token"})
		return
	}
	if qToken.Expire.UTC().Before(time.Now().UTC()) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "token expired"})
		return
	}
	c.Set(qTokenKey, qToken)
	c.Next()
}

func (cont *Controller) UserTokenAuth(c *gin.Context, tokenInfo []string) {
	user, err := cont.data.UserBy("email", tokenInfo[0])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Token Auth Unauthorized 1"})
		return
	}
	if user.Email == "" || user.Token == "" || !user.AllowTokenAuth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Token Auth Unauthorized 3"})
		return
	}
	if user.Token != tokenInfo[1] {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "Token Auth Unauthorized 4"})
		return
	}
	c.Set(responseKey, []byte(tokenInfo[0]))
	c.Next()
}

func (cont *Controller) AccessDataAuth(c *gin.Context) {
	user, err := cont.getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.HTTPError{Error: err.Error()})
		return
	}
	if !user.DataAllowed {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.HTTPError{Error: "user is not allow to access this resource"})
		return
	}
	c.Next()
}
