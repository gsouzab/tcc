<template>
  <v-container fluid grid-list-md>
    <v-layout row wrap>
      <v-flex lg12>
        <v-card>
          <v-card-text class="py-0">
            <plotly-chart :chart="temperatureChart"></plotly-chart>
          </v-card-text>
        </v-card>
      </v-flex>

      <v-flex lg6>
        <v-card>
          <v-card-title>Co2</v-card-title>
          <v-card-text >6</v-card-text>
        </v-card>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex lg6>
        <v-card>
          <v-card-title>Umidade</v-card-title>
          <v-card-text >6</v-card-text>
        </v-card>
      </v-flex>

      <v-flex lg6>
        <v-card>
          <v-card-title>Co2</v-card-title>
          <v-card-text >6</v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import PlotlyChart from './charts/PlotlyChart.vue'

export default {
  name: 'Dashboard',
  components: {
    PlotlyChart
  },
  data () {
    return {
      temperatureChart: {
        uuid: "temp-chart",
        traces: [
          {
            y: [],
            x: [],
            line: {
              color: "#5e9e7e",
              width: 2,
              shape: "spline"
            }
          }
        ],
        layout: {
          title:'Temperatura',
          xaxis: {
            title: 'Data e hora'
          },
          yaxis: {
            title: 'ºC'
          }
        }
      }
    }
  },
  mounted () {
    this.$options.sockets.onmessage = (data) => this.addTelemetryData(JSON.parse(data.data));

    // TODO: fetch initial data
    // for (let i = 0; i < 100; i++) {
    //   this.addData(i);
    // }
  },
  methods: {
    addTelemetryData(data) {
      console.log(data.temp);
      this.temperatureChart.layout.datarevision = new Date().getTime();
      this.temperatureChart.traces[0].y.push(data.temp);
      this.temperatureChart.traces[0].x.push(new Date(data.createdAt));
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
