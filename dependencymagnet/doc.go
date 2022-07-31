// Copyright Contributors to the Open Cluster Management project
//go:build dependencymagnet
// +build dependencymagnet

// go mod won't pull in code that isn't depended upon, but we have some code we don't depend on from code that must be included
// for our build to work.
package dependencymagnet

import (
	_ "github.com/apg/patter"
	_ "github.com/wadey/gocovmerge"
)
