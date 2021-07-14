'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})

addTest('get all clusters in account should return 0 clusters in new account', function (t) {
    return t.get('/v1/clusters')
        .expect(200)
        .expectLen(null, 0)
})

addTest('create a new cluster', function (t) {
    return t.post('/v1/clusters')
        .send({
            name: 'europe-test1',
            accountId: t.state.session.accountId,
        })
        .expect(201)
        .store('clusterId', 'id')
})

addTest('get all clusters in account should return 1 clusters', function (t) {
    return t.get('/v1/clusters')
        .expect(200)
        .expectLen(null, 1)
})

addTest('create a new cluster with same name should fail', function (t) {
    return t.post('/v1/clusters')
        .send({
            name: 'europe-test1',
            accountId: t.state.session.accountId,
        })
        .expect(409)
        .expectField(null, `Key (accountid, name)=(${t.state.session.accountId}, ${'europe-test1'}) already exists.`)
})

addTest('create a new cluster with a different name should succeed', function (t) {
    return t.post('/v1/clusters')
        .send({
            name: 'europe-test2',
            accountId: t.state.session.accountId,
        })
        .expect(201)
})

addTest('get all clusters in account should return both created clusters', function (t) {
    return t.get('/v1/clusters')
        .expect(200)
        .expectLen(null, 2)
})

addTest('delete first cluster', function (t) {
    return t.delete(`/v1/clusters/${t.state.clusterId}`)
        .expect(204)
})

addTest('delete same cluster will return a 404 NOT_FOUND', function (t) {
    return t.delete(`/v1/clusters/${t.state.clusterId}`)
        .expect(404)
})

addTest('get all clusters in account should return one clusters', function (t) {
    return t.get('/v1/clusters')
        .expect(200)
        .expectLen(null, 1)
})