<template>
  <v-container fluid id="bus-map">
    <gmap-map
      :center="center"
      :zoom="zoom"
      :options="{fullscreenControl: false, streetViewControl: false, mapTypeControl: false}"
      slot="activator"
      class="map"
    >
      <div slot="visible">
        <v-toolbar dense floating id="selector-nav" class="mt-3">
          <v-overflow-btn
            :items="dropdown_experiments"
            label="Data experimento"
            class="pa-0"
            hide-details
            v-model="experiment"
            @change="runSnapToRoad"
            target="#bus-map"
          ></v-overflow-btn>

          <v-divider class="mr-2" vertical></v-divider>

          <v-btn-toggle v-model="experiment_type" class="mr-2" @change="runSnapToRoad">
            <v-btn flat value="ground_truth">Real</v-btn>
            <v-btn flat value="estimated">Estimado</v-btn>
          </v-btn-toggle>
        </v-toolbar>
      </div>

      <gmap-polyline
        :key="path.id"
        v-for="path in paths"
        :options="{'strokeColor': path.color, 'strokeWeight': 6}"
        :path="path.snappedCoordinates"
        @mouseover="showLegend($event, path)"
      ></gmap-polyline>

      <gmap-info-window
        :options="infoOptions"
        :position="infoWindowPos"
        :opened="infoWinOpen"
        @closeclick="infoWinOpen=false"
      >
        <div v-html="infoContent"></div>
      </gmap-info-window>

      <GmapMarker
        :key="index"
        v-for="(m, index) in markers"
        :position="m.position"
        :clickable="false"
        :icon="m.icon"
        :label="m.label"
        :draggable="false"
      />

      <div slot="visible">
        <v-btn
          absolute
          dark
          fab
          bottom
          right
          class="add-button"
          color="green"
          @click="showUploadForm=true;"
        >
          <v-icon>add</v-icon>
        </v-btn>
      </div>
    </gmap-map>

    <div>
      <svg id="legend" height="150" width="100" />
    </div>
  </v-container>
</template>

<script>
import axios from "axios";
import { loaded, gmapApi } from "vue2-google-maps";
import Heatmap from "@/components/Heatmap";
import UploadForm from "@/components/forms/UploadForm";
import _ from "lodash";
import * as d3 from "d3";
import { csv } from "d3-fetch";
// import url from '../assets/ida_real.csv';

