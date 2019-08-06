package gobbldflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// Config holds all the information required to query the dialogflow
// v2 api
type Config struct {
	ServiceAccount    string
	PrivateKey        []byte
	ProjectID         string
	MinimumConfidence float64
}

// Response is what dialogflow returns for a detect intent query
type Response struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		QueryText                string            `json:"queryText"`
		Action                   string            `json:"action"`
		Parameters               map[string]string `json:"parameters"`
		AllRequiredParamsPresent bool              `json:"allRequiredParamsPresent"`
		FulfillmentText          string            `json:"fulfillmentText"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
			Platform string `json:"platform,omitempty"`
		} `json:"fulfillmentMessages"`
		Intent struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"intent"`
		IntentDetectionConfidence float64 `json:"intentDetectionConfidence"`
		LanguageCode              string  `json:"languageCode"`
	} `json:"queryResult"`
}

// API is a dialogflow API
type API struct {
	client *http.Client
	config *Config
}

// New creates a new Dialogflow API
func New(config *Config) *API {
	conf := &jwt.Config{
		Email:      config.ServiceAccount,
		PrivateKey: config.PrivateKey,
		TokenURL:   google.JWTTokenURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/dialogflow",
		},
	}

	client := conf.Client(context.Background())
	client.Timeout = 10 * time.Second

	return &API{
		client: client,
		config: config,
	}
}

// Query will query the dialogflow API
func (d *API) Query(config *QueryConfig, sessionID string) (*Response, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	res, err := d.client.Post(
		fmt.Sprintf("https://dialogflow.googleapis.com/v2/projects/%s/agent/sessions/%s:detectIntent", d.config.ProjectID, sessionID),
		"application/json",
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return nil, errors.New(string(body))
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// QueryText will query dialogflow with a query string
func (d *API) QueryText(text, sessionID string) (*Response, error) {
	if len(text) > 255 {
		text = string([]rune(text)[:255])
	}

	return d.Query(&QueryConfig{
		QueryInput: &QueryInput{
			Text: &Text{
				Text:         text,
				LanguageCode: "en-US",
			},
		},
	}, sessionID)
}
