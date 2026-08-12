package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/app"
	internallogger "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// конфигурация
	config := NewConfig()
	absConfigFilePath, _ := filepath.Abs(configFile)
	if err := config.ParseConfigFromFile(absConfigFilePath); err != nil {
		fmt.Printf("Couldn't parse config file %s cause %s", configFile, err)
		os.Exit(1)
	}

	logg := internallogger.New(config.Logger.Level, config.Logger.LogFile)
	defer logg.Close()

	var calendar *app.App
	if config.Storage.InMemory {
		memoryStorage := memorystorage.New()
		calendar = app.New(logg, memoryStorage)
	} else {
		ctx := context.Background()
		sqlStorage := sqlstorage.New()
		sqlStorage.Connect(ctx)
		defer sqlStorage.Close(ctx)
		calendar = app.New(logg, sqlStorage)
	}

	server := internalhttp.NewServerWithConfig(logg, calendar, internalhttp.ServerConfig(config.Server))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
