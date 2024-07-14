package bindlib

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"reflect"
)

func Validate(ptrToStruct any, size, align uintptr, fields ...any) {
	mylog.Call(func() {
		rtype := reflect.TypeOf(ptrToStruct).Elem()

		// Validate size
		if size != rtype.Size() {
			mylog.Check(fmt.Sprintf("Mismatching sizof(%s) 0x%d, expected 0x%x", rtype.Name(), rtype.Size(), size))
		}

		// Validate alignment
		if align != uintptr(rtype.Align()) {
			mylog.Check(fmt.Sprintf("Mismatching alignof(%s) 0x%d, expected 0x%x", rtype.Name(), rtype.Align(), align))
		}

		// Validate fields
		for i := 0; i < len(fields); i += 2 {
			fieldName := fields[i].(string)
			fieldOffset := reflect.ValueOf(fields[i+1]).Int()
			field, ok := rtype.FieldByName(fieldName)
			if !ok {
				panic(fmt.Sprintf("Field %s not found", fieldName))
			}
			if field.Offset != uintptr(fieldOffset) {
				mylog.Check(fmt.Sprintf("Mismatching offsetof(%s::%s): 0x%x, expected 0x%x", rtype.Name(), fieldName, field.Offset, fieldOffset))
			}
		}
	})
}
