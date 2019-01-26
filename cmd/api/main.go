package main

import (
	"github.com/gin-gonic/gin"

	"github.com/bufferserver/api"
	"github.com/bufferserver/config"
)

func main() {
	cfg := config.NewConfig()
	if cfg.GinGonic.IsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	apiServer := api.NewServer(cfg)
	apiServer.Run()
}
