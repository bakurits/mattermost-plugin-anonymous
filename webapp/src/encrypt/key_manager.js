import Client from '../api_client';

const eccrypto = require('eccrypto');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';

// generates ECIES private, public key pair and executes with callback
export async function generateKeyPair() {
    const privateKey = eccrypto.generatePrivate();
    const publicKey = eccrypto.getPublic(privateKey);
    return [privateKey, publicKey];
}

//store private key in a local storage
export async function storeKeyPair(privateKey, publicKey) {
    const pr = JSON.stringify(Array.from(privateKey));
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);
    const response = await Client.storePublicKey(publicKey)
    return response
}

// eslint-disable-next-line no-unused-vars
export async function getKeyPair() {
    const privateKey = localStorage.getItem(LOCAL_STORAGE_KEY);

    //get public key from server
    const publicKey = await Client.retrievePublicKey()
    const pr = Buffer.from(JSON.parse(privateKey));
    return [pr, publicKey]
}
