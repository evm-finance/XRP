<template>
  <div>
    <v-row no-gutters justify="center">
      <v-col cols="12" md="10">
        <v-row justify="center">
          <v-col cols="12">
            <h1 class="text-h4">XRP Trust Lines</h1>
            <p class="text-body-2 grey--text mt-2">
              Manage your Trust lines to hold XRP tokens from different issuers
            </p>
          </v-col>
        </v-row>

        <!-- Wallet Connection Status -->
        <v-row justify="center" class="mb-4">
          <v-col cols="12" lg="6">
            <v-alert
              v-if="!isWalletReady"
              type="warning"
              outlined
            >
              <div class="d-flex align-center">
                <v-icon left>mdi-wallet</v-icon>
                <span>Please connect your GEM wallet to manage Trust lines</span>
                <v-btn
                  color="primary"
                  small
                  class="ml-auto"
                  @click="connectWallet"
                >
                  Connect Wallet
                </v-btn>
              </div>
            </v-alert>
            <v-alert
              v-else
              type="success"
              outlined
            >
              <div class="d-flex align-center">
                <v-icon left>mdi-check-circle</v-icon>
                <span>Connected: {{ formatAddress(address) }}</span>
                <v-btn
                  text
                  small
                  class="ml-auto"
                  @click="copyToClipboard(address)"
                >
                  <v-icon small>mdi-content-copy</v-icon>
                </v-btn>
              </div>
            </v-alert>
          </v-col>
        </v-row>

        <!-- Create New Trust Line -->
        <v-row justify="center" class="mb-6">
          <v-col cols="12" lg="8">
            <v-card tile outlined>
              <v-card-title class="d-flex align-center">
                <v-icon left>mdi-plus-circle</v-icon>
                Create New Trust Line
              </v-card-title>
              <v-card-text>
                <v-form ref="form" v-model="formValid">
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="newTrustLine.currency"
                        label="Currency Code"
                        placeholder="e.g., USDC, BTC, ETH"
                        :rules="[rules.required, rules.currency]"
                        outlined
                        dense
                        required
                      />
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="newTrustLine.issuer"
                        label="Issuer Address"
                        placeholder="rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
                        :rules="[rules.required, rules.address]"
                        outlined
                        dense
                        required
                      />
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="newTrustLine.limit"
                        label="Trust Limit"
                        placeholder="1000"
                        type="number"
                        :rules="[rules.required, rules.positive]"
                        outlined
                        dense
                        required
                      />
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="newTrustLine.flags"
                        label="Trust Line Flags"
                        :items="trustLineFlags"
                        item-text="label"
                        item-value="value"
                        outlined
                        dense
                        multiple
                        chips
                      />
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="12">
                      <v-btn
                        color="primary"
                        :disabled="!formValid || !isWalletReady || loading"
                        :loading="loading"
                        @click="createTrustLine"
                      >
                        <v-icon left>mdi-plus</v-icon>
                        Create Trust Line
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-form>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Existing Trust Lines -->
        <v-row justify="center">
          <v-col cols="12">
            <v-card tile outlined>
              <v-card-title class="d-flex align-center justify-space-between">
                <span>Your Trust Lines</span>
                <v-btn
                  icon
                  @click="refreshTrustLines"
                  :loading="loading"
                >
                  <v-icon>mdi-refresh</v-icon>
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-data-table
                  v-if="!loading"
                  :headers="trustLineHeaders"
                  :items="trustLines"
                  :items-per-page="10"
                  class="elevation-0"
                  mobile-breakpoint="0"
                >
                  <template #[`item.currency`]="{ item }">
                    <div class="d-flex align-center">
                      <v-avatar size="24" class="mr-2">
                        <v-img
                          :src="$imageUrlBySymbol(item.currency.toLowerCase())"
                          :lazy-src="$imageUrlBySymbol(item.currency.toLowerCase())"
                        />
                      </v-avatar>
                      <span class="font-weight-medium">{{ item.currency }}</span>
                    </div>
                  </template>

                  <template #[`item.issuer`]="{ item }">
                    <div class="d-flex align-center">
                      <span class="font-family-mono">{{ formatAddress(item.issuer) }}</span>
                      <v-btn
                        icon
                        x-small
                        class="ml-1"
                        @click="copyToClipboard(item.issuer)"
                      >
                        <v-icon size="16">mdi-content-copy</v-icon>
                      </v-btn>
                    </div>
                  </template>

                  <template #[`item.limit`]="{ item }">
                    <span class="font-weight-medium">{{ formatAmount(item.limit) }}</span>
                  </template>

                  <template #[`item.balance`]="{ item }">
                    <span :class="item.balance > 0 ? 'green--text' : 'grey--text'">
                      {{ formatAmount(item.balance) }}
                    </span>
                  </template>

                  <template #[`item.flags`]="{ item }">
                    <div class="d-flex flex-wrap">
                      <v-chip
                        v-for="flag in item.flags"
                        :key="flag"
                        small
                        outlined
                        class="mr-1 mb-1"
                      >
                        {{ flag }}
                      </v-chip>
                    </div>
                  </template>

                  <template #[`item.actions`]="{ item }">
                    <div class="d-flex">
                      <v-btn
                        icon
                        x-small
                        color="primary"
                        @click="editTrustLine(item)"
                      >
                        <v-icon size="16">mdi-pencil</v-icon>
                      </v-btn>
                      <v-btn
                        icon
                        x-small
                        color="error"
                        @click="deleteTrustLine(item)"
                      >
                        <v-icon size="16">mdi-delete</v-icon>
                      </v-btn>
                    </div>
                  </template>
                </v-data-table>

                <div v-else class="text-center pa-4">
                  <v-progress-circular indeterminate color="primary" />
                  <div class="mt-2">Loading Trust lines...</div>
                </div>

                <div v-if="!loading && trustLines.length === 0" class="text-center pa-4">
                  <v-icon size="48" class="mb-2 grey--text">mdi-account-off</v-icon>
                  <div class="text-h6 grey--text">No Trust lines found</div>
                  <div class="text-body-2 grey--text">Create a Trust line to start holding XRP tokens</div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Popular Tokens -->
        <v-row justify="center" class="mt-6">
          <v-col cols="12">
            <v-card tile outlined>
              <v-card-title>Popular Tokens</v-card-title>
              <v-card-text>
                <v-row>
                  <v-col
                    v-for="token in popularTokens"
                    :key="token.currency"
                    cols="12"
                    sm="6"
                    md="4"
                    lg="3"
                  >
                    <v-card outlined class="pa-3">
                      <div class="d-flex align-center mb-2">
                        <v-avatar size="32" class="mr-2">
                          <v-img
                            :src="$imageUrlBySymbol(token.currency.toLowerCase())"
                            :lazy-src="$imageUrlBySymbol(token.currency.toLowerCase())"
                          />
                        </v-avatar>
                        <div>
                          <div class="font-weight-medium">{{ token.currency }}</div>
                          <div class="text-caption grey--text">{{ token.name }}</div>
                        </div>
                      </div>
                      <div class="text-caption grey--text mb-2">
                        Issuer: {{ formatAddress(token.issuer) }}
                      </div>
                      <v-btn
                        color="primary"
                        small
                        outlined
                        block
                        :disabled="!isWalletReady"
                        @click="quickCreateTrustLine(token)"
                      >
                        Create Trust Line
                      </v-btn>
                    </v-card>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <!-- Edit Trust Line Dialog -->
    <v-dialog v-model="editDialog" max-width="500">
      <v-card>
        <v-card-title>Edit Trust Line</v-card-title>
        <v-card-text>
          <v-form ref="editForm" v-model="editFormValid">
            <v-text-field
              v-model="editingTrustLine.limit"
              label="Trust Limit"
              type="number"
              :rules="[rules.required, rules.positive]"
              outlined
              dense
            />
            <v-select
              v-model="editingTrustLine.flags"
              label="Trust Line Flags"
              :items="trustLineFlags"
              item-text="label"
              item-value="value"
              outlined
              dense
              multiple
              chips
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="editDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :disabled="!editFormValid"
            :loading="loading"
            @click="updateTrustLine"
          >
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>Delete Trust Line</v-card-title>
        <v-card-text>
          Are you sure you want to delete the Trust line for {{ deletingTrustLine.currency }}?
          <br><br>
          <strong>Warning:</strong> This will remove your ability to hold this token.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="deleteDialog = false">Cancel</v-btn>
          <v-btn
            color="error"
            :loading="loading"
            @click="confirmDeleteTrustLine"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, inject, onMounted } from '@nuxtjs/composition-api'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import { setTrustline } from '@gemwallet/api'

