'use strict'

const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'


function randomEmailAddress() {
    `${randomString}@${randomString()}.com`
}

function randomString(stringLength = 10) {
    let str = []
    for (let s = 0; s < stringLength; s++) {
        str.push(randomElement(chars))
    }
    return str.join('')
}

function randomElement(arr) {
    return arr[Math.floor(Math.random() * arr.length)]
}