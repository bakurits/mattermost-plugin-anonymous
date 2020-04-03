import Client from '../api_client';

const eccrypto = require('eccrypto');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair() {
    const privateKey = eccrypto.generatePrivate();
    const publicKey = eccrypto.getPublic(privateKey);
    return [privateKey, publicKey];
}

// generates and stores private and public keys
export async function generateAndStoreKeypair() {
    const keys = generateKeyPair();
    return storeKeyPair(keys[0], keys[1]);
}

//store private key in a local storage
export async function storeKeyPair(privateKey, publicKey) {
    const pr = JSON.stringify(Array.from(privateKey));
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);
    return Client.storePublicKey(publicKey);
}

// eslint-disable-next-line no-unused-vars
export async function getKeyPair() {
    const privateKey = localStorage.getItem(LOCAL_STORAGE_KEY);

    //get public key from server
    const publicKey = await Client.retrievePublicKey();
    const pr = Buffer.from(JSON.parse(privateKey));
    return [pr, publicKey];
}