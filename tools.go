// +build tools

package main

// This file defines tool dependencies for the module.
// See: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/gojuno/minimock/v3/cmd/minimock"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
)
