package glocalize

import (
	"github.com/calebhiebert/gobbl"
)

// ArgModifierFunc takes in a list of supplied arguments and modifies them
// this can be used to add default translation values, etc...
type ArgModifierFunc func(args A, l *Localization) A

// Middleware will set a localization object on the context
// this can later be loaded to do translations.
// This middleware will pull the language from the "lang" flag
// so make sure it's set
func Middleware(config *LocalizationConfig) gbl.MiddlewareFunction {
	return func(c *gbl.Context) {
		var loc *Localization

		if c.HasFlag("lang") {
			loc = MustGetLocalization(c.GetStringFlag("lang"), config.Bundle)
		} else {
			loc = MustGetLocalization("en-US", config.Bundle)
		}

		loc.argModifier = config.ArgModifier

		c.Flag("__localizer", loc)
		c.Next()
	}
}

// GetCurrentLocalization will retrieve the current localization from the gobbl context
func GetCurrentLocalization(c *gbl.Context) *Localization {
	if c.HasFlag("__localizer") {
		return c.GetFlag("__localizer").(*Localization)
	}

	panic("Localizer not present on context")
}
