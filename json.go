package tabi

import (
	"encoding/json"
	"io"
)

type JsonDumpFormat struct{}

func (JsonDumpFormat) Dump(w io.Writer, v any) {
	_ = json.NewEncoder(w).Encode(v)
}
