package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/deathcore666/datagen/jsongen"
)

var (
	enqueued, successes, errors int64
)

func main() {
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup
	config := sarama.NewConfig()
	sarama.MaxRequestSize = 1000000
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Println(err)

	}

	log.Println(sarama.MaxRequestSize)
	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			errors++
		}
	}()

	timer := time.NewTimer(time.Second * 1)
	log.Println("1 secs started yo")

	// go seedMessages(signals, producer, timer)

	for {
		message := &sarama.ProducerMessage{Topic: "test1", Value: jsongen.WrapJSON()}
		select {
		case producer.Input() <- message:
			enqueued++

		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return

		case <-timer.C:
			log.Println("1 secs passed yo")
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return
		}
	}

	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, errors)
}

// func seedMessages(signals chan os.Signal, producer sarama.AsyncProducer, timer *time.Timer) {
// 	for {
// 		message := &sarama.ProducerMessage{Topic: "test1", Value: jsongen.WrapJSON()}
// 		select {
// 		case producer.Input() <- message:
// 			atomic.AddInt64(&enqueued, 1)

// 		case <-signals:
// 			producer.AsyncClose() // Trigger a shutdown of the producer.
// 			return

// 		case <-timer.C:
// 			log.Println("1 secs passed yo")
// 			producer.AsyncClose() // Trigger a shutdown of the producer.
// 			return
// 		}
// 	}
// }
