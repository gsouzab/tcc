<template>
<v-toolbar
  color="grey darken-3"
  dark
  app
  :clipped-left="$vuetify.breakpoint.lgAndUp"
  fixed >

  <v-text-field
    flat
    solo-inverted
    prepend-icon="search"
    label="Search"
    class="hidden-sm-and-down">
  </v-text-field>
  <v-spacer></v-spacer>
  <v-btn icon>
    <v-icon>notifications</v-icon>
  </v-btn>
</v-toolbar>
</template>

<script>
export default {
  name: 'Header',
};
</script>
