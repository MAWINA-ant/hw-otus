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
	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/server/internalgrpc"
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
		sqlStorage.ConnectWithConfig(ctx, sqlstorage.StorageConfig(config.SQLStorage))
		defer sqlStorage.Close(ctx)
		calendar = app.New(logg, sqlStorage)
	}

	httpServer := internalhttp.NewServerWithConfig(logg, calendar, internalhttp.ServerConfig(config.Server))
	grpcServer := internalgrpc.NewServer(logg, calendar, internalgrpc.ServerConfig(config.GRPCServer))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer stopCancel()

		if err := httpServer.Stop(stopCtx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := grpcServer.Stop(stopCtx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	errCh := make(chan error, 2)

	go func() {
		logg.Info("calendar http api is running...")
		if err := httpServer.Start(ctx); err != nil {
			errCh <- fmt.Errorf("failed to start http server: %w", err)
			return
		}
		errCh <- nil
	}()

	go func() {
		logg.Info("calendar grpc api is running...")
		if err := grpcServer.Start(ctx); err != nil {
			errCh <- fmt.Errorf("failed to start grpc server: %w", err)
			return
		}
		errCh <- nil
	}()

	var exitCode int
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			logg.Error(err.Error())
			exitCode = 1
			cancel()
		}
	}

	if exitCode != 0 {
		os.Exit(exitCode) //nolint:gocritic
	}
}
