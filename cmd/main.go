package main

import (
	"flag"
	"fmt"
	"os"
	"sxolla-rest-api/config"
	"sxolla-rest-api/pkg/logging"
	"sxolla-rest-api/pkg/rest"

	"github.com/rs/zerolog/log"
)

var (
	debug = flag.Bool("debug", false, "Use debug mode")
	port  = flag.String("port", "8080", "The server port")
)

func main() {
	flag.Parse()
	loadLogger()
	// initiate Ent Client
	client, err := config.NewEntClient()
	if err != nil {
		log.Printf("err : %s", err)
	}
	defer client.Close()
	rest.RunGinService(*port, client)
}

func loadLogger() {
	logPath := os.Getenv("LOG_DIR")
	app := os.Getenv("APP")
	logConfig := logging.Config{Debug: *debug, EncodeLogsAsJson: true, ConsoleLoggingEnabled: *debug, FileLoggingEnabled: true,
		Directory: logPath, Filename: fmt.Sprintf("%s.log", app), MaxSize: 5, MaxBackups: 5, MaxAge: 10}
	newLog := logging.Configure(logConfig)
	log.Logger = *newLog.Logger
}
