import Client from '../api_client';

import {newFromPrivateKey, Key} from './key';

const NodeRSA = require('node-rsa');
export const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';
const RSA_KEY_SIZE = 512;

/**
 * generates ECIES private, public key pair and executes with callback
 * @returns {NodeRSA} returns newly generated NodeRSA object
 */
export function generateKeyPair() {
    const key = new NodeRSA({b: RSA_KEY_SIZE});
    return key.generateKeyPair();
}

/**
 * generates and stores private and public keys
 * @returns {Object} returns response from api call
 */
export async function generateAndStoreKeyPair() {
    const key = generateKeyPair();
    if (!key) {
        return null;
    }
    return storeKeyPair(key);
}

/**
 * store private key in a local storage
 * @param {NodeRSA} key is nodeRSA key object
 * @returns {Object} returns newly generated NodeRSA object
 */
async function storeKeyPair(key) {
    const privateKey = new Key(null, key);
    storePrivateKey(privateKey);
    const publicKey = new Key(key, null);
    return Client.storePublicKey(publicKey);
}

/**
 *
 * @param {Key | null} key object of Key
 */
export function storePrivateKey(key) {
    if (key === null || key.PrivateKey === null) {
        return;
    }
    localStorage.setItem(LOCAL_STORAGE_KEY, key.PrivateKey);
}

/**
 *
 * @returns {Key | null} returns new key object loaded from localstorage
 * or null if localstorage is empty
 */
export function loadFromLocalStorage() {
    const keyData = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (!keyData) {
        return null;
    }
    const privateKey = newFromPrivateKey(keyData);
    if (!privateKey) {
        return null;
    }
    return privateKey;
}
