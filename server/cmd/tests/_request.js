'use strict'

const assert = require('assert')
const dottie = require('dottie')

const SuperTest = require('supertest').Test

function Request(method, url, state) {
    const serverSocket = '127.0.0.1:8080'
    SuperTest.call(this, serverSocket, method, url)
    this._state = state
    this._shouldDebug = false
    this._debugMessage = null
}

const inherits = require('inherits')

inherits(Request, SuperTest)

Request.prototype.end = function (callback) {
    SuperTest.prototype.end.call(this, callback)

    this.on('response', res => {
        this._debugInfo = {
            request: this.req.method + ' ' + res.req.path,
            requestHeaders: this.req._headers,
            requestBody: this.req.text,
            statusCode: res.statusCode,
            responseHeaders: this.res._headers,
            responseBody: res.body,
        }
        if (this._debug) {
            console.log(this._debugInfo)
            console.log({ state: this._state })
        }
    })
    return this
}

Request.prototype.expect = function (statusCodeOrFunc) {
    if (typeof statusCodeOrFunc === 'function') {
        return SuperTest.prototype.expect.call(this, statusCodeOrFunc)
    } else {
        return SuperTest.prototype.expect.call(this, res => {
            if (res.statusCode != statusCodeOrFunc) {
                throw new Error(`Status ${res.statusCode} != ${statusCodeOrFunc}: ${this._debugInfo}`)
            }
        })
    }
}

Request.prototype.store = function (destKey, srcKey) {
    return this.expect(res => {
        let toSet = dottie.get(res.body, srcKey)
        if (toSet === undefined) {
            throw new Error(`store: path '${path}' not found in response: \n${res.body}`)
        }
        dottie.set(this._state, destKey, toSet)
    })
}

Request.prototype.expectLen = function (path, length) {
    this.expect(function counter(res) {
        var toCheck
        const checkEntireBody = path == null
        if (checkEntireBody) {
            toCheck = res.body
        } else {
            toCheck = dottie.get(res.body, path)
            if (toCheck === undefined) {
                throw new Error(`expectLen: path '${path}' not found in response: \n${res.body}`)
            }
        }
        if (toCheck.length != length) {
            throw new Error(`unexpected length at path '${path}' - expected ${length} but found ${toCheck.length}`)
        }
    })
    return this
}

Request.prototype.expectField = function (path, value) {
    this.expect(function counter(res) {
        var toCheck
        const checkEntireBody = path == null
        if (checkEntireBody) {
            toCheck = res.body
        } else {
            toCheck = dottie.get(res.body, path)
            if (toCheck === undefined) {
                throw new Error(`expectField: path '${path}' not found in response: \n${res.body}`)
            }
        }
        if (toCheck !== value) {
            throw new Error(`unexpected value at path '${path}' -\nexpected:\n\t${value}\nbut found:\n\t${toCheck}`)
        }
    })
    return this
}

Request.prototype.debug = function () {
    this._debug = true
    return this
}

module.exports = Request
