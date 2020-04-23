import Client from '../api_client';

const crypto = require('crypto');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';
const RSA_KEY_SIZE = 4096;

export function privateKeyToString(privateKey) {
    return privateKey.export({
        type: 'pkcs1',
        format: 'pem',
    });
}

export function publicKeyToString(publicKey) {
    return publicKey.export({
        type: 'pkcs1',
        format: 'pem',
    });
}

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair() {
    const {publicKey, privateKey} = crypto.generateKeyPairSync('rsa', {
        modulusLength: RSA_KEY_SIZE,
    });
    if (publicKey === null) {
        return null;
    }
    return {privateKey, publicKey};
}

// generates and stores private and public keys
export async function generateAndStoreKeyPair() {
    const keys = generateKeyPair();

    if (keys === null) {
        return null;
    }
    return storeKeyPair(keys.privateKey, keys.publicKey);
}

//store private key in a local storage
export async function storeKeyPair(privateKey, publicKey) {
    storePrivateKey(privateKey);
    return Client.storePublicKey(publicKey);
}

// get getKeyPair returns key pair
export function getKeyPair() {
    const privateKey = getPrivateKey();
    const publicKey = getPublicKeyFromPrivateKey(privateKey);
    if (privateKey === null) {
        return null;
    }
    return {privateKey, publicKey};
}

// generate public key using private key (buffer)
export function getPublicKeyFromPrivateKey(privateKey) {
    try {
        return crypto.createPublicKey(privateKey);
    } catch (e) {
        return null;
    }
}

// only store private key (buffer)
export function storePrivateKey(privateKey) {
    localStorage.setItem(LOCAL_STORAGE_KEY, privateKeyToString(privateKey));
}

// get private key
export function getPrivateKey() {
    return crypto.createPrivateKey(localStorage.getItem(LOCAL_STORAGE_KEY));
}
