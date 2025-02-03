package tb

import (
	"fmt"
)

func ExpectString(value string) func(string) error {
	return func(actual string) error {
		if actual != value {
			return fmt.Errorf("expected %q, got %q", value, actual)
		}
		return nil
	}
}
