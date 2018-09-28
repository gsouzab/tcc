<template>
<v-content>
  <v-container fluid>
    <router-view></router-view>
  </v-container>
</v-content>
</template>

<script>
export default {
  name: 'Content',
};
</script>

<style>
.content {
  padding: 46px 0px 0px 60px !important;
}
.container.fluid {
  padding: 0;
}
</style>

