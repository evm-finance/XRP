<template>
  <v-row justify="center" no-gutters class="grey--text text--lighten-1 body-2">
    <v-col cols="11" lg="8" md="9">
      <div class="text-h4 py-3 white--text">About Us</div>
      <v-divider />

      <h1 class="text-h6 py-3 white--text">EVM Finance Launches - July 2022!</h1>
      <p>
        We are proud to announce the launch of the EVM.Finance platform, made possible through a grant from the Aave
        Community. EVM.Finance was built by the
        <nuxt-link to="/team" class="white--text text-hover-primary" v-text="`talented team`" />
        that created and supports the innovative
        <a
          target="_blank"
          class="white--text text-hover-primary"
          href="https://quantifycrypto.com"
          v-text="`QuantifyCrypto.com.`"
        />
        We admire the Decentralized Finance (DeFi) projects that have evolved in this space since they represent a
        significant step forward, creating financial systems that are faster, transparent, and more efficient for all.
        DeFi financial services offer higher interest rate returns and advanced hedging techniques not previously
        available for non-accredited investors.
      </p>
      <p>
        The Aave platform and the Aave smart contracts are proven and secure. They achieved a “Process Quality Review”
        score of 96% from an independent audit of the Aave V2 protocol smart contracts, conducted by defisafety.com.
        Very few protocols earn a ranking above 90% and only 9 out of 246 protocols had a higher score than Aave V2. The
        importance of security and safety in the DeFi space cannot be overlooked. Reference the
        <a class="white--text text-hover-primary" target="_blank" href="https://www.defisafety.com/app/pqrs/18">
          link
        </a>
        for more information.
      </p>

      <h2 class="text-h6 py-3 white--text">Aave Professional Interface</h2>
      <p>
        Our first DeFi project for EVM (Ethereum Virtual Machine) Finance is the “Aave Pro Trader” Interface. We were
        attracted to Aave for their advanced and secure crypto loan platform which is backed by robust Web3 technology.
        We are privileged to have our grant proposal accepted, uniting our efforts to advance the Aave platform.
      </p>
      <p>
        The EVM platform uses the identical Aave V2 smart contracts audited by
        <a
          class="white--text text-hover-primary"
          target="_blank"
          href="https://www.defisafety.com/"
          v-text="`defisafety.com.`"
        />
        We offer visual displays, software frameworks, and metrics computed by Quantify Crypto. We have also added
        detailed wallet balances, loan positions and transaction history to benefit users on our platform. Additional
        add on features, from Quantify Crypto metrics, include price trends; support and resistance levels; and live
        prices. Our design philosophy is providing as much quality information as possible on a single screen.
      </p>

      <h2 class="text-h6 py-3 white--text">Key points about EVM Finance</h2>
      <ol class="pl-4 pb-4">
        <li>Uses the same audited V2 smart contracts as the Aave.com platform.</li>
        <li>Adds detailed DeFi wallet balances and transaction history</li>
        <li>Coin pages with live prices, news, trends, price chart and metrics for all tokens on Aave</li>
      </ol>

      <p>
        The best part, there are no fees added by EVM Finance site for this service. A special thanks goes out to the
        <a
          class="white--text text-hover-primary"
          target="_blank"
          href="https://lenster.xyz/u/aavegrants.lens"
          v-text="`AaveGrant.lens`"
        />
        program
      </p>

      <h3 class="py-3 white--text">The EVM Finance Team</h3>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, useRoute } from '@nuxtjs/composition-api'
import { useMetaTags } from '~/composables/useMetaTags'

export default defineComponent({
  setup() {
    useMetaTags('about', useRoute().value.path)
  },
  head: {},
})
</script>
