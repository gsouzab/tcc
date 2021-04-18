<template>
  <v-container fluid grid-list-md>
    <TimePicker @change="onDateRangeChange"></TimePicker>
    <div class="pa-2">

      <v-layout row wrap>
        <v-flex lg12 md12>
          <v-card>
            <v-card-title>
              Presença
              <v-progress-linear indeterminate v-if="loading"></v-progress-linear>
            </v-card-title>
            <v-card-text style="min-height: 200px">
              <div v-if="!loading">
                <p v-if="presenceChartCount == 0">Nenhuma dado disponível.</p>
                <plotly-chart v-if="presenceChartCount > 0" :chart="presenceChart"></plotly-chart>
              </div>
            </v-card-text>
          </v-card>
        </v-flex>
        <v-flex lg12 md12>
          <v-card>
            <v-card-title>
              Temperatura
              <v-progress-linear indeterminate v-if="loading"></v-progress-linear>
            </v-card-title>
            <v-card-text style="min-height: 200px">
              <div v-if="!loading">
                <p v-if="temperatureChartCount == 0">Nenhuma dado disponível.</p>
                <plotly-chart v-if="temperatureChartCount > 0" :chart="temperatureChart"></plotly-chart>
              </div>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>

      <v-layout row wrap>
        <v-flex lg12 md12>
          <v-card>
            <v-card-title>
              Umidade
              <v-progress-linear indeterminate v-if="loading"></v-progress-linear>
            </v-card-title>
            <v-card-text style="min-height: 200px">
              <div v-if="!loading">
                <p v-if="humidityChartCount == 0">Nenhuma dado disponível.</p>
                <plotly-chart v-if="humidityChartCount > 0" :chart="humidityChart"></plotly-chart>
              </div>
            </v-card-text>
          </v-card>
        </v-flex>

        <v-flex lg12 md12>
          <v-card>
            <v-card-title>
              CO2
              <v-progress-linear indeterminate v-if="loading"></v-progress-linear>
            </v-card-title>
            <v-card-text style="min-height: 200px">
              <div v-if="!loading">
                <p v-if="co2ChartCount == 0">Nenhuma dado disponível.</p>
                <plotly-chart :chart="co2Chart" v-if="co2ChartCount > 0"></plotly-chart>
              </div>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </div>
  </v-container>
</template>

<script>
import axios from 'axios';
import * as moment from 'moment';
import _ from 'lodash';
import TimePicker from '@/components/TimePicker';
import PlotlyChart from './charts/PlotlyChart';

const sensorsConfig = {};

