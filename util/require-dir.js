const fs = require('fs')
const path = require('path')
const string = require('./string')

module.exports = (caller, directory, capitalize = false) => {
  const allFiles = {}

  const files = fs.readdirSync(directory)
    .filter(file => file.indexOf('.') !== 0 &&
                                  file !== path.basename(caller) &&
                                  file.indexOf('.js') !== -1)

  files.forEach(file => {
    const contents = require(path.join(directory, file))
    let name = string.camelcase(file.replace('.js', '').replace(/\-/g, ' '))
    if (capitalize) name = string.capitalize(name)
    allFiles[name] = contents
  })

  return allFiles
}
