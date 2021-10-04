package api

import (
	"encoding/base64"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sonntuet1997/medical-chain-utils/cryptography"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/lib/response"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/pb"
	types "github.com/sotanext-team/medical-chain/src/sideservice/src/pb/medichain"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/services"
)

type CosmosApiI interface {
	HandleAccountGet(g *gin.Context)
	HandleUserGet(g *gin.Context)
	HandleServiceUserPost(g *gin.Context)
	HandleCheckSharingGet(g *gin.Context)
}

type CryptoApiI interface {
	HandleSignPost(g *gin.Context)
	HandleVerifySigGet(g *gin.Context)
}

type CosmosApi struct {
	logger        *logrus.Logger
	cosmosService *services.CosmosService
}

func NewCosmosApi(logger *logrus.Logger, cosmosService *services.CosmosService) *CosmosApi {
	return &CosmosApi{logger, cosmosService}
}

func (c *CosmosApi) HandleAccountGet(g *gin.Context) {
	add := g.Param("addr")

	res, err := c.cosmosService.GetAccount(add)
	if err != nil {
		response.ErrInternalServerError(g, err)
		return
	}

	response.Success(g, res)
}

func (c *CosmosApi) HandleUserGet(g *gin.Context) {
	userId := g.Param("userId")

	c.logger.Info(g.Param("userId"))

	req := types.QueryGetUserRequest{
		Index: userId,
	}

	res, err := c.cosmosService.GetUser(&req)

	if err != nil {
		response.ErrInternalServerError(g, err)
		return
	}

	response.Success(g, res)
}

func (c *CosmosApi) HandleServiceUserPost(g *gin.Context) {
	msg64 := g.Query("message")
	serviceUser := g.Query("serviceUser")

	var msg pb.Message
	decMsg, err := base64.RawURLEncoding.DecodeString(msg64)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}
	err = json.Unmarshal(decMsg, &msg)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	var req types.MsgCreateServiceUser
	err = json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	verified, err := c.cosmosService.Verify(msg.Message, msg.Signature, req.UserId)
	if err != nil || !verified {
		response.ErrUnauthorized(g, err)
		return
	}

	req.ServiceUserId = serviceUser
	info, err := c.cosmosService.ShowAccount("admin")
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}
	bech32addr, err := sdk.Bech32ifyAddressBytes("medichain", info.GetAddress().Bytes())
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}
	req.Creator = bech32addr
	if err := req.ValidateBasic(); err != nil {
		c.logger.Errorf("error: %v", err)
	}

	res, err := c.cosmosService.PostCreateServiceUser(&req)
	if err != nil {
		logrus.Error(err)
		response.ErrInternalServerError(g, err)
		return
	}

	response.Success(g, res)
}

func (c *CosmosApi) HandleCheckSharingGet(g *gin.Context) {
	msg64 := g.Query("message")

	var msg pb.Message
	decMsg, err := base64.RawURLEncoding.DecodeString(msg64)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}
	err = json.Unmarshal(decMsg, &msg)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	var req types.QueryCheckSharingRequest
	err = json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	verified, err := c.cosmosService.Verify(msg.Message, msg.Signature, req.ViewerId)
	if err != nil || !verified {
		response.ErrUnauthorized(g, err)
		return
	}

	res, err := c.cosmosService.GetCheckSharing(&req)
	if err != nil {
		response.ErrInternalServerError(g, err)
		return
	}

	response.Success(g, res)
}

func (c *CosmosApi) HandleVerifySigGet(g *gin.Context) {
	msg64 := g.Query("message")

	decMsg, err := base64.RawURLEncoding.DecodeString(msg64)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	var msg pb.Message
	err = json.Unmarshal(decMsg, &msg)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	verified, err := c.cosmosService.Verify(msg.Message, msg.Signature, msg.PubKey)
	if err != nil || !verified {
		response.ErrUnauthorized(g, err)
		return
	}

	response.Success(g, verified)
}

func (c *CosmosApi) HandleSignPost(g *gin.Context) {
	msg64 := g.Query("message")
	key64 := g.Query("key")

	//key64 = "2MhRqfBvzaAZX7bQkR0R4KRQf3PUvd8hwUQLd7Ugp1A="
	//
	//keybz, _ := cryptography.ConvertBase64ToBytes(key64)
	//logrus.Info(base64.RawURLEncoding.EncodeToString(keybz))

	keyBz, err := base64.RawURLEncoding.DecodeString(key64)
	if err != nil {
		response.ErrBadRequest(g, err)
		return
	}

	res, err := c.cosmosService.Sign(msg64, cryptography.ConvertBytesToBase64(keyBz))
	if err != nil {
		response.ErrUnauthorized(g, err)
		return
	}
	response.Success(g, res)
}
