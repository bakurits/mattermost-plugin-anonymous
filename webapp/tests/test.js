/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt} from '../src/encrypt/encrypt';
import {generateKeyPair, getKeyPair, storePrivateKey} from '../src/encrypt/key_manager';

test('should be decrypted same', () => {
    const keys = generateKeyPair();
    const pr = keys.privateKey;
    const pb = keys.publicKey;

    const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];
    testsInput.forEach((test) => {
        encrypt(pb, test, (encrypted) => {
            decrypt(pr, encrypted, (decrypted) => {
                expect(decrypted).toStrictEqual(test);
            });
        });
    });
});

test('storing key in local storage', () => {
    const keys1 = generateKeyPair();
    storePrivateKey(keys1.privateKey);
    const keys2 = getKeyPair();
    expect(keys1.publicKey).toStrictEqual(keys2.publicKey);
    expect(keys1.privateKey).toStrictEqual(keys2.privateKey);
});
