'use strict'

const util = require('util')
const Request = require('./_request')
const AccountTest = require('./_account_test')

const bgRed = '\x1b[41m'
const bgMagenta = "\x1b[45m"
const bgCyan = "\x1b[46m"

const yellow = '\x1b[33m'
const bright = "\x1b[1m"
const black = "\x1b[30m"

const underscore = "\x1b[4m"
const resetColor = '\x1b[0m'

function Test(state) {
  this.state = state
}

// register requests execution methods
['put', 'post', 'get', 'delete'].forEach(function (method) {
  Test.prototype[method] = function (url) {
    return new Request(method, url, this.state)
  }
})

exports.init = function (state = {}) {
  // the test object that will be passed to test functions
  const t = new Test(state)
  const tests = []

  function testAdder(name, body) {
    tests.push([name, body])
  }

  // print a header and start executing the first test
  function runAddedTests() {
    console.log(`${bright}1..${tests.length}${resetColor}`)
    runTest(0)
  }

  // execute the test at index = idx. if there is no error execute the next tests until done
  function runTest(idx) {
    if (!tests[idx]) {
      return
    }
    const [description, fn] = tests[idx]
    fn(t).end(function (err, res) {
      if (err) {
        console.log(`${bgRed}${yellow}not ok:${resetColor} #${idx + 1}: ${description} ${err.stack.replace(/\n/g, '\n# ')}`)
        process.exit(1)
      }

      console.log(`ok: #${idx + 1} - ${description}`)
      runTest(idx + 1)
    })
  }

  // start execution in a bit (after tests are added) async
  // and return the add test function to the caller (so the tests may be added)  
  setTimeout(runAddedTests, 10)

  return testAdder
}

// add functions from other test helper bundles
Test.prototype.createAccount = AccountTest.createAccount