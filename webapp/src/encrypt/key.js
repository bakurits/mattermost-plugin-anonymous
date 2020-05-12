const NodeRSA = require('node-rsa');

export class Key {
    /**
     *
     * @param {NodeRSA | null} publicKey NodeRSA object for public key
     * @param {NodeRSA | null} privateKey NodeRSA object for public key
     */
    constructor(publicKey, privateKey) {
        this.publicKey = publicKey;
        this.privateKey = privateKey;
    }

    /**
     * @returns base64 string of public key
     */
    get PublicKey() {
        if (this.publicKey != null) {
            return this.publicKey.exportKey('public');
        }
        if (this.privateKey != null) {
            return this.privateKey.exportKey('public');
        }
        return null;
    }

    /**
     * @returns base64 string of private key
     */
    get PrivateKey() {
        if (this.privateKey == null) {
            return '';
        }
        return this.privateKey.exportKey('private');
    }

    /**
     *
     * @param {Buffer} data message text
     * @returns {string | null} returns encrypted text or null if public key isn't present
     */
    encrypt(data) {
        if (this.publicKey == null) {
            return null;
        }
        const buffer = Buffer.from(JSON.stringify(data));
        return this.publicKey.encrypt(buffer);
    }

    /**
     *
     * @param {Buffer} data encrypted text
     */
    decrypt(data) {
        if (this.privateKey == null) {
            return null;
        }
        try {
            return this.privateKey.decrypt(data, 'json');
        } catch (e) {
            return null;
        }
    }
}

/**
 * @param {string} publicKeyString NodeRSA object for public key
 * @returns {Key | null} returns new key object generated from private key
 */
export function newFromPublicKey(publicKeyString) {
    const publicKey = keyFromString(publicKeyString);
    if (!publicKey) {
        return null;
    }
    return new Key(publicKey, null);
}

/**
 * @param {string} privateKeyString NodeRSA object for private key
 * @returns {Key | null} returns new key object generated from private key
 */
export function newFromPrivateKey(privateKeyString) {
    const privateKey = keyFromString(privateKeyString);
    if (!privateKey) {
        return null;
    }
    return new Key(null, privateKey);
}

/**
 * @param {string} keyString is string with key data
 * @returns {NodeRSA | null} returns NodeRSA from key string
 */
function keyFromString(keyString) {
    try {
        return new NodeRSA(keyString);
    } catch (e) {
        return null;
    }
}
