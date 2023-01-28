package main

import (
	"context"
	"os"

	"github.com/knadh/koanf"
	"github.com/mikeblum/golang-project-template/conf"
	"github.com/mikeblum/golang-project-template/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log := log.NewLog()
	if file, err := os.Create(conf.ConfFile); err != nil {
		log.WithError(err).Fatal("error creating conf")
	} else {
		defer file.Close()
	}
	cnf, err := conf.NewConf(conf.Provider(conf.ConfFile))
	if err != nil {
		log.WithError(err).Fatal("error loading conf")
	}
	ctx := context.WithValue(context.Background(), conf.Conf{}, cnf)
	log.WithContext(ctx).WithFields(logrus.Fields{
		"ctx": ctx.Value(conf.Conf{}).(*koanf.Koanf) != nil,
	}).Info("DONE")
}
