import paho.mqtt.client as mqtt
from multiprocessing import Pool
import datetime
import random
import time
import json

QOS = 1
SENSORS = [
    'B8:27:EB:F2:64:A9',
    'B8:27:EB:52:43:B7',
]

# The "random" message creation, in json format
def create_json_message(): 
    message = {}
    message['temp'] = random.randint(12,40)
    message['hum'] = random.randint(1,99)
    message['co2'] = random.randint(600,1200)
    message['sensor'] = random.choice(SENSORS)
    message['createdAt'] = int(float(datetime.datetime.now().strftime('%s.%f')) * 1000)
    
    return json.dumps(message)

# The callback for when the client receives a CONNACK response from the server.
def on_connect(client, userdata, flags, rc):
    print("Connected with result code "+str(rc))

# The callback for when a PUBLISH message is received from the server.
def on_message(client, userdata, msg):
    print(msg.topic+" "+str(msg.payload))

def publish(client):
    jsonMessage = create_json_message()
    # print jsonMessage
    client.publish('telemetry', jsonMessage, QOS)

client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect("174.138.126.228", 30000, 60)
client.loop_start()

while True:
    publish(client)
    time.sleep(2)

