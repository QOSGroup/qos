package types

// hooks for mappers
type Hooks interface {
	// return mapper name who use this hooks
	HookMapper() string
}
