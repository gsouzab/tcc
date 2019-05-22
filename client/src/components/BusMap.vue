<template>
  <v-container fluid>
    <gmap-map
      :center="center"
      :zoom="zoom"
      :options="{fullscreenControl: false, streetViewControl: false, mapTypeControl: false}"
      slot="activator"
      class="map">


      <gmap-polyline :key="path.id" v-for="path in paths" :options="{'strokeColor': path.color, 'strokeWeight': 6}" :path="path.snappedCoordinates" ></gmap-polyline>
    </gmap-map>

    <div>
      <svg id="legend" height=150 width=100></svg>
    </div>
  </v-container>
</template>

<script>

import axios from 'axios';
import {loaded, gmapApi} from 'vue2-google-maps';
import Heatmap from '@/components/Heatmap';
import _ from 'lodash';
import * as d3 from 'd3';
import { csv } from 'd3-fetch';
// import url from '../assets/ida_real.csv';

export default {
  async created() {
    await this.runSnapToRoad();
  },
  mounted() {
    this.setLegend();
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
    setLegend() {
      var svg = d3.select("#legend")
        .style("position", "absolute")
        .style("right", "10px")
        // .style("background-color", "white")
        .style("top", "0px");

      // Handmade legend
      svg.append("circle").attr("cx",10).attr("cy",50).attr("r", 6).style("fill", "green")
      svg.append("circle").attr("cx",10).attr("cy",70).attr("r", 6).style("fill", "yellow")
      svg.append("circle").attr("cx",10).attr("cy",90).attr("r", 6).style("fill", "orange")
      svg.append("circle").attr("cx",10).attr("cy",110).attr("r", 6).style("fill", "red")
      svg.append("text").attr("x", 30).attr("y", 50).text("0-9").style("font-size", "14px").attr("alignment-baseline","middle")
      svg.append("text").attr("x", 30).attr("y", 70).text("10-19").style("font-size", "14px").attr("alignment-baseline","middle")
      svg.append("text").attr("x", 30).attr("y", 90).text("20-39").style("font-size", "14px").attr("alignment-baseline","middle")
      svg.append("text").attr("x", 30).attr("y", 110).text(">40").style("font-size", "14px").attr("alignment-baseline","middle")
    },
    async runSnapToRoad(path) {
      let colors = [];
      
      let gpsData = await csv('static/entrada_180_saida_360.csv');
      delete gpsData['columns'];
      
      colors[0] = this.getColorFromValue(gpsData[0].n_sherlock);
      
      var pathQuery = `${gpsData[0].lat},${gpsData[0].lon}`;
      var snappedPoints = [];

      for (let i = 1; i < gpsData.length; i++) {
        colors[i] = this.getColorFromValue(gpsData[i].n_sherlock);
      
        if (pathQuery === "") {
          pathQuery += `${gpsData[i].lat},${gpsData[i].lon}`;
        } else {
          pathQuery += `|${gpsData[i].lat},${gpsData[i].lon}`;
        }

        if (i % 99 === 0) {

          const data = await axios.get('https://roads.googleapis.com/v1/snapToRoads', {
            params: {
              interpolate: true,
              key: "AIzaSyAnY70pmN83mIyNnWRAxqeydbKsuP5ehIE",
              path: pathQuery
            }
          });
          pathQuery = "";
          snappedPoints = _.concat(snappedPoints, data.data.snappedPoints);
        }
      }

      if (pathQuery !== "") {
        const data = await axios.get('https://roads.googleapis.com/v1/snapToRoads', {
          params: {
            interpolate: true,
            key: "AIzaSyAnY70pmN83mIyNnWRAxqeydbKsuP5ehIE",
            path: pathQuery
          }
        });

        snappedPoints = _.concat(snappedPoints, data.data.snappedPoints);
      }

      let currColor = colors[0];
      let paths = [];
      let currCoordinates = [];
      var index = 0;
      _.each(snappedPoints, point => {
        let snappedCoordinate = new google.maps.LatLng(point.location.latitude, point.location.longitude);

        if (point.originalIndex && point.originalIndex > 0) {

          paths.push({id: index, color: currColor, snappedCoordinates: currCoordinates});
          currColor = colors[index];
          index++;
          let lastCoord = currCoordinates[currCoordinates.length-1];
          currCoordinates = [lastCoord];
        }

        currCoordinates.push(snappedCoordinate);

      });

      if (currCoordinates.length) {
        paths.push({id: paths.length+1, color: currColor, snappedCoordinates: currCoordinates});
      }

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

<style>
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
