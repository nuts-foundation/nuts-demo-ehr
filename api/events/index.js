const requireDir = require('../../util/require-dir')
const routes = requireDir(__filename, __dirname)

module.exports = io => {
  Object.keys(routes).forEach(r =>
    routes[r](io.of(`/${r}`))
  )
}
