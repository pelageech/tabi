package tabi

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"text/tabwriter"
)

type TabDumpSliceStruct struct{}

func (TabDumpSliceStruct) Dump(w io.Writer, v any) {
	t := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	defer t.Flush()

	vof := reflect.ValueOf(v)
	if vof.Kind() == reflect.Ptr {
		vof = vof.Elem()
	}
	if vof.Kind() != reflect.Slice {
		fmt.Fprintf(w, "expected slice, got %T\n", v)
		return
	}

	if vof.Len() == 0 {
		fmt.Fprintf(w, "no data\n")
		return
	}

	internal := vof.Index(0)
	if internal.Kind() == reflect.Ptr {
		internal = internal.Elem()
	}
	if internal.Kind() != reflect.Struct {
		fmt.Fprintf(w, "expected slice of struct, got %T\n", v)
		return
	}

	for i := range internal.NumField() {
		colName := internal.Type().Field(i).Tag.Get("col")
		if colName == "" {
			colName = strings.ToUpper(internal.Type().Field(i).Name)
		}
		if i == internal.NumField()-1 {
			fmt.Fprintf(t, "%s\n", colName)
		} else {
			fmt.Fprintf(t, "%s\t", colName)
		}
	}

	for i := 0; i < vof.Len(); i++ {
		var (
			slicesVals []sliceKey
			mapKeys    []mapKey
		)
		internal := vof.Index(i)
		if internal.Kind() == reflect.Ptr {
			internal = internal.Elem()
		}
		if internal.Kind() != reflect.Struct {
			fmt.Fprintf(w, "expected slice of struct, got %T\n", v)
			return
		}

		for j := range internal.NumField() {
			field := internal.Field(j)
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
				if field.Len() == 0 {
					fmt.Fprint(t, "[]\t")
				} else if field.Len() == 1 {
					fmt.Fprintf(t, "%v\t", field.Index(0).Interface())
				} else {
					fmt.Fprintf(t, "%v\t", field.Index(0).Interface())
					slicesVals = append(slicesVals, sliceKey{j, field.Len() - 1, field.Slice(1, field.Len())})
				}
			} else if field.Kind() == reflect.Map {
				if field.Len() == 0 {
					fmt.Fprint(t, "{}\t")
					continue
				}
				keys := field.MapKeys()
				if len(keys) == 1 {
					fmt.Fprintf(t, "%v=%v\t", keys[0].Interface(), field.MapIndex(keys[0]).Interface())
				} else {
					slices.SortFunc(keys, mustCompareOrdered)
					fmt.Fprintf(t, "%v=%v\t", keys[0].Interface(), field.MapIndex(keys[0]).Interface())
					mapKeys = append(mapKeys, mapKey{j, keys[1:], field})
				}
			} else {
				if j == internal.NumField()-1 {
					_, _ = io.WriteString(t, fmt.Sprintf("%v\n", field.Interface()))
				} else {
					_, _ = io.WriteString(t, fmt.Sprintf("%v\t", field.Interface()))
				}
			}
		}
		dumpIterable(t, internal.NumField(), slicesVals, mapKeys)
	}
}

type sliceKey struct {
	pos    int
	length int
	slices reflect.Value
}

type mapKey struct {
	pos  int
	keys []reflect.Value
	mapv reflect.Value
}

func dumpIterable(w io.Writer, cols int, slices []sliceKey, maps []mapKey) {
	maxLen := 0

	for _, slice := range slices {
		if slice.length > maxLen {
			maxLen = slice.length
		}
	}
	for _, key := range maps {
		if len(key.keys) > maxLen {
			maxLen = len(key.keys)
		}
	}

	if maxLen == 0 {
		return
	}

	for i := range maxLen {
		slicesIdx, mapsIdx := 0, 0
		for j := range cols {
			if slicesIdx < len(slices) {
				if s := slices[slicesIdx]; s.pos == j {
					slicesIdx++

					if i >= s.length {
						fmt.Fprintf(w, "\t")
						continue
					}

					fmt.Fprintf(w, "%v", s.slices.Index(i).Interface())
				}
			}
			if mapsIdx < len(maps) {
				if m := maps[mapsIdx]; m.pos == j {
					mapsIdx++
					if i >= len(m.keys) {
						fmt.Fprintf(w, "\t")
						continue
					}
					fmt.Fprintf(w, "%v=%v", m.keys[i].Interface(), m.mapv.MapIndex(m.keys[i]).Interface())
				}
			}
			if j != cols-1 {
				fmt.Fprintf(w, "\t")
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func compareOrdered(a, b reflect.Value) (int, error) {
	if a.Type() != b.Type() {
		return 0, fmt.Errorf("разные типы: %v и %v", a.Type(), b.Type())
	}

	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if a.Int() < b.Int() {
			return -1, nil
		} else if a.Int() > b.Int() {
			return 1, nil
		}
		return 0, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if a.Uint() < b.Uint() {
			return -1, nil
		} else if a.Uint() > b.Uint() {
			return 1, nil
		}
		return 0, nil

	case reflect.Float32, reflect.Float64:
		if a.Float() < b.Float() {
			return -1, nil
		} else if a.Float() > b.Float() {
			return 1, nil
		}
		return 0, nil

	case reflect.String:
		if a.String() < b.String() {
			return -1, nil
		} else if a.String() > b.String() {
			return 1, nil
		}
		return 0, nil

	default:
		return 0, fmt.Errorf("type %v is not comparable", a.Type())
	}
}

func mustCompareOrdered(a, b reflect.Value) int {
	r, err := compareOrdered(a, b)
	if err != nil {
		panic(err)
	}
	return r
}
