import 'babel-polyfill';

import Vue from 'vue';
import Router from 'vue-router';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import * as VueGoogleMaps from 'vue2-google-maps';
import GmapCluster from 'vue2-google-maps/dist/components/cluster'; // replace src with dist if you have Babel issues
import VueNativeSock from 'vue-native-websocket';

import Map from '@/components/Map';
import Dashboard from '@/components/Dashboard';

Vue.use(VueNativeSock, `ws://${process.env.API_HOST}/ws`, {
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
  reconnectionAttempts: 5, // (Number) number of reconnection attempts before giving up (Infinity),
  reconnectionDelay: 3000, // (Number) how long to initially wait before attempting a new (1000)
  format: 'json',
});

Vue.use(VeeValidate);
Vue.use(Router);
Vue.use(Vuetify);

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
