<template>
  <div :ref="chart.uuid"></div>
</template>

<script>
import Plotly from 'plotly.js-dist';

export default {
  props: ['chart', 'loading'],
  mounted() {
    Plotly.plot(
      this.$refs[this.chart.uuid],
      this.chart.traces,
      this.chart.layout,
      { responsive: true },
    );
  },
  watch: {
    chart: {
      handler() {
        Plotly.react(
          this.$refs[this.chart.uuid],
          this.chart.traces,
          this.chart.layout,
        );
      },
      deep: true,
    },
  },
};
</script>
