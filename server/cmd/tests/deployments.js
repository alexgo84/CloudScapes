'use strict'

const addTest = require('./_test').init()

addTest('create a new account', function (t) {
    return t.createAccount()
        .expect(201)
})


addTest('get all deployments in account should return 0 deployments', function (t) {
    return t.get('/v1/deployments')
        .expect(200)
        .expectLen(null, 0)
})

addTest('fail to create a deployment that references a non existing plan', function (t) {
    return t.post('/v1/deployments')
        .send({
            name: 'some deployment',
            accountId: t.state.session.accountId,
            Replicase: 3,
            CPULimit: '100m',
            MemLimit: '100m',
            CPUReq: '10m',
            MemReq: '10m',
            PlanID: 0, // violates foreign key constraint
            DatabaseServiceName: "pg-good-boys",
            DatabaseServiceName: "gcp-europe-west1",
            DatabaseServicePlan: "hobbyist",
            EnvVars: {},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(404)
        .expectField(null, 'Plan with id 0 was not found')
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

addTest('create a deployment', function (t) {
    return t.post('/v1/deployments')
        .send({
            name: 'some deployment',
            accountId: t.state.session.accountId,
            Replicas: 3,
            CPULimit: '100m',
            MemLimit: '100m',
            CPUReq: '10m',
            MemReq: '10m',
            planId: t.state.planId,
            DatabaseServiceName: "pg-good-boys",
            DatabaseServiceName: "gcp-europe-west1",
            DatabaseServicePlan: "hobbyist",
            EnvVars: {},
            CronJobs: [],
            ConfigMaps: [],
        })
        .expect(201)
        .store('deploymentId', 'id')
})

addTest('fail to create a deployment by the same name', function (t) {
    return t.post('/v1/deployments')
        .send({
            name: 'some deployment',
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
        .expectField(null, `Key (accountid, name)=(${t.state.session.accountId}, ${'some deployment'}) already exists.`)
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

addTest('update the deployment to switch over all properties', function (t) {
    const updatedPlan = {
        name: 'some deployment B',
        Replicas: 4,
        CPULimit: '101m',
        MemLimit: '101m',
        CPUReq: '11m',
        MemReq: '11m',
        ClusterID: t.state.clusterId2,
        DatabaseServiceName: "pg-bad-boys",
        DatabaseServiceCloud: "gcp-europe-west2",
        DatabaseServicePlan: "pro",
        EnvVars: {A: "a", Num: 42},
        CronJobs: [],
        ConfigMaps: [],
    }
    return t.put(`/v1/deployments/${t.state.deploymentId}`)
        .send(updatedPlan)
        .expect(200)
        .expectField('id', t.state.deploymentId)
        .expectField('accountId', t.state.session.accountId)
        .expectField('CPULimit', updatedPlan.CPULimit)
        .expectField('memLimit', updatedPlan.MemLimit)
        .expectField('CPUReq', updatedPlan.CPUReq)
        .expectField('memReq', updatedPlan.MemReq)
        .expectField('clusterID', updatedPlan.ClusterID)
        .expectField('databaseServiceName', updatedPlan.DatabaseServiceName)
        .expectField('databaseServiceCloud', updatedPlan.DatabaseServiceCloud)
        .expectField('databaseServicePlan', updatedPlan.DatabaseServicePlan)
        .expectField('envVars.A', updatedPlan.EnvVars.A)
        .expectField('envVars.Num', updatedPlan.EnvVars.Num)
        .expectLen('configMaps', 0)
        .expectLen('cronJobs', 0)
})

addTest('get all deployments in account should return only one deployment', function (t) {
    return t.get('/v1/deployments')
        .expect(200)
        .expectLen(null, 1)
})

addTest('delete first cluster', function (t) {
    return t.delete(`/v1/deployments/${t.state.deploymentId}`)
        .expect(204)
})


addTest('get all deployments in account should return zero deployments', function (t) {
    return t.get('/v1/deployments')
        .expect(200)
        .expectLen(null, 0)
})