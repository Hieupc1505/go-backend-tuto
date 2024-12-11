package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaProducer *kafka.Writer
)

const (
	kafkaURL   = "localhost:9092"
	kafkaTopic = "user_topic_vip"
)

// for producer
func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// for consumer
func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.LastOffset,
	})
}

type StockInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func newStock(msg, typeMsg string) *StockInfo {
	return &StockInfo{
		Message: msg,
		Type:    typeMsg,
	}
}

func actionStock(c *gin.Context) {

	s := newStock(c.Query("msg"), c.Query("type"))
	body := make(map[string]interface{})
	body["action"] = "action"
	body["info"] = s

	jsonBody, _ := json.Marshal(body)

	msg := kafka.Message{
		Key:   []byte("action"),
		Value: []byte(jsonBody),
	}

	err := kafkaProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"err": "",
		"msg": "success",
	})

}

func RegisterConsumerATC(id int) {
	//group consumer
	kafkaGroupId := "consumer-group-"
	reader := getKafkaReader(kafkaURL, kafkaTopic, kafkaGroupId)
	defer reader.Close()

	fmt.Printf("Consumer %d started\n", id)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Consumer %d: %v\n", id, err)
			continue
		}
		fmt.Printf("Consumer(%d), hong topic:%v, partition:%v, offset:%v, time: %d %s = %s\n", id, m.Topic, m.Partition, m.Offset, m.Time.Unix(), string(m.Key), string(m.Value))
	}
}

func main() {
	r := gin.Default()
	kafkaProducer = getKafkaWriter(kafkaURL, kafkaTopic)
	defer kafkaProducer.Close()

	r.POST("action/stock", actionStock)

	//regist hong
	go RegisterConsumerATC(1)
	go RegisterConsumerATC(2)
	r.Run(":8999")
}
