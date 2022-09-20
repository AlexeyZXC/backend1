// Main package of backend that connects all the packages together and starts the app.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler/routerchi"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/server"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/db/postgres"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// config
	config, err := config.NewConfig()
	if err != nil {
		fmt.Println("Failed to read config file: ", err)
	}

	// logger
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeStamp)
	logger, err := conf.Build()
	if err != nil {
		fmt.Println("Creating zap logger error: ", err)
		return
	}
	defer logger.Sync()
	log := logger.Sugar()

	// postgress
	db, err := postgres.NewPgDB()
	if err != nil {
		log.Errorf("Db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	rh := routerchi.NewRouterChi(h, log)

	// for server stopping
	wg := sync.WaitGroup{}
	wg.Add(1)

	srv := server.NewServer(":"+config.Port, rh, &wg, log)

	// stopping stuff
	stopch := make(chan os.Signal, 1)
	signal.Notify(stopch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopch
		log.Infof("Stopping...")
		srv.Stop()
		db.Close()
	}()

	srv.Start(ls)

	wg.Wait()
	log.Infof("--- Stopped ---")
}
