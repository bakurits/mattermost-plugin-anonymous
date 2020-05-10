import Client from '../api_client';

import {newFromPublicKey, newFromPrivateKey} from './key';

const NodeRSA = require('node-rsa');
export const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';
const RSA_KEY_SIZE = 512;

export function keyFromString(keyString) {
    return new NodeRSA(keyString);
}

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair() {
    const key = new NodeRSA({b: RSA_KEY_SIZE});
    return key.generateKeyPair();
}

// generates and stores private and public keys
export async function generateAndStoreKeyPair() {
    const key = generateKeyPair();
    if (!key) {
        return 'error';
    }
    return storeKeyPair(key);
}

//store private key in a local storage
export async function storeKeyPair(key) {
    var privateKey = newFromPrivateKey(key);
    storePrivateKey(privateKey);
    var publicKey = newFromPublicKey(key);
    return Client.storePublicKey(publicKey);
}

/**
 *
 * @param {Key} key object of Key
 */
export function storePrivateKey(key) {
    if (key === null || key.PrivateKey === null) {
        return;
    }
    localStorage.setItem(LOCAL_STORAGE_KEY, key.PrivateKey);
}

// get private key
export function loadKey() {
    const keyData = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (keyData) {
        return keyFromString(keyData);
    }

    return null;
}
