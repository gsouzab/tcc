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
.container.fluid {
  padding: 0;
}
</style>

