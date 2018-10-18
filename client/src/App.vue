<template>
<div id="app">
  <v-app>
    <Sidebar/>
    <Content/>
  </v-app>
</div>
</template>

<script>
import Sidebar from '@/components/layout/Sidebar';
import Content from '@/components/layout/Content';

export default {
  name: 'App',
  components: {
    Content,
    Sidebar,
  },
};
</script>

<style>
html {
  overflow: auto;
}
</style>

