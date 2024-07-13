package clang

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestRecordLayout_UnmarshalString(t *testing.T) {
	after := "[sizeof=32, align=8]"
	Size := 0
	Align := 0
	mylog.Check2(fmt.Sscanf(after, "[sizeof=%d, align=%d]", &Size, &Align))
}

/*
handle https://github.com/HyperDbg/gui/issues/110
handle https://github.com/HyperDbg/gui/issues/111
handle https://github.com/HyperDbg/gui/issues/113
handle https://github.com/HyperDbg/gui/issues/167
handle https://github.com/HyperDbg/gui/issues/173
*/

/*
	0 | struct _CR3_TYPE::(unnamed at ./merged_headers.h:196:9)

0:0-11 |   UINT64 Pcid                 //size 11-0=11 start 0
1:4-39 |   UINT64 PageFrameNumber      //size 39-4=35 start 15
6:0-11 |   UINT64 Reserved1            //size 11-0=11 start 50
7:4-6  |   UINT64 Reserved_2           //size 6-4=2   start 61
7:7-7  |   UINT64 PcidInvalidate       //size 7-7=0   start 63

	| [sizeof=8, align=8]           //size 11-0=11 start
*/
type BitField struct {
	StructName   string
	Name         string
	Type         string
	Start        int
	End          int
	Size         int //End-Start
	IndexForBits int //Size-End
	sizeof       int
	align        int
}

func parseBitFields(text string) ([]BitField, error) {
	lines := strings.Split(text, "\n")
	var bitFields []BitField
	var err error
	parentEnd := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		leftPart := strings.TrimSpace(parts[0])
		rightPart := strings.TrimSpace(parts[1])

		field := BitField{}
		if strings.Contains(rightPart, "struct") && strings.Contains(rightPart, "unnamed at") {
			rightPart = strings.TrimPrefix(rightPart, "struct ")
			before, _, found := strings.Cut(rightPart, "::")
			if found {
				field.StructName = before
				continue
			}
		}
		if strings.HasPrefix(rightPart, "[sizeof=") {
			_, err = fmt.Sscanf(rightPart, "[sizeof=%d, align=%d", &field.sizeof, &field.align)
			if err != nil {
				return nil, err
			}
			continue
		}
		TypeAndName := strings.Split(rightPart, " ")
		startEndParts := strings.Split(leftPart, "-")
		if len(startEndParts) == 1 {
			// 没有位域信息，直接把 | 左边的作为 Start
			fmt.Sscanf(leftPart, "%d", &field.Start)
			field.Name = TypeAndName[1]
			field.Type = TypeAndName[0]
			bitFields = append(bitFields, field)
		} else if len(startEndParts) == 2 {
			// 有位域信息
			split := strings.Split(strings.Split(leftPart, ":")[1], "-")
			fmt.Sscanf(split[0], "%d", &field.Start)
			fmt.Sscanf(split[1], "%d", &field.End)
			if bitFields != nil {
				parentEnd = bitFields[len(bitFields)-1].End
			}
			field.Size = field.End - field.Start + 1
			field.Name = TypeAndName[1]
			field.Type = TypeAndName[0]
			if bitFields != nil {
				field.IndexForBits = bitFields[len(bitFields)-1].IndexForBits + field.Start + parentEnd
			} else {
				field.IndexForBits = field.Start + parentEnd
			}
			bitFields = append(bitFields, field)
		}
	}

	return bitFields, nil
}

func TestParseBitFieldsWithBitRange(t *testing.T) {
	text := `  
         0 | struct _CR3_TYPE::(unnamed at ./merged_headers.h:196:9)
   0:0-11 |   UINT64 Pcid                  
   1:4-39 |   UINT64 PageFrameNumber       
   6:0-11 |   UINT64 Reserved1             
    7:4-6 |   UINT64 Reserved_2            
    7:7-7 |   UINT64 PcidInvalidate        
          | [sizeof=8, align=8]            
`
	bitFields, err := parseBitFields(text)
	if err != nil {
		t.Errorf("Error parsing bit fields: %v", err)
	}
	for _, bf := range bitFields {
		fmt.Printf("StructName: %-10s Name: %-20s Type: %-10s Start: %-5d End: %-5d Size: %-5d IndexForBits: %-5d sizeof: %-5d align: %-5d\n",
			bf.StructName, bf.Name, bf.Type, bf.Start, bf.End, bf.Size, bf.IndexForBits, bf.sizeof, bf.align)
	}
}

func TestParseBitFieldsWithoutBitRange(t *testing.T) {
	text := `
         0 | struct _CR3_TYPE::(unnamed at ./merged_headers.h:196:9)
   11 |   UINT64 Pcid                 //
   39 |   UINT64 PageFrameNumber      //
   11 |   UINT64 Reserved1            //
   6 |   UINT64 Reserved_2           //
   7 |   UINT64 PcidInvalidate       //
          | [sizeof=8, align=8]           //
`

	bitFields, err := parseBitFields(text)
	if err != nil {
		t.Errorf("Error parsing bit fields: %v", err)
	}

	totalBits := 0
	for _, bf := range bitFields {
		// Calculate the byte index and bit index within the byte
		byteIndex := totalBits / 8
		bitIndex := totalBits % 8
		fmt.Printf("Name: %s, Start: %d, End: %d, Byte Index: %d, Bit Index: %d\n", bf.Name, totalBits, totalBits+bf.End-bf.Start, byteIndex, bitIndex)
		totalBits += bf.End - bf.Start + 1
	}
}
