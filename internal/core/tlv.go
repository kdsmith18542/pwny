package core

import (
	"encoding/binary"
	"fmt"
)

// Base TLV Types
const (
	TLVTypeNone       = 0
	TLVTypeString     = 10001
	TLVTypeUint       = 10002
	TLVTypeRaw        = 10003
	TLVTypeBool       = 10004
	TLVTypeQWord      = 10005
	TLVTypeCompressed = 10006
	TLVTypeGroup      = 10007
)

// Specific TLV types
const (
	TLV_TYPE_ANY          = 0
	TLV_TYPE_METHOD       = (TLVTypeString << 16) | 1
	TLV_TYPE_REQUEST_ID   = (TLVTypeString << 16) | 2
	TLV_TYPE_EXCEPTION    = (TLVTypeUint << 16) | 3
	TLV_TYPE_RESULT       = (TLVTypeUint << 16) | 4
	TLV_TYPE_STRING       = (TLVTypeString << 16) | 5
	TLV_TYPE_UINT         = (TLVTypeUint << 16) | 6
	TLV_TYPE_BOOL         = (TLVTypeBool << 16) | 7
	TLV_TYPE_RAW          = (TLVTypeRaw << 16) | 8
	TLV_TYPE_CHANNEL_ID   = (TLVTypeUint << 16) | 9
	TLV_TYPE_CHANNEL_DATA = (TLVTypeRaw << 16) | 10

	// Stdapi types
	TLV_TYPE_PROCESS_NAME      = (TLVTypeString << 16) | 1010
	TLV_TYPE_PROCESS_PATH      = (TLVTypeString << 16) | 1011
	TLV_TYPE_PROCESS_ID        = (TLVTypeUint << 16) | 1012
	TLV_TYPE_PROCESS_PARENT_ID = (TLVTypeUint << 16) | 1013
	TLV_TYPE_PROCESS_USER      = (TLVTypeString << 16) | 1014
	TLV_TYPE_PROCESS_ARCH      = (TLVTypeUint << 16) | 1015

	TLV_TYPE_NETWORK_INTERFACE = (TLVTypeGroup << 16) | 1020
	TLV_TYPE_INTERFACE_NAME    = (TLVTypeString << 16) | 1021
	TLV_TYPE_INTERFACE_ADDR    = (TLVTypeRaw << 16) | 1022
)

type TLV struct {
	Type  uint32
	Value interface{}
}

func (t TLV) Serialize() []byte {
	var val []byte
	switch v := t.Value.(type) {
	case string:
		val = append([]byte(v), 0x00) // Null-terminated
	case uint32:
		val = make([]byte, 4)
		binary.BigEndian.PutUint32(val, v)
	case []byte:
		val = v
	case bool:
		val = []byte{0x00}
		if v {
			val = []byte{0x01}
		}
	}

	res := make([]byte, 8)
	binary.BigEndian.PutUint32(res[0:4], t.Type)
	binary.BigEndian.PutUint32(res[4:8], uint32(len(val)))
	// Note: Meterpreter packet structure is [Type][Length][Value]
	// where Length is the length of the value only, or [Type][Length+8][Value]?
	// Actually, for TLVs it's [Type][Length+8][Value] usually.
	binary.BigEndian.PutUint32(res[4:8], uint32(len(val)+8))
	return append(res, val...)
}

func UnserializeTLV(data []byte) ([]TLV, error) {
	var tlvs []TLV
	offset := 0
	for offset < len(data) {
		if len(data)-offset < 8 {
			break
		}
		t := binary.BigEndian.Uint32(data[offset : offset+4])
		l := binary.BigEndian.Uint32(data[offset+4 : offset+8])
		if offset+int(l) > len(data) {
			return nil, fmt.Errorf("TLV length exceeds data size")
		}

		valData := data[offset+8 : offset+int(l)]
		var val interface{}

		// Determine type from the upper bits
		typeBase := t >> 16
		switch typeBase {
		case uint32(TLVTypeString):
			val = string(valData)
		case uint32(TLVTypeUint):
			if len(valData) >= 4 {
				val = binary.BigEndian.Uint32(valData)
			}
		case uint32(TLVTypeRaw):
			val = valData
		case uint32(TLVTypeBool):
			if len(valData) > 0 {
				val = valData[0] != 0
			}
		}

		tlvs = append(tlvs, TLV{Type: t, Value: val})
		offset += int(l)
	}
	return tlvs, nil
}
