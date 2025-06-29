<template>
  <div v-if="errorState.hasError" class="graphql-error-display">
    <v-alert
      :type="alertType"
      :icon="alertIcon"
      :dismissible="dismissible"
      @click:close="handleDismiss"
      class="mb-4"
    >
      <div class="d-flex align-center">
        <div class="flex-grow-1">
          <div class="text-body-1 font-weight-medium">
            {{ errorState.errorMessage }}
          </div>
          
          <!-- Retry information -->
          <div v-if="canRetry && !errorState.isRetrying" class="text-caption mt-1">
            Retry {{ errorState.retryCount + 1 }} of {{ maxRetries }} attempts
          </div>
          
          <!-- Retrying indicator -->
          <div v-if="errorState.isRetrying" class="text-caption mt-1">
            <v-progress-circular indeterminate size="16" class="mr-2" />
            Retrying... ({{ errorState.retryCount }}/{{ maxRetries }})
          </div>
          
          <!-- Technical details for developers -->
          <div v-if="showTechnicalDetails" class="text-caption mt-2 grey--text">
            <div>Error Type: {{ errorState.errorType }}</div>
            <div v-if="errorState.lastError?.message">
              Technical Details: {{ errorState.lastError.message }}
            </div>
          </div>
        </div>
        
        <!-- Action buttons -->
        <div class="ml-4 d-flex align-center">
          <v-btn
            v-if="canRetry && !errorState.isRetrying"
            small
            color="primary"
            @click="handleRetry"
            :loading="errorState.isRetrying"
          >
            <v-icon left small>mdi-refresh</v-icon>
            Retry
          </v-btn>
          
          <v-btn
            v-if="showTechnicalDetails"
            small
            text
            @click="toggleTechnicalDetails"
            class="ml-2"
          >
            <v-icon small>{{ showTechnicalDetails ? 'mdi-eye-off' : 'mdi-eye' }}</v-icon>
          </v-btn>
        </div>
      </div>
    </v-alert>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from '@nuxtjs/composition-api'

export interface ErrorState {
  hasError: boolean
  errorMessage: string
  errorType: 'network' | 'graphql' | 'validation' | 'unknown'
  retryCount: number
  lastError?: any
  isRetrying: boolean
}

export default defineComponent({
  name: 'GraphQLErrorDisplay',
  props: {
    errorState: {
      type: Object as () => ErrorState,
      required: true
    },
    canRetry: {
      type: Boolean,
      default: true
    },
    maxRetries: {
      type: Number,
      default: 3
    },
    dismissible: {
      type: Boolean,
      default: true
    },
    showTechnicalDetailsByDefault: {
      type: Boolean,
      default: false
    }
  },
  emits: ['retry', 'dismiss'],
  setup(props, { emit }) {
    const showTechnicalDetails = ref(props.showTechnicalDetailsByDefault)

    // Alert type based on error type
    const alertType = computed(() => {
      switch (props.errorState.errorType) {
        case 'network':
          return 'warning'
        case 'graphql':
          return 'error'
        case 'validation':
          return 'info'
        case 'unknown':
        default:
          return 'error'
      }
    })

    // Alert icon based on error type
    const alertIcon = computed(() => {
      switch (props.errorState.errorType) {
        case 'network':
          return 'mdi-wifi-off'
        case 'graphql':
          return 'mdi-database-alert'
        case 'validation':
          return 'mdi-alert-circle'
        case 'unknown':
        default:
          return 'mdi-alert'
      }
    })

    // Handle retry
    const handleRetry = () => {
      emit('retry')
    }

    // Handle dismiss
    const handleDismiss = () => {
      emit('dismiss')
    }

    // Toggle technical details
    const toggleTechnicalDetails = () => {
      showTechnicalDetails.value = !showTechnicalDetails.value
    }

    return {
      alertType,
      alertIcon,
      showTechnicalDetails,
      handleRetry,
      handleDismiss,
      toggleTechnicalDetails
    }
  }
})
</script>

<style scoped>
.graphql-error-display {
  position: relative;
}

.graphql-error-display .v-alert {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.graphql-error-display .v-alert__content {
  padding: 12px 16px;
}
</style> 