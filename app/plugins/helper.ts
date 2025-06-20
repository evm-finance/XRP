 import { defineNuxtPlugin } from '@nuxtjs/composition-api'
import { Context } from '@nuxt/types'
import Vue from 'vue'

// Returns a formatted string with specified minimum and maximum fractional digits
const numFormat = (val: number, minFracDigits = 2, maxFracDigits = 2): string =>
  new Intl.NumberFormat('en', {
    minimumFractionDigits: minFracDigits,
    maximumFractionDigits: maxFracDigits,
  }).format(val)

function applyNumber(val: number, pre: string) {
  return val < 10 ** 3
    ? pre + numFormat(val)
    : 10 ** 3 <= val && val < 10 ** 6
    ? pre + numFormat(val / 10 ** 3) + ' K'
    : 10 ** 6 <= val && val < 10 ** 9
    ? pre + numFormat(val / 10 ** 6) + ' M'
    : pre + numFormat(val / 10 ** 9) + ' B'
}

const numberFormatter = (
  val: number,
  { minDigits = 2, maxDigits = 2, pre = '', after = '', useSymbol = false }
): string => {
  if (useSymbol) {
    return applyNumber(val, pre)
  } else {
    maxDigits = maxDigits <= minDigits ? (maxDigits = minDigits) : maxDigits
    const value = numFormat(val, minDigits, maxDigits)
    return pre + value + after
  }
}

const copyAddressToClipboard = async (value: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(value)
  } catch (e) {}
}

const setAltImageUrl = (event: any) => {
  return (event.target.src = `https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/generic.png`)
}
const imageUrlBySymbol = (symbol: string | null) =>
  symbol
    ? `https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/${symbol.toLowerCase()}.png`
    : `https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/generic.png`

const truncateAddress = (address: string, zeroIndexTo: number, endIndexMinus: number): string =>
  address ? address.slice(0, zeroIndexTo) + '...' + address.slice(address.length - endIndexMinus, address.length) : ''

const navigateToExplorer = (
  address: string,
  type: 'tx' | 'address' = 'address',
  blockExplorerUrl = 'https://etherscan.io/'
): void => {
  const url = `${blockExplorerUrl}${type}/${address}`
  window.open(url)
}

const applyPtcChange = (val: number): { value: string; color: string; icon: string | null } => {
  const rounded = parseFloat((val * 100).toFixed(2))
  if (rounded === 0) {
    return { value: '0,00%', color: 'grey', icon: null }
  }
  if (rounded > 0) {
    return { value: `${rounded}%`, color: 'green', icon: 'mdi-arrow-up' }
  }
  if (rounded < 0) {
    return { value: `${rounded}%`, color: 'red', icon: 'mdi-arrow-down' }
  }
  return { value: '-', color: 'grey', icon: null }
}

export default defineNuxtPlugin((context: Context) => {
  context.$f = numberFormatter
  context.$copyAddressToClipboard = copyAddressToClipboard
  context.$imageUrlBySymbol = imageUrlBySymbol
  context.$applyPtcChange = applyPtcChange

  Vue.prototype.$f = numberFormatter
  Vue.prototype.$copyAddressToClipboard = copyAddressToClipboard
  Vue.prototype.$setAltImageUrl = setAltImageUrl
  Vue.prototype.$imageUrlBySymbol = imageUrlBySymbol
  Vue.prototype.$truncateAddress = truncateAddress
  Vue.prototype.$navigateToExplorer = navigateToExplorer
  Vue.prototype.$applyPtcChange = applyPtcChange
})
