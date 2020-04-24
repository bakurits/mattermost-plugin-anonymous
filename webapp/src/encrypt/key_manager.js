import Client from '../api_client';

const NodeRSA = require('node-rsa');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';
const RSA_KEY_SIZE = 512;

export function privateKeyToString(key) {
    return key.exportKey('private');
}

export function publicKeyToString(key) {
    return key.exportKey('public');
}

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair() {
    return new NodeRSA({b: RSA_KEY_SIZE});
}

// generates and stores private and public keys
export async function generateAndStoreKeyPair() {
    const key = generateKeyPair();
    return storeKeyPair(key);
}

//store private key in a local storage
export async function storeKeyPair(key) {
    storePrivateKey(key);
    return Client.storePublicKey(publicKeyToString(key));
}

// only store private key (buffer)
export function storePrivateKey(key) {
    localStorage.setItem(LOCAL_STORAGE_KEY, privateKeyToString(key));
}

// get private key
export function loadKey() {
    const keyData = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (keyData) {
        return new NodeRSA(keyData);
    }

    return null;
}
