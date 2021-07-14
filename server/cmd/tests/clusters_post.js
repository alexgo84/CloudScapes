'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})

addTest('get all clusters in account should return 0 clusters in new account', function (t) {
    return t.get('/v1/clusters')
        .expect(200)
        .debug()
        .expectLen('', 0)
})

addTest('create a new cluster', function (t) {
    return t.post('/v1/clusters')
        .send({
            name: 'europe-test1',
            accountId: t.state.session.accountId,
        })
        .expect(201)
})