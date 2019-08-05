package types

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
)

var (
	// event type for invariant checking
	EventTypeInvariantCheck = "invariant_check"
)

// An Invariant is a function which tests a particular invariant.
// The invariant returns a descriptive message about what happened
// and total tokens, the mint module always return a negative value
// and a boolean indicating whether the invariant has been broken.
// The simulator will then halt and print the logs.
type Invariant func(ctx context.Context) (string, types.BaseCoins, bool)

// Invariants defines a group of invariants
type Invariants []Invariant

// expected interface for registering invariants
type InvariantRegistry interface {
	RegisterInvarRoute(moduleName, route string, invar Invariant)
}

// FormatInvariant returns a standardized invariant message along with
// a boolean indicating whether the invariant has been broken.
func FormatInvariant(module, name, msg string, coins types.BaseCoins, broken bool) (string, types.BaseCoins, bool) {
	return fmt.Sprintf("%s: %s invariant\n%sinvariant broken: %v\n",
		module, name, msg, broken), coins, broken
}

// invariant route
type InvarRoute struct {
	ModuleName string
	Route      string
	Invar      Invariant
}

// NewInvarRoute - create an InvarRoute object
func NewInvarRoute(moduleName, route string, invar Invariant) InvarRoute {
	return InvarRoute{
		ModuleName: moduleName,
		Route:      route,
		Invar:      invar,
	}
}

// get the full invariance route
func (i InvarRoute) FullRoute() string {
	return i.ModuleName + "/" + i.Route
}
