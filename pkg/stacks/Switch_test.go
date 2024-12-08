package stacks

import "testing"

func TestSwitch(t *testing.T) {

	type executionLog struct {
		name  string
		param string
	}

	var logs []executionLog

	var s1 Stackable[string] = func(next Handler[string]) Handler[string] {
		return func(scenario string) error {
			logs = append(logs, executionLog{name: "s1", param: scenario})
			return next(scenario)
		}
	}

	var s2 Stackable[string] = func(next Handler[string]) Handler[string] {
		return func(scenario string) error {
			logs = append(logs, executionLog{name: "s2", param: scenario})
			return next(scenario)
		}
	}

	var stackMap = map[string]Stackable[string]{
		"s1": s1,
		"s2": s2,
	}

	s := Switch[string](func(scenario string) string {
		if scenario == "test" {
			return "s1"
		} else {
			return "s2"
		}
	}, stackMap)

	err := s(func(scenario string) error {
		return nil
	})(`test`)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 logs, got %d", len(logs))
	}

	if logs[0].name != "s1" {
		t.Errorf("Expected s1, got %s", logs[0].name)
	}

	if logs[0].param != "test" {
		t.Errorf("Expected test, got %s", logs[0].param)
	}

	err = s(func(scenario string) error {
		return nil
	})(`test-other`)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	if logs[1].name != "s2" {
		t.Errorf("Expected s2, got %s", logs[1].name)
	}

	if logs[1].param != "test-other" {
		t.Errorf("Expected test-other, got %s", logs[1].param)
	}
}
