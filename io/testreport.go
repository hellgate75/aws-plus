package io

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func EncodeTextFormatSummary(in interface{}) ([]byte, error) {
	var err error
	var out = make([]byte, 0)

	if strings.Contains(fmt.Sprintf("%T", in), "[]") {
		// It's an array of elements
		list, err := splitInterfaceArray(in)
		if err != nil {
			var text = fmt.Sprintf("Error: %v", err)
			out = append(out, []byte(text)...)
		} else {
			var text = ""
			if len(list) == 0 {
				text += "No results ...\n"
			}
			for idx, elem := range list {
				header, values, err := decomposeElement(elem)
				//b, _ := EncodeValue(&header, EncodingYaml)
				//fmt.Println("Index:", idx, "Header:", string(b))
				//b, _ = EncodeValue(&values, EncodingYaml)
				//fmt.Println("Index:", idx, "Values:", string(b))
				if err != nil {
					text += fmt.Sprintf("Error in line %v: %v\n", idx, err)
				} else {
					if idx == 0 {
						text += calculateHeaderLines(header)
					}
					text += calculateValueLines(values)
				}
			}
			out = append(out, []byte(text)...)
		}
	} else {
		// It's a single object
		header, values, err := decomposeElement(in)
		//b, _ := EncodeValue(&header, EncodingYaml)
		//fmt.Println("Header:", string(b))
		//b, _ = EncodeValue(&values, EncodingYaml)
		//fmt.Println("Values:", string(b))
		if err != nil {
			var text = fmt.Sprintf("Error: %v\n", err)
			out = append(out, []byte(text)...)
		} else {
			var text = calculateHeaderLines(header)
			text += calculateValueLines(values)
			out = append(out, []byte(text)...)
		}
	}
	return out, err
}

func splitInterfaceArray(in interface{}) ([]interface{}, error) {
	var err error
	var out = make([]interface{}, 0)
	var kind = reflect.TypeOf(in).Kind()
	if kind == reflect.Slice ||  kind == reflect.Array {
		v := reflect.Indirect(reflect.ValueOf(&in)).Elem()
		var length = v.Len()
		for i := 0; i < length; i++ {
			out = append(out, v.Index(i).Interface())
		}
	} else {
		err = errors.New(fmt.Sprintf("Invalid slice or list type: %s", kind.String()))
	}
	return out, err
}

