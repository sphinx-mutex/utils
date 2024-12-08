package booter

import "fmt"

// Bootable is a function that bootstraps a service
// insde that function, you can use the Booter to get other services if needed
type Bootable[Service interface{}] func(booter *Booter) Service

// Booter is a service bootstrapper that allows for lazy loading of services
type Booter struct {
	// registry is a map of service names to bootable functions
	registry map[string]Bootable[any]

	// cached is a map of service names to resolved services
	cached map[string]any

	// bootSeq is a slice of service names that currently being booted
	// using this bootSeq, we can detect circular dependencies
	//
	// if the booter trying to boot a service that is already in this list
	// it means there is a circular dependency
	bootSeq []string
}

// NewBooter creates a new Booter
func NewBooter(registry map[string]Bootable[any]) *Booter {
	return &Booter{
		registry: registry,
		cached:   make(map[string]any),
		bootSeq:  make([]string, 0),
	}
}

// Register registers a service with the Booter
func (b *Booter) Register(name string, bootable Bootable[any]) {
	b.registry[name] = bootable
}

// Get gets a service by name
func (b *Booter) Get(name string) (any, bool) {

	// if the service is already cached, return it
	if cached, ok := b.cached[name]; ok {
		return cached, true
	}

	// if the service is not registered, return false
	if bootable, ok := b.registry[name]; ok {

		// check for circular dependencies
		for _, booting := range b.bootSeq {

			// if the service is already being booted, it means there is a circular dependency
			if booting == name {
				fullSeq := append(b.bootSeq, name)
				// panic with the circular dependency and the boot sequence
				panic(fmt.Sprintf("circular dependency detected for %s, boot sequnce %v", name, fullSeq))
			}
		}

		// add the service to the boot sequence
		b.bootSeq = append(b.bootSeq, name)

		// boot the service
		b.cached[name] = bootable(b)

		// remove the service from the boot sequence
		b.bootSeq = b.bootSeq[:len(b.bootSeq)-1]

		// return the service
		return b.cached[name], true
	}

	return nil, false
}

// MustGet gets a service by name and panics if it is not found
func (b *Booter) MustGet(name string) any {
	if service, ok := b.Get(name); ok {
		return service
	}

	panic("service not found")
}

// Cache caches a service by name
func (b *Booter) Cache(name string, service any) {
	b.cached[name] = service
}

// Resolved returns a map of all services that have been resolved
func (b *Booter) Resolved() map[string]any {
	return b.cached
}
