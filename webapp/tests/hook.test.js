
import 'mattermost-webapp/tests/setup';
import {generateKeyPair, LOCAL_STORAGE_KEY} from '../src/encrypt/key_manager';
import {Key} from '../src/encrypt/key';

import Hooks from '../src/hook/hook';

import Client4 from '../__mocks__/Client4';
import Client from '../__mocks__/Client';

const hooks = new Hooks(null, null, Client4, Client);

test('hook encrypt/decrypt test', () => {
    const key = generateKeyPair();
    const publicKey = new Key(key, null).PublicKey;
    const privateKey = new Key(null, key).PrivateKey;
    localStorage.setItem(LOCAL_STORAGE_KEY, privateKey);

    Client.retrievePublicKey.mockImplementationOnce(() =>
        Promise.resolve({public_keys: [publicKey]})
    );

    const message = 'message';

    hooks.encryptMessage('someChannel', message).then(
        (enc) => {
            expect(enc.success).toStrictEqual(true);
            const post = {
                message: enc.message,
            };

            const dec = hooks.decryptMessage(post);

            expect(dec).toStrictEqual(message);
        }
    );
});
