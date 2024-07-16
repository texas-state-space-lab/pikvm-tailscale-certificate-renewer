package helpers_test

import (
	"regexp"
	"testing"

	"github.com/nateinaction/pikvm-tailscale-cert-renewer/internal/helpers"
)

func TestSetLine(t *testing.T) {
	t.Parallel()

	contents := []string{
		"line 1",
		"line 2",
		"line 3",
	}

	regex := regexp.MustCompile(`^line \d$`)
	certLine := "new line"

	expected := []string{
		"new line",
		"line 2",
		"line 3",
	}

	result := helpers.SetLine(contents, regex, certLine)

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, but got %d", len(expected), len(result))
	}

	for i, line := range result {
		if line != expected[i] {
			t.Errorf("Expected line '%s', but got '%s'", expected[i], line)
		}
	}
}
