const eccrypto = require('eccrypto');

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

