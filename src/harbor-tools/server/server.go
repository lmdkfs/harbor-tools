package server

import (
	"fmt"
	"harbor-tools/harbor-tools/config"
	"harbor-tools/harbor-tools/router"
)

func NewServer() func() {
	cfg := config.NewConfig()
	return func() {
		addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
		GinServer := router.NewRouter()
		err := GinServer.Run(addr)
		if err != nil {
			panic(err)
		}
	}
}
