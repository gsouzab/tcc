<template>
  <v-container fluid grid-list-md>
    <SensorForm :visible="showSensorForm"
                :initLat="currLat"
                :initLng="currLng"
                @close="showSensorForm=false"
                @onSave="addSensor" />

    <gmap-map
      :center="center"
      :zoom="15"
      @rightclick="this.showMenu"
      class="map">

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
          @click="toggleInfoWindow(s, index)">

        </gmap-marker>
      </gmap-cluster>

    </gmap-map>
    <v-speed-dial
      :bottom=true
      :right=true
      :open-on-hover=true>
      <v-btn
        slot="activator"
        color="blue darken-2"
        dark
        fab
        hover>
        <v-icon>add</v-icon>
        <v-icon>close</v-icon>
      </v-btn>
      <v-btn
        fab
        dark
        small
        color="green"
        @click="showSensorForm=true">
        <v-icon>settings_input_antenna</v-icon>
      </v-btn>
      <v-btn
        fab
        dark
        small
        color="indigo">
        <v-icon>layers</v-icon>
      </v-btn>
    </v-speed-dial>
  </v-container>
</template>

<script>

import axios from 'axios';
import {loaded} from 'vue2-google-maps';
import SensorForm from '@/components/forms/SensorForm';

export default {
  created() {
    this.getSensors();
  },
  data() {
    return {
      center: { lat: -22.8617784, lng: -43.2296038 },
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
    showMenu(event) {
      this.currLat = event.latLng.lat();
      this.currLng = event.latLng.lng();

      this.showSensorForm = true;
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
    toggleInfoWindow: function(s, idx) {
      let position = {lat: parseFloat(s.latitude), lng: parseFloat(s.longitude)};

      this.center = position;
      this.infoWindowPos = position;
      this.infoContent = `${s.name}<br/>${s.description}`;
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
  width: calc(100% - 80px);
  height: calc(100% - 64px);
  position: absolute;
  left:80px;
}

.speed-dial {
  position: absolute;
  right: 48px;
}

.btn--floating {
  position: relative;
}
</style>
