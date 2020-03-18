
const LOCAL_STORAGE_KEY = 'anonymous_plugin_private_key';

//store private key in a local storage
export function storeKeypair(privateKey, publicKey, callback) {
    const pr = JSON.stringify(Array.from(privateKey));
    const pb = JSON.stringify(Array.from(publicKey));
    localStorage.setItem(LOCAL_STORAGE_KEY, pr);

    //post a request to server to store public key

    callback(0);
}

// eslint-disable-next-line no-unused-vars
export function getKeypair(callback) {
    const privateKey = localStorage.getItem(LOCAL_STORAGE_KEY);

    //get public key from server
    const publicKey = '[]';
    const pr = Buffer.from(JSON.parse(privateKey));
    const pb = Buffer.from(JSON.parse(publicKey));
    callback(pr, pb);
}
