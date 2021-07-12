package wire

type ConfigMap struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}
