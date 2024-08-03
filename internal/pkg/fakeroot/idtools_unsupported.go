//go:build !linux || !libsubid || !cgo
// +build !linux !libsubid !cgo

package fakeroot

import (
	"github.com/apptainer/apptainer/internal/pkg/util/user"
)

func readSubuid(user *user.User) ([]*Entry, error) {
	return make([]*Entry, 0), nil
}

func readSubgid(user *user.User) ([]*Entry, error) {
	return make([]*Entry, 0), nil
}
