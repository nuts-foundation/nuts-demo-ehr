const requireDir = require('../util/require-dir');
const routes     = requireDir(__filename, __dirname);
const router     = require('express').Router();

Object.keys(routes).forEach(r =>
  router.use(`/${r}`, routes[r]));

module.exports = router;
