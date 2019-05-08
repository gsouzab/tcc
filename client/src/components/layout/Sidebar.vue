<template>
  <v-navigation-drawer
    v-model="drawer"
    mini-variant
    dark
    permanent
    :mini-variant-width="60"
    app>

    <v-toolbar flat>
      <v-list>
        <v-list-tile>
          <v-list-tile-title class="title">

          </v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-toolbar>

    <v-divider></v-divider>

    <v-list dense>

      <v-list-tile
        v-for="item in items"
        :key="item.title"
        @click="$router.push({name: item.route})"
      >
        <v-list-tile-action>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-tile-action>

      </v-list-tile>
    </v-list>

  </v-navigation-drawer>

</template>

<script>
export default {
  name: 'Sidebar',
  data () {
    return {
      drawer: true,
      items: [
        { title: 'Bus', icon: 'map', route: 'map' },
        { title: 'Map', icon: 'directions_bus', route: 'bus' },
        { title: 'Dashboard', icon: 'dashboard', route: 'dashboard' }
      ]
    }
  }
};
</script>

<style>
.v-navigation-drawer--mini-variant .v-list__tile__action, .v-navigation-drawer--mini-variant .v-list__tile__avatar {
  min-width: 28px;
}
</style>
