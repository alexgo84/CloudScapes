'use strict'

const addTest = require('./_test').init()

addTest('call healthcheck endpoint', function (t) {
    return t.get('/v1/status/health')
    .expect(200)
})