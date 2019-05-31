package gobbldflow

type QueryConfig struct {
	QueryParams *QueryParams `json:"queryParams,omitempty"`
	QueryInput  *QueryInput  `json:"queryInput,omitempty"`
}

type QueryInput struct {
	Text  *Text  `json:"text,omitempty"`
	Event *Event `json:"event,omitempty"`
}

type Event struct {
	Name         string                 `json:"name"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	LanguageCode string                 `json:"languageCode"`
}

type Text struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type QueryParams struct {
	TimeZone      string       `json:"timeZone,omitempty"`
	GeoLocation   *GeoLocation `json:"geoLocation,omitempty"`
	Contexts      []Context    `json:"contexts,omitempty"`
	ResetContexts bool         `json:"resetContexts"`
	Payload       interface{}  `json:"payload"`
}

type Context struct {
	Name          string `json:"name"`
	LifespanCount int    `json:"lifespanCount,omitempty"`
	Parameters    map[string]interface{}
}

type GeoLocation struct {
	Lat float64
	Lng float64
}
