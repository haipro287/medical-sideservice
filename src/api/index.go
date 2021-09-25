package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	g    *gin.Engine
	addr string
}

func (s *Service) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.g)
}

func NewApiService(addr string, indexApi *IndexApi, cosmosApi *CosmosApi) *Service {
	g := gin.Default()
	// Index
	index := g.Group("/")
	index.GET("/", indexApi.HandleIndexGet)
	index.GET("/kill", indexApi.HandleKillGet)

	v1 := g.Group("/medichain")

	v1.GET("/account/:addr", cosmosApi.HandleAccountGet)
	v1.GET("/user/:userId", cosmosApi.HandleUserGet)
	v1.GET("/sharing/validate", cosmosApi.HandleCheckSharingGet)

	v1.POST("/service/user", cosmosApi.HandleServiceUserPost)

	v1.GET("/crypto/verify", cosmosApi.HandleVerifySigGet)
	v1.POST("/crypto/sign", cosmosApi.HandleSignPost)

	return &Service{
		g:    g,
		addr: addr,
	}
}
