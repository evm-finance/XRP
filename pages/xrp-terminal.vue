<template>
  <div>
    <v-container fluid>
      <v-row>
        <v-col>
          <v-card>
            <v-card-title>XRP Terminal</v-card-title>
            <v-card-text>
              <v-alert type="info">
                XRP Terminal component is currently disabled. Please check back later.
              </v-alert>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  name: 'XrpTerminalPage'
})
</script> 