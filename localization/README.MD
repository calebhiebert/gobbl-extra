# GOBBL Localization

This package helps with localizing chatbots built with GOBBL.
It's built using the wonderful https://github.com/nicksnyder/go-i18n

## Usage

First, a bundle must be created
```go

import (
  "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"  
)

func main() {
  bundle := &i18n.Bundle{DefaultLanguage: language.English}
  bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

  // Load language files
  bundle.MustLoadMessageFile("assets/lang/en-US.json")
}

```

Then, it must be required as middleware

```go
  gobblr := gbl.New()
  gobblr.Use(glocalize.Middleware(bundle))
```

## Translation in handlers

Translation happens via a few different functions, the simplest one: T

```go
func HExampleHandler(c *gbl.Context) {
  l := glocalize.GetCurrentLocalization(c)

  fmt.Println(l.T("language.key"))
}
```