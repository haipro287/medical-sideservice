package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/lib/response"
	"syscall"
)

type IndexApi struct {
	logger *logrus.Logger
}

func NewIndexApi(logger *logrus.Logger) *IndexApi {
	return &IndexApi{logger}
}
func (u *IndexApi) HandleIndexGet(g *gin.Context) {
	response.Success(g, "Hello from Medical-chain-server")
}

func (u *IndexApi) HandleKillGet(g *gin.Context) {
	response.Success(g, "Bye bye from Medical-chain-server")
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
