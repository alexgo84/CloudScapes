'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})

addTest('get all accounts', function (t) {
    return t.get('/v1/accounts')
        .expect(200)
})

addTest('create a new user in same account', function (t) {
    return t.post('/v1/users')
        .send({
            name: 'alex',
            email: 'gordonsnif@gmail.com',
            password: 'my-secure-password',
            accountId: t.state.session.accountId,
        })
        .expect(201)
})

