import Client from '../api_client';

const eccrypto = require('eccrypto');
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';

// generates ECIES private, public key pair and executes with callback
export function generateKeyPair(callback) {
    const privateKey = eccrypto.generatePrivate();
    const publicKey = eccrypto.getPublic(privateKey);
    callback(privateKey, publicKey);
}

//store private key in a local storage
export function storeKeyPair(privateKey, publicKey, callback) {
    const pr = JSON.stringify(Array.from(privateKey));
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);
    Client.storePublicKey(publicKey).then((response) => {
        callback(response);
    });
}

// eslint-disable-next-line no-unused-vars
export function getKeyPair(callback) {
    const privateKey = localStorage.getItem(LOCAL_STORAGE_KEY);

    //get public key from server
    Client.retrievePublicKey().then((publicKey) => {
        const pr = Buffer.from(JSON.parse(privateKey));
        callback(pr, publicKey);
    });
}
