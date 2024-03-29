package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Log(index string, logger interface{}) error {
	log.SetFlags(0)

	var (
		wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	urlElastic := validateEnv(os.Getenv("ELASTIC_URL"))

	cfg := elasticsearch.Config{
		Addresses: []string{
			urlElastic,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return err
	}

	_, err = es.Info()
	if err != nil {
		log.Printf("Setup ElasticSearch: failed cause, %s", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Set up the request object.
		body, err := json.Marshal(logger)
		if err != nil {
			log.Fatal("Failed @marshal")
		}

		req := esapi.IndexRequest{
			Index:      index,
			DocumentID: strconv.Itoa(int(time.Now().UnixNano())),
			Body:       bytes.NewReader(body),
			Refresh:    "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
	}()
	wg.Wait()
	return err
}

func validateEnv(url string) string {
	if url == "" {
		return "http://localhost:9200/"
	} else {
		return url
	}
}
