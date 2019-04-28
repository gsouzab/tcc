<script>
import {MapElementFactory} from 'vue2-google-maps'

export default MapElementFactory({
  name: 'groundOverlay',
  mappedProps: {
    'opacity': {}
  },
  props: {
    'source': {type: String},
    'bounds': {type: Object},
  },
  events: ['click', 'dblclick', 'rightclick'],
  ctr: () => google.maps.GroundOverlay,
  ctrArgs: (options, {source, bounds}) => [source, bounds, options],
})
</script>
