'use strict'

const xtend = require('xtend')
const faker = require('faker')

module.exports =
{
    createAccount: function (suffix = '', params = {}) {
        return this.post('/v1/accounts')
            .send(
                xtend(
                    {
                        companyName: faker.company.companyName(),
                        email: faker.internet.email(),
                        password: faker.internet.password()
                    },
                    params
                )
            )
            .store('session' + suffix, [])
    }
}