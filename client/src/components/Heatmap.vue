<script>
import {MapElementFactory} from 'vue2-google-maps'

export default MapElementFactory({
  name: 'heatmap',
  ctr: () => google.maps.visualization.HeatmapLayer,
  events: [],
  mappedProps: {
    data: {type: Array},
    options: {type: Object},
  }
})
</script>