export default {
  async created() {
    // await this.runSnapToRoad();
  },
  mounted() {
    this.setLegend();
  },
  data() {
    return {
      showUploadForm: false,
      infoContent: "",
      infoWindowPos: null,
      infoWinOpen: false,
      infoOptions: {
      },
      center: { lat: -22.861659764789806, lng: -43.22840063788988 },
      zoom: 13,
      experiment: "Volta 29/03/2019",
      experiment_type: "ground_truth",
      paths: [],
      dropdown_experiments: ["Volta 29/03/2019"],
      markers: [],
      files: {
        "Volta 29/03/2019": {
          estimated: "static/entrada_180_saida_360.csv",
          ground_truth: "static/volta_real.csv"
        }
      }
    };
  },
  computed: {
    google: gmapApi
  },
  methods: {
    showLegend(e, path) {
      let position = { lat: e.latLng.lat(), lng: e.latLng.lng() };

      this.$nextTick(() => {
        // this.center = position;
        this.infoContent = `Quantidade estimada: ${path.n_sherlock}
        <br>Quantidade real: 12
        <br>Velocidade: 26 km/h
        <br>Data e hora: ${new Date(path.timestamp).toLocaleString('pt-BR')}`;
        this.infoWindowPos = position;
        this.infoWinOpen = true;
      })
      // console.log(path.timestamp, path.n_sherlock);
      // this.infoWinOpen = position;
      // this.infoWinOpen = true;
    },
    setLegend() {
      var svg = d3
        .select("#legend")
        .style("position", "absolute")
        .style("right", "10px")
        .style("background-color", "white")
        .style("border-radius", "6px")
        .style("height", "80px")
        .style("width", "70px")
        .style("top", "16px");

      // Handmade legend
      svg
        .append("circle")
        .attr("cx", 10)
        .attr("cy", 10)
        .attr("r", 6)
        .style("fill", "green");
      svg
        .append("circle")
        .attr("cx", 10)
        .attr("cy", 30)
        .attr("r", 6)
        .style("fill", "yellow");
      svg
        .append("circle")
        .attr("cx", 10)
        .attr("cy", 50)
        .attr("r", 6)
        .style("fill", "orange");
      svg
        .append("circle")
        .attr("cx", 10)
        .attr("cy", 70)
        .attr("r", 6)
        .style("fill", "red");
      svg
        .append("text")
        .attr("x", 25)
        .attr("y", 10)
        .text("0-9")
        .style("font-size", "14px")
        .attr("alignment-baseline", "middle");
      svg
        .append("text")
        .attr("x", 25)
        .attr("y", 30)
        .text("10-19")
        .style("font-size", "14px")
        .attr("alignment-baseline", "middle");
      svg
        .append("text")
        .attr("x", 25)
        .attr("y", 50)
        .text("20-39")
        .style("font-size", "14px")
        .attr("alignment-baseline", "middle");
      svg
        .append("text")
        .attr("x", 25)
        .attr("y", 70)
        .text(">40")
        .style("font-size", "14px")
        .attr("alignment-baseline", "middle");
    },
    clearMarkers() {
      this.markers = [];
    },
    addEndMarker(position) {
      this.markers.push({
        position,
        icon: {
          url: "http://maps.google.com/mapfiles/ms/icons/red-dot.png",
          labelOrigin: new google.maps.Point(15, -10)
        },
        label: "Fim"
      });
    },
    addStartMarker(position) {
      this.markers.push({
        position,
        icon: {
          url: "http://maps.google.com/mapfiles/ms/icons/green-dot.png",
          labelOrigin: new google.maps.Point(15, -10)
        },
        label: "Início"
      });
    },
    async runSnapToRoad() {
      let file = this.files[this.experiment][this.experiment_type];

      if (!file) {
        return;
      }

      let colors = [];

      let gpsData = await csv(file);
      delete gpsData["columns"];

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
          const data = await axios.get(
            "https://roads.googleapis.com/v1/snapToRoads",
            {
              params: {
                interpolate: true,
                key: "AIzaSyAnY70pmN83mIyNnWRAxqeydbKsuP5ehIE",
                path: pathQuery
              }
            }
          );
          pathQuery = "";
          snappedPoints = _.concat(snappedPoints, data.data.snappedPoints);
        }
      }

      if (pathQuery !== "") {
        const data = await axios.get(
          "https://roads.googleapis.com/v1/snapToRoads",
          {
            params: {
              interpolate: true,
              key: "AIzaSyAnY70pmN83mIyNnWRAxqeydbKsuP5ehIE",
              path: pathQuery
            }
          }
        );

        snappedPoints = _.concat(snappedPoints, data.data.snappedPoints);
      }

      let currColor = colors[0];
      const paths = [];
      let currCoordinates = [];

      _.each(snappedPoints, point => {
        let snappedCoordinate = new google.maps.LatLng(
          point.location.latitude,
          point.location.longitude
        );

        if (point.originalIndex && point.originalIndex > 0) {
          let orgIndex = point.originalIndex;
          paths.push({
            id: orgIndex,
            color: currColor,
            snappedCoordinates: currCoordinates,
            n_sherlock: gpsData[orgIndex].n_sherlock,
            timestamp: gpsData[orgIndex].time
          });
          currColor = colors[orgIndex];
          let lastCoord = currCoordinates[currCoordinates.length - 1];
          currCoordinates = [lastCoord];
        }

        currCoordinates.push(snappedCoordinate);
      });

      if (currCoordinates.length) {
        paths.push({
          id: 11231231,
          color: currColor,
          snappedCoordinates: currCoordinates
        });
      }

      this.$nextTick(() => {
        this.paths = null;
        this.paths = paths;
        this.center = paths[0].snappedCoordinates[0];
      })

      this.clearMarkers();
      this.addStartMarker(paths[0].snappedCoordinates[0]);
      this.addEndMarker(paths[paths.length - 1].snappedCoordinates[0]);
    },

    getColorFromValue(n) {
      if (n < 10) {
        return "green";
      } else if (n < 20) {
        return "yellow";
      } else if (n < 30) {
        return "orange";
      }

      return "red";
    }
  },
  components: {
    Heatmap,
    UploadForm
  }
};
</script>

<style>
.v-toolbar__content {
  height: auto !important;
  padding: 0;
}
.v-select__slot {
  min-width: 220px;
}
</style>

<style scoped>
.map,
.container {
  width: 100%;
  height: 100%;
}
.container.fluid,
.container {
  padding: 0 !important;
}
.add-button {
  right: 80px;
  bottom: 36px;
}
</style>
