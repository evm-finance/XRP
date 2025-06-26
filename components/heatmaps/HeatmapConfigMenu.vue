<template>
  <div>
    <!-- XRP Heatmap Configuration Menu -->
    <v-menu offset-y>
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          color="primary"
          dark
          v-bind="attrs"
          v-on="on"
          icon
        >
          <v-icon>mdi-cog</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item @click="configureHeatmap">
          <v-list-item-title>Configure Heatmap</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  name: 'HeatmapConfigMenu',
  setup() {
    const configureHeatmap = () => {
      // XRP heatmap configuration logic
      console.log('Configure XRP heatmap')
    }

    return {
      configureHeatmap
    }
  }
})
</script> 