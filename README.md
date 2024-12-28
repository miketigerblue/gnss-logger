                .-------------------.
    (NAV-SAT)  |   [ Space GNSS ]  | 
                '-------------------'
                     \
                      \    (NMEA & UBX)
                       \   ~~~~~~~~~~~~~
                        \
    .----------------.   \
    |   GNSS Logger  |----O----> [ Influx ]
    '----------------'           [ Grafana ]

# GNSS Logger (NMEA + UBX) with InfluxDB & Grafana

A Raspberry Pi–based GNSS logging application that:
- **Reads** mixed NMEA (ASCII) and UBX (binary) data from a u-blox ZED-F9 module.
- **Parses** NMEA sentences (GGA, RMC, etc.) and UBX messages (NAV-PVT, NAV-POSLLH, etc.).
- **Stores** parsed fields (latitude, longitude, altitude, speed, etc.) in InfluxDB.
- **Visualizes** real-time or historical GNSS data in Grafana.

This project demonstrates how to handle interleaved ASCII + binary streams reliably (via a ring buffer), then build a small IoT pipeline using Docker Compose for local or edge deployments.

---

## Features

1. **Robust Parsing**  
   - **NMEA** via [go-nmea](https://github.com/adrianmo/go-nmea)  
   - **UBX** via custom [\`ubx\`](./ubx/) package
2. **InfluxDB Integration**  
   - Writes lat/lon, altitude, speed, fix type, etc. into InfluxDB for time-series storage.
3. **Grafana Dashboards**  
   - Visualize location, speed, altitude over time, or use the Geomap panel to see lat/lon.
4. **Docker Compose**  
   - Spins up **InfluxDB**, **Grafana**, and **gnss-logger** containers together.

---

## Hardware Requirements

- **Raspberry Pi** (or similar Linux SBC)  
- **ZED-F9** module (e.g. SparkFun ZED-F9P, ZED-F9R, or equivalent)  
- USB or UART cable connection to Pi  
- (Optional) Internet connectivity if you want to access Grafana from outside

---

## Directory Layout

    gnss-logger/ 
        ├─ .env                  # environment variables for Docker Compose 
        ├─ .gitignore 
        ├─ LICENSE               # e.g. MIT license 
        ├─ README.md 
        ├─ Dockerfile            # Multi-stage build for the Go logger 
        ├─ docker-compose.yml    # InfluxDB, Grafana, gnss-logger 
        ├─ go.mod 
        ├─ go.sum 
        ├─ main.go               # Entrypoint for reading & parsing data 
        └─ ubx/ 
            ├─ ubx.go            # Common UBX constants/defs 
            ├─ parser.go         # ParseUBX dispatcher 
            ├─ nav_pvt.go        # Example UBX parse for NAV-PVT 
            ├─ nav_status.go     # Example UBX parse for NAV-STATUS 
            └─ nav_posllh.go     # Example UBX parse for NAV-POSLLH


---

## Getting Started

**Clone this repo**:

    git clone https://github.com/miketigerblue/gnss-logger.git
    cd gnss-logger

### (Optional) Configure environment in .env:

    INFLUXDB_USER=admin
    INFLUXDB_PASS=adminpassword
    INFLUXDB_TOKEN=my-secret-token
    INFLUXDB_ORG=my-org
    INFLUXDB_BUCKET=gnss


### Build & run with Docker Compose:

    docker-compose build
    docker-compose up -d

![Docker Compose](/screenshots/docker-compose-build-up.png)

**This launches:**

    influxdb at http://rpi:8086
    grafana at http://rpi:3000
    gnss-logger container, reading from /dev/ttyACM0 or /dev/ttyUSB0 (adjust in docker-compose.yml).

**Check logs:**
    
    docker-compose logs -f gnss-logger
    Look for parsed NMEA lines like GGA => lat=..., lon=..., alt=... or UBX lines like NAV-PVT => fixType=..., lat=..., lon=....

## Access InfluxDB & Grafana:

    InfluxDB: http://rpi:8086 (login with INFLUXDB_USER/PASS)
    Grafana: http://rpi:3000 (default is admin/admin unless overridden)

## Usage / Examples

**Editing Baud Rate**

If your ZED-F9 is configured for 921600 or 115200, update baudRate in main.go or your environment.

**Disabling Extra Messages**

In u-center or via UBX CFG-MSG commands, uncheck unwanted messages (GNGST, GNGRS, etc.) if you don’t want parse errors for them.

**Extending UBX Parsing**

Add new parse functions in `ubx/parser.go` and create a file like nav_XXXX.go for each UBX message you need (e.g., NAV-SAT).

## Demo Screenshots

**Check Kernel Ring Buffer and USB Bus**

![Check KRB and USB Bus](/screenshots/check-usb-u-blox.png)

**UBX and NMEA on the same serial port**

![UBX and NMEA on the same serial port](/screenshots/ubx-and-nmea-raw-serial.png)

**InfluxDB Data Explorer**

![InfluxDB Data Explorer](/screenshots/influxdb-data-explorer.png)

**Geomap**

![Grafana - Geomap - Lat Long](/screenshots/grafana-geomap-lat-lon.png)

**Groud Speed over time**

![Grafana - Ground Speed](/screenshots/grafana-ground-speed.png)




## Development


**Local Go build (without Docker):**
    
    go mod tidy
    go build -o gnss-logger main.go
    ./gnss-logger
    

This will read from /dev/ttyACM0 and try to connect to Influx at http://localhost:8086 (depending on your config).

**Testing additional UBX messages:**

Add test code in ubx/..._test.go.

You can feed sample UBX payloads to ensure your parse functions work.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

**Fork this repo**

Create a feature branch (git checkout -b new-feature)
Commit changes & push
Submit a Pull Request
Bug reports / feature requests are welcome via Issues.

## Contact

Author: @miketigerblue
Open an Issue if you have questions or suggestions.
Enjoy logging your GNSS data with full NMEA + UBX parsing!