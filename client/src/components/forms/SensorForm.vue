<template>
  <v-dialog v-model="show" persistent max-width="500px">
    <v-card>
      <v-card-title
        class="grey lighten-4 py-4 title">
        {{isEdit ? 'Editar' : 'Adicionar'}} Sensor
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
            :disabled="isEdit"
            :error-messages="errors.collect('mac')"
            mask="NN:NN:NN:NN:NN:NN"
            return-masked-value
            v-validate="{ required: true, regex: /^[a-fA-F0-9:]{17}|[a-fA-F0-9]{12}$/ }"
            data-vv-name="mac"
            required>
          </v-text-field>
          </v-flex>
          <v-flex xs12>
            <v-textarea
              prepend-icon="notes"
              v-model="description"
              data-vv-name="description"
              label="Descrição"
              :error-messages="errors.collect('description')"
            ></v-textarea>
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
  props: ['visible', 'initData', 'isEdit'],
  watch: {
    initData(val) { // watch it
      if (val) {
        this.name = val.name;
        this.mac = val.mac;
        this.latitude = val.latitude;
        this.longitude = val.longitude;
        this.description = val.description;
      }
    },
  },
  data() {
    return {
      name: '',
      mac: '',
      latitude: '',
      longitude: '',
      description: '',
    };
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

        let response = await axios.post(`https://${process.env.API_HOST}/sensors`, formData);
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
