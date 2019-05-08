<template>
  <v-container fluid>
    <gmap-map
      :center="center"
      :zoom="zoom"
      :options="{fullscreenControl: false, streetViewControl: false, mapTypeControl: false}"
      slot="activator"
      class="map">


      <gmap-polyline :key="path.id" v-for="path in paths" :options="{'strokeColor': path.color, 'strokeWeight': 6}" :path="path.snappedCoordinates" ></gmap-polyline>
      <!-- <heatmap :data="snappedCoordinates" ></heatmap> -->
    </gmap-map>

  </v-container>
</template>

<script>

import axios from 'axios';
import {loaded, gmapApi} from 'vue2-google-maps';
import Heatmap from '@/components/Heatmap';
import _ from 'lodash';
import * as d3 from 'd3';

export default {
  async created() {
    await this.runSnapToRoad();
  //   await this.getLastTelemetry();
  //   this.connectWS();
  //   this.setLegend();
  },
  data() {
    return {
      center: { lat: -22.861659764789806, lng: -43.22840063788988 },
      zoom: 13,
      paths: [],
    };
  },
  computed: {
    google: gmapApi,
  },
  methods: {
    async runSnapToRoad(path) {
      var pathValues = [
        "-22.91786373,-43.24705558",
        "-22.91765307,-43.24534777",
        "-22.91734733,-43.24391027",
        "-22.91705271,-43.242477",
        "-22.91704565,-43.24245743",
        "-22.91677556,-43.24134835",
        "-22.91647701,-43.23904644",
        "-22.91614004,-43.23718893",
        "-22.91597462,-43.23587119",
        "-22.91596102,-43.23586924",
        "-22.9146078,-43.23376945",
        "-22.91448082,-43.23360262",
        "-22.91433112,-43.23340576",
        "-22.91538595,-43.23234093",
        "-22.91578144,-43.23174912",
        "-22.91574339,-43.23157724",
        "-22.91619641,-43.23027063",
        "-22.91620987,-43.23028371",
        "-22.91617081,-43.22829075",
        "-22.9164546,-43.2265166",
        "-22.91754207,-43.22487945",
        "-22.9181644,-43.22442526",
        "-22.92050683,-43.22277233",
        "-22.92158869,-43.22025502",
        "-22.92121071,-43.21931594",
        "-22.92125678,-43.21914435",
        "-22.92070835,-43.21772955",
        "-22.92030096,-43.21647528",
        "-22.91966809,-43.21484186",
        "-22.91857924,-43.21283215",
        "-22.91834687,-43.21246911",
        "-22.91833423,-43.21242836",
        "-22.91804398,-43.2119079",
        "-22.91773052,-43.21145367"];

      const people = [
        5,
        5,
        5,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        4,
        5,
        5,
        5,
        5,
        5,
        5,
        12,
        12,
        12,
        21,
        21,
        21,
        50,
        50,
        50,
        5,
        5,
        5,
        5,
      ]

      let colors = [];
      
      _.each(people, (p, index) => {
        colors[index] = this.getColorFromValue(p);
      });

      const data = await axios.get('https://roads.googleapis.com/v1/snapToRoads', {
        params: {
          interpolate: true,
          key: "AIzaSyAnY70pmN83mIyNnWRAxqeydbKsuP5ehIE",
          path: pathValues.join('|')
        }
      });

      const snappedPoints = data.data.snappedPoints;

      let currColor = colors[0];
      let paths = [];
      let currCoordinates = [];

      _.each(snappedPoints, point => {
        let snappedCoordinate = new google.maps.LatLng(point.location.latitude, point.location.longitude);

        if (point.originalIndex && point.originalIndex > 0) {
        
          console.log(point.originalIndex, currCoordinates.length);
          paths.push({id: point.originalIndex - 1, color: currColor, snappedCoordinates: currCoordinates});

          currColor = colors[point.originalIndex];

          let lastCoord = currCoordinates[currCoordinates.length-1];
          currCoordinates = [lastCoord];
        }

        currCoordinates.push(snappedCoordinate);

      });

      if (currCoordinates.length) {
        paths.push({id: paths.length+1, color: currColor, snappedCoordinates: currCoordinates});
      }

      console.log(paths);
      this.paths = paths;
      this.center = paths[0].snappedCoordinates[0];
    },

    getColorFromValue(n) {
      if (n < 10) {
        return 'green';
      } else if (n < 20) {
        return 'yellow';
      } else if (n < 30) {
        return 'orange';
      } 

      return 'red';
    }
  },
  components: {
    Heatmap
  },
};
</script>

<style scoped>
.map, .container {
  width: 100%;
  height: 100%;
}
.container.fluid, .container {
  padding: 0 !important;
}
</style>
