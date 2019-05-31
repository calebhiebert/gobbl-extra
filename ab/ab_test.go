package ab

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/calebhiebert/gobbl"
)

func TestHashToInt(t *testing.T) {

	h1 := func(c *gbl.Context) {
		c.Info("Handler 1")
	}

	h2 := func(c *gbl.Context) {
		c.Info("Handler 2")
	}

	h3 := func(c *gbl.Context) {
		c.Info("Handler 3")
	}

	ids := []string{}

	for i := 0; i < 100; i++ {
		ids = append(ids, strconv.Itoa(i))
	}

	for _, id := range ids {

		ab := New()

		ab.Register(
			Test{Type: "menu-order", Variations: TestList{h1}},
			Test{Type: "mailing-list", Variations: TestList{h1, h2}},
			Test{Type: "prob", Variations: TestList{h1, h2, h3}, Probability: []float64{0.5, 0.5, 0.01}},
		)

		suite, err := genSuite(ab, id)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(suite)
	}
}
