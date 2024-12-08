package stacks

import "testing"

func TestStackup(t *testing.T) {

	type executionLog struct {
		name  string
		param string
	}

	var logs []executionLog

	var s1 Stackable[string] = func(next Handler[string]) Handler[string] {
		return func(scenario string) error {
			logs = append(logs, executionLog{name: "s1", param: scenario})
			return next(scenario + "s1")
		}
	}

	var s2 Stackable[string] = func(next Handler[string]) Handler[string] {
		return func(scenario string) error {
			logs = append(logs, executionLog{name: "s2", param: scenario})
			return next(scenario + "s2")
		}
	}

	var s3 Stackable[string] = func(next Handler[string]) Handler[string] {
		return func(scenario string) error {
			logs = append(logs, executionLog{name: "s3", param: scenario})
			return next(scenario + "s3")
		}
	}

	err := Stackup(s1, s2, s3)(func(scenario string) error {
		logs = append(logs, executionLog{name: "handler", param: scenario})
		return nil
	})(`test`)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(logs) != 4 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}

	if logs[0].name != "s1" {
		t.Errorf("Expected s1, got %s", logs[0].name)
	}

	if logs[0].param != "test" {
		t.Errorf("Expected test, got %s", logs[0].param)
	}

	if logs[1].name != "s2" {
		t.Errorf("Expected s2, got %s", logs[1].name)
	}

	if logs[1].param != "tests1" {
		t.Errorf("Expected tests1, got %s", logs[1].param)
	}

	if logs[2].name != "s3" {
		t.Errorf("Expected s3, got %s", logs[2].name)
	}

	if logs[2].param != "tests1s2" {
		t.Errorf("Expected tests1s2, got %s", logs[2].param)
	}

	if logs[3].name != "handler" {
		t.Errorf("Expected handler, got %s", logs[3].name)
	}

	if logs[3].param != "tests1s2s3" {
		t.Errorf("Expected tests1s2s3, got %s", logs[3].param)
	}
}
