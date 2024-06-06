package responses

import "github.com/NorskHelsenett/ror/pkg/apicontracts"

type Clusters struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	Data []apicontracts.Cluster `json:"data"`

	TotalCount int `json:"totalCount"`
}

type Projects struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	Data []apicontracts.Project `json:"data"`

	TotalCount int `json:"totalCount"`
}
