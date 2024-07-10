package conf

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

func InitElasticsearch() {
	esHost := os.Getenv("ES_HOSTS")
	esUser := os.Getenv("ES_USER")
	esPassword := os.Getenv("ES_PASSWORD")
	cfg := elasticsearch.Config{
		Addresses: []string{
			esHost,
		},
		Username: esUser,
		Password: esPassword,
	}

	var err error
	EsClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	log.Println("Elasticsearch client initialized successfully")
}
