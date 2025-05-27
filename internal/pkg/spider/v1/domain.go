// internal/pkg/spider/v1/domain.go - Enhanced domain model with JSON support
package spider

import "encoding/json"

type Domain struct {
	URL    string `json:"url"`
	Name   string `json:"name"`
	TLD    string `json:"tld"`
	Status string `json:"status"`
}

func (d Domain) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d Domain) CSVRow() []string {
	return []string{d.URL, d.Name, d.TLD, d.Status}
}
