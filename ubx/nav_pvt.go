// nav_pvt.go
package ubx

import (
	"encoding/binary"
	"fmt"
)

// NavPVT struct for the standard 84-byte version of NAV-PVT
// (ZED-F9 might have 92 bytes if it has extra fields).
type NavPVT struct {
	// Put all relevant fields
	ITOW       uint32
	Year       uint16
	Month      uint8
	Day        uint8
	Hour       uint8
	Min        uint8
	Sec        uint8
	Valid      uint8
	TAcc       uint32
	Nano       int32
	FixType    uint8
	Flags      uint8
	Flags2     uint8
	NumSV      uint8
	Lon        int32
	Lat        int32
	Height     int32
	HMSL       int32
	HAcc       uint32
	VAcc       uint32
	VelN       int32
	VelE       int32
	VelD       int32
	GSpeed     int32
	Heading    int32
	SAcc       uint32
	HeadingAcc uint32
	PDOP       uint16
	Flags3     uint8
	Reserved1  uint8
	Reserved2  uint32

	// If length >= 92, parse additional fields: headVeh, magDec, etc.
}

func (p *NavPVT) ClassID() uint8 { return ClassNAV }
func (p *NavPVT) MsgID() uint8   { return MsgIDNAVPVT }
func (p *NavPVT) String() string {
	return fmt.Sprintf("NavPVT fixType=%d, #SV=%d, lat=%d, lon=%d", p.FixType, p.NumSV, p.Lat, p.Lon)
}

// ParseNavPVT parses the UBX NAV-PVT payload into a NavPVT struct
func ParseNavPVT(payload []byte) (*NavPVT, error) {
	if len(payload) < 84 {
		return nil, fmt.Errorf("NAV-PVT payload too short: %d bytes", len(payload))
	}
	nav := &NavPVT{}
	nav.ITOW = binary.LittleEndian.Uint32(payload[0:4])
	nav.Year = binary.LittleEndian.Uint16(payload[4:6])
	nav.Month = payload[6]
	nav.Day = payload[7]
	nav.Hour = payload[8]
	nav.Min = payload[9]
	nav.Sec = payload[10]
	nav.Valid = payload[11]
	nav.TAcc = binary.LittleEndian.Uint32(payload[12:16])
	nav.Nano = int32(binary.LittleEndian.Uint32(payload[16:20]))
	nav.FixType = payload[20]
	nav.Flags = payload[21]
	nav.Flags2 = payload[22]
	nav.NumSV = payload[23]
	nav.Lon = int32(binary.LittleEndian.Uint32(payload[24:28]))
	nav.Lat = int32(binary.LittleEndian.Uint32(payload[28:32]))
	nav.Height = int32(binary.LittleEndian.Uint32(payload[32:36]))
	nav.HMSL = int32(binary.LittleEndian.Uint32(payload[36:40]))
	nav.HAcc = binary.LittleEndian.Uint32(payload[40:44])
	nav.VAcc = binary.LittleEndian.Uint32(payload[44:48])
	nav.VelN = int32(binary.LittleEndian.Uint32(payload[48:52]))
	nav.VelE = int32(binary.LittleEndian.Uint32(payload[52:56]))
	nav.VelD = int32(binary.LittleEndian.Uint32(payload[56:60]))
	nav.GSpeed = int32(binary.LittleEndian.Uint32(payload[60:64]))
	nav.Heading = int32(binary.LittleEndian.Uint32(payload[64:68]))
	nav.SAcc = binary.LittleEndian.Uint32(payload[68:72])
	nav.HeadingAcc = binary.LittleEndian.Uint32(payload[72:76])
	nav.PDOP = binary.LittleEndian.Uint16(payload[76:78])
	nav.Flags3 = payload[78]
	nav.Reserved1 = payload[79]
	nav.Reserved2 = binary.LittleEndian.Uint32(payload[80:84])

	// If len(payload) >= 92, parse additional fields...
	return nav, nil
}
