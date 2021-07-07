'use strict'

const addTest = require('./_test').init()

addTest('create account', function (t) {
    return t.createAccount()
})

addTest('create account 2', function (t) {
    return t.createAccount()
})