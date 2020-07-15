const CryptoJS = require('crypto-js');

const AES_KEY_SIZE = 20;

/**
 * @param {string} data message that needs to be encrypted
 * @returns {{message: string, key: string}} returns encrypted data and randomly generated key
 */
export function encrypt(data) {
    const key = CryptoJS.lib.WordArray.random(AES_KEY_SIZE).toString();
    return {
        message: CryptoJS.AES.encrypt(data, key).toString(),
        key,
    };
}

/**
 * @param {string} data encrypted data
 * @param {string} key encryption key
 * @returns {string} returns original data before encryption
 */
export function decrypt(data, key) {
    const bytes = CryptoJS.AES.decrypt(data, key);
    return bytes.toString(CryptoJS.enc.Utf8);
}
