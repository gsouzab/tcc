<template>
  <v-dialog v-model="show">
    <v-card>
      <v-card-title
        class="grey lighten-4 py-4 title">
        Adicionar Sensor
      </v-card-title>
      <v-container grid-list-sm class="pa-4">
        <v-layout row wrap>
          <v-flex xs12>
          <v-text-field
            prepend-icon="settings_input_antenna"
            v-model="name"
            label="Nome"
            :error-messages="errors.collect('name')"
            v-validate="'required'"
            data-vv-name="name"
            required>
          </v-text-field>
          </v-flex>
          <v-flex xs12>
          <v-text-field
            prepend-icon="settings_input_antenna"
            v-model="mac"
            label="MAC"
            :error-messages="errors.collect('mac')"
            mask="NN:NN:NN:NN:NN:NN"
            v-validate="{ required: true, regex: /^[a-fA-F0-9:]{17}|[a-fA-F0-9]{12}$/ }"
            data-vv-name="mac"
            required>
          </v-text-field>
          </v-flex>
          <v-flex xs6>
            <v-text-field
              prepend-icon="add_location"
              v-model="latitude"
              data-vv-name="latitude"
              label="Latitude"
              :error-messages="errors.collect('latitude')"
              v-validate="'required'"
              required
            ></v-text-field>
          </v-flex>
          <v-flex xs6>
            <v-text-field
              v-model="longitude"
              data-vv-name="longitude"
              label="Longitude"
              :error-messages="errors.collect('longitude')"
              v-validate="'required'"
              required
            ></v-text-field>
          </v-flex>
          <v-flex xs12>
            <v-text-field
              prepend-icon="notes"
              v-model="description"
              data-vv-name="description"
              label="Descrição"
              :error-messages="errors.collect('description')"
              v-validate="'required'"
              multi-line
              required
            ></v-text-field>
          </v-flex>
        </v-layout>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn flat @click="show = false">Cancelar</v-btn>
        <v-btn flat color="primary" @click="saveSensor()">Salvar</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

import axios from 'axios';

export default {
  props: ["visible", "initLat", "initLng"],
  watch: {
    initLat(lat, oldVal) { // watch it
      this.latitude = lat;
    },
    initLng(lng, oldVal) { // watch it
      this.longitude = lng;
    }
  },
  data() {
    return {
      name: '',
      mac: '',
      latitude: '',
      longitude: '',
      description: '',
    }
  },
  events: {
    setLatLng(lat, lng) {
      this.latitude = lat;
      this.longitude = lon;
    }
  },
  computed: {
    show: {
      get() {
        return this.visible;
      },
      set(value) {
        if (!value) {
          this.clearData();
          this.$emit("close");
        }
      }
    }
  },
  methods: {
    async saveSensor() {

      let valid = await this.$validator.validateAll();

      if (valid) {
        let formData = {
          name:         this.name,
          mac:          this.mac,
          latitude:     this.latitude.toString(),
          longitude:    this.longitude.toString(),
          description:  this.description
        };

        let response = await axios.post(`http://${window.location.hostname}:8000/sensors`, formData);
        this.$emit("onSave", formData);
        this.show = false;
      }

    },
    clearData() {
      this.name         = '';
      this.latitude     = '';
      this.mac          = '';
      this.longitude    = '';
      this.description  = '';
      this.errors.clear();
    }
  }
};
</script>
