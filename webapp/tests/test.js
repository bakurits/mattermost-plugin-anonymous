import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt, generateKeyPair} from '../src/encrypt/encrypt';

test('should be decrypted same', () => {
    generateKeyPair((privateKey, publicKey) => {
        const pb = publicKey;
        const pr = privateKey;

        // eslint-disable-next-line no-magic-numbers
        const testsInput = [[1, 3, 123], {key: 'value'}, 'bakuri'];

        // eslint-disable-next-line max-nested-callbacks
        testsInput.forEach((test) => {
            // eslint-disable-next-line max-nested-callbacks
            encrypt(pb, test, (encrypted) => {
                // eslint-disable-next-line max-nested-callbacks
                decrypt(pr, encrypted, (decrypted) => {
                    expect(decrypted).toStrictEqual(test);
                });
            });
        });
    });
});
