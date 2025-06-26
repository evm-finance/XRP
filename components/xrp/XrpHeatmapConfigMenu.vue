<template>
  <div class="xrp-heatmap-config-menu">
    <v-menu v-model="menuOpen" offset-y>
      <template #activator="{ on, attrs }">
        <v-btn
          icon
          v-bind="attrs"
          v-on="on"
        >
          <v-icon>mdi-cog</v-icon>
        </v-btn>
      </template>
      
      <v-card>
        <v-card-title>XRP Heatmap Configuration</v-card-title>
        <v-card-text>
          <p class="caption grey--text">XRP heatmap configuration options will be added here</p>
        </v-card-text>
      </v-card>
    </v-menu>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'

export default defineComponent({
  name: 'XrpHeatmapConfigMenu',
  setup() {
    const menuOpen = ref(false)
    
    return {
      menuOpen
    }
  }
})
</script> 