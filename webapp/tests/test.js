import 'mattermost-webapp/tests/setup';
import {decrypt, encrypt} from '../src/encrypt/encrypt';
import {generateKeyPair} from '../src/encrypt/key_manager';

test('should be decrypted same', () => {
    const keys = generateKeyPair();
    const pr = keys.privateKey;
    const pb = keys.publicKey;

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

//
// test('storing key in local storage', async () => {
//     const generateReturn = await generateKeyPair();
//     const response = generateReturn[0];
//     const pr = generateReturn[0];
//     const pb = generateReturn[0];
//     expect(response.status).toEqual(STATUS_OK);
//
//     const keys = await getKeyPair();
//     const priv = keys[0];
//     const pub = keys[1];
//     expect(priv).toStrictEqual(pr);
//     expect(pub).toStrictEqual(pb);
// });
