package kafkax

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestCm(t *testing.T) {
	NewKafkaConsumer(
		[]string{"localhost:9092"},
		"commentGroup",
		"comment",
		10*time.Second,
	)
}

func TestP(t *testing.T) {
	writer := NewKafkaAsyncProducer([]string{"localhost:9092"})
	for i := 0; i < 100; i++ {
		m := map[string]string{
			"11111111111111111": strconv.Itoa(i),
		}
		bytes, _ := json.Marshal(m)
		writer.Send(Message{
			Data:  bytes,
			Topic: "comment",
		})
	}
	time.Sleep(time.Second * 100)
}
