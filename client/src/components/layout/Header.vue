<template>
<v-toolbar
  color="grey darken-3"
  dark
  app
  compact
  :clipped-left="$vuetify.breakpoint.lgAndUp"
  fixed >


  <v-spacer></v-spacer>
  <v-text-field
    flat
    solo-inverted
    prepend-icon="search"
    label="Search"
    class="hidden-sm-and-down">
  </v-text-field>
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

<style >
.toolbar__content {
  height: 46px !important;
}

.input-group.input-group--solo {
  height: 32px ;
  min-height: 32px;
}

.input-group label {
  height: 16px;
  line-height: 16px;
}

.input-group.input-group--solo .input-group__input {
  padding: 2px 16px;
}
</style>
