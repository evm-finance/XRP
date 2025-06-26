import { computed, useContext, useStore } from '@nuxtjs/composition-api'
import { mdiContentCopy, mdiStar, mdiExternalLink, mdiSwapHorizontal } from '@mdi/js'
import { useXrpFormatters } from './useXrpFormatters'

export interface XrpGridRendererParams {
  data: any
  value: any
  column: any
  node: any
}

export function useXrpGridRenderers() {
  const { $f, $imageUrlBySymbol } = useContext()
  const { state, dispatch } = useStore()
  const { 
    formatXrpPrice, 
    formatIssuerAddress, 
    formatTrustLine, 
    formatAmmPool,
    formatPercentageChange,
    formatMarketCap,
    formatVolume,
    formatTransactionHash,
    formatTimestamp,
    formatXrpAmount
  } = useXrpFormatters()

  // Favorite tokens management
  const favoriteTokens = computed<string[]>({
    get: () => state.xrp?.favorites || [],
    set: (value: string[]) => dispatch('xrp/setFavorites', value)
  })

  const setFavoriteToken = (tokenId: string, status: boolean) => {
    const favorites = [...favoriteTokens.value]
    if (!status) {
      favorites.push(tokenId)
    } else {
      const index = favorites.indexOf(tokenId)
      if (index > -1) {
        favorites.splice(index, 1)
      }
    }
    favoriteTokens.value = favorites
  }

  /**
   * Token name with icon and link renderer
   */
  const tokenNameRenderer = (params: XrpGridRendererParams) => {
    const token = params.data
    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex align-center'
    iDiv.style.cssText = 'gap: 8px; cursor: pointer;'

    // Token icon
    const image = document.createElement('img')
    image.src = $imageUrlBySymbol ? $imageUrlBySymbol(token.currency) : '/img/default-token.svg'
    image.className = 'token-icon'
    image.style.cssText = 'width: 24px; height: 24px; border-radius: 50%;'

    // Token name and currency
    const textDiv = document.createElement('div')
    textDiv.className = 'd-flex flex-column'
    
    const nameSpan = document.createElement('span')
    nameSpan.textContent = token.name || token.currency
    nameSpan.style.cssText = 'font-weight: 500; font-size: 14px;'
    
    const currencySpan = document.createElement('span')
    currencySpan.textContent = token.currency
    currencySpan.style.cssText = 'font-size: 12px; color: #888;'

    textDiv.appendChild(nameSpan)
    textDiv.appendChild(currencySpan)

    iDiv.appendChild(image)
    iDiv.appendChild(textDiv)

    // Click handler for navigation
    iDiv.addEventListener('click', () => {
      window.$nuxt.$router.push(`/token/${token.currency}`)
    })

    return iDiv
  }

  /**
   * Issuer address with copy functionality
   */
  const issuerAddressRenderer = (params: XrpGridRendererParams) => {
    const address = params.value
    if (!address) return '-'

    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex align-center'
    iDiv.style.cssText = 'gap: 4px;'

    const addressSpan = document.createElement('span')
    addressSpan.textContent = formatIssuerAddress(address)
    addressSpan.style.cssText = 'font-family: monospace; font-size: 12px;'

    const copyBtn = document.createElement('button')
    copyBtn.className = 'v-btn v-btn--icon v-btn--small'
    copyBtn.innerHTML = `<v-icon small>${mdiContentCopy}</v-icon>`
    copyBtn.addEventListener('click', (e) => {
      e.stopPropagation()
      navigator.clipboard.writeText(address)
      // Show toast notification
      window.$nuxt.$emit('show-toast', { message: 'Address copied to clipboard', type: 'success' })
    })

    iDiv.appendChild(addressSpan)
    iDiv.appendChild(copyBtn)

    return iDiv
  }

  /**
   * Price with change percentage
   */
  const priceRenderer = (params: XrpGridRendererParams) => {
    const token = params.data
    const price = token.price_usd || 0
    const change = token.price_change_24h || 0

    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex flex-column'
    iDiv.style.cssText = 'gap: 2px;'

    const priceSpan = document.createElement('span')
    priceSpan.textContent = formatXrpPrice(price)
    priceSpan.style.cssText = 'font-weight: 500;'

    const changeSpan = document.createElement('span')
    changeSpan.textContent = formatPercentageChange(change)
    changeSpan.style.cssText = `font-size: 12px; color: ${change >= 0 ? '#4caf50' : '#f44336'};`

    iDiv.appendChild(priceSpan)
    iDiv.appendChild(changeSpan)

    return iDiv
  }

  /**
   * Market cap with formatting
   */
  const marketCapRenderer = (params: XrpGridRendererParams) => {
    const marketCap = params.value
    const span = document.createElement('span')
    span.textContent = formatMarketCap(marketCap)
    span.style.cssText = 'font-weight: 500;'
    return span
  }

  /**
   * Volume with formatting
   */
  const volumeRenderer = (params: XrpGridRendererParams) => {
    const volume = params.value
    const span = document.createElement('span')
    span.textContent = formatVolume(volume)
    span.style.cssText = 'font-weight: 500;'
    return span
  }

  /**
   * Trust line status renderer
   */
  const trustLineRenderer = (params: XrpGridRendererParams) => {
    const token = params.data
    const trustLine = token.trustLine

    if (!trustLine) {
      const span = document.createElement('span')
      span.textContent = 'No Trust Line'
      span.style.cssText = 'color: #888; font-style: italic;'
      return span
    }

    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex flex-column'
    iDiv.style.cssText = 'gap: 2px;'

    const statusSpan = document.createElement('span')
    statusSpan.textContent = trustLine.status
    statusSpan.style.cssText = `font-weight: 500; color: ${trustLine.status === 'Active' ? '#4caf50' : '#f44336'};`

    const balanceSpan = document.createElement('span')
    balanceSpan.textContent = trustLine.balance
    balanceSpan.style.cssText = 'font-size: 12px; color: #888;'

    iDiv.appendChild(statusSpan)
    iDiv.appendChild(balanceSpan)

    return iDiv
  }

  /**
   * AMM pool information renderer
   */
  const ammPoolRenderer = (params: XrpGridRendererParams) => {
    const pool = params.data.ammPool
    if (!pool) {
      const span = document.createElement('span')
      span.textContent = 'No AMM Pool'
      span.style.cssText = 'color: #888; font-style: italic;'
      return span
    }

    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex align-center'
    iDiv.style.cssText = 'gap: 8px;'

    const icon = document.createElement('v-icon')
    icon.innerHTML = mdiSwapHorizontal
    icon.style.cssText = 'color: #536af6;'

    const textDiv = document.createElement('div')
    textDiv.className = 'd-flex flex-column'
    
    const pairSpan = document.createElement('span')
    pairSpan.textContent = `${pool.asset1}/${pool.asset2}`
    pairSpan.style.cssText = 'font-weight: 500; font-size: 12px;'
    
    const liquiditySpan = document.createElement('span')
    liquiditySpan.textContent = pool.liquidity
    liquiditySpan.style.cssText = 'font-size: 11px; color: #888;'

    textDiv.appendChild(pairSpan)
    textDiv.appendChild(liquiditySpan)

    iDiv.appendChild(icon)
    iDiv.appendChild(textDiv)

    return iDiv
  }

  /**
   * Favorite button renderer
   */
  const favoriteRenderer = (params: XrpGridRendererParams) => {
    const tokenId = params.data.currency
    const isFavorite = favoriteTokens.value.includes(tokenId)

    const button = document.createElement('button')
    button.className = 'v-btn v-btn--icon v-btn--small'
    button.style.cssText = 'color: transparent;'
    button.innerHTML = `<v-icon small style="color: ${isFavorite ? '#ff9800' : '#888'}">${mdiStar}</v-icon>`
    
    button.addEventListener('click', (e) => {
      e.stopPropagation()
      setFavoriteToken(tokenId, isFavorite)
    })

    return button
  }

  /**
   * Transaction hash with link renderer
   */
  const transactionHashRenderer = (params: XrpGridRendererParams) => {
    const hash = params.value
    if (!hash) return '-'

    const iDiv = document.createElement('div')
    iDiv.className = 'd-flex align-center'
    iDiv.style.cssText = 'gap: 4px;'

    const hashSpan = document.createElement('span')
    hashSpan.textContent = formatTransactionHash(hash)
    hashSpan.style.cssText = 'font-family: monospace; font-size: 12px; cursor: pointer;'

    const linkBtn = document.createElement('button')
    linkBtn.className = 'v-btn v-btn--icon v-btn--small'
    linkBtn.innerHTML = `<v-icon small>${mdiExternalLink}</v-icon>`
    linkBtn.addEventListener('click', (e) => {
      e.stopPropagation()
      window.open(`/xrp-explorer/tx/${hash}`, '_blank')
    })

    iDiv.appendChild(hashSpan)
    iDiv.appendChild(linkBtn)

    return iDiv
  }

  /**
   * Timestamp renderer
   */
  const timestampRenderer = (params: XrpGridRendererParams) => {
    const timestamp = params.value
    const span = document.createElement('span')
    span.textContent = formatTimestamp(timestamp)
    span.style.cssText = 'font-size: 12px; color: #888;'
    return span
  }

  /**
   * XRP amount renderer (drops to XRP)
   */
  const xrpAmountRenderer = (params: XrpGridRendererParams) => {
    const drops = params.value
    const span = document.createElement('span')
    span.textContent = formatXrpAmount(drops)
    span.style.cssText = 'font-weight: 500;'
    return span
  }

  /**
   * Ledger index renderer
   */
  const ledgerIndexRenderer = (params: XrpGridRendererParams) => {
    const index = params.value
    const span = document.createElement('span')
    span.textContent = formatLedgerIndex(index)
    span.style.cssText = 'font-family: monospace; font-size: 12px;'
    return span
  }

  /**
   * Status indicator renderer
   */
  const statusRenderer = (params: XrpGridRendererParams) => {
    const status = params.value
    const span = document.createElement('span')
    span.textContent = status
    span.style.cssText = `
      padding: 2px 8px;
      border-radius: 12px;
      font-size: 11px;
      font-weight: 500;
      text-transform: uppercase;
      background-color: ${status === 'Success' ? '#4caf50' : status === 'Failed' ? '#f44336' : '#ff9800'};
      color: white;
    `
    return span
  }

  return {
    tokenNameRenderer,
    issuerAddressRenderer,
    priceRenderer,
    marketCapRenderer,
    volumeRenderer,
    trustLineRenderer,
    ammPoolRenderer,
    favoriteRenderer,
    transactionHashRenderer,
    timestampRenderer,
    xrpAmountRenderer,
    ledgerIndexRenderer,
    statusRenderer,
    setFavoriteToken
  }
} 