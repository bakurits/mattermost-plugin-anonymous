/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt} from '../src/encrypt/encrypt';
import {
    generateKeyPair,
    getKeyPair,
    getPublicKeyFromPrivateKey,
    storePrivateKey,
} from '../src/encrypt/key_manager';

test('should be decrypted same', () => {
    const keys = generateKeyPair();
    const pr = keys.privateKey;
    const pb = keys.publicKey;

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = encrypt(pb, test);
        const dec = decrypt(pr, enc);
        expect(dec).toStrictEqual(test);
    });
});

test('storing key in local storage', () => {
    const keys1 = generateKeyPair();
    storePrivateKey(keys1.privateKey);
    const keys2 = getKeyPair();
    expect(keys1.publicKey).toStrictEqual(keys2.publicKey);
});

test('test get public key from private', () => {
    const keys = generateKeyPair();
    const a = getPublicKeyFromPrivateKey(keys.privateKey);
    expect(a).toStrictEqual(keys.publicKey);
});

test('should be decrypted same with stored keys', () => {
    const keys = generateKeyPair();
    storePrivateKey(keys.privateKey);

    const storedKeys = getKeyPair();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = encrypt(keys.publicKey, test);
        const dec = decrypt(storedKeys.privateKey, enc);
        expect(dec).toStrictEqual(test);
    });
});