interface TrustLine {
  currency: string
  issuer: string
  limit: number
  balance: number
  flags: string[]
}

interface PopularToken {
  currency: string
  name: string
  issuer: string
}

export default defineComponent({
  setup() {
    const xrpClient = inject(XRP_PLUGIN_KEY) as XrpClient
    
    // State
    const loading = ref(false)
    const formValid = ref(false)
    const editFormValid = ref(false)
    const editDialog = ref(false)
    const deleteDialog = ref(false)
    const trustLines = ref<TrustLine[]>([])
    const editingTrustLine = ref<TrustLine | null>(null)
    const deletingTrustLine = ref<TrustLine | null>(null)
    
    // Form data
    const newTrustLine = ref({
      currency: '',
      issuer: '',
      limit: '',
      flags: [] as string[]
    })

    // Computed
    const isWalletReady = computed(() => xrpClient.isWalletReady.value)
    const address = computed(() => xrpClient.address.value)

    // Trust line flags
    const trustLineFlags = [
      { label: 'Set Authorization', value: 'tfSetfAuth' },
      { label: 'Set No Ripple', value: 'tfSetNoRipple' },
      { label: 'Clear No Ripple', value: 'tfClearNoRipple' },
      { label: 'Set Freeze', value: 'tfSetFreeze' },
      { label: 'Clear Freeze', value: 'tfClearFreeze' }
    ]

    // Popular tokens
    const popularTokens: PopularToken[] = [
      {
        currency: 'USDC',
        name: 'USD Coin',
        issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q'
      },
      {
        currency: 'USDT',
        name: 'Tether USD',
        issuer: 'rchGBxcD1A1C2tdxF6papQYZ8kjRKMYcL'
      },
      {
        currency: 'BTC',
        name: 'Bitcoin',
        issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B'
      },
      {
        currency: 'ETH',
        name: 'Ethereum',
        issuer: 'rcA8X3TVMST1n3CJeAdGk1RdRCHii7N2h'
      }
    ]

    // Table headers
    const trustLineHeaders = [
      { text: 'Currency', value: 'currency', width: '150' },
      { text: 'Issuer', value: 'issuer', width: '200' },
      { text: 'Limit', value: 'limit', width: '120' },
      { text: 'Balance', value: 'balance', width: '120' },
      { text: 'Flags', value: 'flags', width: '150' },
      { text: 'Actions', value: 'actions', width: '100', sortable: false }
    ]

    // Form validation rules
    const rules = {
      required: (v: any) => !!v || 'This field is required',
      currency: (v: string) => /^[A-Z]{3,4}$/.test(v) || 'Currency must be 3-4 uppercase letters',
      address: (v: string) => /^r[a-zA-Z0-9]{25,34}$/.test(v) || 'Invalid XRP address format',
      positive: (v: number) => v > 0 || 'Value must be positive'
    }

    // Methods
    const connectWallet = () => {
      xrpClient.connectWallet()
    }

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // Could add toast notification here
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }

    const formatAddress = (address: string) => {
      return `${address.slice(0, 10)}...${address.slice(-10)}`
    }

    const formatAmount = (amount: number) => {
      return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 6,
      }).format(amount)
    }

    const loadTrustLines = async () => {
      if (!isWalletReady.value) return
      
      loading.value = true
      try {
        // Mock data for development - would be replaced with actual API call
        trustLines.value = [
          {
            currency: 'USDC',
            issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
            limit: 10000,
            balance: 500,
            flags: ['tfSetfAuth']
          },
          {
            currency: 'BTC',
            issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B',
            limit: 100,
            balance: 0,
            flags: []
          }
        ]
      } catch (error) {
        console.error('Failed to load trust lines:', error)
      } finally {
        loading.value = false
      }
    }

    const createTrustLine = async () => {
      if (!isWalletReady.value) return
      
      loading.value = true
      try {
        const payload = {
          limitAmount: {
            currency: newTrustLine.value.currency,
            issuer: newTrustLine.value.issuer,
            value: newTrustLine.value.limit.toString()
          },
          flags: newTrustLine.value.flags.reduce((acc, flag) => {
            acc[flag] = true
            return acc
          }, {} as any)
        }

        const response = await setTrustline(payload)
        console.log('Trust line created:', response)
        
        // Reset form
        newTrustLine.value = {
          currency: '',
          issuer: '',
          limit: '',
          flags: []
        }
        
        // Reload trust lines
        await loadTrustLines()
      } catch (error) {
        console.error('Failed to create trust line:', error)
      } finally {
        loading.value = false
      }
    }

    const quickCreateTrustLine = (token: PopularToken) => {
      newTrustLine.value.currency = token.currency
      newTrustLine.value.issuer = token.issuer
      newTrustLine.value.limit = '1000'
      newTrustLine.value.flags = []
    }

    const editTrustLine = (trustLine: TrustLine) => {
      editingTrustLine.value = { ...trustLine }
      editDialog.value = true
    }

    const updateTrustLine = async () => {
      if (!editingTrustLine.value) return
      
      loading.value = true
      try {
        // Mock update - would be replaced with actual API call
        const index = trustLines.value.findIndex(tl => 
          tl.currency === editingTrustLine.value?.currency && 
          tl.issuer === editingTrustLine.value?.issuer
        )
        
        if (index !== -1) {
          trustLines.value[index] = { ...editingTrustLine.value }
        }
        
        editDialog.value = false
        editingTrustLine.value = null
      } catch (error) {
        console.error('Failed to update trust line:', error)
      } finally {
        loading.value = false
      }
    }

    const deleteTrustLine = (trustLine: TrustLine) => {
      deletingTrustLine.value = trustLine
      deleteDialog.value = true
    }

    const confirmDeleteTrustLine = async () => {
      if (!deletingTrustLine.value) return
      
      loading.value = true
      try {
        // Mock delete - would be replaced with actual API call
        trustLines.value = trustLines.value.filter(tl => 
          !(tl.currency === deletingTrustLine.value?.currency && 
            tl.issuer === deletingTrustLine.value?.issuer)
        )
        
        deleteDialog.value = false
        deletingTrustLine.value = null
      } catch (error) {
        console.error('Failed to delete trust line:', error)
      } finally {
        loading.value = false
      }
    }

    const refreshTrustLines = () => {
      loadTrustLines()
    }

    // Lifecycle
    onMounted(() => {
      if (isWalletReady.value) {
        loadTrustLines()
      }
    })

    return {
      // State
      loading,
      formValid,
      editFormValid,
      editDialog,
      deleteDialog,
      trustLines,
      editingTrustLine,
      deletingTrustLine,
      newTrustLine,
      
      // Computed
      isWalletReady,
      address,
      
      // Data
      trustLineFlags,
      popularTokens,
      trustLineHeaders,
      rules,
      
      // Methods
      connectWallet,
      copyToClipboard,
      formatAddress,
      formatAmount,
      createTrustLine,
      quickCreateTrustLine,
      editTrustLine,
      updateTrustLine,
      deleteTrustLine,
      confirmDeleteTrustLine,
      refreshTrustLines
    }
  },
  head: {
    title: 'XRP Trust Lines',
  },
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style> 