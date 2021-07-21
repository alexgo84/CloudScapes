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
            email: 'gordonsniff@gmail.com',
            password: 'my-secure-password',
            accountId: t.state.session.accountId,
        })
        .expect(201)
})

addTest('verify account has two users', function (t) {
    return t.get('/v1/users')
        .expectField('0.email', 'gordonsniff@gmail.com')
        .expectField('0.name', 'alex')
        .expect(200)
})
