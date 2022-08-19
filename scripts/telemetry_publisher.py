from random import randint
import paho.mqtt.client as mqtt
import time
import json

sensors_file = open("sensors.txt", "r")
N = 10
QOS = 1
SENSORS = sensors_file.read().splitlines()

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
    for sensor in SENSORS[0:N]:
        jsonMessage = create_json_message(sensor)
        print("enviando mensagem"+ jsonMessage)
        client.publish('telemetry', jsonMessage, QOS)

client = mqtt.Client(client_id="python-script", transport="websockets")
client.connect("mqtt.sensein.tech", 9001, 10)
client.loop_start()

while True:
    publish(client)
    time.sleep(1)
