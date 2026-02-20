package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"NotificationService/internal/app"
	"NotificationService/internal/http"
	"NotificationService/internal/infrastructure/handler"
	q "NotificationService/internal/infrastructure/queue"
	"NotificationService/internal/infrastructure/repository"
	"NotificationService/internal/infrastructure/store"
	"github.com/gofiber/fiber/v3"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
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
	conn, err := q.NewConn(viper.GetString("RABBITMQ_URL"))
	if err != nil {
		log.Printf("%v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%v", err)
	}
	idem_store, repo, queue := create(ctx, ch)
	defer idem_store.Close()
	defer repo.CloseDB()
	defer queue.CloseChannel()
	defer queue.CloseConnection(conn)
	handler := &handler.Dispatcher{}

	router := fiber.New()
	acceptEvent := &app.AcceptEvent{
		IdemStore: idem_store,
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
	go func() {
		log.Fatalf("%v", router.Listen(":3000"))
	}()
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
}

func create(
	ctx context.Context,
	ch *amqp091.Channel,
) (*store.RedisStore, *repository.Repository, *q.RabbitQueue) {
	viper.AutomaticEnv()
	idem_store, err := store.NewRedisClient(viper.GetString("REDIS_URL"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	repo, err := repository.NewPoolPG(ctx, viper.GetString("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	queue, err := q.NewRabbitQueue(ch, "EventQueue")
	if err != nil {
		log.Fatalf("%v", err)
	}
	return idem_store, repo, queue
}
