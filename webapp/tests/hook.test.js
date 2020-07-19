import 'mattermost-webapp/tests/setup';
import {generateKeyPair, LOCAL_STORAGE_KEY} from '../src/encrypt/key_manager';
import {Key} from '../src/encrypt/key';

import Hooks from '../src/hook/hook';

import Client from '../__mocks__/Client';
import Store from '../__mocks__/Store';

const hooks = new Hooks(Store, null, Client);

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

test('hook handlePostCommand error test', () => {
    Client.retrievePublicKey.mockImplementationOnce(() =>
        Promise.resolve({public_keys: null})
    );

    const commands = ['test', 'command'];

    hooks.handlePostCommand(commands, {}).then(
        (result) => {
            expect(result.error).toEqual('could not encrypt message properly');
        }
    );
});

test('hook handlePostCommand success test', () => {
    Client.retrievePublicKey.mockImplementationOnce(() =>
        Promise.resolve({public_keys: []})
    );

    const commands = ['test', 'command'];

    hooks.handlePostCommand(commands, {}).then(
        (result) => {
            expect(result).toEqual({});
        }
    );
});

test('hook messageWillBePostedHook test', () => {
    const post = {channel_id: '123', message: 'test message'};
    const key = generateKeyPair();
    const publicKey = new Key(key, null).PublicKey;
    const privateKey = new Key(null, key).PrivateKey;
    localStorage.setItem(LOCAL_STORAGE_KEY, privateKey);

    Client.retrievePublicKey.mockImplementationOnce(() =>
        Promise.resolve({public_keys: [publicKey]})
    );

    const props = {encrypted: true};

    hooks.messageWillBePostedHook(post).then(
        (result) => {
            expect(result.post.channel_id).toEqual('123');
            expect(result.post.props).toEqual(props);
        }
    );
});
