'use strict'

const xtend = require('xtend')
const faker = require('faker')

module.exports =
{
    createAccount: function (suffix = '', params = {}) {
        return this.post('/accounts')
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
            .expect(200)
            .store('session' + suffix, [])
    }
}