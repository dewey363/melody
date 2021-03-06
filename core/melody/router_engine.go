package melody

import (
	"io"
	"melody/config"
	"melody/logging"

	botmonitor "melody/middleware/melody-botmonitor/gin"
	cors "melody/middleware/melody-cors/gin"
	httpsecure "melody/middleware/melody-httpsecure/gin"

	"github.com/gin-gonic/gin"
)

// NewEngine 返回一个基于gin的默认Engine
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, gelf io.Writer) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: gelf}), gin.Recovery())

	// 默认 重定向全部打开
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true
	// 注册跨域middleware
	if mw := cors.New(cfg.ExtraConfig); mw != nil {
		engine.Use(mw)
	}
	//http secure middleware
	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Warning(err)
	}
	//TODO lua register
	//botmonitor middleware
	botmonitor.Register(cfg, logger, engine)
	return engine
}
