//go:build linux && cgo && libsubid
// +build linux,cgo,libsubid

package fakeroot

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/apptainer/apptainer/internal/pkg/util/user"
)

/*
#cgo LDFLAGS: -l subid

#include <shadow/subid.h>
#include <stdlib.h>
#include <stdio.h>

struct subid_range apptainer_get_range(struct subid_range *ranges, int i)
{
	return ranges[i];
}

#if !defined(SUBID_ABI_MAJOR) || (SUBID_ABI_MAJOR < 4)
# define subid_get_uid_ranges get_subuid_ranges
# define subid_get_gid_ranges get_subgid_ranges
#endif
*/
import "C"

func readSubid(user *user.User, isUser bool) ([]*Entry, error) {
	ret := make([]*Entry, 0)
	uidstr := fmt.Sprintf("%d", user.UID)

	if user.Name == "ALL" {
		return nil, errors.New("username ALL not supported")
	}

	cUsername := C.CString(user.Name)
	defer C.free(unsafe.Pointer(cUsername))

	cuidstr := C.CString(uidstr)
	defer C.free(unsafe.Pointer(cuidstr))

	var nRanges C.int
	var cRanges *C.struct_subid_range
	if isUser {
		nRanges = C.subid_get_uid_ranges(cUsername, &cRanges)
		if nRanges <= 0 {
			nRanges = C.subid_get_uid_ranges(cuidstr, &cRanges)
		}
	} else {
		nRanges = C.subid_get_gid_ranges(cUsername, &cRanges)
		if nRanges <= 0 {
			nRanges = C.subid_get_gid_ranges(cuidstr, &cRanges)
		}
	}
	if nRanges < 0 {
		return nil, errors.New("cannot read subids")
	}
	defer C.free(unsafe.Pointer(cRanges))

	for i := 0; i < int(nRanges); i++ {
		r := C.apptainer_get_range(cRanges, C.int(i))
		line := fmt.Sprintf("%d:%d:%d", user.UID, r.start, r.count)
		ret = append(
			ret,
			&Entry{
				UID:      user.UID,
				Start:    uint32(r.start),
				Count:    uint32(r.count),
				disabled: false,
				line:     line,
			})
	}
	return ret, nil
}

func readSubuid(user *user.User) ([]*Entry, error) {
	return readSubid(user, true)
}

func readSubgid(user *user.User) ([]*Entry, error) {
	return readSubid(user, false)
}
