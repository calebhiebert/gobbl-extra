package gobbldflow

import (
	"fmt"
	"math/rand"

	"github.com/calebhiebert/gobbl"
)

// Middleware will return a gobbl compatable middleware
func Middleware(dflow *API) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		var sessionID string

		if c.User.ID != "" {
			sessionID = c.User.ID
		} else {
			sessionID = fmt.Sprintf("%f", rand.Float64())
		}

		res, err := dflow.QueryText(c.Request.Text, sessionID)
		if err != nil {
			c.Errorf("DFLOW %v", err)
			c.Next()
			return
		}

		c.Flag("dflow", res)

		if res.QueryResult.IntentDetectionConfidence >= dflow.config.MinimumConfidence {
			c.Flag("intent", res.QueryResult.Intent.DisplayName)
			c.Flag("intent:score", res.QueryResult.IntentDetectionConfidence)

			for k, v := range res.QueryResult.Parameters {
				c.Flag("dflow:p:"+k, v)
			}
		}

		c.Next()
	}
}
