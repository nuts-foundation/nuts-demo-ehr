module.exports = {

  camelcase: (str) => {
    return str.replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => {
      return index == 0 ? word.toLowerCase() : word.toUpperCase()
    }).replace(/\s+/g, '')
  },

  capitalize: (str) => {
    return str.charAt(0).toUpperCase() + str.slice(1)
  }

}
