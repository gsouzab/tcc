import pandas as pd
import matplotlib.pyplot as plt

estimated = pd.read_csv('entrada_180_saida_360.csv', sep=',', parse_dates=['time'])
ground_truth = pd.read_csv('files/exp_20190329/volta_real.csv', sep=',', parse_dates=['time'])

# estimated['time'] = estimated['time'] - pd.Timedelta(minutes=10)

ax = estimated.plot(x='time', y='estimated')
ax = ground_truth.plot(ax=ax, x='time', y='n_sherlock')

plt.show()