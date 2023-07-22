package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"route256/libs/cache"
	"route256/libs/logger"
	"route256/libs/postgres/tx"
	"route256/libs/tracer"
	"route256/notifications/internal/config"
	"route256/notifications/internal/pkg/notifications"
	"route256/notifications/internal/repository/postgres"
	"route256/notifications/internal/service"
	"route256/notifications/internal/telegram"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func main() {
	err := godotenv.Load("/.env")
	if err != nil {
		logger.Fatal("ERR loading .env file")
	}
	err = config.Init()
	if err != nil {
		logger.Fatal("ERR config.Init: ", err)
	}

	flag.Parse()
	// Init logger for environment
	logger.SetLoggerByEnvironment(*environment)

	// Init tracer
	if err = tracer.InitGlobal(config.AppConfig.ServiceName); err != nil {
		logger.Fatal(err)
	}

	// создаёт бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		logger.Panic(err)
	}

	telegramBot := telegram.NewBot(bot)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connection to db
	pool, err := pgxpool.Connect(ctx, os.Getenv("NOTIFICATIONS_DATABASE_URL"))
	if err != nil {
		logger.Fatalf("connect to db: %w", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)
	cache := cache.NewLRU(config.AppConfig.CacheCapacity)

	// создаёт сервис нотификаций
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.GrpcPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	service := service.NewNotificationService(provider, repo, cache, config.AppConfig.Brokers, telegramBot)
	notifications.RegisterNotificationsServer(s, service)

	logger.Infof("server listening at %v", lis.Addr())

	go func() {
		if err = s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	service.ConsumerGroupStatuses(config.AppConfig.Brokers, telegramBot, repo)
}
