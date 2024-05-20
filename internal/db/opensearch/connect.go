package connect

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/shaik80/ODIW/config"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	lp "github.com/shaik80/ODIW/utils/logger"
)

var Client *opensearchapi.Client

// InitOpenSearchClient initializes the global OpenSearch client
func InitOpenSearchClient(cfg config.Config) error {
	opensearchCfg := opensearch.Config{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}, Addresses: []string{
			fmt.Sprintf("http://%s:%s", cfg.OpenSearch.Host, cfg.OpenSearch.Port),
		},
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	}
	client, err := opensearchapi.NewClient(opensearchapi.Config{Client: opensearchCfg})
	if err != nil {
		return fmt.Errorf("error creating OpenSearch client: %w", err)
	}
	lp.Logs.Infof("openseach connected successfully")
	Client = client

	ctx := context.Background()
	// Print OpenSearch version information on console.
	infoResp, err := client.Info(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Cluster INFO:\n  Cluster Name: %s\n  Cluster UUID: %s\n  Version Number: %s\n", infoResp.ClusterName, infoResp.ClusterUUID, infoResp.Version.Number)

	return nil
}
