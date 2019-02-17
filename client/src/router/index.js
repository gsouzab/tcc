import 'babel-polyfill';

import Vue from 'vue';
import Router from 'vue-router';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import * as VueGoogleMaps from 'vue2-google-maps';
import GmapCluster from 'vue2-google-maps/dist/components/cluster'; // replace src with dist if you have Babel issues
// import VueNativeSock from 'vue-native-websocket';

import Map from '@/components/Map';
import Dashboard from '@/components/Dashboard';

// Vue.use(VueNativeSock, `ws://${process.env.API_HOST}/ws/probe`, {
//   reconnection: true, // (Boolean) whether to reconnect automatically (false)
//   reconnectionAttempts: 5, // (Number) number of reconnection attempts before giving up (Infinity),
//   reconnectionDelay: 3000, // (Number) how long to initially wait before attempting a new (1000)
//   format: 'json',
// });

Vue.use(VeeValidate);
Vue.use(Router);
Vue.use(Vuetify);

const connectTelemetryWS = () => {
  console.log('CONNECTING TELEMERY SOCKET...');
  Vue.prototype.$telemetryWS = new WebSocket(`ws://${process.env.API_HOST}/ws/telemetry`);

  Vue.prototype.$telemetryWS.onclose = function(e) {
    console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
    setTimeout(function() {
      connectTelemetryWS();
    }, 1000);
  };

  Vue.prototype.$telemetryWS.onerror = function(err) {
    console.error('Socket encountered error: ', err.message, 'Closing socket');
    Vue.prototype.$telemetryWS.close();
  };
}

const connectProbesWS = () => {
  console.log('CONNECTING PROBE SOCKET...');
  Vue.prototype.$probesWS = new WebSocket(`ws://${process.env.API_HOST}/ws/probe`);

  Vue.prototype.$probesWS.onclose = function(e) {
    console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
    setTimeout(function() {
      connectProbesWS();
    }, 1000);
  };

  Vue.prototype.$probesWS.onerror = function(err) {
    console.error('Socket encountered error: ', err.message, 'Closing socket');
    Vue.prototype.$probesWS.close();
  };
}

connectTelemetryWS();
connectProbesWS();

Vue.use(VueGoogleMaps, {
  load: {
    key: process.env.MAPS_API_KEY,
    libraries: 'visualization'
  },
});

Vue.component('GmapCluster', GmapCluster);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'map',
      component: Map,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: Dashboard,
    },
  ],
});
