package tabi

import (
	"bytes"
	"strings"
	"testing"
)

func TestDumpTab(t *testing.T) {

	s := []struct {
		Name  string
		Slice []int
		Map   map[string]int
		Map2  map[string]int
		Map3  map[string]int
		Int   int `col:"MY_COL" json:"int"`
		Bool  bool
		Any   any
	}{
		{
			Name:  "test",
			Slice: []int{1, 2},
			Map: map[string]int{
				"one": 1,
				"two": 2,
			},
			Map2: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
			Map3: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
				"four":  4,
			},
			Int:  42,
			Bool: false,
			Any: struct {
				Test string
			}{
				Test: "qwerty",
			},
		},
		{
			Name:  "test2",
			Slice: []int{1, 2, 5, 10, 20},
			Map: map[string]int{
				"one": 1,
			},
			Int:  0,
			Bool: true,
		},
	}

	t.Run("tab", func(t *testing.T) {
		exp := strings.TrimSpace(`
NAME   SLICE  MAP    MAP2     MAP3     MY_COL  BOOL   ANY
test   1      one=1  one=1    four=4   42      false  {qwerty}
       2      two=2  three=3  one=1                   
                     two=2    three=3                 
                              two=2                   
test2  1      one=1  {}       {}       0       true   <nil>
       2                                              
       5                                              
       10                                             
       20
`)
		d := TabDumpSliceStruct{}
		buf := bytes.NewBuffer(nil)
		d.Dump(buf, s)

		act := strings.TrimSpace(buf.String())

		if exp != act {
			t.Errorf("\nexp:\n%s\ngot:\n%s", exp, act)
		}
	})

	t.Run("json", func(t *testing.T) {
		exp := `[{"Name":"test","Slice":[1,2],"Map":{"one":1,"two":2},"Map2":{"one":1,"three":3,"two":2},"Map3":{"four":4,"one":1,"three":3,"two":2},"int":42,"Bool":false,"Any":{"Test":"qwerty"}},{"Name":"test2","Slice":[1,2,5,10,20],"Map":{"one":1},"Map2":null,"Map3":null,"int":0,"Bool":true,"Any":null}]`
		d := JsonDumpFormat{}
		buf := bytes.NewBuffer(nil)
		d.Dump(buf, s)

		if strings.TrimSpace(exp) != strings.TrimSpace(buf.String()) {
			t.Errorf("\nexp:\n%s\ngot:\n%s", exp, buf.String())
		}
	})

}
