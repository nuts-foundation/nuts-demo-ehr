const consents = []

module.exports = {
  get: (match) => {
    return Promise.resolve(consents.filter((c) => {
      for (const prop in match) {
        if (c[prop] !== match[prop]) { return false }
      }
      return true
    }))
  },

  store: (consent) => {
    consent.id = Math.max(...consents.map(o => o.id), -1) + 1
    consents.push(consent)
    return Promise.resolve(consent)
  }
}
