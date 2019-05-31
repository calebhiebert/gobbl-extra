# GOBBL AB

This is a package with some utilities that make adding AB testing to a GOBBL project easier.

## Usage

```go

// Create a new tester instance
tester := ab.New()

// Register all test types
tester.Register(
  // This line notes that the menu-order ab test has just 1 variation
  Test{Type: "menu-order", Variations: TestList{h1}},

  // This line notes that the mailing-list ab test has 2 variations
  // each variation has a 50% chance of being chosen
  Test{Type: "mailing-list", Variations: TestList{h1, h2}},

  // This line notes that the greeting ab test has 3 variations
  // the first variation has a 50% chance of being chosen
  // the second variation has a 49% chance of being chosen
  // and the last variation only has a 1% chance of being chosen 
  Test{Type: "greeting", Variations: TestList{h1, h2, h3}, Probability: []float64{0.5, 0.49, 0.01}},
)

// This will automatically select which variation of menu-order to use based on the user
gobblr.Use(tester.AB("menu-order"))
```

## Custom ID Selection

It is possible to change the parameter used to select a user's ab suite. By default the context `User.ID` parameter is used. To change the id, just use `NewCustom` instead of `New` when creating a tester.

```go
// This example will use the user's first an last name to generate their ab suite
// this means that users with the same first and last name will be ab
// tested the exact same way
ab := NewCustom(func(c *gbl.Context) string {
  return c.User.FirstName + c.User.LastName
})
```

## Retrieve the AB suite

It's possible to generate the ab suite manually, for storage or caching.

```go
ab := New()

// Pass in the ID property manually to generate the suite
suite := ab.GetSuite("123456789")
```