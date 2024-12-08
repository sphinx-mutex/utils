package stacks

// Stackup stacks the handlers in the order they are passed
func Stackup[T Scenario](stackables ...Stackable[T]) Stackable[T] {
	return func(handler Handler[T]) Handler[T] {
		for i := len(stackables) - 1; i >= 0; i-- {
			handler = stackables[i](handler)
		}
		return handler
	}
}
