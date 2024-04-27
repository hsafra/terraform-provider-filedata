// Copyright (c) Harel Safra
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
)

func withEnvironmentOverrideString(currentValue, envOverrideKey string) string {
	envValue, ok := os.LookupEnv(envOverrideKey)
	if ok {
		return envValue
	}

	return currentValue
}
