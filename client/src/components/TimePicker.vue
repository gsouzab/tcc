<template>
  <v-expansion-panel v-model="panel" expand>
    <v-expansion-panel-content>
      <div slot="header">
        <v-layout row align-center justify-space-between>
          <div>
            <v-icon>schedule</v-icon> <span class="pl-2">{{selectedTimeText}}</span>
          </div>

          <div>
            <v-btn small depressed v-on:click.stop="setRealTime" color="primary" :disabled="startDate === null && endDate === null"><v-icon style="color: #4caf50 !important;" small v-if="startDate === null && endDate === null">fiber_manual_record</v-icon>Tempo real</v-btn>
          </div>
        </v-layout>
      </div>
      <v-card>
        <v-card-text>
          <v-layout row wrap>
            <v-flex>
              <v-list>
                <v-subheader>Acesso rápido</v-subheader>
                <v-list-tile
                  v-for="(preset, index) in presets1"
                  :key="index"
                  @click="onPresetSelect(preset)"
                >
                  <v-list-tile-content>
                    {{ preset.title }}
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
            </v-flex>

            <v-flex>
              <v-list>
                <v-subheader></v-subheader>
                <v-list-tile
                  v-for="(preset, index) in presets2"
                  :key="index"
                  @click="onPresetSelect(preset)"
                >
                  <v-list-tile-content>
                    {{ preset.title }}
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
            </v-flex>

            <v-flex>
              <v-list>
                <v-subheader></v-subheader>
                <v-list-tile
                  v-for="(preset, index) in presets3"
                  :key="index"
                  @click="onPresetSelect(preset)"
                >
                  <v-list-tile-content>
                    {{ preset.title }}
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
            </v-flex>

            <v-flex>
              <v-text-field
                v-model="startDateTime"
                :label="'Data inicial'"
                name="startDateTime"
                class="date-textfield"
                prepend-icon="event"
                @change="onDateRangeChange"
              />
              <v-date-picker v-model="startDate" no-title reactive locale="pt-br" :max="endDate" @change="onDateRangeChange"></v-date-picker>
            </v-flex>

            <v-flex>
              <v-text-field
                v-model="endDateTime"
                :label="'Data final'"
                name="endDateTime"
                class="date-textfield"
                @change="onDateRangeChange"
              />
              <v-date-picker v-model="endDate" no-title reactive locale="pt-br" :min="startDate" @change="onDateRangeChange"></v-date-picker>
            </v-flex>
            <v-btn color="success" large outline style="width: 44px; min-width: 44px;" @click="onDateRangeSelect"><v-icon>search</v-icon></v-btn>
          </v-layout>

        </v-card-text>
      </v-card>
    </v-expansion-panel-content>
  </v-expansion-panel>
</template>


<script>
import * as moment from "moment";

export default {
  name: "TimePicker",
  data() {
    return {
      selectedTimeText: "Agora",
      startDate: null,
      endDate: null,
      startDateTime: null,
      endDateTime: null,
      panel: [false],
      presets1: [
        {
          title: "Hoje",
          startDate: () => moment().startOf("day"),
          endDate: () => moment().endOf("day")
        },
        {
          title: "Essa semana",
          startDate: () => moment().startOf("week"),
          endDate: () => moment().endOf("week")
        },
        {
          title: "Esse mês",
          startDate: () => moment().startOf("month"),
          endDate: () => moment().endOf("month")
        },
        {
          title: "Esse ano",
          startDate: () => moment().startOf("year"),
          endDate: () => moment().endOf("year")
        },
        {
          title: "Desde o início da semana",
          startDate: () => moment().startOf("week"),
          endDate: () => moment()
        },
        {
          title: "Desde o início do mês",
          startDate: () => moment().startOf("month"),
          endDate: () => moment()
        }
      ],
      presets2: [
        {
          title: "Ontem",
          startDate: () =>
            moment()
              .subtract(1, "days")
              .startOf("day"),
          endDate: () =>
            moment()
              .subtract(1, "days")
              .endOf("day")
        },
        {
          title: "Semana passada",
          startDate: () =>
            moment()
              .subtract(1, "week")
              .startOf("week"),
          endDate: () =>
            moment()
              .subtract(1, "week")
              .endOf("week")
        },
        {
          title: "Mês passado",
          startDate: () =>
            moment()
              .subtract(1, "month")
              .startOf("month"),
          endDate: () =>
            moment()
              .subtract(1, "month")
              .endOf("month")
        },
        {
          title: "Ano passado",
          startDate: () =>
            moment()
              .subtract(1, "year")
              .startOf("year"),
          endDate: () =>
            moment()
              .subtract(1, "year")
              .endOf("year")
        },
        {
          title: "Últimos 7 dias",
          startDate: () =>
            moment()
              .subtract(7, "days")
              .startOf("day"),
          endDate: () => moment().endOf("day")
        },
        {
          title: "Últimos 30 dias",
          startDate: () =>
            moment()
              .subtract(30, "days")
              .startOf("day"),
          endDate: () => moment().endOf("day")
        }
      ],
      presets3: [
        {
          title: "Últimos 15 minutos",
          startDate: () => moment().subtract(15, "minutes"),
          endDate: () => moment()
        },
        {
          title: "Últimos 30 minutos",
          startDate: () => moment().subtract(30, "minutes"),
          endDate: () => moment()
        },
        {
          title: "Última hora",
          startDate: () => moment().subtract(1, "hour"),
          endDate: () => moment()
        },
        {
          title: "Últimas 4 horas",
          startDate: () => moment().subtract(4, "hours"),
          endDate: () => moment()
        },
        {
          title: "Últimas 12 horas",
          startDate: () => moment().subtract(12, "hours"),
          endDate: () => moment()
        },
        {
          title: "Últimas 24 horas",
          startDate: () => moment().subtract(24, "hours"),
          endDate: () => moment()
        }
      ]
    };
  },
  methods: {
    onPresetSelect(preset) {
      this.selectedTimeText = preset.title;
      this.startDateTime = preset.startDate().format("YYYY-MM-DD\\THH:mm:ss-02:00");
      this.endDateTime = preset.endDate().format("YYYY-MM-DD\\THH:mm:ss-02:00");

      this.startDate = preset.startDate().format("YYYY-MM-DD");
      this.endDate = preset.endDate().format("YYYY-MM-DD");

      this.onDateRangeSelect();
    },
    onDateRangeChange(data) {
      this.startDateTime = this.startDate ? this.startDate + "T00:00:00-02:00" : null;
      this.endDateTime = this.endDate ? this.endDate + "T23:59:59-02:00" : null;

      if (this.endDate === null && this.startDate === null) {
        this.selectedTimeText = 'Agora';
      } else {
        this.selectedTimeText = `${this.startDateTime} à ${this.endDateTime}`;
      }
    },
    onDateRangeSelect() {
      this.$nextTick(() => {
        this.panel = [false];
      });

      this.$emit('change', this.startDateTime, this.endDateTime)
    },
    setRealTime() {
      this.startDate = null;
      this.endDate = null;
      this.onDateRangeChange();
      this.onDateRangeSelect();
    }
  }
};
</script>

<style scoped>
.date-textfield {
  width: 290px;
}
</style>
