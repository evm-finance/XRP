<template>
  <v-row justify="center">
    <v-col cols="12" lg="10">
      <v-row justify="center" class="pt-md-10 pt-3 pb-5">
        <v-col cols="11" lg="4">
          <h1 class="text-h4">Our Team</h1>
          <p class="mt-2 grey--text text--lighten-1">
            Our team's goal is empowering cryptocurrency enthusiasts and traders for success. We know from our own
            trading experiences that the crypto market is fast moving which demands reliable data and tools. Our
            platforms provide content-rich resources that inform, track, and simplify the trading process with superior
            visuals displays. While each member has a varied skill set, we collaborate to offer a fast platform with
            advanced interactive features. We are a new site in active development, so we encourage feedback from users
            regarding features they desire.
          </p>
        </v-col>

        <v-col cols="11" lg="8" class="pl-lg-10">
          <v-row class="my-5">
            <v-col v-for="member in team" :key="member.id" cols="12" md="6" width="200">
              <team-member-profile-dialog :member-info="member"></team-member-profile-dialog>

              <div class="d-flex justify-left">
                <v-avatar size="80" class="ml-5">
                  <img
                    :src="`https://quantifycrypto.s3.us-west-2.amazonaws.com/pictures/team/${member.thumbnailUrl}.jpg`"
                    :alt="`${member.name} avatar`"
                  />
                </v-avatar>

                <div class="ml-6 mt-2">
                  <h4 class="text-h6 v-chip--clickable text-hover-primary" @click="toggleBio(member.id, false, true)">
                    {{ member.name }}
                  </h4>
                  <div class="text--secondary">
                    {{ member.jobTitle }}
                  </div>

                  <v-btn
                    v-for="social in member.socials"
                    :key="social.link"
                    class="mt-1 mr-2 text-hover-primary"
                    small
                    icon
                    :href="social.link"
                    target="_blank"
                  >
                    <v-icon size="20">mdi-{{ social.icon }}</v-icon>
                  </v-btn>
                </div>
              </div>
            </v-col>
          </v-row>
        </v-col>
      </v-row>

      <v-divider></v-divider>

      <v-row justify="center" class="pt-10">
        <v-col cols="11" lg="4">
          <h2 class="text-h4">Our Advisors</h2>
          <p class="mt-2 grey--text text--lighten-1">
            Our Advisors are some of the top performers to work and interface with financial trading systems. All have
            been significant contributors at the New York Stock Exchange, they fully understand how trading systems work
            and the features traders expect.
          </p>
        </v-col>

        <v-col cols="11" lg="8" class="pl-lg-10">
          <v-row class="my-5">
            <v-col v-for="member in advisors" :key="member.id" cols="12" md="6" width="200">
              <team-member-profile-dialog :member-info="member"></team-member-profile-dialog>

              <div class="d-flex justify-left">
                <v-avatar size="80" class="ml-5">
                  <img
                    :src="`https://quantifycrypto.s3.us-west-2.amazonaws.com/pictures/team/${member.thumbnailUrl}.jpg`"
                    :alt="`${member.name} avatar`"
                  />
                </v-avatar>

                <div class="ml-6 mt-2">
                  <h4 class="text-h6 v-chip--clickable text-hover-primary" @click="toggleBio(member.id, true, true)">
                    {{ member.name }}
                  </h4>
                  <div class="text--secondary">
                    {{ member.jobTitle }}
                  </div>

                  <v-btn
                    v-for="social in member.socials"
                    :key="social.link"
                    class="mt-1 mr-2 text-hover-primary"
                    small
                    icon
                    :href="social.link"
                    target="_blank"
                  >
                    <v-icon size="18">mdi-{{ social.icon }}</v-icon>
                  </v-btn>
                </div>
              </div>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, ref, useRoute } from '@nuxtjs/composition-api'
import { useMetaTags } from '~/composables/useMetaTags'
import TeamMemberProfileDialog from '~/components/team/TeamMemberProfileDialog.vue'
import { TeamMember } from '~/types/team'

