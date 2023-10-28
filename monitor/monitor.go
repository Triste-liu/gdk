package monitor

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/triste-liu/gdk/log"
)

type Config struct {
	Port   int
	Prefix string
}

func Run(c Config) {
	addr := fmt.Sprintf(":%d", c.Port)
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	pprof.Register(g, c.Prefix)
	log.Info("monitor route: %s", addr+c.Prefix)
	err := g.Run(addr)
	if err != nil {
		log.Warning("monitor run error:%s", err)
	}
}
