<template>
  <div>
    <v-container fluid class="pa-0">
      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-title class="py-3">
              <div class="d-flex align-center">
                <v-icon large class="mr-3" color="primary">mdi-monitor-dashboard</v-icon>
                <div>
                  <h1 class="text-h4 font-weight-bold">Terminal Dashboard</h1>
                  <p class="text-subtitle-1 grey--text mb-0">
                    Real-time market data and analytics
                  </p>
                </div>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>

      <!-- XRP Heatmap Section -->
      <v-row no-gutters class="mb-4">
        <v-col cols="12" md="6" class="pa-1">
          <v-card tile outlined>
            <v-card-title class="py-3">
              <span class="text-h6">XRP Token Heatmap</span>
              <v-spacer />
              <v-btn
                text
                color="primary"
                to="/xrp-heatmap"
              >
                View Full
                <v-icon right>mdi-arrow-right</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text class="pa-0">
              <xrp-token-heatmap :height="400" :user-can-access-trend="true" />
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" md="6" class="pa-1">
          <v-card tile outlined>
            <v-card-title class="py-3">
              <span class="text-h6">XRP AMM Heatmap</span>
              <v-spacer />
              <v-btn
                text
                color="primary"
                to="/xrp-amm-heatmap"
              >
                View Full
                <v-icon right>mdi-arrow-right</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text class="pa-0">
              <xrp-amm-heatmap :height="400" :user-can-access-trend="true" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Quick Stats Section -->
      <v-row no-gutters class="mb-4">
        <v-col cols="12" md="3" class="pa-1">
          <v-card tile outlined class="text-center pa-4">
            <v-icon large color="primary" class="mb-2">mdi-currency-usd</v-icon>
            <div class="text-h6 font-weight-bold">$2.5B</div>
            <div class="text-caption grey--text">Total Market Cap</div>
          </v-card>
        </v-col>
        <v-col cols="12" md="3" class="pa-1">
          <v-card tile outlined class="text-center pa-4">
            <v-icon large color="success" class="mb-2">mdi-trending-up</v-icon>
            <div class="text-h6 font-weight-bold text-success">+5.2%</div>
            <div class="text-caption grey--text">24h Change</div>
          </v-card>
        </v-col>
        <v-col cols="12" md="3" class="pa-1">
          <v-card tile outlined class="text-center pa-4">
            <v-icon large color="info" class="mb-2">mdi-chart-line</v-icon>
            <div class="text-h6 font-weight-bold">$180M</div>
            <div class="text-caption grey--text">24h Volume</div>
          </v-card>
        </v-col>
        <v-col cols="12" md="3" class="pa-1">
          <v-card tile outlined class="text-center pa-4">
            <v-icon large color="warning" class="mb-2">mdi-fire</v-icon>
            <div class="text-h6 font-weight-bold">1,247</div>
            <div class="text-caption grey--text">Active Tokens</div>
          </v-card>
        </v-col>
      </v-row>

      <!-- Quick Actions Section -->
      <v-row no-gutters>
        <v-col cols="12" md="6" class="pa-1">
          <v-card tile outlined class="pa-4">
            <v-card-title class="text-h6 pb-2">Quick Actions</v-card-title>
            <v-row>
              <v-col cols="6">
                <v-btn
                  block
                  color="primary"
                  to="/xrp-screener"
                  class="mb-2"
                >
                  <v-icon left>mdi-view-list</v-icon>
                  XRP Screener
                </v-btn>
              </v-col>
              <v-col cols="6">
                <v-btn
                  block
                  color="secondary"
                  to="/swap"
                  class="mb-2"
                >
                  <v-icon left>mdi-swap-horizontal</v-icon>
                  Swap Tokens
                </v-btn>
              </v-col>
              <v-col cols="6">
                <v-btn
                  block
                  color="info"
                  to="/xrp-trust-lines"
                  class="mb-2"
                >
                  <v-icon left>mdi-handshake</v-icon>
                  Trust Lines
                </v-btn>
              </v-col>
              <v-col cols="6">
                <v-btn
                  block
                  color="success"
                  to="/xrp-token-mints"
                  class="mb-2"
                >
                  <v-icon left>mdi-cube-scan</v-icon>
                  Token Mints
                </v-btn>
              </v-col>
            </v-row>
          </v-card>
        </v-col>
        <v-col cols="12" md="6" class="pa-1">
          <v-card tile outlined class="pa-4">
            <v-card-title class="text-h6 pb-2">Recent Activity</v-card-title>
            <v-list dense>
              <v-list-item v-for="(activity, index) in recentActivities" :key="index">
                <v-list-item-icon>
                  <v-icon :color="activity.color">{{ activity.icon }}</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title class="text-body-2">{{ activity.title }}</v-list-item-title>
                  <v-list-item-subtitle class="text-caption">{{ activity.time }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'
import XrpTokenHeatmap from '~/components/xrp/XrpTokenHeatmap.vue'
import XrpAmmHeatmap from '~/components/xrp/XrpAmmHeatmap.vue'

export default defineComponent({
  components: {
    XrpTokenHeatmap,
    XrpAmmHeatmap,
  },
  setup() {
    const recentActivities = ref([
      {
        title: 'New token minted: SOL',
        time: '2 minutes ago',
        icon: 'mdi-plus-circle',
        color: 'success'
      },
      {
        title: 'Large swap: 10,000 USDC → XRP',
        time: '5 minutes ago',
        icon: 'mdi-swap-horizontal',
        color: 'primary'
      },
      {
        title: 'Trust line created for USDT',
        time: '12 minutes ago',
        icon: 'mdi-handshake',
        color: 'info'
      },
      {
        title: 'Market cap increased by 2.3%',
        time: '15 minutes ago',
        icon: 'mdi-trending-up',
        color: 'success'
      }
    ])

    return {
      recentActivities,
    }
  },
  head() {
    return {
      title: 'Terminal Dashboard - EVM Finance',
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: 'Real-time market data, XRP token heatmap, and analytics dashboard for the XRP Ledger.',
        },
      ],
    }
  },
})
</script> 