<template>
  <v-container fluid grid-list-md>
    <SensorForm :visible="showSensorForm"
                :initData="sensorData"
                :isEdit="isEdit"
                @close="showSensorForm=false"
                @onSave="addSensor" />

    <gmap-map
      :center="center"
      :zoom="15"
      @rightclick="this.showMenu"
      slot="activator"
      class="map">

      <v-menu
        v-model="showContextMenu"
        :position-x="x"
        :position-y="y"
        absolute
        offset-y
      >
        <v-list>
          <v-list-tile
            @click="showSensorForm = true"
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

      <gmap-info-window :options="infoOptions" :position="infoWindowPos" :opened="infoWinOpen" @closeclick="infoWinOpen=false">
        <div v-html="infoContent"></div>
      </gmap-info-window>

      <gmap-cluster :zoomOnClick="true">
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

    </gmap-map>
  </v-container>
</template>

<script>

import axios from 'axios';
import {loaded} from 'vue2-google-maps';
import SensorForm from '@/components/forms/SensorForm';
import _ from 'lodash';

export default {
  created() {
    this.getSensors();
    this.$options.sockets.onmessage = (data) => this.addTelemetryData(JSON.parse(data.data));
  },
  data() {
    return {
      center: { lat: -22.8617784, lng: -43.2296038 },
      showContextMenu: false,
      isSensorMenu: false,
      isEdit: false,
      sensorData: null,
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
          height: -35
        }
      },
    };
  },
  methods: {
    addSensor(sensorData) {
      this.sensors.push(sensorData);
    },
    addTelemetryData(data) {
      const i = _.findIndex(this.sensors, {mac: data.sensor});
      this.sensors[i].telemetry = data;

      if (this.infoWinOpen && this.infoWindowSensorMac == data.sensor) {
        this.updateInfoContent(this.sensors[i]);
      }
    },
    openMenu(e) {
      this.showContextMenu = false;

      this.x = e.wa.clientX;
      this.y = e.wa.clientY;
      this.$nextTick(() => {
        this.showContextMenu = true;
      })
    },
    showSensorMenu(e, sensor, index) {
      this.showContextMenu = false;
      this.isSensorMenu = true;

      this.sensorData = null;
      this.$nextTick(() => {
        this.sensorData = sensor;
      })

      this.openMenu(e);
    },
    showMenu(e) {

      this.isSensorMenu = false;

      this.sensorData = {
        mac: '',
        name: '',
        description: '',
        latitude: e.latLng.lat(),
        longitude: e.latLng.lng()
      }
      this.openMenu(e);
    },
    async removeSensor() {
      let sensor = this.sensorData;
      try {
        let response = await axios.delete(`http://${window.location.hostname}:8000/sensors/${sensor.mac}`);
        if (response.status == 200) {
          this.sensors.splice(sensor, 1);
        }
      } catch (error) {
        console.log(error);
      }
    },
    async getSensors() {
      this.waitingSensors = true;

      try {
        let response = await axios.get(`http://${window.location.hostname}:8000/sensors`);
        if (response.status == 200) {
          this.sensors = response.data.data;
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
    },
    updateInfoContent(sensor) {
      this.infoContent = `
        <h3>${sensor.name}</h3>
        ${sensor.description}`;

      if (sensor.telemetry) {
        this.infoContent += `
          <br>
          Temperatura: ${sensor.telemetry.temp}º C <br>
          CO2: ${sensor.telemetry.co2}<br>
          Umidade: ${sensor.telemetry.hum}% <br><br>
          Atualizado em: ${new Date(sensor.telemetry.createdAt).toLocaleString("pt-BR")}
        `;
      }
    },
    toggleInfoWindow: function(s, idx) {
      let position = {lat: parseFloat(s.latitude), lng: parseFloat(s.longitude)};

      this.infoWindowPos = position;
      this.infoWindowSensorMac = s.mac;
      this.updateInfoContent(s);
      //check if its the same marker that was selected if yes toggle
      if (this.currentMidx == idx) {
        this.infoWinOpen = !this.infoWinOpen;
      } else {
        //if different marker set infowindow to open and reset current marker index
        this.infoWinOpen = true;
        this.currentMidx = idx;
      }
    },
  },
  components: {
    SensorForm
  }
};
</script>

<style scoped>
.map {
  width: calc(100% - 60px);
  height: calc(100% - 46px);
  position: absolute;
  left: 60px;
}


.container.fluid, .container {
  padding: 0 !important;
}
</style>
