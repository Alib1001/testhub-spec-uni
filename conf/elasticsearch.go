package conf

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: "elastic",
		Password: "10011001",
	}

	var err error
	EsClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	log.Println("Elasticsearch client initialized successfully")
}
