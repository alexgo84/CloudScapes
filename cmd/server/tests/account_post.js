'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})

addTest('get an invalid account', function (t) {
    return t.get('/v1/accounts/asdf')
        .expect(400)
})