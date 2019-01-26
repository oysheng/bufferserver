package main

import (
	"github.com/gin-gonic/gin"

	"github.com/blockcenter/api"
	"github.com/blockcenter/config"
)

func main() {
	cfg := config.NewConfig()
	if cfg.GinGonic.IsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	apiServer := api.NewServer(cfg)
	apiServer.Run()
}
