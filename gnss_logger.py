#!/usr/bin/env python3

import serial
import time
import datetime
import os

LOG_DIR = "/home/pi/gnss_logs"
SERIAL_PORT = "/dev/ttyACM0"
BAUD_RATE = 921600

def main():
    if not os.path.exists(LOG_DIR):
        os.makedirs(LOG_DIR)

    # Create a timestamped log file
    log_filename = f"gnss_{datetime.datetime.now().strftime('%Y%m%d_%H%M%S')}.log"
    logfile_path = os.path.join(LOG_DIR, log_filename)

    with serial.Serial(SERIAL_PORT, BAUD_RATE, timeout=1) as ser, \
         open(logfile_path, 'w') as logfile:
        
        while True:
            try:
                line = ser.readline().decode('ascii', errors='replace').strip()
                if line:
                    print(line)
                    logfile.write(line + "\n")
                time.sleep(0.1)
            except KeyboardInterrupt:
                print("Exiting...")
                break
            except Exception as e:
                print(f"Error: {e}")
                time.sleep(1)

if __name__ == "__main__":
    main()

