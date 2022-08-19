<template>
  <v-container fluid>
    <SensorForm :visible="showSensorForm"
                :initData="sensorData"
                :isEdit="isEdit"
                @close="showSensorForm=false"
                @onSave="addSensor" />

    <gmap-map
      :center="center"
      :zoom="zoom"
      :options="{fullscreenControl: false, streetViewControl: false, mapTypeControl: false}"
      @rightclick="this.showMenu"
      slot="activator"
      class="map">

      <!-- <div slot="visible">
        <v-toolbar
          dense
          floating
          class="mt-3"
        >
          <v-flex xs12 sm6>
            <v-btn-toggle v-model="viewMode" mandatory class="transparent">
              <v-btn :value="SENSOR_VIEW_MODE" flat>
                Sensores
              </v-btn>
              <v-btn :value="TEMPERATURE_VIEW_MODE" flat>
                Temperatura
              </v-btn>
              <v-btn :value="CO2_VIEW_MODE" flat>
                CO2
              </v-btn>
              <v-btn :value="HUMIDITY_VIEW_MODE" flat>
                Umidade
              </v-btn>
            </v-btn-toggle>
          </v-flex>
        </v-toolbar>
      </div> -->

      <v-menu
        v-model="showContextMenu"
        :position-x="x"
        :position-y="y"
        absolute
        offset-y
      >
        <v-list>
          <v-list-tile
            @click="showSensorForm = true;  isEdit = false;"
            v-if="!isSensorMenu"
          >
            <v-list-tile-title>Adicionar sensor aqui</v-list-tile-title>
          </v-list-tile>

          <v-list-tile
            @click="showSensorForm = true; isEdit = true;"
            v-if="isSensorMenu"
          >
            <v-list-tile-title>Editar sensor</v-list-tile-title>
          </v-list-tile>

          <v-list-tile
            @click="removeSensor()"
            v-if="isSensorMenu"
          >
            <v-list-tile-title>Remover sensor</v-list-tile-title>
          </v-list-tile>
        </v-list>
      </v-menu>

      <section v-if="viewMode === TEMPERATURE_VIEW_MODE">
        <gmap-circle
          :key="index"
          v-for="(m, index) in measurementsArray"
          :radius="m.radius"
          :center="m.center"
          :options="m.options">
        </gmap-circle>
      </section>>

      <section v-if="viewMode === SENSOR_VIEW_MODE">
        <heatmap :data="heatmapDataArray" :options="heatmapOptions"></heatmap>

        <gmap-info-window :options="infoOptions" :position="infoWindowPos" :opened="infoWinOpen" @closeclick="infoWinOpen=false">
          <div v-html="infoContent"></div>
        </gmap-info-window>

        <gmap-cluster :zoomOnClick="true" :maxZoom=21>
          <gmap-marker
            :key="index"
            v-for="(s, index) in sensors"
            :icon="require('../assets/sensor.png')"
            :position="{lat: parseFloat(s.latitude), lng: parseFloat(s.longitude)}"
            :clickable="true"
            @click="toggleInfoWindow(s, index)"
            @rightclick="showSensorMenu($event, s, index)">

          </gmap-marker>
        </gmap-cluster>
      </section>

      <ground-overlay :source="require('../assets/lab_layer.png')" :bounds="{
          north: -22.861479854128124,
          south: -22.861765421786675,
          east: -43.22791502677535,
          west: -43.22829133787252,}"
          @rightclick="showMenu">
      </ground-overlay>


    </gmap-map>

    <div id="legend" v-show="viewMode !== SENSOR_VIEW_MODE"> </div>
  </v-container>
</template>

<script>

import axios from 'axios';
import { gmapApi } from 'vue2-google-maps';
import SensorForm from '@/components/forms/SensorForm';
import Heatmap from '@/components/Heatmap';
import GroundOverlay from '@/components/GroundOverlay';
import DatetimeFieldPicker from '@/components/DatetimeFieldPicker';
import _ from 'lodash';
import * as d3 from 'd3';
import * as moment from 'moment';

const sensorsConfig = {};

