package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"pillars-backend/src/constants"
	"pillars-backend/src/db"
	"pillars-backend/src/errors"
	"pillars-backend/src/types"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

func InsertSingleEarthquake(c echo.Context) error {
	var body types.EarthquakeDB
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, types.Error{
			Error: fmt.Sprintf("Invalid request body: %s", err.Error()),
			Code:  errors.InvalidRequestBody,
		})
	}
	if _, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).InsertOne(context.Background(), body); err != nil {
		return c.JSON(500, types.Error{
			Error: "Failed to insert earthquake. Error: " + err.Error(),
			Code:  errors.FailedToInsertEarthquake,
		})
	}
	return c.NoContent(204)
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func GetAllEarthquakes() {
	var msg kafka.Message
	writer := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  []string{constants.KAFKA_ADDR},
			Topic:    constants.KAFKA_EARTHQUAKE_STREAM_TOPIC,
			Balancer: &kafka.LeastBytes{},
		},
	)
	defer writer.Close()
	cursor, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).Find(context.Background(), bson.D{})
	if err != nil {
		errors.WriteKafkaProducerError(writer, err, context.Background())
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var earthquake types.EarthquakeDB
		if err := cursor.Decode(&earthquake); err != nil {
			errors.WriteKafkaProducerError(writer, err, context.Background())
			return
		}
		msg = kafka.Message{
			Key:   []byte(earthquake.ID.Hex()),
			Value: []byte(cursor.Current.String()),
		}
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			errors.WriteKafkaProducerError(writer, err, context.Background())
			return
		}
	}

}

func GetLatestEarthquakes() {
	var msg kafka.Message
	writer := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  []string{constants.KAFKA_ADDR},
			Topic:    constants.KAFKA_EARTHQUAKE_STREAM_TOPIC,
			Balancer: &kafka.LeastBytes{},
		},
	)
	defer writer.Close()
	cursor, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).Watch(context.Background(), mongo.Pipeline{})
	if err != nil {
		errors.WriteKafkaProducerError(writer, err, context.Background())
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		elems, err := cursor.Current.Elements()
		if err != nil {
			errors.WriteKafkaProducerError(writer, err, context.Background())
			return
		}
		// iterate over the cursor's current elements
		// which corresponds to earthquakes grouped by grid
		for _, elem := range elems {
			// there is two fields in each group
			// one is _id and the other is earthquakes
			if elem.Key() == "fullDocument" {
				// get the earthquakes array
				earthquake := elem.Value().Document()
				// validate the earthquakes array, if not a BSON array return an error
				if err := earthquake.Validate(); err != nil {
					errors.WriteKafkaProducerError(writer, err, context.Background())
					return
				}
				msg = kafka.Message{
					Key:   []byte(earthquake.Lookup("_id").String()),
					Value: []byte(earthquake.String()),
				}
				if err := writer.WriteMessages(context.Background(), msg); err != nil {
					errors.WriteKafkaProducerError(writer, err, context.Background())
					return
				}
			}
		}
	}
}
func GetEarthquakeStream(c echo.Context) error {
	var err error
	// grid_radius := c.QueryParam("grid_radius")
	// var grid_radius_f64 float64 = constants.STREAM_EARTHQUAKE_DEFAULT_GRID_SIZE
	// if grid_radius != "" {
	// 	grid_radius_f64, err = strconv.ParseFloat(grid_radius, 64)
	// 	if err != nil {
	// 		return c.JSON(400, types.Error{
	// 			Error: "Invalid grid_radius. Error: " + err.Error(),
	// 			Code:  errors.InvalidRequestBody,
	// 		})
	// 	}
	// }
	// Kafka connection and consumer setup
	// Replace with your specific Kafka configuration
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{constants.KAFKA_ADDR},
			Topic:     constants.KAFKA_EARTHQUAKE_STREAM_TOPIC,
			Partition: constants.KAFKA_EARTHQUAKE_PARTITION,
			MaxBytes:  constants.KAFKA_READER_MAX_BYTES, // 10 MB
		},
	)
	reader.SetOffset(kafka.LastOffset)
	go GetLatestEarthquakes()
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	var sentIDs []primitive.ObjectID
	// we will determine to show the earthquake or not based on the magnitude
	for {
		// no TO ctx, because we want to keep the connection open
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read message from Kafka. Error: %s", err.Error()))
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return err
		}
		slog.Info(fmt.Sprintf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value)))
		if string(m.Key) == constants.KAFKA_PRODUCER_ERROR_KEY {
			slog.Error(fmt.Sprintf("Kafka producer error: %s", string(m.Value)))
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return err
		}
		var earthquake types.EarthquakeDB
		err = bson.UnmarshalExtJSON(m.Value, true, &earthquake)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to unmarshal message to EarthquakeDB. Error: %s", err.Error()))
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return err
		}
		if !slices.Contains(sentIDs, earthquake.ID) {
			marshalled_bytes, err := json.Marshal(earthquake)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to marshal message to JSON. Error: %s", err.Error()))
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return err
			}
			if err := ws.WriteMessage(websocket.BinaryMessage, marshalled_bytes); err != nil {
				slog.Error(fmt.Sprintf("Failed to write message to websocket. Error: %s", err.Error()))
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return err
			}
			sentIDs = append(sentIDs, earthquake.ID)
		}

	}
}

