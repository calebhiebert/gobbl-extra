package gobbldflow

import (
	"encoding/json"
	"io/ioutil"
)

// ServiceKey represents a google account service key file
type ServiceKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// LoadServiceKeyFile will load a google account service key from a file
// and return the proper dialogflow config
func LoadServiceKeyFile(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return LoadServiceKeyJSONBytes(file)
}

// LoadServiceKeyJSONBytes will create a dialogflow config from some json
// this json should be a google service account key file
func LoadServiceKeyJSONBytes(jsonBytes []byte) (*Config, error) {
	var serviceKey ServiceKey

	err := json.Unmarshal(jsonBytes, &serviceKey)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServiceAccount: serviceKey.ClientEmail,
		PrivateKey:     []byte(serviceKey.PrivateKey),
		ProjectID:      serviceKey.ProjectID,
	}, nil
}
