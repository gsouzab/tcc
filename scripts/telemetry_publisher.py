from random import randint
import paho.mqtt.client as mqtt
import time
import json

QOS = 1
SENSORS = [
    'B8:27:EB:F2:64:A9',
    'B8:27:EB:52:43:B7',
]

# The "random" message creation, in json format
def create_json_message(sensor): 
    message = {}
    message['temp'] = randint(38,40)
    message['hum'] = randint(65,70)
    message['co2'] = randint(1100,1110)
    message['sensor'] = sensor
    message['createdAt'] = int(time.time() * 1000)
    return json.dumps(message)

def publish(client):
    for sensor in SENSORS:
        jsonMessage = create_json_message(sensor)
        print("enviando mensagem"+ jsonMessage)
        client.publish('telemetry', jsonMessage, QOS)

client = mqtt.Client()
client.connect("localhost", 1883, 60)
client.loop_start()

while True:
    publish(client)
    time.sleep(1.5)