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
  </gmap-map>
</template>

<script>

  export default {
    data () {
      return {
        center: {lat: -22.8617784, lng: -43.2296038},
        markers: [],
      }
    }
  }
</script>
