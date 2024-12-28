package main

import (
	"context"
	"encoding/binary"
	"log"
	"strings"
	"time"

	nmea "github.com/adrianmo/go-nmea"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	"go.bug.st/serial"

	// Adjust to your actual module name or path
	"github.com/miketigerblue/gnss-logger/ubx"
)

// Adjust these as needed
const (
	serialPort   = "/dev/ttyACM0"
	baudRate     = 115200
	influxURL    = "http://influxdb:8086"
	influxToken  = "rhYPk3z6ewCsDwWy59iYfvvbCZTwOgNbHP2x2eZl5keNrmxgVCgL5HhJ_J1PkFPSuEs0Akr7w0dCoEUUiLXuSg=="
	influxOrg    = "my-org"
	influxBucket = "gnss"
)

func main() {
	// 1) Open serial
	port, err := serial.Open(serialPort, &serial.Mode{BaudRate: baudRate})
	if err != nil {
		log.Fatalf("Could not open serial port %s: %v", serialPort, err)
	}
	defer port.Close()

	// 2) Influx client
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()
	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)

	// 3) Persistent buffer for data
	buffer := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)

	log.Printf("Starting GNSS logger on %s at %d baud...", serialPort, baudRate)

	for {
		// Read from serial
		n, err := port.Read(tmp)
		if err != nil {
			log.Printf("Read error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if n == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		// Append to buffer
		buffer = append(buffer, tmp[:n]...)

		// Parse from buffer
		parseBuffer(&buffer, writeAPI)
	}
}

// parseBuffer handles NMEA lines (split by newline) and UBX packets
func parseBuffer(buf *[]byte, writeAPI api.WriteAPIBlocking) {
	// 1) NMEA lines
	for {
		newlineIdx := indexOfByte(*buf, '\n')
		if newlineIdx == -1 {
			break
		}
		line := (*buf)[:newlineIdx]
		*buf = (*buf)[newlineIdx+1:]

		lineStr := strings.TrimSpace(string(line))
		if strings.HasPrefix(lineStr, "$") {
			parseNMEALine(lineStr, writeAPI)
		} else {
			// Possibly partial UBX or random data
			// We just skip it here because UBX parse is below
		}
	}

	// 2) UBX packets
	for {
		if len(*buf) < 6 {
			break
		}
		idx := findUBXSync(*buf)
		if idx == -1 {
			break
		}
		// Discard any leading junk before sync
		if idx > 0 {
			*buf = (*buf)[idx:]
		}
		if len(*buf) < 6 {
			break
		}
		length := binary.LittleEndian.Uint16((*buf)[4:6])
		totalNeeded := 6 + int(length) + 2 // header + payload + 2 checksum
		if len(*buf) < totalNeeded {
			break
		}
		// Extract full packet
		packet := (*buf)[:totalNeeded]
		*buf = (*buf)[totalNeeded:]
		parseUBXPacket(packet, writeAPI)
	}
}

// parseNMEALine handles standard NMEA like GGA, RMC, GSA, GSV, etc.
func parseNMEALine(line string, writeAPI api.WriteAPIBlocking) {
	msg, err := nmea.Parse(line)
	if err != nil {
		log.Printf("[NMEA parse error] %v => '%s'", err, line)
		return
	}

	switch typed := msg.(type) {
	case nmea.GGA:
		lat, lon, alt := typed.Latitude, typed.Longitude, typed.Altitude
		log.Printf("GGA => lat=%.6f, lon=%.6f, alt=%.2f", lat, lon, alt)
		p := influxdb2.NewPoint("nmea",
			map[string]string{"type": "GGA"},
			map[string]interface{}{
				"latitude":  lat,
				"longitude": lon,
				"altitude":  alt,
			},
			time.Now())
		writeIfErr(writeAPI.WritePoint(context.Background(), p))

	case nmea.RMC:
		lat, lon, spd := typed.Latitude, typed.Longitude, typed.Speed
		log.Printf("RMC => lat=%.6f, lon=%.6f, speed=%.2f", lat, lon, spd)
		p := influxdb2.NewPoint("nmea",
			map[string]string{"type": "RMC"},
			map[string]interface{}{
				"latitude":  lat,
				"longitude": lon,
				"speed":     spd,
			},
			time.Now())
		writeIfErr(writeAPI.WritePoint(context.Background(), p))

	case nmea.GSA, nmea.GSV, nmea.VTG, nmea.GLL:
		// other standard messages recognized by go-nmea
		log.Printf("Other NMEA => %s => %+v", typed.DataType(), typed)
		// You can parse out fields if you want, or just skip

	default:
		// Unhandled or custom message. Possibly a prefix like PUBX that isn't recognized
		log.Printf("Unhandled NMEA => %s => %v", msg.DataType(), line)
	}
}

// parseUBXPacket delegates to ubx.ParseUBX(...) and logs/writes recognized messages
func parseUBXPacket(packet []byte, writeAPI api.WriteAPIBlocking) {
	// UBX layout: [0..1]=sync, [2]=class, [3]=id, [4..5]=len, [6..(6+len-1)]=payload, [end-2..end-1]=ck
	if len(packet) < 8 {
		log.Printf("UBX packet too short (%d bytes)", len(packet))
		return
	}
	classID := packet[2]
	msgID := packet[3]
	length := binary.LittleEndian.Uint16(packet[4:6])
	payload := packet[6 : 6+length]

	parsed, err := ubx.ParseUBX(classID, msgID, payload)
	if err != nil {
		log.Printf("UBX parse error: %v", err)
		return
	}

	switch msg := parsed.(type) {
	case *ubx.NavPVT:
		lat := float64(msg.Lat) / 1e7
		lon := float64(msg.Lon) / 1e7
		log.Printf("NAV-PVT => fixType=%d, #SV=%d, lat=%.7f, lon=%.7f", msg.FixType, msg.NumSV, lat, lon)
		p := influxdb2.NewPoint("ubx",
			map[string]string{"message": "NAV-PVT"},
			map[string]interface{}{
				"latitude": lat, "longitude": lon,
				"fixType": msg.FixType, "numSV": msg.NumSV,
			},
			time.Now())
		writeIfErr(writeAPI.WritePoint(context.Background(), p))

	case *ubx.NavStatus:
		// Example: nav_status has fixType, flags, etc.
		log.Printf("NAV-STATUS => fixType=%d, flags=0x%02X", msg.FixType, msg.Flags)
		p := influxdb2.NewPoint("ubx",
			map[string]string{"message": "NAV-STATUS"},
			map[string]interface{}{
				"fixType": msg.FixType,
				"flags":   msg.Flags,
			},
			time.Now())
		writeIfErr(writeAPI.WritePoint(context.Background(), p))

	case *ubx.NavPosLLH:
		// Example: lat/lon in 1e-7, height in mm
		lat := float64(msg.Lat) / 1e7
		lon := float64(msg.Lon) / 1e7
		alt := float64(msg.Height) / 1000
		log.Printf("NAV-POSLLH => lat=%.7f, lon=%.7f, alt=%.2f m", lat, lon, alt)
		p := influxdb2.NewPoint("ubx",
			map[string]string{"message": "NAV-POSLLH"},
			map[string]interface{}{
				"latitude":  lat,
				"longitude": lon,
				"altitude":  alt,
			},
			time.Now())
		writeIfErr(writeAPI.WritePoint(context.Background(), p))

	default:
		// Fallback for unhandled UBX
		log.Printf("Unhandled UBX class=0x%02X, id=0x%02X => %#v", classID, msgID, parsed)
	}
}

// indexOfByte finds the first occurrence of 'b' in 'data', or -1 if not found
func indexOfByte(data []byte, b byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == b {
			return i
		}
	}
	return -1
}

// findUBXSync looks for 0xB5,0x62 in the data; returns index or -1
func findUBXSync(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0xB5 && data[i+1] == 0x62 {
			return i
		}
	}
	return -1
}

// writeIfErr logs any write errors to Influx
func writeIfErr(err error) {
	if err != nil {
		log.Printf("Influx write error: %v", err)
	}
}
