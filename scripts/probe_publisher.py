from random import randint
import paho.mqtt.client as mqtt
import time
import json

QOS = 1
N = 100
sensors_file = open("sensors.txt", "r")
SENSORS = sensors_file.read().splitlines()

NEARBY_DEVICES = [
    '11:11:11:11:11:11',
    '22:22:22:22:22:21',
    '22:22:22:22:22:25',
    '22:22:22:22:22:26',
    '22:22:22:22:22:28',
    '22:22:22:22:22:29',
    '22:22:22:22:22:30',
    '22:22:22:22:22:31',
]

# The "random" message creation, in json format
def create_json_message(srcMac, sensor): 
    message = {}
    message['rssi'] = randint(-75,-20)
    message['sensor'] = sensor
    message['srcMac'] = srcMac
    message['createdAt'] = int(time.time() * 1000)
    return json.dumps(message)

def publish(client):
    for sensor in SENSORS[0:N]:
        for srcMac in NEARBY_DEVICES:
            jsonMessage = create_json_message(srcMac, sensor)
            print(jsonMessage)
            client.publish('telemetry/probes', jsonMessage, QOS)

client = mqtt.Client()
client.connect("localhost", 1883, 60)
client.loop_start()

while True:
    publish(client)
    time.sleep(1)

