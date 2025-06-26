export const useMetaTags = (options: {
  title?: string
  description?: string
  image?: string
  url?: string
}) => {
  const { title, description, image, url } = options

  return {
    title: title || 'XRP Finance',
    meta: [
      {
        hid: 'description',
        name: 'description',
        content: description || 'XRP Finance - DeFi analytics and trading platform'
      },
      {
        hid: 'og:title',
        property: 'og:title',
        content: title || 'XRP Finance'
      },
      {
        hid: 'og:description',
        property: 'og:description',
        content: description || 'XRP Finance - DeFi analytics and trading platform'
      },
      {
        hid: 'og:image',
        property: 'og:image',
        content: image || '/img/logo/evmfinance-logo.svg'
      },
      {
        hid: 'og:url',
        property: 'og:url',
        content: url || process.env.BASE_URL
      },
      {
        hid: 'twitter:card',
        name: 'twitter:card',
        content: 'summary_large_image'
      },
      {
        hid: 'twitter:title',
        name: 'twitter:title',
        content: title || 'XRP Finance'
      },
      {
        hid: 'twitter:description',
        name: 'twitter:description',
        content: description || 'XRP Finance - DeFi analytics and trading platform'
      },
      {
        hid: 'twitter:image',
        name: 'twitter:image',
        content: image || '/img/logo/evmfinance-logo.svg'
      }
    ]
  }
} 