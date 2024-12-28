// parser.go
package ubx

import (
	"encoding/binary"
	"fmt"
)

// ParseUBX is a general entry point: you pass in the UBX class, ID, payload, etc.
// If it's recognized (e.g. NAV-PVT), we parse it further and return something useful.
func ParseUBX(classID, msgID uint8, payload []byte) (interface{}, error) {
	switch {
	case classID == 0x01 && msgID == 0x07:
		// NAV-PVT
		navpvt, err := ParseNavPVT(payload)
		if err != nil {
			return nil, err
		}
		return navpvt, nil

	// Add more cases for other UBX messages you want to parse
	// e.g. NAV-STATUS, NAV-SAT, etc.

	default:
		// If we don't have a parser, just log or return an error
		return nil, fmt.Errorf("unknown UBX message class=0x%02X, id=0x%02X", classID, msgID)
	}
}

// Example helper to read the sync bytes and length, for a fully robust approach:
func IsUBX(data []byte) bool {
	return len(data) >= 2 && data[0] == 0xB5 && data[1] == 0x62
}

// Extract length (2 bytes after class,ID) in little-endian
func UBXPayloadLength(data []byte) (uint16, error) {
	if len(data) < 6 {
		return 0, fmt.Errorf("UBX data too short for length extraction")
	}
	return binary.LittleEndian.Uint16(data[4:6]), nil
}
