package tabi

import (
	"io"
	"strings"
)

type Dumper interface {
	Dump(io.Writer, any)
}

func GetDumper(s string) Dumper {
	switch strings.ToLower(s) {
	case "json":
		return JsonDumpFormat{}
	default:
		return TabDumpSliceStruct{}
	}
}
