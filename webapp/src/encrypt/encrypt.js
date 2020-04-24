// encrypts message with public key and executes callback with encrypted message
export function encrypt(key, message) {
    return key.encrypt(Buffer.from(JSON.stringify(message)));
}

// decrypts cipher text and executes callback with plain text
export function decrypt(key, encrypted) {
    return key.decrypt(encrypted, 'json');
}

