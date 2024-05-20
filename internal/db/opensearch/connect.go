package connect

import (
	"fmt"
	"shaik80/ODIW/config"

	"github.com/opensearch-project/opensearch-go"
)

var Client *opensearch.Client

// InitOpenSearchClient initializes the global OpenSearch client
func InitOpenSearchClient() error {
	cfg := config.GetConfig()
	opensearchCfg := opensearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", cfg.OpenSearch.Host, cfg.OpenSearch.Port),
		},
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	}
	client, err := opensearch.NewClient(opensearchCfg)
	if err != nil {
		return fmt.Errorf("error creating OpenSearch client: %w", err)
	}

	Client = client
	return nil
}
