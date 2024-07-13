package conf

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

var EsClient *elasticsearch.Client

func InitElasticsearch() {
	esHost := os.Getenv("ES_HOSTS")
	if esHost == "" {
		log.Fatalf("ES_HOSTS environment variable not set")
	}

	log.Printf("Using Elasticsearch host: %s", esHost)

	cfg := elasticsearch.Config{
		Addresses: []string{
			esHost,
		},
	}

	var err error
	EsClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	log.Println("Elasticsearch client initialized successfully")
}
