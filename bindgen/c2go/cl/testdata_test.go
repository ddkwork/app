package cl_test

import (
	"testing"

	"github.com/ddkwork/app/bindgen/c2go"
	"github.com/ddkwork/app/bindgen/c2go/cl"
)

// -----------------------------------------------------------------------------

func TestFromTestdata(t *testing.T) {
	t.Skip()
	cl.SetDebug(0)
	defer cl.SetDebug(cl.DbgFlagAll)

	c2go.Run("", "../testdata/...", c2go.FlagRunTest|c2go.FlagTestMain|c2go.FlagFailFast, nil)
}

// -----------------------------------------------------------------------------
