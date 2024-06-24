package translate

import (
	"fmt"
	"testing"
	"utils"
)

func TestTypeParse(t *testing.T) {
	utils.SetLogLevel(utils.LL_Verbose)
	src := "int (const char *, ...)"
	ps := NewTypeParse(src)
	ps.Next()
	tp := ps.parseType()
	fmt.Println(tp)
}
