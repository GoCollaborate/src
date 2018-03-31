package ioHelper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/GoCollaborate/src/constants"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type status int

const (
	StatusNormal status = iota
	StatusError
)

type source struct {
	// the status of data source1
	Status status
	Reader io.Reader
}

type CSVOperator struct {
	s *source // input source
}

func (s *source) NewCSVOperator() *CSVOperator {
	return &CSVOperator{s}
}

func FromURL(url string) *source {
	resp, err := http.Get(url)
	if err != nil {
		return &source{StatusError, nil}
	}
	defer resp.Body.Close()
	return FromIOReader(resp.Body)
}

func FromPath(path string) *source {
	f, err := os.Open(path)
	if err != nil {
		return &source{StatusError, nil}
	}
	return FromFile(f)
}

func FromFile(f *os.File) *source {
	return &source{StatusNormal, f}
}

func FromString(s string) *source {
	return FromBytes([]byte(s))
}

func FromBytes(bs []byte) *source {
	return FromIOReader(bytes.NewReader(bs))
}

func FromIOReader(reader io.Reader) *source {
	return &source{StatusNormal, reader}
}

// fill up the given struct from pre-defined data source
func (op *CSVOperator) Fill(v interface{}) error {
	switch op.s.Status {
	case StatusNormal:
		return Decode(op.s.Reader, v)
	}
	return constants.ErrInputStreamCorrupted
}

// The code below is inspired and adapted from radioinmyhead/csv (https://github.com/radioinmyhead/csv/blob/master/csv.go)
// Updates:
// - Float64 type supported
// - Float32 type supported
// - Bool type supported

// decode
func Decode(in io.Reader, v interface{}) (err error) {
	m, err := csv2map(in)
	if err != nil {
		return
	}
	return map2list(m, v)
}

// io => []map[string]string
func csv2map(in io.Reader) ([]map[string]string, error) {
	var (
		r       = csv.NewReader(in)
		records [][]string
		ret     []map[string]string
		err     error
	)

	records, err = r.ReadAll()
	if err != nil {
		return ret, err
	}
	name := records[0]
	for i := 1; i < len(records); i++ {
		m := make(map[string]string, len(records))
		for j, k := range name {
			value := records[i][j]
			if value == "" {
				continue
			}
			if strings.HasPrefix(value, "#") {
				continue
			}
			if k == "" {
				continue
			}
			m[k] = value
		}
		if len(m) == 0 {
			continue
		}
		ret = append(ret, m)
	}
	return ret, nil
}

// []map[string]string => []interface{}
func map2list(m []map[string]string, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return constants.ErrIODecodePointerRequired
	}

	for {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			continue
		}
		break
	}
	if rv.Kind() != reflect.Slice {
		return constants.ErrIODecodeSliceRequired
	}
	ret := reflect.MakeSlice(rv.Type(), 0, 0)
	for _, i := range m {
		d := reflect.New(rv.Type().Elem())
		if err := map2struct(i, d); err != nil {
			return err
		}
		ret = reflect.Append(ret, d.Elem())
	}
	rv.Set(ret)
	return
}

// map[string]string => struct
func map2struct(m map[string]string, v reflect.Value) (err error) {
	v = v.Elem()
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.New(v.Type().Elem()))
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return constants.ErrIODecodeStructRequired
	}
	return value(m, v)
}

func getMapValue(m map[string]string, key string) (ret string) {
	k := strings.ToLower(key)
	ret, _ = m[k]
	return
}

func value(m map[string]string, v reflect.Value) error {
	var (
		err error
	)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		t := v.Type().Field(i).Tag.Get("csv")
		k := v.Type().Field(i).Name
		if !f.CanSet() {
			continue
		}
		d := getMapValue(m, k)
		switch f.Kind() {
		default:
			err = constants.ErrInputStreamNotSupported
		case reflect.Struct:
			if t == "extends" {
				err = value(m, f)
			} else {
				err = valueStruct(d, f)
			}
		case reflect.String:
			err = valueString(d, f)
		case reflect.Int:
			err = valueInt(d, f)
		case reflect.Int32:
			err = valueIntX(d, f, 32)
		case reflect.Int64:
			err = valueIntX(d, f, 64)
		case reflect.Bool:
			err = valueBool(d, f)
		case reflect.Float32:
			err = valueFloatX(d, f, 32)
		case reflect.Float64:
			err = valueFloatX(d, f, 64)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func valueStruct(s string, v reflect.Value) error {
	var (
		b   = []byte(s)
		err error
	)
	err = json.Unmarshal(b, v.Addr().Interface())
	return err
}

func valueString(s string, v reflect.Value) error {
	v.Set(reflect.ValueOf(s))
	return nil
}

func valueInt(s string, v reflect.Value) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(n))
	return nil
}

func valueIntX(s string, v reflect.Value, x int) error {
	n, err := strconv.ParseInt(s, 10, x)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(n))
	return nil
}

func valueFloatX(s string, v reflect.Value, x int) error {
	n, err := strconv.ParseFloat(s, x)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(n))
	return nil
}

func valueBool(s string, v reflect.Value) error {
	n, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(n))
	return nil
}
