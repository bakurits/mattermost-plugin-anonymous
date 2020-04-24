/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt} from '../src/encrypt/encrypt';
import {
    generateKeyPair, loadKey, publicKeyToString,
    storePrivateKey,
} from '../src/encrypt/key_manager';

test('should be decrypted same', () => {
    const key = generateKeyPair();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = encrypt(key, test);

        const dec = decrypt(key, enc);

        expect(dec).toStrictEqual(test);
    });
});

test('storing key in local storage', () => {
    const key1 = generateKeyPair();
    storePrivateKey(key1);
    const key2 = loadKey();
    const s1 = publicKeyToString(key1);
    const s2 = publicKeyToString(key2);
    expect(s1).toStrictEqual(s2);
});

test('should be decrypted same with stored keys', () => {
    const key = generateKeyPair();
    storePrivateKey(key);

    const storedKey = loadKey();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = encrypt(key, test);
        const dec = decrypt(storedKey, enc);
        expect(dec).toStrictEqual(test);
    });
});
