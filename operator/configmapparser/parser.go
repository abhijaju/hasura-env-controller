package parser

import "encoding/json"

type ConfigMapData struct {
	DeploymentName string            `json:"deploymentName"`
	EnvVars        map[string]string `json:"envVars"`
}

func NewConfigMapData(content string) (*ConfigMapData, error) {
	cmData := ConfigMapData{}
	if err := json.Unmarshal([]byte(content), &cmData); err != nil {
		return &cmData, err
	}
	return &cmData, nil
}
