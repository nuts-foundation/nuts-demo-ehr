const requireDir = require('../../util/require-dir')
const routes = requireDir(__filename, __dirname)
const router = require('express').Router()


module.exports = (organisation) => {
    // Every endpoint gets the path of /api/filename/methodName
    Object.keys(routes).forEach(r => {
        router.use(`/${r}`, routes[r](organisation));
    });
    return router
};
