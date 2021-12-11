package main

import (
	"github.com/foryouandyourcustomers/helm-keyvault/internal/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// setup logging
	log.SetOutput(os.Stdout)
	//log.SetLevel(config.Cfg.LogLevel)
}

func main() {
	if len(os.Args[1:]) != 4 {
		log.Fatal("Wrong amount of parameters given")
	}

	cmd.Download(os.Args[4])
}
