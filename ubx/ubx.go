// ubx.go
package ubx

import (
	"encoding/binary"
	"fmt"
)

// UBX Sync Characters
const (
	UBXSync1 = 0xB5
	UBXSync2 = 0x62
)

// Some Common Classes
const (
	ClassNAV = 0x01 // Navigation
	ClassRXM = 0x02 // Receiver Manager
	// ... add more as needed
)

// Some Common NAV Message IDs
const (
	MsgIDNAVPVT    = 0x07 // NAV-PVT
	MsgIDNAVSTATUS = 0x03 // NAV-STATUS
	MsgIDNAVPOSLLH = 0x02 // NAV-POSLLH
	// etc.
)

// UBXMessage is an optional interface if you want typed messages
type UBXMessage interface {
	ClassID() uint8
	MsgID() uint8
	String() string
}

// Helper to read length
func PayloadLength(data []byte) (uint16, error) {
	if len(data) < 6 {
		return 0, fmt.Errorf("UBX data too short to extract length")
	}
	return binary.LittleEndian.Uint16(data[4:6]), nil
}
