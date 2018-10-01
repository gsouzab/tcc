<template>
<div id="app">
  <v-app>
    <Sidebar/>
    <Header/>
    <Content/>
  </v-app>
</div>
</template>

<script>
import Sidebar from '@/components/layout/Sidebar';
import Content from '@/components/layout/Content';
import Header from '@/components/layout/Header';

export default {
  name: 'App',
  components: {
    Header,
    Content,
    Sidebar
  },
};
</script>

<style>
html {
  overflow: auto;
}
</style>

