package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/config"
	"github.com/bufferserver/database"
)

type Server struct {
	db     *database.DB
	cache  *database.RedisDB
	cfg    *config.Config
	engine *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	db, err := database.NewMySQLDB(cfg.MySQL, cfg.API.MySQLConnCfg)
	if err != nil {
		log.WithField("err", err).Panic("initialize mysql db error")
	}

	return NewServerWithDB(cfg, db)
}

func NewServerWithDB(cfg *config.Config, db *database.DB) *Server {
	cache, err := database.NewRedisDB(cfg.Redis)
	if err != nil {
		log.WithField("err", err).Panic("initialize redis error")
	}

	server := &Server{
		db:    db,
		cache: cache,
		cfg:   cfg,
	}
	setupRouter(server)
	return server
}

func (s *Server) Run() {
	s.engine.Run(fmt.Sprintf(":%d", s.cfg.GinGonic.ListeningPort))
}

func (s *Server) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.ServerLabel, s)
		c.Set(common.DBLabel, s.db)
		c.Set(common.CacheLabel, s.cache)
		c.Next()
	}
}

func (s *Server) Head(_ *gin.Context) error {
	return nil
}

func setupRouter(apiServer *Server) {
	r := gin.Default()
	r.Use(apiServer.Middleware())
	r.HEAD("/dapp", handlerMiddleware(apiServer.Head))
	v1 := r.Group("/dapp")
	v1.POST("/update-base", handlerMiddleware(apiServer.UpdateBase))
	v1.POST("/list-utxos", handlerMiddleware(apiServer.ListUtxos))
	v1.POST("/list-balances", handlerMiddleware(apiServer.ListBalances))
	//v1.POST("/update-utxo", handlerMiddleware(apiServer.UpdateUtxo))
	//v1.POST("/update-balance", handlerMiddleware(apiServer.UpdateBalance))

	apiServer.engine = r
}

func handlerMiddleware(handleFunc interface{}) func(*gin.Context) {
	if err := common.ValidateFuncType(handleFunc); err != nil {
		panic(err)
	}

	return func(context *gin.Context) {
		common.HandleRequest(context, handleFunc)
	}
}
