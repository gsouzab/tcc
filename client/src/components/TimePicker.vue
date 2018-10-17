<template>
  <v-expansion-panel>
    <v-expansion-panel-content>
      <div slot="header">
        <v-layout row align-center>
         <v-icon >schedule</v-icon> <span class="pl-2">{{selectedTimeText}}</span>
        </v-layout>
      </div>
      <v-card>
        <v-card-text>
          <v-layout row justify-end>
            <v-date-picker v-model="startDatePicker" no-title reactive locale="pt-br" :events="selectedStartRange"></v-date-picker>

            <v-date-picker v-model="endDatePicker" no-title reactive locale="pt-br" :events="selectedEndRange" :min="minEndDate"></v-date-picker>
          </v-layout>

        </v-card-text>
      </v-card>
    </v-expansion-panel-content>
  </v-expansion-panel>
</template>


<script>
export default {
  name: 'TimePicker',
  data() {
    return {
      selectedTimeText: 'últimos 15 minutos',
      startDatePicker: null,
      endDatePicker: null,
      selectedStartTime: null,
      selectedEndTime: null,
      minEndDate: null
    }
  },
  methods: {
    selectedStartRange(date) {
      const [,, day] = date.split('-')
      return parseInt(day, 10) % 3 === 0
    },
    selectedEndRange(date) {
      const [,, day] = date.split('-')
      return parseInt(day, 10) % 3 === 0
    },
  }
}
</script>
