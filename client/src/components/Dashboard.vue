<template>
  <v-container fluid grid-list-md>
    <TimePicker></TimePicker>
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
import PlotlyChart from './charts/PlotlyChart.vue'
import TimePicker from '@/components/TimePicker';
import axios from 'axios';
import _ from 'lodash';

var sensorsConfig = {};

export default {
  name: 'Dashboard',
  components: {
    PlotlyChart,
    TimePicker
  },
  data () {
    return {
      temperatureChart: {
        uuid: "temp-chart",
        traces: [],
        layout: {
          title:'Temperatura',
          xaxis: {
            title: 'Data e hora'
          },
          yaxis: {
            title: 'ºC'
          }
        }
      },
      co2Chart: {
        uuid: "co2-chart",
        traces: [],
        layout: {
          title:'Co2',
          xaxis: {
            title: 'Data e hora'
          },
          yaxis: {
            title: 'ppm'
          }
        }
      },
      humidityChart: {
        uuid: "hum-chart",
        traces: [],
        layout: {
          title:'Umidade',
          xaxis: {
            title: 'Data e hora'
          },
          yaxis: {
            title: '%'
          }
        }
      }
    }
  },
  async mounted() {
    let sensors = await this.getSensors();

    _.forEach(sensors, (value, key) => {
      sensorsConfig[value.mac] = {index: key, color: '#5e9e7e'};
      this.temperatureChart.traces.push({x: [], y: [], name: value.name});
      this.co2Chart.traces.push({x: [], y: [], name: value.name});
      this.humidityChart.traces.push({x: [], y: [], name: value.name});
    });

    this.$options.sockets.onmessage = (data) => this.addTelemetryData(JSON.parse(data.data));

    // TODO: fetch initial data
    // for (let i = 0; i < 100; i++) {
    //   this.addData(i);
    // }
  },
  methods: {
    async getSensors() {
      try {
        let response = await axios.get(`http://${process.env.API_HOST}/sensors`);
        if (response.status == 200) {
          return response.data.data;
        }
      } catch (error) {
        console.log(error);
      }

      return [];
    },
    addTelemetryData(data) {
      let dateTime = new Date(data.createdAt);
      let traceIndex = sensorsConfig[data.sensor].index;

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
    }
  }
};
</script>

<style scoped>
</style>
