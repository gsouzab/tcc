import random
import requests
import json

N_SENSORS = 1000
CENTER_LAT = -22.861592962437157
CENTER_LON = -43.22808397926099

def create_sensor(i):
    mac = "02:00:00:%02X:%02X:%02X" % (random.randint(0, 255),
                             random.randint(0, 255),
                             random.randint(0, 255))

    return {
        "mac": mac,
        "name": "Sensor %d" % i,
        "description": "Sensor #%d" % i,
        "latitude": "%s" % (CENTER_LAT + (random.uniform(-1, 1) / 100)),
        "longitude": "%s" % (CENTER_LON + (random.uniform(-1, 1) / 100))
    }

for i in range(1, N_SENSORS + 1):
    sensor = create_sensor(i)
    requests.post("http://localhost:8000/sensors", data=json.dumps(sensor))

    print(sensor["mac"])
