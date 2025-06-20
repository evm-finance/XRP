export interface TeamMember {
  id: string
  name: string
  jobTitle: string
  thumbnailUrl: string
  bioIsVisible: boolean
  socials: { icon: string; link: string }[]
  bio: string
}
