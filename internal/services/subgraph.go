package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/etherfi/eth-dashboard/internal/models"
)

type SubgraphClient struct {
	etherfiURL string
	lidoURL    string
	client     *http.Client
}

func NewSubgraphClient(etherfiURL, lidoURL string) *SubgraphClient {
	return &SubgraphClient{
		etherfiURL: etherfiURL,
		lidoURL:    lidoURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *SubgraphClient) Query(ctx context.Context, url string, query string, variables map[string]interface{}, result interface{}) error {
	payload := models.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("subgraph request failed: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result); err != nil {
		return err
	}
	return nil
}

const etherfiStakeQuery = `
query GetUserStakes($user: String!) {
  stakers(where: { id: $user }) {
    id
    totalStaked
    totalWithdrawn
    stakes {
      amount
      timestamp
      validator {
        id
        status
      }
    }
  }
}
`

const protocolTVLQuery = `
query GetProtocolTVL($startTime: Int!) {
  protocolMetrics(
    where: { timestamp_gt: $startTime }
    orderBy: timestamp
    orderDirection: asc
  ) {
    timestamp
    totalValueLocked
    totalStakers
  }
}
`

const lidoStakeQuery = `
query GetLidoStake($user: String!) {
  users(where: { id: $user }) {
    id
    totalStaked
    totalClaimed
  }
}
`

type etherfiStakeResponse struct {
	Data struct {
		Stakers []struct {
			ID            string `json:"id"`
			TotalStaked   string `json:"totalStaked"`
			TotalWithdrawn string `json:"totalWithdrawn"`
			Stakes        []struct {
				Amount    string `json:"amount"`
				Timestamp string `json:"timestamp"`
				Validator struct {
					ID     string `json:"id"`
					Status string `json:"status"`
				} `json:"validator"`
			} `json:"stakes"`
		} `json:"stakers"`
	} `json:"data"`
}

type lidoStakeResponse struct {
	Data struct {
		Users []struct {
			ID          string `json:"id"`
			TotalStaked string `json:"totalStaked"`
			TotalClaimed string `json:"totalClaimed"`
		} `json:"users"`
	} `json:"data"`
}

type protocolTVLResponse struct {
	Data struct {
		ProtocolMetrics []struct {
			Timestamp       string `json:"timestamp"`
			TotalValueLocked string `json:"totalValueLocked"`
			TotalStakers    string `json:"totalStakers"`
		} `json:"protocolMetrics"`
	} `json:"data"`
}

func (c *SubgraphClient) FetchEtherFiStake(ctx context.Context, user string) (etherfiStakeResponse, error) {
	var resp etherfiStakeResponse
	err := c.Query(ctx, c.etherfiURL, etherfiStakeQuery, map[string]interface{}{"user": user}, &resp)
	return resp, err
}

func (c *SubgraphClient) FetchLidoStake(ctx context.Context, user string) (lidoStakeResponse, error) {
	var resp lidoStakeResponse
	err := c.Query(ctx, c.lidoURL, lidoStakeQuery, map[string]interface{}{"user": user}, &resp)
	return resp, err
}

func (c *SubgraphClient) FetchEtherFiTVL(ctx context.Context, startTime int) (protocolTVLResponse, error) {
	var resp protocolTVLResponse
	err := c.Query(ctx, c.etherfiURL, protocolTVLQuery, map[string]interface{}{"startTime": startTime}, &resp)
	return resp, err
}

func (c *SubgraphClient) FetchLidoTVL(ctx context.Context, startTime int) (protocolTVLResponse, error) {
	var resp protocolTVLResponse
	err := c.Query(ctx, c.lidoURL, protocolTVLQuery, map[string]interface{}{"startTime": startTime}, &resp)
	return resp, err
}