export default {
  name: 'Dashboard',
  components: {
    PlotlyChart,
    TimePicker,
  },
  data() {
    return {
      sensors: [],
      meanInterval: null,
      loading: false,
      temperatureChartCount: 0,
      co2ChartCount: 0,
      humidityChartCount: 0,
      presenceChartCount: 0,
      temperatureChart: {
        uuid: 'temp-chart',
        traces: [],
        layout: {
          margin: {
            t: '20px',
          },
          xaxis: {
            type: 'date',
          },
          yaxis: {
            title: 'ºC',
            range: [-10, 50],
          },
        },
      },
      co2Chart: {
        uuid: 'co2-chart',
        traces: [],
        layout: {
          margin: {
            t: '20px',
          },
          xaxis: {
            type: 'date',
          },
          yaxis: {
            title: 'ppm',
          },
        },
      },
      humidityChart: {
        uuid: 'hum-chart',
        traces: [],
        layout: {
          margin: {
            t: '20px',
          },
          xaxis: {
            type: 'date',
          },
          yaxis: {
            title: '%',
            range: [0, 100],
          },
        },
      },
      presenceChart: {
        uuid: 'pres-chart',
        traces: [],
        layout: {
          barmode: 'stack',
          margin: {
            t: '20px',
          },
          xaxis: {
            type: 'date',
          },
          yaxis: {
            title: '# de dispositivos',
          },
        },
      },
    };
  },
  async mounted() {
    this.sensors = await this.getSensors();

    this.loading = true;
    await this.initRealtimeData();
    this.loading = false;
  },
  methods: {
    async getSensors() {
      try {
        const response = await axios.get(`http://${process.env.API_HOST}/sensors`);
        if (response.status === 200) {
          return response.data.data;
        }
      } catch (error) {
        console.error(error);
      }

      return [];
    },
    initTraces(sensors) {
      this.temperatureChart.traces = [];
      this.co2Chart.traces = [];
      this.humidityChart.traces = [];
      this.presenceChart.traces = [];

      this.temperatureChartCount = 0;
      this.co2ChartCount = 0;
      this.humidityChartCount = 0;
      this.presenceChartCount = 0;

      _.forEach(sensors, (value, key) => {
        sensorsConfig[value.mac] = { index: key, color: '#5e9e7e' };
        this.temperatureChart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
        this.co2Chart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
        this.humidityChart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
        this.presenceChart.traces.push({ x: [], y: [], name: value.name, type: 'scatter', fill: 'tozeroy' });
      });
    },

    addTelemetryData(data) {
      const dateTime = new Date(data.createdAt);
      const traceIndex = sensorsConfig[data.sensor].index;

      this.co2Chart.traces[traceIndex].y.push(data.co2);
      this.co2Chart.traces[traceIndex].x.push(dateTime);

      this.temperatureChart.traces[traceIndex].y.push(data.temp);
      this.temperatureChart.traces[traceIndex].x.push(dateTime);

      this.humidityChart.traces[traceIndex].y.push(data.hum);
      this.humidityChart.traces[traceIndex].x.push(dateTime);

      this.temperatureChart.layout.datarevision = dateTime.getTime();
      this.co2Chart.layout.datarevision = dateTime.getTime();
      this.humidityChart.layout.datarevision = dateTime.getTime();

      this.temperatureChartCount++;
      this.co2ChartCount++;
      this.humidityChartCount++;
    },
    addProbeData(data) {
      const dateTime = new Date(data.createdAt);
      const traceIndex = sensorsConfig[data.sensor].index;

      this.presenceChart.traces[traceIndex].x.push(dateTime);
      this.presenceChart.traces[traceIndex].y.push(data.count);
      this.presenceChart.layout.datarevision = dateTime.getTime();

      this.presenceChartCount += 1;
    },
    addOnMessageListener() {
      this.$telemetryWS.onmessage = (data) => {
        const measurements = data.data.split('\n');
        measurements.forEach((m) => {
          if (m !== '') this.addTelemetryData(JSON.parse(m));
        });
      };

      this.$probesWS.onmessage = (data) => {
        const measurements = data.data.split('\n');
        for (const m of measurements) {
          if (m !== '') this.addProbeData(JSON.parse(m));
        }
      };
    },
    async getPresenceData(whereStartTime, whereEndTime) {
      const duration = moment.duration(moment(whereEndTime).diff(moment(whereStartTime))).as('minutes');
      const maxPoints = 1000;
      const selectMeanInterval = Math.ceil(duration / maxPoints);
      const query = `http://174.138.126.228:30004/query?pretty=true&db=tcc_data&q=SELECT COUNT(DISTINCT("srcMac")) FROM "probe_data" WHERE time >= '${moment(whereStartTime).toISOString()}' AND time <= '${moment(whereEndTime).toISOString()}' GROUP BY time(${selectMeanInterval}m), sensor`;

      try {
        const response = await axios.get(query);
        if (response.status === 200) {
          if (response.data.results[0].series) {
            return response.data.results[0].series;
          }
          return [];
        }
      } catch (error) {
        console.log(error);
      }
    },
    async getTelemetryData(whereStartTime, whereEndTime, sensors, selectMeanField) {
      const duration = moment.duration(moment(whereEndTime).diff(moment(whereStartTime))).as('minutes');

      const maxPoints = 1440 * 3;
      const selectMeanInterval = Math.ceil(duration / maxPoints);

      this.meanInterval = selectMeanInterval;
      try {
        const response = await axios.post(`http://${process.env.API_HOST}/telemetry/query`, {
          whereStartTime,
          whereEndTime,
          selectMeanField,
          selectMeanInterval: selectMeanInterval.toString(),
          GroupByTag: 'sensor',
        });

        if (response.status === 200) {
          return response.data.data;
        }
      } catch (error) {
        console.log(error);
      }

      return [];
    },
    async onDateRangeChange(startDate, endDate) {
      this.loading = true;

      if (startDate === null && endDate === null) {
        await this.initRealtimeData();
      } else {
        this.$telemetryWS.onmessage = null;
        this.$probesWS.onmessage = null;

        this.initTraces(this.sensors);

        const temperatureData = await this.getTelemetryData(startDate, endDate, this.sensors, 'temp');
        const co2Data = await this.getTelemetryData(startDate, endDate, this.sensors, 'co2');
        const humData = await this.getTelemetryData(startDate, endDate, this.sensors, 'hum');
        const presenceData = await this.getPresenceData(startDate, endDate, this.sensors);

        this.fillTelemetryMeasurements(
          { chart: 'temperatureChart', data: temperatureData },
          { chart: 'co2Chart', data: co2Data },
          { chart: 'humidityChart', data: humData });

        this.fillPresenceMeasurements(presenceData);
      }
      this.loading = false;
    },
    async initRealtimeData() {
      const now = moment();
      const endDate = now.format('YYYY-MM-DD\\THH:mm:ssZ');
      const startDate = now.subtract(1, 'minute').format('YYYY-MM-DD\\THH:mm:ssZ');

      const temperatureData = await this.getTelemetryData(startDate, endDate, this.sensors, 'temp');
      const co2Data = await this.getTelemetryData(startDate, endDate, this.sensors, 'co2');
      const humData = await this.getTelemetryData(startDate, endDate, this.sensors, 'hum');
      const presenceData = await this.getPresenceData(startDate, endDate, this.sensors);

      this.initTraces(this.sensors);

      this.fillTelemetryMeasurements(
        { chart: 'temperatureChart', data: temperatureData },
        { chart: 'co2Chart', data: co2Data },
        { chart: 'humidityChart', data: humData });

      this.fillPresenceMeasurements(presenceData);
      this.addOnMessageListener();
    },
    fillTelemetryMeasurements(...telemetryMeasurements) {
      for (const measurement of telemetryMeasurements) {
        _.forEach(measurement.data, (groupedData) => {
          if (!sensorsConfig[groupedData.tags.sensor]) return;

          const traceIndex = sensorsConfig[groupedData.tags.sensor].index;

          _.forEach(groupedData.values, (data) => {
            if (data[1] !== null) {
              this[measurement.chart].traces[traceIndex].x.push(new Date(data[0]));
              this[measurement.chart].traces[traceIndex].y.push(data[1]);

              this[`${measurement.chart}Count`]++;
            }
          });

          this[measurement.chart].layout.datarevision = new Date().getTime();
        });
      }
    },
    fillPresenceMeasurements(presenceData) {
      for (const sensorData of presenceData) {
        if (!sensorsConfig[sensorData.tags.sensor]) return;

        const traceIndex = sensorsConfig[sensorData.tags.sensor].index;
        _.forEach(sensorData.values, (data) => {
          this.presenceChart.traces[traceIndex].x.push(new Date(data[0]));
          this.presenceChart.traces[traceIndex].y.push(data[1]);
          this.presenceChartCount++;
        });
      }
      this.presenceChart.layout.datarevision = new Date().getTime();
    },
  },
};
</script>

<style scoped>
  .container {
    padding: 0px;
  }
</style>
