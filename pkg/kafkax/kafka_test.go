package kafkax

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
	"time"
)

func TestCm(t *testing.T) {
	NewKafkaConsumer(
		[]string{"localhost:9092"},
		"commentGroup",
		"comments",
		1*time.Second,
	)
	select {}
}

func TestP(t *testing.T) {
	writer := NewKafkaProducer([]string{"localhost:9092"}, "comments")
	for i := 0; i < 10; i++ {
		event := Event{
			LifeLogId: 123456,
			UserId:    int64(i),
		}
		marshal, err := json.Marshal(event)
		if err != nil {
			log.Fatal(err)
		}
		message := kafka.Message{
			Value: marshal,
		}
		writer.Send(message)
	}

	time.Sleep(time.Second * 100)
}

func TestBCm(t *testing.T) {
	consumer := NewKafkaAsyncBatchConsumer[Event](
		[]string{"localhost:9092"},
		"commentGroup",
		"comment111",
		10*time.Second,
		10,
		CreateHandler(),
	)
	consumer.ReadAndProcessMsg()
	select {}
}

// 定义事件结构体
type Event struct {
	LifeLogId int64 `json:"lifeLogId"`
	UserId    int64 `json:"userId"`
}

// 定义 handler 函数
func CreateHandler() func(vals []Event) error {
	return func(vals []Event) error {
		log.Printf("createHandler处理函数")
		for _, val := range vals {
			log.Printf("lifeLogId:%d,userId:%d", val.LifeLogId, val.UserId)
		}
		return nil
	}
}