export default {
  async created() {
    await this.getSensors();
    await this.getLastTelemetry();
    this.connectWS();
    this.setLegend();
    setInterval(() => {
      this.clearOldProbes();
    }, 5000);
  },
  data() {
    return {
      center: { lat: -22.861659764789806, lng: -43.22840063788988 },
      zoom: 20,
      viewMode: 1,
      SENSOR_VIEW_MODE: 1,
      TEMPERATURE_VIEW_MODE: 2,
      CO2_VIEW_MODE: 3,
      HUMIDITY_VIEW_MODE: 4,
      heatmapData: [],
      heatmapOptions: {
        radius: 80,
        opacity: 0.8,
      },
      measurements: new Map(),
      measurementsTracker: 0,
      probes: new Map(),
      probesTracker: 0,
      showContextMenu: false,
      isSensorMenu: false,
      isEdit: false,
      sensorData: null,
      colorScale: null,
      x: 0,
      y: 0,
      sensors: [],
      currLat: null,
      currLng: null,
      showSensorForm: false,
      googleMapsInitialized: false,
      infoContent: '',
      infoWindowPos: null,
      infoWinOpen: false,
      currentMidx: null,
      infoOptions: {
        pixelOffset: {
          width: 0,
          height: -35,
        },
      },
    };
  },
  computed: {
    google: gmapApi,
    measurementsArray() {
      return (this.measurementsTracker && Array.from(this.measurements.values())) || [];
    },
    heatmapDataArray() {
      return (this.probesTracker && Array.from(this.probes.values())) || [];
    },
  },
  methods: {
    setLegend() {
      const w = 300,
        h = 40;

      const key = d3.select('#legend')
        .append('svg')
        .attr('width', w)
        .attr('height', h)
        .style('position', 'absolute')
        .style('right', '10px')
        .style('top', '0px');

      const legend = key.append('defs')
        .append('svg:linearGradient')
        .attr('id', 'gradient')
        .attr('x1', '0%')
        .attr('y1', '100%')
        .attr('x2', '100%')
        .attr('y2', '100%')
        .attr('spreadMethod', 'pad');

      legend.selectAll('stop')
        .data([
          { offset: '0%', color: '#2c7bb6' },
          { offset: '12.5%', color: '#00a6ca' },
          { offset: '25%', color: '#00ccbc' },
          { offset: '37.5%', color: '#90eb9d' },
          { offset: '50%', color: '#ffff8c' },
          { offset: '62.5%', color: '#f9d057' },
          { offset: '75%', color: '#f29e2e' },
          { offset: '87.5%', color: '#e76818' },
          { offset: '100%', color: '#d7191c' },
        ])
        .enter().append('stop')
        .attr('offset', d => d.offset)
        .attr('stop-color', d => d.color);


      key.append('rect')
        .attr('width', w)
        .attr('height', h - 25)
        .style('fill', 'url(#gradient)')
        .attr('transform', 'translate(0,5)');

      const y = d3.scaleLinear()
        .domain([10, 35])
        .range([10, 280]);

      const yAxis = d3.axisBottom()
        .scale(y)
        .tickFormat(d => `${d}°C`)
        .ticks(5);

      key.append('g')
        .attr('class', 'y axis')
        .attr('transform', 'translate(0,20)')
        .call(yAxis)
        .append('text')
        .attr('transform', 'rotate(-90)')
        .attr('y', 0)
        .attr('dy', '.71em')
        .style('text-anchor', 'end')
        .text('axis title');

      this.colorScale = d3.scaleLinear()
        .domain([10, 20, 35])
        .range(['#2c7bb6', '#90eb9d', '#d7191c'])
        .interpolate(d3.interpolateHcl);
    },
    addSensor(sensorData) {
      if (this.isEdit) {
        this.sensors[this.sensorData.index] = sensorData;
        this.updateInfoContent(sensorData);
      } else {
        this.sensors.push(sensorData);
      }
    },
    addTelemetryData(data) {
      const i = _.findIndex(this.sensors, { mac: data.sensor });
      this.sensors[i].telemetry = data;

      if (this.infoWinOpen && this.infoWindowSensorMac === data.sensor) {
        this.updateInfoContent(this.sensors[i]);
      }

      const m = {
        radius: 30,
        center: sensorsConfig[data.sensor].position,
        options: {
          strokeColor: this.colorScale(data.temp),
          strokeOpacity: 0,
          strokeWeight: 0,
          fillColor: this.colorScale(data.temp),
          fillOpacity: 0.7,
        },
      };

      this.measurements.set(data.sensor, m);
      this.measurementsTracker += 1;
    },
    addProbeData(data) {
      const i = _.findIndex(this.sensors, { mac: data.sensor });
      const heatmapData = {
        location: new google.maps.LatLng(this.sensors[i].latitude, this.sensors[i].longitude),
        weight: data.count,
      };

      this.sensors[i].probes = data;

      if (this.infoWinOpen && this.infoWindowSensorMac === data.sensor) {
        this.updateInfoContent(this.sensors[i]);
      }

      this.probes.set(data.sensor, heatmapData);
      this.probesTracker += 1;
    },
    clearOldProbes() {
      const currentTs = moment();
      _.forEach(this.sensors, (sensor) => {
        if ((sensor.probes && sensor.probes.count > 0) && currentTs.diff(moment(sensor.probes.createdAt), 'seconds') > 10) {
          sensor.probes.count = 0;
          sensor.probes.createdAt = currentTs.toDate();
          this.probes.delete(sensor.mac)

          if (this.infoWinOpen && this.infoWindowSensorMac === sensor.mac) {
            this.updateInfoContent(sensor);
          }
        }
      });
      this.probesTracker += 1;
    },
    openMenu(e) {
      this.showContextMenu = false;
      const mouseEvent = _.find(e, attr => attr.clientX && attr.clientY);

      this.x = mouseEvent.clientX;
      this.y = mouseEvent.clientY;
      this.$nextTick(() => {
        this.showContextMenu = true;
      });
    },
    showSensorMenu(e, sensor, index) {
      this.showContextMenu = false;
      this.isSensorMenu = true;

      this.sensorData = null;
      this.$nextTick(() => {
        this.sensorData = sensor;
        this.sensorData.index = index;
      });

      this.openMenu(e);
    },
    showMenu(e) {
      this.isSensorMenu = false;

      this.sensorData = {
        mac: '',
        name: '',
        description: '',
        latitude: e.latLng.lat(),
        longitude: e.latLng.lng(),
      };
      this.openMenu(e);
    },
    async removeSensor() {
      const sensor = this.sensorData;
      try {
        const response = await axios.delete(`${process.env.API_HOST}/sensors/${sensor.mac}`);
        if (response.status === 200) {
          this.sensors.splice(sensor.index, 1);
        }
      } catch (error) {
        console.error(error);
      }
    },
    async getSensors() {
      this.waitingSensors = true;

      try {
        const response = await axios.get(`${process.env.API_HOST}/sensors`);
        if (response.status === 200) {
          this.sensors = response.data.data;

          _.forEach(this.sensors, (s, index) => {
            sensorsConfig[s.mac] = { position: { lat: parseFloat(s.latitude), lng: parseFloat(s.longitude) }, index };
          });
        }
      } catch (error) {
        console.error(error);
      }

      this.waitingSensors = false;
    },
    async getLastTelemetry() {
      this.waitingSensors = true;

      try {
        const response = await axios.post(`${process.env.API_HOST}/telemetry/query`, {
          whereLastXMinutes: '5',
          selectMeanField: 'temp',
          GroupByTag: 'sensor',
        });
        if (response.status === 200) {
          return response.data.data;
        }
      } catch (error) {
        console.log(error);
      }

      this.waitingSensors = false;
    },
    getLastProbes() {
      return [];
    },
    connectWS() {
      /**
       * Method to connect to the WS instance
       */

      this.$probesWS.onmessage = (data) => {
        const measurements = data.data.split('\n');
        for (const m of measurements) {
          if (m !== '') this.addProbeData(JSON.parse(m));
        }
      };

      this.$telemetryWS.onmessage = (data) => {
        const measurements = data.data.split('\n');

        for (const m of measurements) {
          if (m !== '') this.addTelemetryData(JSON.parse(m));
        }
      };
    },
    updateInfoContent(sensor) {
      this.infoContent = `
        <h3>${sensor.name}</h3>
        <p style="white-space: pre;"><b>MAC</b>: ${sensor.mac} </br>${sensor.description || ''}</p>`;

      if (sensor.telemetry) {
        this.infoContent += `
          <p>
            Temperatura: ${sensor.telemetry.temp}º C <br>
            CO2: ${sensor.telemetry.co2} ppm<br>
            Umidade: ${sensor.telemetry.hum}% <br><br>
            Atualizado em: ${new Date(sensor.telemetry.createdAt).toLocaleString('pt-BR')}
          </p>`;
      }

      if (sensor.probes) {
        this.infoContent += `
          <p>
            Dispositivos nas proximidades: ${sensor.probes.count} <br><br>
            Atualizado em: ${new Date(sensor.probes.createdAt).toLocaleString('pt-BR')}
          </p>`;
      }
    },
    toggleInfoWindow(s, idx) {
      const position = { lat: parseFloat(s.latitude), lng: parseFloat(s.longitude) };

      this.updateInfoContent(s);
      // check if its the same marker that was selected if yes toggle
      if (this.currentMidx === idx) {
        this.infoWinOpen = !this.infoWinOpen;
      } else {
        // if different marker set infowindow to open and reset current marker index
        this.infoWinOpen = true;
        this.currentMidx = idx;
      }

      this.$nextTick(() => {
        this.center = position;
        this.infoWindowPos = position;
        this.infoWindowSensorMac = s.mac;
      });
    },
  },
  components: {
    SensorForm,
    Heatmap,
    GroundOverlay,
    DatetimeFieldPicker,
  },
};
</script>

<style>
.axis path,
.axis tick,
.axis line {
  fill: none;
  stroke: none;
}

.v-toolbar__content {
  height: auto !important;
  padding: 0;
}
</style>

<style scoped>
.map, .container {
  width: 100%;
  height: 100%;
}
.container.fluid, .container {
  padding: 0 !important;
}
</style>
