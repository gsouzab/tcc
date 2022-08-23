from random import randint
import paho.mqtt.client as mqtt
import json
import time

MAX_NEW_NEARBY = 5
MIN_PRESENCE_TIME = 3

QOS = 1
sensors_file = open("sensors.txt", "r")
SENSORS = sensors_file.read().splitlines()


def random_new_nearby():
    new_nearby_count = max(0, randint(-MAX_NEW_NEARBY, MAX_NEW_NEARBY))

    return generate_devices(new_nearby_count)


def generate_mac_address():
    return "02:00:00:%02X:%02X:%02X" % (randint(0, 255),
                                        randint(0, 255),
                                        randint(0, 255))


def generate_devices(n):
    devices = list()

    for _ in range(n):
        device = {}
        device['mac'] = generate_mac_address()
        device['presence'] = randint(1, MIN_PRESENCE_TIME)
        devices.append(device)
    
    return devices


def create_json_message(srcMac, sensor): 
    message = {}
    message['rssi'] = randint(-75,-20)
    message['sensor'] = sensor
    message['srcMac'] = srcMac
    message['createdAt'] = int(time.time() * 1000)
    return json.dumps(message)


def publish(client, sensors):
    for sensor in sensors:
        for device in sensors[sensor]:
            jsonMessage = create_json_message(device['mac'], sensor)
            client.publish('probes', jsonMessage, QOS)


def refresh_sensors(sensors):
    for sensor in sensors:
        devices_refreshed = list()

        for device in sensors[sensor]:
            device['presence'] -= randint(0, 1)

            if(device['presence'] > 0):
                devices_refreshed.append(device)

        sensors[sensor] = devices_refreshed

client = mqtt.Client(client_id="python-script", transport="websockets")
client.connect("mqtt.sensein.tech", 9001, 10)
client.loop_start()

sensors = {}

for sensor in SENSORS:
    sensors[sensor] = random_new_nearby()

while True:
    publish(client, sensors)

    refresh_sensors(sensors)

    for sensor in sensors:
        sensors[sensor] += random_new_nearby()

    time.sleep(5)
