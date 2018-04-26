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
          v-validate="{ required: true, regex: /^[a-fA-F0-9:]{17}|[a-fA-F0-9]{12}$/ }"
          data-vv-name="mac"
          required>
        </v-text-field>
        </v-flex>
        <v-flex xs6>
          <v-text-field
            prepend-icon="add_location"
            v-model="latitude"
            label="Latitude"
            :rules="[() => latitude.length > 0 || 'Este campo é obrigatório']"
            required
          ></v-text-field>
        </v-flex>
        <v-flex xs6>
          <v-text-field
            v-model="longitude"
            label="Longitude"
            :rules="[() => longitude.length > 0 || 'Este campo é obrigatório']"
            required
          ></v-text-field>
        </v-flex>
        <v-flex xs12>
          <v-text-field
            prepend-icon="notes"
            v-model="description"
            label="Descrição"
            :rules="[() => description.length > 0 || 'Este campo é obrigatório']"
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
export default {
  props: ["visible"],
  data: function() {
    return {
      name: '',
      mac: '',
      latitude: '',
      longitude: '',
      description: '',
    }
  },
  computed: {
    show: {
      get() {
        return this.visible;
      },
      set(value) {
        if (!value) {
          // this.clearData();
          this.$emit("close");
        }
      }
    }
  },
  methods: {
    saveSensor: function() {
      let formData = {
        name:         this.name,
        mac:          this.mac,
        latitude:     this.latitude,
        longitude:    this.longitude,
        description:  this.description
      };

      this.$emit("onSave", formData);
      this.clearData();
      this.show = false;
    },
    clearData: function() {
      this.name         = '';
      this.latitude     = '';
      this.mac          = '';
      this.longitude    = '';
      this.description  = '';
    }
  }
};
</script>
