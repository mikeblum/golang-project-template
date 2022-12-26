package main

import (
	"github.com/mikeblum/golang-project-template/conf"
	"github.com/mikeblum/golang-project-template/log"
)

func main() {
	log := log.NewLog()
	config, err := conf.NewConfig()
	if err != nil {
		log.Fatalf("error loading conf: %v", err)
	}
	log.Infof("golang-project-template: %t", config != nil)
}
