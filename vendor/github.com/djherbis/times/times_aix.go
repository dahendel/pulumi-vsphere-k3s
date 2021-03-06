// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://golang.org/src/os/stat_aix.go

package times

import (
	"os"
	"syscall"
	"time"
)

// HasChangeTime and HasBirthTime are true if and only if
// the target OS supports them.
const (
	HasChangeTime = true
	HasBirthTime  = false
)

type timespec struct {
	atime
	mtime
	ctime
	nobtime
}

func timespecToTime(ts syscall.StTimespec_t) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func getTimespec(fi os.FileInfo) (t timespec) {
	stat := fi.Sys().(*syscall.Stat_t)
	t.atime.v = timespecToTime(stat.Atim)
	t.mtime.v = timespecToTime(stat.Mtim)
	t.ctime.v = timespecToTime(stat.Ctim)
	return t
}
