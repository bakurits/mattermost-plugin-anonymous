/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {newFromPrivateKey, newFromPublicKey, loadFromLocalStorage} from '../src/encrypt/key';

import {
    generateKeyPair,
    storePrivateKey,
} from '../src/encrypt/key_manager';

test('should be decrypted same', () => {
    const key = generateKeyPair();
    const keyPrivate = newFromPrivateKey(key);
    const keyPublic = newFromPublicKey(key);

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyPrivate.decrypt(enc);

        expect(dec).toStrictEqual(test);
    });
});

test('storing key in local storage', () => {
    const key1 = generateKeyPair();
    var privateKey = newFromPrivateKey(key1);
    storePrivateKey(privateKey);
    const s1 = privateKey.PrivateKey;
    const s2 = loadFromLocalStorage().PrivateKey;
    expect(s1).toStrictEqual(s2);
});

test('should be decrypted same with stored keys', () => {
    const key = generateKeyPair();
    var privateKey = newFromPrivateKey(key);
    storePrivateKey(privateKey);

    const keyPublic = newFromPublicKey(key);
    const keyLocalStorage = loadFromLocalStorage();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyLocalStorage.decrypt(enc);

        expect(dec).toStrictEqual(test);
    });
});

test('should not be encrypted with private key', () => {
    const key = generateKeyPair();
    const keyPrivate = newFromPrivateKey(key);

    const keyLocalStorage = loadFromLocalStorage();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPrivate.encrypt(test);

        const dec = keyPrivate.decrypt(enc);

        expect(dec).toEqual(null);
    });

    testsInput.forEach((test) => {
        const enc = keyLocalStorage.encrypt(test);

        const dec = keyLocalStorage.decrypt(enc);

        expect(dec).toEqual(null);
    });
});

test('should not be decrypted with public key', () => {
    const key = generateKeyPair();
    const keyPublic = newFromPublicKey(key);

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyPublic.decrypt(enc);

        expect(dec).toEqual(null);
    });
});