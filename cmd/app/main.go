package main

import (
	"NotificationService/internal/app"
	"NotificationService/internal/infrastructure/queue"
	"NotificationService/internal/infrastructure/repository"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	processEvent := &app.ProcessEvent{
		Repo:  repo,
		Queue: queue,
	}
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