export default defineComponent({
  name: 'Team',
  components: {
    TeamMemberProfileDialog,
  },
  setup() {
    const team = ref<TeamMember[]>([
      {
        id: 'john_b',
        name: 'John Barry',
        jobTitle: 'CEO / Co-Founder',
        thumbnailUrl: 'john',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/johnbarryit/',
          },
          {
            icon: 'twitter',
            link: 'https://twitter.com/Quantify_Crypto',
          },
        ],
        bio: 'John is an Information Technology developer with over 30 years of executive experience in financial services, including 23 years at the New York Stock Exchange. He flourished as lead developer for the vital SuperDot order flow system that tracked every order entering the NYSE computer systems. Later, John was responsible for the Capacity Planning of all NYSE order flow system, earning the Chairman’s Award for his project management success. As the CEO of Quantify Crypto, John combines his deep understanding of the securities industry and project leadership expertise to form the company vision.',
      },
      {
        id: 'alex_b',
        name: 'Alexander Barry',
        jobTitle: 'CTO / Co-Founder',
        thumbnailUrl: 'alex',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/alexander-barry-56bb13119/',
          },
          {
            icon: 'twitter',
            link: 'https://twitter.com/Alexclarkbarry',
          },
        ],
        bio: 'Alex bought his first Bitcoins in 2016 while in college studying information security as a Computer Science major. Since then, he immersed himself in everything crypto and is an accomplished developer of automated modal bots. Shortly after making his first trades on Uniswap during the DeFi Summer of 2020, he dove headfirst into programming for Ethereum and is the driving force behind integrating Ethereum based products into our platform.',
      },
      {
        id: 'maksym_z',
        name: 'Max Zhakun',
        jobTitle: 'Software Engineer',
        thumbnailUrl: 'maksym',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/maksym-zhakun-938362a2/',
          },
          {
            icon: 'twitter',
            link: 'https://twitter.com/maxzhakun',
          },
        ],
        bio: 'Max is a brilliant developer with complete usage of latest software components and tools.  He is lead system architect of the Quantify Crypto platform, deciding the necessary technology and techniques the system requires.  His full knowledge of databases, AWS servers, networking, programming languages, communication protocols and standards is demonstrated on our platform.',
      },
    ])
    const advisors = ref<TeamMember[]>([
      {
        id: 'jerry_r',
        name: 'Jerry Raio',
        jobTitle: 'Advisor',
        thumbnailUrl: 'jerry',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/jerry-raio-5556a814/',
          },
        ],
        bio: `<div>Jerry is the Founder, President and CEO of Arbor Lane Advisors, Inc., a financial services consulting firm with expertise in product development, product management, marketing, and distribution of investment products. His career on Wall Street spans over 30 years, mostly focused on the distribution of new issue Equity Capital Markets product to Financial Advisors and retail investors. Jerry was Head of Retail Origination at Wells Fargo Securities where he worked closely with asset management clients to structure, market and execute new closed end fund offerings. Jerry also had significant roles at Wachovia (now Wells Fargo), Morgan Stanley and Citigroup.  He serves on the board/advisory board of several firms including RiverNorth Capital Management closed-end funds, FLX Distribution and Qudos Technologies. </div>
        <div class='mt-3'>Jerry holds multiple FINRA licenses including the following: Series 7 (General Securities Representative), 24 (General Principal) and 79 (Investment Banking Representative)</div>`,
      },

      {
        id: 'hardol_b',
        name: 'Harold Bott',
        jobTitle: 'Advisor',
        thumbnailUrl: 'harold',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/haroldbott/',
          },
          {
            icon: 'twitter',
            link: 'https://twitter.com/haroldbott',
          },
        ],
        bio: "Harold’s career started in 1993 at Goldman Sachs while a graduate student at the Columbia University School of Engineering and Applied Science. At Goldman, he worked with the government bond and repo desks to develop modal platforms. He was an angel investor in a Goldman colleague's startup, eGain Communications, now a public company. As Chief Information Officer at specialist firm Bear Wagner, lead the development of electronic modal for NYSE and ARCA equities. Promoted by JP Morgan (after acquiring Bear Wagner) to Executive Director before leaving to serve as a consultant for high frequency modal firms, including Getco. Harold received his BA in Music cum laude and throughout his career has worked as a professional musician.",
      },
      {
        id: 'warren_r',
        name: 'Warren Rosenbaum',
        jobTitle: 'Advisor',
        thumbnailUrl: 'warren',
        bioIsVisible: false,
        socials: [
          {
            icon: 'linkedin',
            link: 'https://www.linkedin.com/in/warrenrosenbaum/',
          },
        ],
        bio: 'Warren was the chief architect of the Technology and Operations divisions of the New York Stock Exchange. A senior Vice President leader with full responsibility for the software and hardware architecture of the NYSE enterprise systems.  Warren provided the vision and innovation guidance keeping NYSE as a technology leader. During his tenure, over 25+ years, the NYSE systems consistently had operational uptimes of over 99.99% due to his successful.',
      },
    ])

    function toggleBio(memberId: string, isAdvisor: boolean, bioIsVisible: boolean): void {
      let teamMembers = team

      if (isAdvisor) {
        teamMembers = advisors
      }

      const member = teamMembers.value.find((obj) => obj.id === memberId)

      if (member) {
        member.bioIsVisible = bioIsVisible
      }
    }

    // META TAGS
    useMetaTags('teams', useRoute().value.path)

    return { team, advisors, toggleBio }
  },
  head: {},
})
</script>
