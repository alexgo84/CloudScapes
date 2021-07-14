package wire

import (
	"encoding/json"
	"errors"
)

type ConfigMap struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type ConfigMaps []ConfigMap

func (cm *ConfigMaps) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &cm)
}
