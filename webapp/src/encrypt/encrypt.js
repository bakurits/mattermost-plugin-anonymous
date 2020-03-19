const eccrypto = require('eccrypto');

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair(callback) {
    const privateKey = eccrypto.generatePrivate();
    const publicKey = eccrypto.getPublic(privateKey);
    callback(privateKey, publicKey);
}

// encrypts message with public key and executes callback with encrypted message
export function encrypt(publicKey, message, callback) {
    eccrypto.encrypt(publicKey, Buffer.from(JSON.stringify(message)), {}).then((encrypted) => {
        callback(encrypted);
    });
}

// decrypts cipher text and executes callback with plain text
export function decrypt(privateKey, encrypted, callback) {
    eccrypto.decrypt(privateKey, encrypted).then((plaintext) => {
        callback(JSON.parse(plaintext));
    });
}

