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
          placeholder="Nome">
        </v-text-field>
        </v-flex>
        <v-flex xs6>
          <v-text-field
            prepend-icon="add_location"
            placeholder="Latitude"
          ></v-text-field>
        </v-flex>
        <v-flex xs6>
          <v-text-field
            placeholder="Longitude"
          ></v-text-field>
        </v-flex>
        <v-flex xs12>
          <v-text-field
            prepend-icon="notes"
            placeholder="Descrição"
          ></v-text-field>
        </v-flex>
      </v-layout>
    </v-container>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn flat @click="show = false">Cancelar</v-btn>
      <v-btn flat color="primary" @click="this.saveSensor">Salvar</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
</template>

<script>
export default {
  props: ["visible"],
  computed: {
    show: {
      get() {
        return this.visible;
      },
      set(value) {
        if (!value) {
          this.$emit("close");
        }
      }
    }
  },
  methods: {
    saveSensor: () => {
      this.show = false;
    }
  }
};
</script>
