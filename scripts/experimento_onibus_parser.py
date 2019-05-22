import pandas as pd
from math import log10
from datetime import datetime

gps_file = "files/exp_20190329/gps.csv"
probe_file = "files/exp_20190329/probe_volta.csv"
real_file = "files/exp_20190329/pessoas_volta.csv"

out_file_probe = "files/exp_20190329/volta.csv"
out_file_real = "files/exp_20190329/ida_volta.csv"

def dateparse (time_in_secs):    
    return datetime.utcfromtimestamp(float(time_in_secs))

def dateparse_mili (time_in_mili):    
    return datetime.utcfromtimestamp(float(time_in_mili) / 1000)

def dateparse_format (date_string):    
    return datetime.strptime(date_string, '%Y-%m-%dT%H:%M:%S.%fZ')

def int_len(n):
    return int(log10(abs(n))) + 1
dfGPS = pd.read_csv(gps_file, usecols=[0,1,2], index_col=[0], names=["time", "lat", "lon"], parse_dates=["time"], date_parser=dateparse_format)
dfProbe = pd.read_csv(probe_file, usecols=[0,1,4], index_col=[0], names=["time", "rssi", "mac"], parse_dates=True, date_parser=dateparse)
dfSherlock = pd.read_csv(real_file, index_col=[0], names=["time", "n_sherlock"], parse_dates=True, date_parser=dateparse_mili)

dfGPS.sort_index(inplace=True)

dfProbe = dfProbe.reset_index().drop_duplicates(subset=['time','mac']).set_index('time')
dfProbe.sort_index(inplace=True)

out = pd.merge_asof(dfProbe, dfGPS, left_index=True, right_index=True)

out.lat = out.lat.apply(lambda lat: lat / float(10** (int_len(lat) - 2)))
out.lon = out.lon.apply(lambda lon: lon / float(10** (int_len(lon) - 2)))

out.to_csv(out_file_probe)

out_real = pd.merge_asof(dfSherlock, dfGPS, left_index=True, right_index=True)
out_real.lat = out_real.lat.apply(lambda lat: lat / float(10** (int_len(lat) - 2)))
out_real.lon = out_real.lon.apply(lambda lon: lon / float(10** (int_len(lon) - 2)))

out_real.to_csv(out_file_real)