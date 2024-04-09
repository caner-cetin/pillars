package errors

import (
	"context"
	"pillars-backend/src/constants"

	"github.com/segmentio/kafka-go"
)

func WriteKafkaProducerError(w *kafka.Writer, err error, to context.Context) {
	w.WriteMessages(to, kafka.Message{
		Key:   []byte(constants.KAFKA_PRODUCER_ERROR_KEY),
		Value: []byte(err.Error()),
	})
}
