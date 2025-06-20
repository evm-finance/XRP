<template>
  <v-row justify="center" no-gutters class="grey--text text--lighten-1 body-2">
    <v-col cols="11" lg="8" md="9">
      <h1 class="text-h4 py-3 white--text">Frequently Asked Questions about Aave and EVM Finance</h1>
      <v-divider />

      <h2 class="text-h6 py-3 white--text">What does EVM mean?</h2>
      <p>
        EVM is the acronym for Ethereum Virtual Machine. EVM is the software platform where all Ethereum Blockchain
        accounts are stored, smart contracts and blockchain application are executed on.
      </p>

      <h2 class="text-h6 py-3 white--text">What is a smart contract?</h2>
      <p>
        Smart contracts are programmed code designed to execute on the blockchain. Smart contracts provide a secure,
        transparent and efficient way for DeFi transactions to occur without the need for a “middle man” found in
        traditional financial systems.
      </p>

      <h2 class="text-h6 py-3 white--text">What is Aave?</h2>
      <p>
        Aave is a liquidity protocol built for the Ethereum Virtual Machine that allows for decentralized lending and
        borrowing of assets. You can lend assets into the protocol using the deposit function and earn fees when other
        users borrow those assets, additionally your deposited assets can be used as collateral to borrow other crypto
        currencies from the protocol.
      </p>

      <h2 class="text-h6 py-3 white--text">What can I do with Aave?</h2>
      <p>
        Aave’s lending protocol allows you to have multiple assets as both collateral and debt, which can be used to
        create leveraged long and short positions across all supported assets. While using leverage can be risky, it is
        useful to have short positions during long bear markets where the price of Bitcoin and Ethereum are known to
        drop by over 85% in value.
      </p>

      <h2 class="text-h6 py-3 white--text">What makes Aave safer than centralized lending platforms?</h2>
      <p>
        The blockchain has 100% uptime and you can withdraw your funds from the protocol at any time. All of Aave’s
        assets and liabilities are recorded publicly on the blockchain so users can easily verify the protocol’s balance
        sheet is healthy. Aave has a decentralized liquidation system, which means any user can liquidate an account
        that has become delinquent on its loans (health factor below 1) ensuring the protocol’s solvency! These features
        make the system far safer than centralized lending platforms, such as Celsius and Block Fi, which are not
        transparent about how they lend out funds or the process of liquidating bad loans, which has resulted in losses
        of Billions of dollars of user deposits.
      </p>

      <h2 class="text-h6 py-3 white--text">Which smart contracts are used by EVM Finance?</h2>
      <p>
        We use the Aave V2 Smart Contracts, which scored 96% in a defisafety.com audit. We have modified the user
        interface for the action buttons used on the EVM Finance platform. Reference the attached link for more info
        <a href="https://www.defisafety.com/app/pqrs/18" target="_blank" class="white--text text-hover-primary"
          >https://www.defisafety.com/app/pqrs/18</a
        >
      </p>

      <h2 class="text-h6 py-3 white--text">Does EVM Finance collect fees?</h2>
      <p>No. The EVM Finance Aave pro interface is made possible from an Aave grant.</p>

      <h2 class="text-h6 py-3 white--text">Who has custody of your crypto funds?</h2>
      <p>
        Almost all DeFi products, including EVM finance, do not take custody of your funds. Metamask or another
        cryptocurrency wallet is used for custody of your funds.
      </p>

      <h2 class="text-h6 py-3 white--text">What are the key benefits for using EVM Finance?</h2>
      <p>
        We provide detailed token pages with live updates, consolidated displays with quick action buttons with balances
        and transaction history.
      </p>

      <h2 class="text-h6 py-3 white--text">Does the Ethereum gas fees change dynamically?</h2>
      <p>
        Yes, gas fees are based on how active the Ethereum network is. Ethereum gas prices (called gwei) goes up when
        activity increases and decreases when traffic declines. The current gwei gas amount is updating automatically in
        the upper high hand corner of the screen.
      </p>

      <h2 class="text-h6 py-3 white--text">How do I create a long position using EVM Finance?</h2>
      <ol class="pl-4 pb-4">
        <li>Use the “Deposit” function for the asset you wish to be long (ETH, AAVE, UNI) as collateral to Aave.</li>
        <li>Use the “Borrow” function to receive stablecoins that will become debt.</li>
        <li>Use borrowed stablecoins to buy more of the asset you want. Using a DEX like Uniswap.</li>
        <li>(Optional) Deposit tokens back to Aave to increase health factor and avoid liquidation.</li>
        <li>When your assets go up, close out all or part of the position using the “Repay” function.</li>
      </ol>

      <h2 class="text-h6 py-3 white--text">How do I create a short position using EVM Finance?</h2>
      <ol class="pl-4 pb-4">
        <li>Use the “Deposit” function for a stablecoin that can be used as collateral (USDC or DAI, not USDT)</li>
        <li>Use the “Borrow” function for the token you wish to short.</li>
        <li>Sell the borrowed tokens for stablecoins. Using a DEX like Uniswap.</li>
        <li>(Optional) Deposit tokens back to Aave to increase health factor and avoid liquidation.</li>
        <li>
          When the assets you borrowed goes down, close out all or part of the position using the “Repay” function.
        </li>
      </ol>

      <div class="pt-3">
        Have an additional question – Email us at
        <a class="text-hover-primary white--text" href="mailto:admin@quantifycrypto.net">admin@quantifycrypto.net</a>
      </div>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, useRoute } from '@nuxtjs/composition-api'
import { useMetaTags } from '~/composables/useMetaTags'
export default defineComponent({
  setup() {
    useMetaTags('faq', useRoute().value.path)
  },
  head: {},
})
</script>
