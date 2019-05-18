import pandas as pd
from datetime import datetime
import math

N = 5
M = 5
P = 5

N_PESSOAS = 0

class Device:
    states = ['possible_presence', 'confirmed_presence', 'possible_absence', 'confirmed_absence']

    def __eq__(self, other):
        if isinstance(other, Device):
            return self.mac == other.mac

        return False

    def __ne__(self, other):
        return not self.__eq__(other)

    def __init__(self, mac, last_reading):
        self.mac = mac
        self.state = self.states[0]
        self.last_reading = last_reading

    def updateReading(self, last_reading):
        self.last_reading = last_reading

    def updateState(self, state):
        self.state = state

devices = dict()

def parse_date(str_date):
    return datetime.strptime(str_date, '%Y-%m-%d %H:%M:%S')

def processReading(reading, device):
    global N_PESSOAS
    time_diff = reading['time'] - device.last_reading

    if time_diff.seconds < 5:
        if device.state == 'possible_presence':
            N_PESSOAS += 1
            device.updateState('confirmed_presence')
        elif device.state == 'possible_absence':
            device.updateState('confirmed_presence')

    elif time_diff.seconds > 5:
        if device.state == 'confirmed_presence':
            device.updateState('possible_absence')
        elif device.state == 'possible_absence':
            device.updateState('confirmed_absence')
            N_PESSOAS -= 1

def processDevices(reading):
    mac = reading['mac']
    
    try:
        device = devices[mac]
        processReading(reading, device)
        device.updateReading(reading['time'])
    except KeyError:
    	devices[mac] = Device(mac, reading['time'])
    
    for m, d in devices.items():
        if m != mac:
            processReading(reading, d)
    
    print(reading['time'], N_PESSOAS)
    
dfIda = pd.read_csv('files/ida_preprocessada.csv', sep=',', date_parser=parse_date, parse_dates=['time'])

for index, row in dfIda.iterrows():
    if row['rssi'] > -90:
        processDevices(row)

