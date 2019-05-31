package ab

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/calebhiebert/gobbl"
)

// Test is an object that stores information about a single AB test
type Test struct {

	// Type is the type of AB test taking place, this string should be
	// unique for each AB test performed
	Type string

	// Variations is a slice of handlers that represents all the different variations
	// that a single test can have
	Variations TestList

	// Probability is a []float of numbers between 0 and 1 that notes the probability of
	// each variation ocurring. If Probability is nil, each variation has an equal chance
	Probability []float64
}

// Tester stores all the AB tests for a single project
type Tester struct {
	testArr []Test
	idgen   IDGenFunc
}

// TestList is a list of gobbl handlers to be treated as variations for a test
type TestList []gbl.MiddlewareFunction

// IDGenFunc is a function that generates an id string from a gobbl context
type IDGenFunc func(c *gbl.Context) string

// TestSuite is a map of test types with the value being the test variation
// chosen for the user
type TestSuite map[string]int

// New will create a new ABTester instance
func New() *Tester {
	return &Tester{
		testArr: []Test{},
	}
}

// NewCustom returns a tester that uses a custom
func NewCustom(idgen IDGenFunc) *Tester {
	return &Tester{
		testArr: []Test{},
		idgen:   idgen,
	}
}

// Register will register a new AB test
func (t *Tester) Register(tests ...Test) {
	for _, test := range tests {
		if test.Variations == nil {
			panic("no variations supplied")
		}

		if test.Type == "" {
			panic("cannot supply blank test type")
		}

		if test.Probability != nil && len(test.Probability) != len(test.Variations) {
			panic("different number of tests and probabilities")
		}

		t.testArr = append(t.testArr, test)
	}
}

// AB will return a gobbl middleware generated for the given test type.
// when the bot calls this middleware, it will decide which handler to call based on the user
func (t *Tester) AB(testType string) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		var abSuite TestSuite

		if c.HasFlag("__abSuite") {
			abSuite = c.GetFlag("__abSuite").(TestSuite)
			c.Tracef("Loaded AB Suite %+v", abSuite)
		} else {

			id := c.User.ID

			if t.idgen != nil {
				id = t.idgen(c)
			}

			suite, err := genSuite(t, id)
			if err != nil {
				c.Errorf("Error generating AB suite for user %v", err)
			}

			abSuite = suite

			c.Flag("__abSuite", abSuite)
			c.Tracef("Generated AB Suite %+v", abSuite)
		}

		var test Test

		for _, searchTest := range t.testArr {
			if searchTest.Type == testType {
				test = searchTest
				break
			}
		}

		if test.Type == "" {
			panic("invalid test type " + testType)
		}

		variationIdx, exists := abSuite[testType]
		if !exists {
			panic(fmt.Sprintf("test type %s was not registered", testType))
		}
		// TODO cache generated suite

		test.Variations[variationIdx](c)
	}
}

// GetSuite returns a suite of ab test selections for a given id
func (t *Tester) GetSuite(id string) TestSuite {
	suite, err := genSuite(t, id)
	if err != nil {
		panic(err)
	}

	return suite
}

// genSuite will generate the entire suite of test variations for a given id
func genSuite(t *Tester, id string) (TestSuite, error) {
	rand, err := getRand(id)
	if err != nil {
		return nil, err
	}

	suite := TestSuite{}

	for _, test := range t.testArr {
		if test.Probability == nil {
			suite[test.Type] = rand.Intn(len(test.Variations))
		} else {
			if len(test.Probability) != len(test.Variations) {
				panic("incorrect number of probabilities provided " + test.Type)
			}

			var totalWeight float64

			for _, prob := range test.Probability {
				if prob < 0 {
					panic("probability cannot be a negative number")
				}

				totalWeight += prob
			}

			if totalWeight <= 0 {
				panic("probabilities must add up to a positive number")
			}

			randomNumber := rand.Float64() * totalWeight

			for i := 0; i < len(test.Variations); i++ {
				randomNumber -= test.Probability[i]

				if randomNumber <= 0 {
					suite[test.Type] = i
					break
				}
			}
		}
	}

	return suite, nil
}

// getRand will return a rand object that is seeded
func getRand(id string) (*rand.Rand, error) {
	i, err := hashToInt(id)
	if err != nil {
		return nil, err
	}

	return rand.New(rand.NewSource(i)), nil
}

func hashToInt(id string) (int64, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(id))
	if err != nil {
		return 0, err
	}

	hashBytes := hash.Sum(nil)

	i := int64(binary.LittleEndian.Uint64(hashBytes[:8]))

	return i, nil
}
