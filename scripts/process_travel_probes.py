import pandas as pd
from datetime import datetime
import math
import matplotlib.pyplot as plt
from matplotlib.ticker import StrMethodFormatter
from datetime import timedelta  

MIN_DURATION = 180

# N = 60
M = 300
P = 60

N_PESSOAS = 0

probe_file = "files/exp_20190329/volta.csv"

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
        self.first_reading = last_reading
        self.seen_count = 1
        self.last_reading = last_reading
        self.duration = 0

    def updateReading(self, last_reading):
        self.last_reading = last_reading
        self.duration = (self.last_reading - self.first_reading).seconds

    def updateState(self, state):
        self.state = state

    def increase_seen_count(self):
        self.seen_count += 1

    def reset_seen_count(self):
        self.seen_count = 1

devices = dict()

def parse_date(str_date):
    return datetime.strptime(str_date, '%Y-%m-%d %H:%M:%S')

def processArrival(time, reading, device):
    global N_PESSOAS
    time_diff = time - device.last_reading
    device.updateReading(time)
    device.increase_seen_count()

    if device.duration >= MIN_DURATION and device.state == 'possible_presence':        
        N_PESSOAS += 1
        print("%s, %f, %f, %d" % (time, reading['lat'], reading['lon'], N_PESSOAS))
        device.updateState('confirmed_presence')
    elif time_diff.seconds <= P and device.state == 'possible_absence':
        device.updateState('confirmed_presence')

    
def processDeparture(time, device, lat, lon):
    global N_PESSOAS
    time_diff = time - device.last_reading

    if time_diff.seconds > M and device.state == 'confirmed_presence':
        device.updateState('possible_absence')
    elif time_diff.seconds > P and device.state == 'possible_absence':
        device.updateState('confirmed_absence')
        N_PESSOAS -= 1
        print("%s, %f, %f, %d" % (time, lat, lon, N_PESSOAS))

def processReading(time, reading):
    mac = reading['mac']
    
    try:
        device = devices[mac]
        processArrival(time, reading, device)
    except KeyError:
    	devices[mac] = Device(mac, time)
    
dfIda = pd.read_csv(probe_file, date_parser=parse_date, parse_dates=['time'], index_col=[0])

min_time = dfIda.index.min()
max_time = dfIda.index.max() + timedelta(minutes=5)

while min_time < max_time:
    start_range = min_time
    min_time += timedelta(seconds=30)

    lat = lon = 0
    for time, row in dfIda.loc[start_range:min_time].iterrows():
        processReading(time, row)

        lat = row['lat']
        lon = row['lon']

    for m, d in devices.iteritems():
        processDeparture(min_time, d, lat, lon)