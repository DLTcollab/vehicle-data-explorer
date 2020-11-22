package elasticsearch

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/DLTcollab/vehicle-data-explorer/models/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type DeviceLog struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	DeviceID  string `json:"deviceID"`
}

var (
	ElasticClient *elasticsearch.Client
)

func init() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.ConfigHandler.GetString("ELASTICSEARCH_HOST"),
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	ElasticClient = client
}

func InsertDeviceLog(deviceID string, message string) bool {
	CurrentTime := time.Now().Unix()

	deviceLog := DeviceLog{
		Timestamp: CurrentTime,
		Message:   message,
		DeviceID:  deviceID,
	}

	// Marshal the struct to JSON and check for errors
	jsonStr, err := json.Marshal(deviceLog)
	if err != nil {
		log.Println("json.Marshal ERROR:", err)
		return false
	}

	documentID := strconv.FormatInt(CurrentTime, 10)
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      deviceID,
		DocumentID: documentID,
		Body:       strings.NewReader(string(jsonStr)),
		Refresh:    "true",
	}

	// Return an API response object from request
	res, err := req.Do(context.Background(), ElasticClient)
	if err != nil {
		log.Fatalf("IndexRequest ERROR: %s", err)
		return false
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("%s ERROR indexing document ID=%s", res.Status(), documentID)
		return false
	}

	return true
}

func QueryDeviceLog(deviceID string) []DeviceLog {

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var (
		r envelopeResponse
	)

	// Perform the search request.
	res, err := ElasticClient.Search(
		ElasticClient.Search.WithContext(context.Background()),
		ElasticClient.Search.WithIndex(deviceID),
		ElasticClient.Search.WithTrackTotalHits(true),
		ElasticClient.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Error getting response: %s", err)
		return nil
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); res.IsError() || err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil
	}

	var logs []DeviceLog
	for _, hit := range r.Hits.Hits {
		deviceLog := DeviceLog{}
		if err := json.Unmarshal(hit.Source, &deviceLog); err != nil {
			continue
		}
		logs = append(logs, deviceLog)
	}
	return logs
}
