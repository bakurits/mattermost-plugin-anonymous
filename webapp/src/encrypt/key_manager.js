import Client from '../api_client';

const eccrypto = require('eccrypto');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair() {
    const privateKey = eccrypto.generatePrivate();
    const publicKey = getPublicKeyFromPrivateKey(privateKey);
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
    return storeKeyPair(keys.privateKey, keys.privateKey);
}

//store private key in a local storage
export async function storeKeyPair(privateKey, publicKey) {
    const pr = privateKey.toString('base64');
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);
    return Client.storePublicKey(publicKey);
}

// get getKeyPair returns key pair
export function getKeyPair() {
    const pr = localStorage.getItem(LOCAL_STORAGE_KEY);
    const privateKey = Buffer.from(pr, 'base64');
    const publicKey = getPublicKeyFromPrivateKey(privateKey);
    if (privateKey === null) {
        return null;
    }
    return {privateKey, publicKey};
}

// generate public key using private key (buffer)
export function getPublicKeyFromPrivateKey(privateKey) {
    try {
        return eccrypto.getPublic(privateKey);
    } catch (e) {
        return null;
    }
}

// only store private key (buffer)
export function storePrivateKey(privateKey) {
    const pr = privateKey.toString('base64');
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);
}

// get private key
export function getPrivateKey() {
    const privateKey = localStorage.getItem(LOCAL_STORAGE_KEY);
    return Buffer.from(privateKey, 'base64');
}
