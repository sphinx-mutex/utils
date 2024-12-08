package booter

import (
	"fmt"
	"testing"
)

func TestBooter_Register_Get(t *testing.T) {
	booter := NewBooter(make(map[string]Bootable[any]))
	testBootingCount := 0
	booter.Register("test", func(booter *Booter) any {
		testBootingCount++
		return "test"
	})

	booter.Register("test2", func(booter *Booter) any {
		return "test2" + booter.MustGet("test").(string)
	})

	service, ok := booter.Get("test")
	if !ok {
		t.Error("Expected true, got false")
	}

	if service != "test" {
		t.Errorf("Expected test, got %v", service)
	}

	service, ok = booter.Get("test2")
	if !ok {
		t.Error("Expected true, got false")
	}

	if service != "test2test" {
		t.Errorf("Expected test2test, got %v", service)
	}

	service, ok = booter.Get("notfound")

	if ok {
		t.Error("Expected false, got true")
	}

	if service != nil {
		t.Errorf("Expected nil, got %v", service)
	}
}

func TestBooter_MustGet(t *testing.T) {
	booter := NewBooter(make(map[string]Bootable[any]))
	booter.Register("test", func(booter *Booter) any {
		return "test"
	})

	service := booter.MustGet("test")

	if service != "test" {
		t.Errorf("Expected test, got %v", service)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, got nil")
		}
	}()
	booter.MustGet("notfound")
}

func TestBooter_Cache(t *testing.T) {
	booter := NewBooter(make(map[string]Bootable[any]))
	booter.Cache("test", "test")

	service, ok := booter.Get("test")
	if !ok {
		t.Error("Expected true, got false")
	}

	if service != "test" {
		t.Errorf("Expected test, got %v", service)
	}
}

func TestBooter_CallOnce(t *testing.T) {

	booter := NewBooter(make(map[string]Bootable[any]))
	testBootingCount := 0
	booter.Register("test", func(b *Booter) any {
		testBootingCount++
		return "test"
	})

	booter.Get("test")
	booter.Get("test")
	booter.Get("test")

	if testBootingCount != 1 {
		t.Errorf("Expected 1, got %d", testBootingCount)
	}
}

func TestBooter_HandleCyclicDependencies(t *testing.T) {
	booter := NewBooter(make(map[string]Bootable[any]))

	booter.Register("test", func(b *Booter) any {
		return b.MustGet("test3")
	})
	booter.Register("test2", func(b *Booter) any {
		return b.MustGet("test")
	})
	booter.Register("test3", func(b *Booter) any {
		return b.MustGet("test2")
	})

	defer func() {
		if r := recover(); r == nil {
			fmt.Println(r)
			t.Error("Expected panic, got nil")
		}
	}()
	booter.Get("test3")
}