func calculateHeaderLine(header HeaderSet) (string, []HeaderSet, []string, bool) {
	var hs = make([]HeaderSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for _, hl := range header.Columns {
		var length = len(hl.Name)
		out += hl.Name
		hs = append(hs, HeaderSet{hl.Columns})
		if len(hl.Columns) > 0 {
			hasMore = true
		}
		spc = append(spc, strings.Repeat(" ", length))
	}
	return out, hs, spc, hasMore
}

func calculateSubHeaderLine(headers []HeaderSet, spaces []string) (string, []HeaderSet, []string, bool) {
	var hs = make([]HeaderSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for idx, hsi := range headers {
		var spcLen = 0
		if len(hsi.Columns) == 0 {
			out += spaces[idx]
			spcLen = len(spaces[idx])
			hs = append(hs, HeaderSet{hsi.Columns})
			spc = append(spc, strings.Repeat(" ", spcLen))

		} else {
			line, hsC, spC, more := calculateHeaderLine(hsi)
			out += line
			hs = append(hs, hsC...)
			spc = append(spc, spC...)
			if more {
				hasMore = true
			}
		}
	}
	return out, hs, spc, hasMore
}

func calculateHeaderLines(header HeaderSet) string {
	var linesArr = make([]string, 0)
	line, hs, spc, more := calculateHeaderLine(header)
	//fmt.Println(more)
	//fmt.Println(hs)
	//fmt.Println(spc)
	//fmt.Println(line)
	linesArr = append(linesArr, line)
	for more {
		line, hs, spc, more = calculateSubHeaderLine(hs, spc)
		linesArr = append(linesArr, line)
		//fmt.Println(more)
		//fmt.Println(hs)
		//fmt.Println(spc)
		//fmt.Println(line)
	}
	return strings.Join(linesArr, "\n") + "\n"
}

func calculateValueLine(value ValuesSet) (string, []ValuesSet, []string, bool) {
	var hs = make([]ValuesSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for _, hl := range value.Values {
		var length = len(hl.Value)
		out += hl.Value
		hs = append(hs, ValuesSet{hl.SubValues})
		if len(hl.SubValues) > 0 {
			hasMore = true
		}
		spc = append(spc, strings.Repeat(" ", length))
	}
	return out, hs, spc, hasMore
}

func calculateSubValueLine(values []ValuesSet, spaces []string) (string, []ValuesSet, []string, bool) {
	var hs = make([]ValuesSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for idx, hsi := range values {
		var spcLen = 0
		if len(hsi.Values) == 0 {
			out += spaces[idx]
			spcLen = len(spaces[idx])
			hs = append(hs, ValuesSet{hsi.Values})
			spc = append(spc, strings.Repeat(" ", spcLen))

		} else {
			line, hsC, spC, more := calculateValueLine(hsi)
			out += line
			hs = append(hs, hsC...)
			spc = append(spc, spC...)
			if more {
				hasMore = true
			}
		}
	}
	return out, hs, spc, hasMore
}


func calculateValueLines(values ValuesSet) string {
	var linesArr = make([]string, 0)
	line, hs, spc, more := calculateValueLine(values)
	linesArr = append(linesArr, line)
	for more {
		line, hs, spc, more = calculateSubValueLine(hs, spc)
		linesArr = append(linesArr, line)
	}
	return strings.Join(linesArr, "\n") + "\n"
}

var runeA = byte('a')
var diff = byte('a') - runeA

func formatColumn(s string) string {
	var out = ""
	for idx, c := range s {
		if idx == 0 {
			if byte(c) >= runeA {
				c = rune(byte(c) - diff)
			}
			out += fmt.Sprintf("%c", c)
		} else {
			if byte(c) < runeA {
				out += " "
			}
			out += fmt.Sprintf("%c", c)
		}
	}
	return out
}

func trimLen(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else if len(s) < n {
		return s + strings.Repeat(" ", n-len(s)-1)
	}
	return s
}


func decomposeElement(in interface{}) (HeaderSet, ValuesSet, error) {
	var hSet = HeaderSet{make([]*HeaderElem, 0)}
	var vSet = ValuesSet{make([]*ValueItem, 0)}
	var err error

	k := reflect.ValueOf(in).Type().Kind()
	if k == reflect.Ptr {
		in = reflect.Indirect(reflect.ValueOf(in)).Interface()
	}
	if reflect.ValueOf(in).Type().Kind() == reflect.Struct {
		var e = reflect.ValueOf(in)
		var typeOfT = e.Type()
		numFields := e.NumField()
		for i := 0; i < numFields; i++ {
			f := e.Field(i)
			name := typeOfT.Field(i).Name
			fName := formatColumn(name)
			var value interface{}
			if e.FieldByName(name).Type().Kind() != reflect.Interface || ! e.FieldByName(name).IsNil() {
				value = e.FieldByName(name).Interface()
			}
			fValue := ""
			if value != nil {
				fValue = fmt.Sprintf("%v", value)
			}
			vType := f.Type()
			var length = len(fValue)
			if length > 0 && len(fName) > len(fValue) {
				length = len(fName)
				fName = trimLen(fName, length + 2)
				fValue = trimLen(fValue, length + 2)
			} else if length > 0 {
				fName = trimLen(fName, length + 2)
				fValue = trimLen(fValue, length + 2)
			}
			length += 2
			headItem := HeaderElem{
				Value: name,
				Name: fName,
				Columns: make([]*HeaderElem, 0),
				Size: length,
			}
			valueItem := ValueItem{
				Value: fValue,
				SubValues: make([]*ValueItem, 0),
			}
			if vType.Kind() == reflect.Struct || vType.Kind() == reflect.Interface {
				var hs HeaderSet
				var vs ValuesSet
				hs, vs, err := decomposeElement(value)
				if err == nil {
					length = 0
					for _, c := range hs.Columns {
						length += c.Size
						headItem.Columns = append(headItem.Columns, c)
					}
					headItem.Size = length
					for _, c := range vs.Values {
						valueItem.SubValues = append(valueItem.SubValues, c)
					}
					valueItem.Value = strings.Repeat(" ", headItem.Size)
				}
			}
			hSet.Columns = append(hSet.Columns, &headItem)
			vSet.Values = append(vSet.Values, &valueItem)
		}
	} else {
		err = errors.New(fmt.Sprintf("Cannot decompose non strcuture element %+v", in))
	}
	return hSet, vSet, err
}


type HeaderElem struct {
	Name		string
	Value		string
	Size		int
	Columns		[]*HeaderElem
}

type HeaderSet struct {
	Columns		[]*HeaderElem
}

type ValueItem struct {
	Value		string
	SubValues	[]*ValueItem
}

type ValuesSet struct {
	Values		[]*ValueItem
}