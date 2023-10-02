import { ethers } from 'ethers'
import { ChainChangeParamInterface, ConnectorInterface, Web3ErrorInterface } from './connector'
import { Chain } from '~/types/apollo/main/types'
declare const window: any

// TODO: Add network  change
export class MetamaskConnector implements ConnectorInterface {
  private idVal: string = 'MetamaskInjector'
  private providerVal: ethers.providers.Web3Provider | null = null
  private accountVal: string | null = null
  private chainIdVal: number | null = null
  private activeVal: boolean = false
  private errorVal: Web3ErrorInterface = { status: false, message: null }

  get id() {
    return this.idVal
  }

  get provider() {
    return this.providerVal
  }

  get account() {
    return this.accountVal
  }

  get chainId() {
    return this.chainIdVal
  }

  get active() {
    return this.activeVal
  }

  get error() {
    return this.errorVal
  }

  async connect(): Promise<{ account: string | null; error: Web3ErrorInterface }> {
    const isMetamaskInstalled = this.isMetaMaskInstalled()
    if (isMetamaskInstalled) {
      const provider = new ethers.providers.Web3Provider(window.ethereum)
      this.providerVal = provider
      this.activeVal = true

      let accounts: string[] = []
      let chainId = null
      try {
        accounts = await provider.send('eth_requestAccounts', [])
        chainId = await provider.getNetwork()
        if (accounts.length > 0 && chainId) {
          this.accountVal = accounts[0]
          this.chainIdVal = chainId.chainId
          this.registerListeners()
        }
      } catch (error: any) {
        this.errorVal.status = true
        this.errorVal.message = `Metamask Connection Error: ${error.message}`
      }
    } else {
      this.errorVal.status = true
      this.errorVal.message = 'Tried to connect to MetaMask but it was not detected. Please install MetaMask.'
    }
    return {
      account: this.account,
      error: this.error,
    }
  }

  resetErrors() {
    this.error.status = false
    this.error.message = null
  }

  registerListeners() {
    if (!this.provider) {
      throw new Error('No provider registered for this connector. Something is very wrong.')
    }
    window.ethereum.on('accountsChanged', this.handleAccountsChanged)
    window.ethereum.on('chainChanged', this.handleChainChanged)
  }

  isMetaMaskInstalled() {
    // Have to check the ethereum binding on the window object to see if it's installed
    const ethereum = window.ethereum
    return Boolean(ethereum && ethereum.isMetaMask)
  }

  handleAccountsChanged = async (accounts: string[]) => {
    if (accounts.length === 0) {
      await this.handleDisconnect()
    } else {
      this.accountVal = accounts[0]
    }
  }

  handleChainChanged = async () => {
    // Reset Provider if chain changed
    this.providerVal = new ethers.providers.Web3Provider(window.ethereum)
    const chainId = await this.providerVal.getNetwork()
    this.chainIdVal = chainId.chainId
  }

  handleDisconnect() {
    // reset everything
    if (window.ethereum?.removeAllListeners) window.ethereum.removeAllListeners()
    this.providerVal?.removeAllListeners()
    this.providerVal = null
    this.accountVal = null
    this.chainIdVal = null
    this.activeVal = false
    if (window.ethereum.disconnect) {
      window.ethereum.disconnect()
    }
    if (window.ethereum.close) {
      window.ethereum.close()
    }
  }

  async handleChanChange(chain: Chain) {
    let chainId: string | number = chain.chainIdentifier

    if (chainId === 1) {
      chainId = '0x' + chain.chainIdentifier
    } else if (chainId === 1337) {
      chainId = '0x539'
    } else if (chainId === 10) {
      chainId = '0xA'
    } else {
      chainId = ethers.utils.hexlify(chain.chainIdentifier)
    }

    const params: ChainChangeParamInterface = {
      chainId,
      chainName: chain.name,
      nativeCurrency: {
        name: chain.name,
        symbol: chain.symbol, // 2-6 characters long
        decimals: 18,
      },
      blockExplorerUrls: [chain.blockExplorerUrl],
      rpcUrls: [chain.rpcUrl],
    }
    try {
      await window.ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId }], // chainId must be in hexadecimal numbers
      })
    } catch (error: any) {
      if (error.code === 4902) {
        try {
          await window.ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [params],
          })
        } catch (error) {
          console.log(error)
        }
      }
    }
  }

  async importTokenToMetamask({
    address,
    symbol,
    decimals,
    image,
  }: {
    address: string
    symbol: string
    decimals: number
    image: string
  }) {
    try {
      await window.ethereum.request({
        method: 'wallet_watchAsset',
        params: { type: 'ERC20', options: { address, symbol, decimals, image } },
      })
    } catch (error) {
      console.log(error)
    }
  }
}
