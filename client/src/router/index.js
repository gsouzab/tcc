import 'babel-polyfill';

import Vue from 'vue';
import Router from 'vue-router';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import Map from '@/components/Map';
import * as VueGoogleMaps from 'vue2-google-maps';

Vue.use(VeeValidate);
Vue.use(Router);
Vue.use(Vuetify);
Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyBUyTr9f-xkXsS4pMsww0ICCK1O-OR33sg',
  },
});

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Map',
      component: Map,
    },
  ],
});
