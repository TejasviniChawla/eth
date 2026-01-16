package models

import "time"

type Protocol string

const (
	ProtocolEtherFi   Protocol = "ether.fi"
	ProtocolLido      Protocol = "lido"
	ProtocolRocketPool Protocol = "rocketpool"
)

type StakingPosition struct {
	Protocol     Protocol `json:"protocol"`
	StakedAmount float64  `json:"stakedAmount"`
	Rewards      float64  `json:"rewards"`
	CurrentValue float64  `json:"currentValue"`
}

type YieldPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Rewards   float64   `json:"rewards"`
	Protocol  Protocol  `json:"protocol"`
}

type ProtocolStats struct {
	Protocol  Protocol  `json:"protocol"`
	CurrentAPY float64  `json:"currentApy"`
	TVL       float64   `json:"tvl"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
