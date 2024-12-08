package stacks

// SwitchReducer is the type of the function that decides which stackable to call by returning the name of the stackable
type SwitchReducer[T Scenario] func(scenario T) string

// Switch is a stackable that switches between stackables based on the switcher
// It is acting like a Router in http context
func Switch[T Scenario](switcher SwitchReducer[T], stackables map[string]Stackable[T]) Stackable[T] {
	return func(next Handler[T]) Handler[T] {
		return func(scenario T) error {
			return stackables[switcher(scenario)](next)(scenario)
		}
	}
}
