<template>
  <v-navigation-drawer
    v-model="drawer"
    :clipped="$vuetify.breakpoint.lgAndUp"
    mini-variant
    dark
    app>

    <v-toolbar flat class="transparent">
      <v-list class="pt-0">
        <v-divider></v-divider>

        <v-list-tile
          v-for="item in items"
          :key="item.title"
          @click="$router.push({name: item.route})"
        >
            <v-list-tile-action>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-tile-action>

            <v-list-tile-content>
              <v-list-tile-title>{{ item.title }}</v-list-tile-title>
            </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-toolbar>

  </v-navigation-drawer>

</template>

<script>
export default {
  name: 'Sidebar',
  data () {
    return {
      drawer: true,
      items: [
        { title: 'Map', icon: 'map', route: 'map' },
        { title: 'Dashboard', icon: 'dashboard', route: 'dashboard' }
      ]
    }
  }
};
</script>

<style >
.navigation-drawer {
  width: 60px !important;
}
.navigation-drawer--mini-variant .list__tile__action, .navigation-drawer--mini-variant .list__tile__avatar {
  min-width: 26px !important;
}
</style>
