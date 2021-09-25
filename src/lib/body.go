package lib

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func GetBody(g *gin.Context, target interface{}) error {
	var bbody []byte
	var err error
	ibody, ok := g.Get("body")
	if !ok {
		// Not get in middleware
		bbody, err = ioutil.ReadAll(g.Request.Body)
		if err != nil {
			return err
		}
	} else {
		bbody = ibody.([]byte)
	}
	return json.Unmarshal(bbody, target)
}
