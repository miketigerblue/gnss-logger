// nav_posllh.go
package ubx

import (
	"encoding/binary"
	"fmt"
)

// NavPosLLH (class=0x01, id=0x02) => 28 bytes
type NavPosLLH struct {
	ITOW   uint32 // ms
	Lon    int32  // 1e-7 deg
	Lat    int32  // 1e-7 deg
	Height int32  // mm
	HMSL   int32  // mm
	HAcc   uint32 // mm
	VAcc   uint32 // mm
}

func (p *NavPosLLH) ClassID() uint8 { return ClassNAV }
func (p *NavPosLLH) MsgID() uint8   { return MsgIDNAVPOSLLH }
func (p *NavPosLLH) String() string {
	return fmt.Sprintf("NavPosLLH lat=%d lon=%d h=%d", p.Lat, p.Lon, p.Height)
}

func ParseNavPosLLH(payload []byte) (*NavPosLLH, error) {
	if len(payload) < 28 {
		return nil, fmt.Errorf("NAV-POSLLH payload too short: %d bytes", len(payload))
	}
	llh := &NavPosLLH{}
	llh.ITOW = binary.LittleEndian.Uint32(payload[0:4])
	llh.Lon = int32(binary.LittleEndian.Uint32(payload[4:8]))
	llh.Lat = int32(binary.LittleEndian.Uint32(payload[8:12]))
	llh.Height = int32(binary.LittleEndian.Uint32(payload[12:16]))
	llh.HMSL = int32(binary.LittleEndian.Uint32(payload[16:20]))
	llh.HAcc = binary.LittleEndian.Uint32(payload[20:24])
	llh.VAcc = binary.LittleEndian.Uint32(payload[24:28])
	return llh, nil
}
