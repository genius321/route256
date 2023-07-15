package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"route256/libs/cache"
	"route256/libs/logger"
	"route256/libs/tracer"
	"route256/notifications/internal/pkg/notifications"
	"route256/notifications/internal/repository/postgres"
	"route256/notifications/internal/repository/postgres/tx"
	"route256/notifications/internal/service"
	"route256/notifications/internal/telegram"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

const grpcPort = 50053
const capacity = 3

func main() {
	flag.Parse()
	// Init logger for environment
	logger.SetLoggerByEnvironment(*environment)

	// Init tracer
	if err := tracer.InitGlobal("NOTIFICATIONS"); err != nil {
		logger.Fatal(err)
	}

	// создаёт бота
	bot, err := tgbotapi.NewBotAPI("6376011717:AAFrInwAk9DH2NjwbRXuBCslp9D6QVHsPew")
	if err != nil {
		logger.Panic(err)
	}

	telegramBot := telegram.NewBot(bot)

	// список брокеров
	var brokers = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connection to db
	pool, err := pgxpool.Connect(ctx, os.Getenv("NOTIFICATIONS_DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to db: %w", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)
	cache := cache.NewLRU(capacity)

	// создаёт сервис нотификаций
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	service := service.NewNotificationService(provider, repo, cache, brokers, telegramBot)
	notifications.RegisterNotificationsServer(s, service)

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	service.ConsumerGroupStatuses(brokers, telegramBot, repo)
}
