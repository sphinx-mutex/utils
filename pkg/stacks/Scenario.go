package stacks

// Scenario is the type of the scenario that is passed to the handler
type Scenario interface{}

// Handler is the type of the function that is called when the scenario is triggered
type Handler[T Scenario] func(scenario T) error

// Stackable is the type of the function that stacks the handlers
type Stackable[T Scenario] func(next Handler[T]) Handler[T]
