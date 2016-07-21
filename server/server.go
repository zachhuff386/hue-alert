package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/handlers"
	"net/http"
	"time"
)

func Server() (err error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())

	handlers.Register(router)

	addr := fmt.Sprintf(":%d", config.Config.ServerPort)

	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"address": addr,
	}).Info("server: Starting oauth server")

	err = server.ListenAndServe()
	if err != nil {
		return
	}

	return
}
