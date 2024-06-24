package utils

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

const ident = "  "

type dumpTask struct {
	w     io.Writer
	name  string
	depth int
	done  map[uintptr]bool
}

func Dump(name string, v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	task := &dumpTask{w: buf, name: name, done: make(map[uintptr]bool)}

	if f, ok := v.(reflect.Value); ok {
		mylog.Check(task.dump(f))
	} else {
		mylog.Check(task.dump(reflect.ValueOf(v)))
	}

	return buf.Bytes(), nil
}

func (s *dumpTask) write(v any) error {
	pre := strings.Repeat(ident, s.depth)
	_ := mylog.Check2(fmt.Fprintf(s.w, "%s%s: %v\n", pre, s.name, v))
	return err
}

func (s *dumpTask) dump(v reflect.Value) (err error) {
	if !v.IsValid() {
		return s.write(nil)
	}
	switch v.Kind() {
	case reflect.Bool:
		mylog.Check(s.write(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		mylog.Check(s.write(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		mylog.Check(s.write(v.Uint()))
	case reflect.Float32, reflect.Float64:
		mylog.Check(s.write(v.Float()))
	case reflect.Complex64, reflect.Complex128:
		mylog.Check(s.write(v.Complex()))
	case reflect.String:
		mylog.Check(s.write("\"" + html.EscapeString(v.String()) + "\""))
	case reflect.Slice:
		if v.IsNil() {
			mylog.Check(s.write(nil))
			break
		}
		fallthrough
	case reflect.Array:
		mylog.Check(s.dumpSlice(v))
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			mylog.Check(s.write(nil))
		} else {
			mylog.Check(s.dump(v.Elem()))
		}
	case reflect.Map:
		mylog.Check(s.dumpMap(v))
	case reflect.Struct:
		mylog.Check(s.dumpStruct(v))
	default:
		mylog.Check(s.write(v.Interface()))
	}
	return
}

func (s *dumpTask) dumpSlice(v reflect.Value) (err error) {
	if s.done[v.UnsafeAddr()] {
		return s.write("<cycle via>")
	}
	s.done[v.UnsafeAddr()] = true
	mylog.Check(s.write(v.Type().String()))

	s.depth++
	l := v.Len()
	for i := 0; i < l; i++ {
		s.name = strconv.Itoa(i)
		mylog.Check(s.dump(v.Index(i)))

	}
	s.depth--
	s.done[v.UnsafeAddr()] = false
	return
}

func (s *dumpTask) dumpMap(v reflect.Value) (err error) {
	if v.IsNil() {
		return s.write(nil)
	}
	if s.done[v.UnsafeAddr()] {
		return s.write("<cycle via>")
	}
	s.done[v.UnsafeAddr()] = true
	mylog.Check(s.write(v.Type().String()))

	s.depth++
	iter := v.MapRange()
	for iter.Next() {
		s.name = fmt.Sprint(iter.Key())
		mylog.Check(s.dump(iter.Value()))

	}
	s.depth--
	s.done[v.UnsafeAddr()] = false
	return
}

func (s *dumpTask) dumpStruct(v reflect.Value) (err error) {
	if s.done[v.UnsafeAddr()] {
		return s.write("<cycle via>")
	}
	s.done[v.UnsafeAddr()] = true
	mylog.Check(s.write(v.Type().String()))

	s.depth++
	t := v.Type()
	l := v.NumField()
	for i := 0; i < l; i++ {
		sv := v.Field(i)
		st := t.Field(i)
		if !st.IsExported() {
			continue
		}

		s.name = st.Name
		if st.Anonymous {
			s.name = st.Type.Name()
		}
		mylog.Check(s.dump(sv))

	}
	s.depth--
	s.done[v.UnsafeAddr()] = false
	return
}
