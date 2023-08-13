package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mikeblum/golang-project-template/conf"
	"github.com/mikeblum/golang-project-template/log"
)

func main() {
	log := log.NewLog()
	if file, err := os.Create(conf.ConfFile); err != nil {
		log.WithError(err).Error("error creating conf")
	} else {
		defer file.Close()
	}
	cnf, err := conf.NewConf(conf.Provider(conf.ConfFile))
	if err != nil {
		log.WithError(err).Error("error loading conf")
	}
	ctx := context.WithValue(context.Background(), conf.Conf{}, cnf)
	log.LogAttrs(
		ctx,
		slog.LevelInfo,
		"example",
	)

}
