<template>
  <v-container fluid grid-list-md>
    <TimePicker @change="onDateRangeChange"></TimePicker>
    <v-layout row wrap>
      <v-flex lg12>
        <v-card>
          <v-card-text class="py-0 px-0">
            <plotly-chart :chart="temperatureChart"></plotly-chart>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex lg6 md12>
        <v-card>
          <v-card-text class="py-0 px-0">
            <plotly-chart :chart="humidityChart"></plotly-chart>
          </v-card-text>
        </v-card>
      </v-flex>

      <v-flex lg6 md12>
        <v-card>
          <v-card-text class="py-0" px-0>
            <plotly-chart :chart="co2Chart"></plotly-chart>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
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
      temperatureChart: {
        uuid: 'temp-chart',
        traces: [],
        layout: {
          title: 'Temperatura',
          xaxis: {
            title: 'Data e hora',
          },
          yaxis: {
            title: 'ºC',
          },
        },
      },
      co2Chart: {
        uuid: 'co2-chart',
        traces: [],
        layout: {
          title: 'Co2',
          xaxis: {
            title: 'Data e hora',
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
          title: 'Umidade',
          xaxis: {
            title: 'Data e hora',
          },
          yaxis: {
            title: '%',
          },
        },
      },
    };
  },
  async mounted() {
    this.sensors = await this.getSensors();

    this.initTraces(this.sensors);
    this.addOnMessageListener();
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
      _.forEach(sensors, (value, key) => {
        sensorsConfig[value.mac] = { index: key, color: '#5e9e7e' };
        this.temperatureChart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
        this.co2Chart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
        this.humidityChart.traces.push({ x: [], y: [], name: value.name, line: { shape: 'spline' } });
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
    },
    addData(i) {
      this.temperatureChart.layout.datarevision = new Date().getTime();
      this.temperatureChart.traces[0].y.push(Math.random());
      this.temperatureChart.traces[0].x.push(i);
    },
    addOnMessageListener() {
      this.$options.sockets.onmessage = (data) => {
        const measurements = data.data.split('\n');
        measurements.forEach((m) => {
          if (m !== '') this.addTelemetryData(JSON.parse(m));
        });
      };
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
          GroupByTag: 'sensor'
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
      if (startDate === null && endDate === null) {
        this.addOnMessageListener();
        this.initTraces(this.sensors);
      } else {
        delete this.$options.sockets.onmessage;

        this.initTraces(this.sensors);

        const temperatureData = await this.getTelemetryData(startDate, endDate, this.sensors, "temp");
        const co2Data = await this.getTelemetryData(startDate, endDate, this.sensors, "co2");
        const humData = await this.getTelemetryData(startDate, endDate, this.sensors, "hum");

        this.fillTelemetryMeasurements(
          {chart: "temperatureChart", data: temperatureData},
          {chart: "co2Chart", data: co2Data},
          {chart: "humidityChart", data: humData})
      }
    },
    fillTelemetryMeasurements(...telemetryMeasurements) {

      for (let measurement of telemetryMeasurements) {
        _.forEach(measurement.data, (groupedData) => {
          if (!sensorsConfig[groupedData.tags.sensor]) return;

          let traceIndex = sensorsConfig[groupedData.tags.sensor].index;

          _.forEach(groupedData.values, (data) => {
            if (data[1] !== null) {
              this[measurement.chart].traces[traceIndex].x.push(new Date(data[0]))
              this[measurement.chart].traces[traceIndex].y.push(data[1])
            }

          });

          this[measurement.chart].layout.datarevision = new Date().getTime();
        });
      }

    }
  }
};
</script>

<style scoped>
  .container {
    padding: 0px;
  }
</style>
