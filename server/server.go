package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	APIRouter    *gin.Engine
	APPRouter    *gin.Engine
	TokenSetTime time.Time
}

func New() *Server {
	return &Server{
		APIRouter: gin.Default(),
		APPRouter: gin.Default(),
	}
}

func (server *Server) RunAPP() {
	// start the server
	log.Fatal(server.APPRouter.Run("http://127.0.0.1:9090"))
}
