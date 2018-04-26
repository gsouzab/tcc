<template>
  <v-container fluid grid-list-md>
    <SensorForm :visible="showSensorForm" @close="showSensorForm=false" @onSave="addSensor" />

    <gmap-map
      :center="center"
      :zoom="15"
      style="width: 100%; height: calc(100% - 64px); position: absolute; left:0;">

      <gmap-marker
        :key="index"
        v-for="(s, index) in sensors"
        :position="{lat: s.latitude, lon: s.longitude}"
        :clickable="true"
        @click="center={lat: s.latitude, lon: s.longitude}">

      </gmap-marker>

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
import SensorForm from "@/components/forms/SensorForm";

export default {
  data() {
    return {
      center: { lat: -22.8617784, lng: -43.2296038 },
      markers: [],
      sensors: [],
      showSensorForm: false
    };
  },
  methods: {
    addSensor: function(sensorData) {
      console.log(sensorData);
      this.sensors.push(sensorData);
    },
    getSensors: () => {
      return [];
    },
    getLastProbes: () => {
      return [];
    },
    connectWS: () => {
      /**
       * Method to connect to the WS instance
       */
    }
  },
  components: {
    SensorForm
  }
};
</script>

<style scoped>
.speed-dial {
  position: absolute;
  right: 48px;
}

.btn--floating {
  position: relative;
}
</style>
