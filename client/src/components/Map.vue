<template>
  <gmap-map
    :center="center"
    :zoom="15"
    style="width: 100%; height: calc(100% - 64px); position: absolute; left:0;">

    <gmap-marker
      :key="index"
      v-for="(m, index) in markers"
      :position="m.position"
      :clickable="true"
      @click="center=m.position">

    </gmap-marker>
    <v-speed-dial
      v-model="fab"
      :top="top"
      :bottom="bottom"
      :right="right"
      :left="left"
      :direction="direction"
      :open-on-hover="hover"
      :transition="transition"
    >
      <v-btn
        slot="activator"
        color="blue darken-2"
        dark
        fab
        hover
        v-model="fab"
      >
        <v-icon>account_circle</v-icon>
        <v-icon>close</v-icon>
      </v-btn>
      <v-btn
        fab
        dark
        small
        color="green"
      >
        <v-icon>edit</v-icon>
      </v-btn>
      <v-btn
        fab
        dark
        small
        color="indigo"
      >
        <v-icon>add</v-icon>
      </v-btn>
      <v-btn
        fab
        dark
        small
        color="red"
      >
        <v-icon>delete</v-icon>
      </v-btn>
    </v-speed-dial>
  </gmap-map>
</template>

<script>
import SensorForm from "@/components/forms/SensorForm";

export default {
  data() {
    return {
      fab: false,
      center: { lat: -22.8617784, lng: -43.2296038 },
      markers: [],
      sensors: [],
      direction: "top",
      fab: false,
      fling: false,
      hover: false,
      tabs: null,
      top: false,
      right: true,
      bottom: true,
      left: false,
      transition: "slide-y-reverse-transition"
    };
  },
  methods: {
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
  computed: {
    activeFab() {
      switch (this.tabs) {
        case "one":
          return { class: "purple", icon: "account_circle" };
        case "two":
          return { class: "red", icon: "edit" };
        case "three":
          return { class: "green", icon: "keyboard_arrow_up" };
        default:
          return {};
      }
    }
  },

  watch: {
    top(val) {
      this.bottom = !val;
    },
    right(val) {
      this.left = !val;
    },
    bottom(val) {
      this.top = !val;
    },
    left(val) {
      this.right = !val;
    }
  },
  components: {
    SensorForm
  }
};
</script>
