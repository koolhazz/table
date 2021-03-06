// Copyright 2017 modood. All rights reserved.
// license that can be found in the LICENSE file.

// Package table produces a string that represents slice of structs data in a text table
package table

import (
	"fmt"
	"reflect"
)

// Output formats slice of structs data and writes to standard output.
func Output(slice interface{}) {
	fmt.Println(Table(slice))
}

// Table formats slice of structs data and returns the resulting string.
func Table(slice interface{}) string {
	var coln []string   // name of columns
	var colw []int      // width of columns
	var rows [][]string // rows of content

	for i, u := range sliceconv(slice) {
		v := reflect.ValueOf(u)
		t := reflect.TypeOf(u)
		if v.Kind() != reflect.Struct {
			panic("Table: items of slice should be on struct value")
		}
		var row []string

		m := 0 // count of unexported field
		for n := 0; n < v.NumField(); n++ {
			if t.Field(n).PkgPath != "" {
				m++
				continue
			}
			cn := t.Field(n).Name
			cv := fmt.Sprintf("%+v", v.FieldByName(cn).Interface())

			if i == 0 {
				coln = append(coln, cn)
				colw = append(colw, len(cn))
			}
			if colw[n-m] < len(cv) {
				colw[n-m] = len(cv)
			}

			row = append(row, cv)
		}
		rows = append(rows, row)
	}
	table := table(coln, colw, rows)
	return table
}

func table(coln []string, colw []int, rows [][]string) (table string) {
	head := [][]rune{[]rune{'┌'}, []rune{'│'}, []rune{'├'}}
	bttm := []rune{'└'}
	for i, v := range colw {
		head[0] = append(head[0], []rune(repeat(v+2, '─')+"┬")...)
		head[1] = append(head[1], []rune(" "+coln[i]+repeat(v-len(coln[i]), ' ')+" │")...)
		head[2] = append(head[2], []rune(repeat(v+2, '─')+"┼")...)
		bttm = append(bttm, []rune(repeat(v+2, '─')+"┴")...)
	}
	head[0][len(head[0])-1] = '┐'
	head[2][len(head[2])-1] = '┤'
	bttm[len(bttm)-1] = '┘'

	var body [][]rune
	for _, r := range rows {
		row := []rune{'│'}
		for i, v := range colw {
			// handle non-ascii character
			lb := len(r[i])
			lr := len([]rune(r[i]))

			row = append(row, []rune(" "+r[i]+repeat(v-lb+(lb-lr)/2, ' ')+" │")...)
		}
		body = append(body, row)
	}

	for _, v := range head {
		table += string(v) + "\n"
	}
	for _, v := range body {
		table += string(v) + "\n"
	}
	table += string(bttm)
	return table
}

func sliceconv(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("sliceconv: param \"slice\" should be on slice value")
	}

	l := v.Len()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

func repeat(time int, char rune) string {
	var s = make([]rune, time)
	for i := range s {
		s[i] = char
	}
	return string(s)
}
