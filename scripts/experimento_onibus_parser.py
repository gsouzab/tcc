import pandas as pd
from math import log10
from datetime import datetime

def dateparse (time_in_secs):    
    return datetime.utcfromtimestamp(float(time_in_secs))

def dateparse_mili (time_in_mili):    
    return datetime.utcfromtimestamp(float(time_in_mili) / 1000)

def dateparse_format (date_string):    
    return datetime.strptime(date_string, '%Y-%m-%dT%H:%M:%S.%fZ')

def int_len(n):
    return int(log10(abs(n))) + 1
dfGPS = pd.read_csv("dados_experimento/gps.csv", usecols=[0,1,2], sep=";", index_col=[0], parse_dates=["time"], date_parser=dateparse_format)
dfProbe = pd.read_csv("dados_experimento/probe_volta.csv", usecols=[0,1,4], sep=";", index_col=[0], names=["time", "rssi", "mac"], parse_dates=True, date_parser=dateparse)
# dfSherlock = pd.read_csv("dados_experimento/passageiros_ida.txt", index_col=[0], names=["time", "n_sherlock"], parse_dates=True, date_parser=dateparse_mili)

dfGPS.sort_index(inplace=True)
# dfProbe = dfProbe[dfProbe.rssi > -60]
# dfProbe = dfProbe.groupby(pd.TimeGrouper(freq='1S')).mac.nunique().to_frame()
dfProbe = dfProbe.reset_index().drop_duplicates(subset=['time','mac']).set_index('time')
# dfProbe.drop_duplicates(inplace=True,subset=["mac"])

# dfProbe = dfProbe.groupby(dfProbe.index.to_period('T')).mac.nunique().to_frame()
# dfProbe = dfProbe.groupby(dfProbe.index.to_period('T')).mac.nunique().to_frame()

# teste = pd.merge_asof(dfProbe, dfSherlock, left_index=True, right_index=True)
teste = pd.merge_asof(dfProbe, dfGPS, left_index=True, right_index=True)

# print dfProbe.groupby(dfProbe.index.to_period('T')).mac.nunique().to_frame()

# teste.to_csv("sem_filtro_rssi_30seg.csv")
teste.lat = teste.lat.apply(lambda lat: lat / float(10** (int_len(lat) - 2)))
teste.lon = teste.lon.apply(lambda lon: lon / float(10** (int_len(lon) - 2)))

teste.to_csv("volta_preprocessada.csv")

# print dfProbe.groupby("time").mac.nunique()