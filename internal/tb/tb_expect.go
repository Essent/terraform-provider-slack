// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func ExpectBool(value bool) func(string) error {
	return func(actual string) error {
		if actual != fmt.Sprintf("%t", value) {
			return fmt.Errorf("expected %t, got %q", value, actual)
		}
		return nil
	}
}
