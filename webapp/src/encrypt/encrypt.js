const eccrypto = require('eccrypto');

// encrypts message with public key and executes callback with encrypted message
export function encrypt(publicKey, message) {
    return eccrypto.encrypt(publicKey, Buffer.from(JSON.stringify(message)), {});
}

// decrypts cipher text and executes callback with plain text
export async function decrypt(privateKey, encrypted) {
    return eccrypto.decrypt(privateKey, encrypted);
}

