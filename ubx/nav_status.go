// nav_status.go
package ubx

import (
	"encoding/binary"
	"fmt"
)

// NavStatus for UBX-NAV-STATUS (class=0x01, id=0x03)
type NavStatus struct {
	ITOW    uint32 // GPS time of week (ms)
	FixType uint8
	Flags   uint8
	FixStat uint8
	Flags2  uint8
	// More fields if needed, total ~16 bytes
}

// Ensure NavStatus implements UBXMessage (optional)
func (s *NavStatus) ClassID() uint8 { return ClassNAV }
func (s *NavStatus) MsgID() uint8   { return MsgIDNAVSTATUS }
func (s *NavStatus) String() string {
	return fmt.Sprintf("NavStatus fixType=%d, flags=0x%02X", s.FixType, s.Flags)
}

// ParseNavStatus parses a 16-byte payload
func ParseNavStatus(payload []byte) (*NavStatus, error) {
	if len(payload) < 16 {
		return nil, fmt.Errorf("NAV-STATUS payload too short: %d bytes", len(payload))
	}
	st := &NavStatus{}
	st.ITOW = binary.LittleEndian.Uint32(payload[0:4])
	st.FixType = payload[4]
	st.Flags = payload[5]
	st.FixStat = payload[6]
	st.Flags2 = payload[7]
	// Reserved8..12, TTFF, MSSS if needed for the rest
	return st, nil
}
