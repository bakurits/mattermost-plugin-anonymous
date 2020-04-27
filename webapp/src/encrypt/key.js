import {loadKey} from './key_manager'

export class Key {
    /**
     *
     * @param {string | null} publicKey base64 string of public key or null
     * @param {string | null} privateKey base64 string of private key or null
     */
    constructor(publicKey, privateKey) {
        this.publicKey = publicKey;
        this.privateKey = privateKey;
    }

    /**
     * @returns base64 string of public key
     */
    get PublicKey() {
        if (this.publicKey == null)
            return this.privateKey.exportKey('public');
        return this.publicKey.exportKey('public');
    }

    /**
     * @returns base64 string of private key
     */
    get PrivateKey() {
        if (this.privateKey == null)
            return '';
        return this.privateKey.exportKey('private');
    }

    /**
     *
     * @param {string} data message text
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
     * @param {string} data encrypted text
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

export function newFromPublicKey(publicKey) {
    return new Key(publicKey, null);
}

export function newFromPrivateKey(privateKey) {
    return new Key(null, privateKey);
}

export function loadFromLocalStorage() {
    const privateKey = loadKey();
    if (!privateKey) {
        return null;
    }
    return new Key(null, privateKey);
}