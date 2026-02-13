package main

import (
	"NotificationService/internal/app"
	"NotificationService/internal/http"
	"NotificationService/internal/infrastructure/handler"
	"NotificationService/internal/infrastructure/queue"
	"NotificationService/internal/infrastructure/repository"
	"NotificationService/internal/infrastructure/store"
	"NotificationService/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func main() {
	if err := utils.Init(); err != nil {
		log.Fatalf("%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Printf("Signal shutdown recieved")
		cancel()
	}()
	store, err := store.NewRedisClient(store.GetUrlString())
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer store.Close()
	repo, err := repository.NewPoolPG(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer repo.CloseDB()

	conn, err := queue.NewConn(viper.GetString("RABBITMQ_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%v", err)
	}
	queue, err := queue.NewRabbitQueue(ch, "EventQueue")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func() {
		queue.CloseChannel()
		queue.CloseConnection(conn)
	}()
	handler := &handler.Dispatcher{}

	router := fiber.New()
	acceptEvent := &app.AcceptEvent{
		IdemStore: store,
		Repo:      repo,
		Publisher: queue,
	}
	request := http.NewRequest(acceptEvent)
	router.Post("/events", request.CreateEvent)

	processEvent := &app.ProcessEvent{
		Repo:    repo,
		Queue:   queue,
		Handler: handler,
	}
	log.Println("Server start on port 3000")
	log.Fatalf("%v", router.Listen(":3000"))
	go func() {
		log.Printf("Worker started")
		for {
			select {
			case <-ctx.Done():
				log.Printf("worker stopped")
				return
			default:
				err := processEvent.Execute(ctx)
				if err != nil {
					log.Println("process error: ", err)
				}
			}
		}
	}()
}
