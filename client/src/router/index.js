import 'babel-polyfill';

import Vue from 'vue';
import Router from 'vue-router';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import * as VueGoogleMaps from 'vue2-google-maps';
import GmapCluster from 'vue2-google-maps/dist/components/cluster'; // replace src with dist if you have Babel issues

import Map from '@/components/Map';
import Dashboard from '@/components/Dashboard';

Vue.use(VeeValidate);
Vue.use(Router);
Vue.use(Vuetify);

Vue.prototype.$connectTelemetryWS = () => {
  console.log('CONNECTING TELEMERY SOCKET...');
  Vue.prototype.$telemetryWS = new WebSocket(`wss://${process.env.API_HOST}/ws/telemetry`);

  Vue.prototype.$telemetryWS.onclose = (e) => {
    console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
    setTimeout(() => {
      Vue.prototype.$connectTelemetryWS();
    }, 1000);
  };

  Vue.prototype.$telemetryWS.onerror = (err) => {
    console.error('Socket encountered error: ', err.message, 'Closing socket');
    Vue.prototype.$telemetryWS.close();
  };
};

Vue.prototype.$connectProbesWS = () => {
  console.log('CONNECTING PROBE SOCKET...');
  Vue.prototype.$probesWS = new WebSocket(`wss://${process.env.API_HOST}/ws/probe`);

  Vue.prototype.$probesWS.onclose = (e) => {
    console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
    setTimeout(() => {
      Vue.prototype.$connectProbesWS();
    }, 1000);
  };

  Vue.prototype.$probesWS.onerror = (err) => {
    console.error('Socket encountered error: ', err.message, 'Closing socket');
    Vue.prototype.$probesWS.close();
  };
};

Vue.prototype.$connectTelemetryWS();
Vue.prototype.$connectProbesWS();

Vue.use(VueGoogleMaps, {
  load: {
    key: process.env.MAPS_API_KEY,
    libraries: 'visualization',
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
