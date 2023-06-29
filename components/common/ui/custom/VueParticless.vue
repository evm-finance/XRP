<!--
<template>
  <client-only>
    <Particles id="tsparticles" :particles-init="particlesInit" :options="options" />
  </client-only>
</template>

<script lang="ts">
import Particles from 'particles.vue'
import { loadFull } from 'tsparticles'
import { defineComponent } from '@nuxtjs/composition-api'
import { Engine } from 'tsparticles-engine'
import Vue from 'vue'

Vue.use(Particles)

export default defineComponent({
  name: 'VueParticles',
  setup() {
    const particlesInit = async (engine: Engine) => {
      await loadFull(engine)
    }

    const particleImages = [
      { src: '/img/logo/evmfinance-logo.svg' },
      { src: 'https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/wbtc.png' },
      { src: 'https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/eth.png' },
      { src: 'https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/bnb.png' },
      { src: 'https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/ftm.png' },
      { src: 'https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/matic.png' },
    ]

    const options = {
      autoPlay: true,
      background: {
        color: {
          value: '#000000',
        },
        image: '',
        position: '50% 50%',
        repeat: 'no-repeat',
        size: 'cover',
        opacity: 1,
      },
      backgroundMask: {
        composite: 'destination-out',
        cover: {
          color: {
            value: '#fff',
          },
          opacity: 1,
        },
        enable: false,
      },
      fullScreen: {
        enable: true,
        zIndex: 1,
      },
      detectRetina: true,
      duration: 0,
      fpsLimit: 60,
      modes: {
        attract: {
          distance: 200,
          duration: 0.4,
          easing: 'ease-out-quad',
          factor: 1,
          maxSpeed: 50,
          speed: 1,
        },
      },
      particles: {
        bounce: {
          horizontal: {
            random: {
              enable: true,
              minimumValue: 0.1,
            },
            value: 1,
          },
        },
        // The lower the value the smaller the stars
        groups: {
          z5000: {
            number: {
              value: 70,
            },
            zIndex: {
              value: 60,
            },
          },
          z7500: {
            number: {
              value: 80,
            },
            zIndex: {
              value: 75,
            },
          },
          z2500: {
            number: {
              value: 50,
            },
            zIndex: {
              value: 30,
            },
          },
          z1000: {
            number: {
              value: 40,
            },
            zIndex: {
              value: 10,
            },
          },
        },

        move: {
          angle: {
            offset: 0,
            value: 10,
          },

          center: {
            x: 50,
            y: 50,
            radius: 0,
          },
          decay: 0,
          distance: {},
          direction: 'right',
          drift: 0,
          enable: true,
        },
        // the density of the stars
        number: {
          density: {
            enable: false,
            area: 800,
            factor: 1000,
          },
          limit: 0,
          value: 200,
        },
        opacity: {
          random: {
            enable: false,
            minimumValue: 0.1,
          },
          value: 1,
          animation: {
            count: 0,
            enable: false,
            speed: 3,
            sync: false,
            destroy: 'none',
            startValue: 'random',
            minimumValue: 0.1,
          },
        },
        orbit: {
          animation: {
            count: 0,
            enable: false,
            speed: 1,
            sync: false,
          },
          enable: false,
          opacity: 1,
          rotation: {
            random: {
              enable: false,
              minimumValue: 0,
            },
            value: 45,
          },
          width: 1,
        },

        roll: {
          darken: {
            enable: false,
            value: 0,
          },
          enable: false,
          enlighten: {
            enable: false,
            value: 0,
          },
          mode: 'vertical',
          speed: 25,
        },
        rotate: {
          random: {
            enable: false,
            minimumValue: 0,
          },
          value: 0,
          animation: {
            enable: false,
            speed: 0,
            sync: false,
          },
          direction: 'clockwise',
          path: false,
        },
        shadow: {
          blur: 0,
          color: {
            value: '#000',
          },
          enable: false,
          offset: {
            x: 0,
            y: 0,
          },
        },
        shape: {
          options: {},
          type: 'circle',
        },
        size: {
          random: {
            enable: false,
            minimumValue: 1,
          },
          value: 3,
          animation: {
            count: 0,
            enable: false,
            speed: 5,
            sync: false,
            destroy: 'none',
            startValue: 'random',
          },
        },
      },
      zLayers: 100,
      emitters: {
        autoPlay: true,
        fill: true,
        rate: {
          quantity: 1,
          delay: 5,
        },
        shape: 'square',
        startCount: 0,
        size: {
          mode: 'percent',
          height: 0,
          width: 0,
        },
        direction: 'none',
        particles: {
          shape: {
            type: 'images',
            options: {
              images: particleImages,
            },
          },
          size: {
            value: 20,
          },
          move: {
            speed: 7,
            outModes: {
              default: 'none',
              right: 'destroy',
            },
            straight: true,
          },
          zIndex: {
            value: 0,
          },
          rotate: {
            value: {
              min: 0,
              max: 360,
            },
            // rotation speed
            animation: {
              enable: true,
              speed: 10,
              sync: true,
            },
          },
        },
        position: {
          x: -5,
          y: 55,
        },
      },
    }

    return { particlesInit, options }
  },
})
</script>
-->
