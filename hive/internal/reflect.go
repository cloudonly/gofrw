// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var (
	baseNameRegex = regexp.MustCompile(`github\.com/luomu/gofrw/[\w\/]+/`)
)

func TrimName(name string) string {
	return string(baseNameRegex.ReplaceAll([]byte(name), []byte{}))
}

func PrettyType(x any) string {
	return TrimName(fmt.Sprintf("%T", x))
}

func FuncNameAndLocation(fn any) string {
	f := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	file, line := f.FileLine(f.Entry())
	name := TrimName(f.Name())
	name = strings.TrimSuffix(name, "-fm")
	if file != "<autogenerated>" {
		return fmt.Sprintf("%s (%s:%d)", name, usefulPathSegment(file), line)
	}
	return name
}

// Purely a heuristic.
var commonRoots = map[string]struct{}{
	"pkg": {},
	"cmd": {},
}

func usefulPathSegment(file string) string {
	p := filepath.Clean(file)
	segs := strings.Split(p, string(os.PathSeparator))
	for i := len(segs) - 1; i > 0; i-- {
		if _, ok := commonRoots[segs[i]]; ok {
			segs = segs[i:]
			break
		}
	}
	return filepath.Join(segs...)
}
