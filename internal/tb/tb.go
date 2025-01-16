// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tb

import "time"

func RandBool() bool {
	return time.Now().UnixNano()%2 == 0
}
