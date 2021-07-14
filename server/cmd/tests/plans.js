'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})


addTest('get all plans in account should return 0 plans', function (t) {
    return t.get('/v1/plans')
        .expect(200)
        .expectLen(null, 0)
})

addTest('fail to create a plan that references a non existing cluster', function (t) {
    return t.post('/v1/plans')
        .send({
            name: 'good plan',
            accountId: t.state.session.accountId,
            Replicase: 3,
            CPULimit: '100m',
            MemLimit: '100m',
            CPUReq: '10m',
            MemReq: '10m',
            ClusterID: 1, // violates foreign key constraint
            DatabaseServiceName: "pg-good-boys",
            DatabaseServiceName: "gcp-europe-west1",
            DatabaseServicePlan: "hobbyist",
            EnvVars: {},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(404)
        .expectField(null, 'Cluster with id 1 was not found')
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

addTest('create a plan', function (t) {
    return t.post('/v1/plans')
        .send({
            name: 'good plan',
            accountId: t.state.session.accountId,
            Replicas: 3,
            CPULimit: '100m',
            MemLimit: '100m',
            CPUReq: '10m',
            MemReq: '10m',
            ClusterID: t.state.clusterId,
            DatabaseServiceName: "pg-good-boys",
            DatabaseServiceName: "gcp-europe-west1",
            DatabaseServicePlan: "hobbyist",
            EnvVars: {},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(201)
        .store('planId', 'id')
})

addTest('fail to create a plan by the same name', function (t) {
    return t.post('/v1/plans')
        .send({
            name: 'good plan',
            accountId: t.state.session.accountId,
            Replicas: 3,
            CPULimit: '100m',
            MemLimit: '100m',
            CPUReq: '10m',
            MemReq: '10m',
            ClusterID: t.state.clusterId,
            DatabaseServiceName: "pg-good-boys",
            DatabaseServiceName: "gcp-europe-west1",
            DatabaseServicePlan: "hobbyist",
            EnvVars: {},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(409)
})

addTest('create another cluster', function (t) {
    return t.post('/v1/clusters')
        .send({
            name: 'europe-test2',
            accountId: t.state.session.accountId,
        })
        .expect(201)
        .store('clusterId2', 'id')
})

addTest('update the plan to switch over all properties', function (t) {
    return t.put(`/v1/plans/${t.state.planId}`)
        .send({
            name: 'good plan B',
            accountId: t.state.session.accountId, // this should not change
            Replicas: 4,
            CPULimit: '101m',
            MemLimit: '101m',
            CPUReq: '11m',
            MemReq: '11m',
            ClusterID: t.state.clusterId2,
            DatabaseServiceName: "pg-bad-boys",
            DatabaseServiceName: "gcp-europe-west2",
            DatabaseServicePlan: "pro",
            EnvVars: {A: "a", Num: 42},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(404)
})