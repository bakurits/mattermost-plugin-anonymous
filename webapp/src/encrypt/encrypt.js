const crypto = require('crypto');

// encrypts message with public key and executes callback with encrypted message
export function encrypt(publicKey, message) {
    return crypto.publicEncrypt(publicKey, Buffer.from(JSON.stringify(message)));
}

// decrypts cipher text and executes callback with plain text
export function decrypt(privateKey, encrypted) {
    return JSON.parse(crypto.privateDecrypt(privateKey, encrypted).toString());
}

