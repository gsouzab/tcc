<template>
<div id="app">
  <v-app>
    <Header/>
    <Content/>
  </v-app>
</div>
</template>

<script>
import Content from '@/components/layout/Content';
import Header from '@/components/layout/Header';

export default {
  name: 'App',
  components: {
    Header,
    Content,
  },
};
</script>
