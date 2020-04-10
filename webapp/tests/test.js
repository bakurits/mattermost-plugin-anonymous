/* eslint-disable no-magic-numbers,max-nested-callbacks */
import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt} from '../src/encrypt/encrypt';
import {generateKeyPair, getKeyPair, getPublicKeyFromPrivateKey, storePrivateKey} from '../src/encrypt/key_manager';

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

test('test get public key from private', () => {
    const a = getPublicKeyFromPrivateKey('123');
    expect(a).toBeNull();

    const b = getPublicKeyFromPrivateKey('74PCP0wLKoBApiR1iOBiNGAE+WUScXr40bjWfUHtB8Y=');
    expect(b).toBeNull();

    const c = getPublicKeyFromPrivateKey(Buffer.from('74PCP0wLKoBApiR1iOBiNGAE+WUScXr40bjWfUHtB8Y=', 'base64'));
    expect(c).toStrictEqual(Buffer.from([4, 230, 209, 139, 9, 141, 217, 63, 192, 48, 95, 207, 193, 98, 113, 155, 196, 143, 20, 107, 69, 250, 237, 169, 144, 53, 122, 84, 76, 170, 18, 34, 43, 80, 124, 44, 75, 96, 95, 230, 43, 144, 157, 119, 90, 188, 98, 27, 60, 170, 99, 52, 162, 246, 100, 223, 160, 93, 58, 148, 40, 126, 65, 139, 120]));

    const d = getPublicKeyFromPrivateKey();
    expect(d).toBeNull();
});
