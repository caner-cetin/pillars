package constants

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var MONGODB *mongo.Client

const DATABASE_NAME = "pillars"

const EARTHQUAKE_COLLECTION_NAME = "earthquakes"

const KAFKA_EARTHQUAKE_STREAM_TOPIC = "earthquake-stream"
const KAFKA_EARTHQUAKE_PAGE_LOAD_TOPIC = "earthquake-page-load"

const KAFKA_EARTHQUAKE_PARTITION = 0
const KAFKA_ADDR = "kafka:9092"
const KAFKA_CONNECTION_METHOD = "tcp"
const KAFKA_CONNECTION_TIMEOUT = 10

const KAFKA_READER_MAX_BYTES = 10e6
const KAFKA_PRODUCER_ERROR_KEY = "producer_error"
const KAFKA_NEW_EARTHQUAKE_BATCH_START_KEY = "new_earthquake_batch_start"
const KAFKA_NEW_EARTHQUAKE_BATCH_END_KEY = "new_earthquake_batch_end"

const STREAM_EARTHQUAKE_DEFAULT_GRID_SIZE = 0.05

const MONGODB_PAGINATION_DEFAULT_LIMIT = 5
