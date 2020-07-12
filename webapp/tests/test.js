/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {Key} from '../src/encrypt/key';
import {decrypt, encrypt} from '../src/encrypt/aes';

import {generateKeyPair, loadKeyFromLocalStorage, storePrivateKey} from '../src/encrypt/key_manager';
import LFUCache from '../src/cache/lfu_cache';

test('should be decrypted same', () => {
    const key = generateKeyPair();
    const keyPrivate = new Key(null, key);
    const keyPublic = new Key(key, null);

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyPrivate.decrypt(enc);

        expect(dec).toStrictEqual(test);
    });
});

test('storing key in local storage', () => {
    const key1 = generateKeyPair();
    const privateKey = new Key(null, key1);
    storePrivateKey(privateKey);
    const s1 = privateKey.PrivateKey;
    const s2 = loadKeyFromLocalStorage().PrivateKey;
    expect(s1).toStrictEqual(s2);
});

test('should be decrypted same with stored keys', () => {
    const key = generateKeyPair();
    const privateKey = new Key(null, key);
    storePrivateKey(privateKey);

    const keyPublic = new Key(key, null);
    const keyLocalStorage = loadKeyFromLocalStorage();

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyLocalStorage.decrypt(enc);

        expect(dec).toStrictEqual(test);
    });
});

test('should not be encrypted with private key', () => {
    const key = generateKeyPair();
    const keyPrivate = new Key(null, key);

    const keyLocalStorage = loadKeyFromLocalStorage();

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
    const keyPublic = new Key(key, null);

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        const enc = keyPublic.encrypt(test);

        const dec = keyPublic.decrypt(enc);

        expect(dec).toEqual(null);
    });
});

test('aes tests', () => {
    const testsInput = ['[1, 3, 123]', '{key: \'value\'}', 'bakuri'];
    testsInput.forEach((test) => {
        const data = encrypt(test);
        const g = decrypt(data.message, data.key);
        expect(g).toEqual(test);
    });
});

test('cache test', () => {
    const cache = new LFUCache(2);
    cache.put('key1', 'value1');
    cache.put('key2', 'value2');
    expect(cache.get('key1')).toEqual('value1');
    expect(cache.get('key2')).toEqual('value2');
    cache.put('key3', 'value3');
    expect(cache.get('key2')).toEqual('value2');
    expect(cache.get('key3')).toEqual('value3');
    expect(cache.get('key1')).toBeNull();
    expect(cache.getLFU()).toEqual('key3');
});