func GetEarthquakeCount(c echo.Context) error {
	filter_param := c.QueryParam("filter")
	var filter bson.D
	if filter_param != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: "poll coffee"}}}}
	} else {
		filter = bson.D{}
	}
	count, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).CountDocuments(context.Background(), filter)
	if err != nil {
		return c.JSON(500, types.Error{
			Error: "Failed to get earthquake count. Error: " + err.Error(),
			Code:  errors.FailedToGetEarthquakeCountFromMongo,
		})
	}
	return c.JSON(200, types.EarthquakeCount{
		Count: count,
	})
}

// not stream this time
func GetEarthquakesPaged(c echo.Context) error {
	var err error
	limit := c.QueryParam("limit")
	var limit_int int
	if limit == "" {
		limit_int = constants.MONGODB_PAGINATION_DEFAULT_LIMIT
	} else {
		limit_int, err = strconv.Atoi(limit)
		if err != nil {
			return c.JSON(400, types.Error{
				Error: "Invalid page. Error: " + err.Error(),
				Code:  errors.InvalidRequestBody,
			})
		}
	}
	offset := c.QueryParam("offset")
	var offset_int int
	if offset == "" {
		offset_int = 0
	} else {
		offset_int, err = strconv.Atoi(offset)
		if err != nil {
			return c.JSON(400, types.Error{
				Error: "Invalid page_size. Error: " + err.Error(),
				Code:  errors.InvalidRequestBody,
			})
		}
	}
	filter_param := c.QueryParam("filter")
	var pipeline mongo.Pipeline
	if filter_param != "" {
		pipeline = append(pipeline,
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "$text",
							Value: bson.D{
								{Key: "$search",
									Value: filter_param},
							},
						},
					},
				},
			},
		)
	}
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: offset_int}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit_int}})
	var earthquakes []types.EarthquakeDB
	cursor, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.JSON(500, types.Error{
			Error: "Failed to get earthquakes. Error: " + err.Error(),
			Code:  errors.FailedToGetEarthquakes,
		})
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var earthquake types.EarthquakeDB
		if err := cursor.Decode(&earthquake); err != nil {
			return c.JSON(500, types.Error{
				Error: "Failed to decode data. Error: " + err.Error(),
				Code:  errors.FailedToDecodeData,
			})
		}
		earthquakes = append(earthquakes, earthquake)
	}
	return c.JSON(200, earthquakes)
}

func InsertBulkEarthquakes(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["usgs-dump"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		data, err := io.ReadAll(src)
		if err != nil {
			return err
		}
		var body types.USGSEarthquake
		if err := json.Unmarshal(data, &body); err != nil {
			return c.JSON(400, types.Error{
				Error: fmt.Sprintf("Invalid request body: %s", err.Error()),
				Code:  errors.InvalidRequestBody,
			})
		}
		var earthquakes []types.EarthquakeDB
		for _, feature := range body.Features {
			earthquake_time, err := strconv.ParseInt(fmt.Sprint(feature.Properties.Time), 10, 64)
			if err != nil {
				return c.JSON(400, types.Error{
					Error: fmt.Sprintf("Invalid time: %s", err.Error()),
					Code:  errors.InvalidRequestBody,
				})
			}
			earthquake := types.EarthquakeDB{
				ID:    primitive.NewObjectID(),
				Mag:   feature.Properties.Mag,
				Place: feature.Properties.Place,
				Time:  primitive.DateTime(earthquake_time),
				Location: types.EarthquakeLocation{
					Type:        "Point",
					Coordinates: []float64{feature.Geometry.Coordinates[0], feature.Geometry.Coordinates[1]},
				},
			}
			earthquakes = append(earthquakes, earthquake)
		}
		var documents []interface{}
		for _, earthquake := range earthquakes {
			documents = append(documents, earthquake)
		}
		if _, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).InsertMany(context.Background(), documents); err != nil {
			return c.JSON(500, types.Error{
				Error: "Failed to insert earthquakes. Error: " + err.Error(),
				Code:  errors.FailedToInsertEarthquake,
			})
		}
	}
	return c.NoContent(204)
}

func WipeAllEarthquakes(c echo.Context) error {
	if _, err := db.GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).DeleteMany(context.Background(), bson.D{}); err != nil {
		return c.JSON(500, types.Error{
			Error: "Failed to wipe all earthquakes. Error: " + err.Error(),
			Code:  errors.FailedToDeleteDataFromMongo,
		})
	}
	return c.NoContent(204)
}
