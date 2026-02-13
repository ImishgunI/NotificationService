package main

import (
	"NotificationService/internal/app"
	"NotificationService/internal/http"
	"NotificationService/internal/infrastructure/handler"
	"NotificationService/internal/infrastructure/queue"
	"NotificationService/internal/infrastructure/repository"
	"NotificationService/internal/infrastructure/store"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Printf("Signal shutdown recieved")
		cancel()
	}()
	store, err := store.NewRedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Printf("%v", err)
		cancel()
	}
	defer store.Close()
	repo, err := repository.NewPoolPG(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Printf("%v", err)
		cancel()
	}
	defer repo.CloseDB()

	conn := queue.NewConn(os.Getenv("RABBITMQ_CONNECTION_STRING"))
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%v", err)
		cancel()
	}
	queue, err := queue.NewRabbitQueue(ch, "EventQueue")
	if err != nil {
		log.Printf("%v", err)
		cancel()
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
