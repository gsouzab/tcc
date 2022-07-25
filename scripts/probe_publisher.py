from random import randint
from random import shuffle
import paho.mqtt.client as mqtt
import time
import json

QOS = 1
MAX_SENSORS = 500
MAX_NEARBY_DEVICES = 100

sensors_file = open("sensors.txt", "r")
SENSORS = sensors_file.read().splitlines()

def generate_mac_address():
    return "02:00:00:%02X:%02X:%02X" % (randint(0, 255),
                                        randint(0, 255),
                                        randint(0, 255))

def generate_nearby_devices(n):
    devices = list()

    for _ in range(n):
        devices.append(generate_mac_address())
    
    return devices

# The "random" message creation, in json format
def create_json_message(srcMac, sensor): 
    message = {}
    message['rssi'] = randint(-75,-20)
    message['sensor'] = sensor
    message['srcMac'] = srcMac
    message['createdAt'] = int(time.time() * 1000)
    return json.dumps(message)

def publish(client, sensors, nearby_devices):
    for sensor in sensors:
        for srcMac in nearby_devices:
            jsonMessage = create_json_message(srcMac, sensor)
            # print(jsonMessage)
            client.publish('telemetry/probes', jsonMessage, QOS)

client = mqtt.Client()
client.connect("localhost", 1883, 60)
client.loop_start()

i = 1
j = 1
sensors_count        = 10
curr_sensors         = SENSORS[0:sensors_count]
nearby_devices_count = 1
nearby_devices       = generate_nearby_devices(nearby_devices_count)

print("[Iteração %d] Número de sensores: %d, número de dispositivos na proximidade: %d" % (j, sensors_count, nearby_devices_count))
while True:
    if i % 10 == 0:
        j += 1
        if sensors_count < MAX_SENSORS:
            sensors_count += 10

        if nearby_devices_count < MAX_NEARBY_DEVICES:
            nearby_devices_count += 1

        shuffle(SENSORS)
        print("[Iteração %d] Número de sensores: %d, número de dispositivos na proximidade: %d" % (j, sensors_count, nearby_devices_count))
        curr_sensors   = SENSORS[0:sensors_count]
        nearby_devices = generate_nearby_devices(nearby_devices_count)
    
    publish(client, curr_sensors, nearby_devices)
    time.sleep(1)
    i += 1
