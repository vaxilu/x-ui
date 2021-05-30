package xray

import (
	"encoding/json"
	"x-ui/util/json_util"
)

type InboundConfig struct {
	Listen         json.RawMessage `json:"listen"` // listen 不能为空字符串
	Port           int             `json:"port"`
	Protocol       string          `json:"protocol"`
	Settings       json.RawMessage `json:"settings"`
	StreamSettings json.RawMessage `json:"streamSettings"`
	Tag            string          `json:"tag"`
	Sniffing       json.RawMessage `json:"sniffing"`
}

func (i *InboundConfig) MarshalJSON() ([]byte, error) {
	return json_util.MarshalJSON(i)
}
