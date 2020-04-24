import {Client4} from 'mattermost-redux/client';

import {
    generateAndStoreKeyPair,
    keyFromString, loadKey, publicKeyToString, privateKeyToString,
    storePrivateKey,
} from '../encrypt/key_manager';
import {sendEphemeralPost} from '../actions/actions';
import Client from '../api_client';
import {decrypt, encrypt} from '../encrypt/encrypt';

export default class Hooks {
    constructor(store, settings) {
        this.store = store;
        this.settings = settings;
    }

    handleKeyPair = async (commands, args) => {
        let key;
        let response;
        switch (commands[0]) {
        case '--generate':
            response = await generateAndStoreKeyPair();
            if (response.status !== 'OK') {
                return Promise.resolve({error: {message: 'Error occurred while trying to store key on a server'}});
            }
            this.store.dispatch(sendEphemeralPost('Your keys were successfully generated and stored', args.channel_id));
            return Promise.resolve({});

        case '--overwrite':
            if (commands.length < 2) {
                return Promise.resolve({error: {message: "Private key isn't specified"}});
            }
            if (commands.length > 2) {
                return Promise.resolve({error: {message: 'Too many arguments'}});
            }

            key = keyFromString(commands[1]);
            if (!key) {
                return Promise.resolve({error: {message: 'Invalid private key'}});
            }

            storePrivateKey(key);
            return Promise.resolve({message: '/anonymous keypair --overwrite ' + publicKeyToString(key), args});

        case '--export':
            key = loadKey();
            this.store.dispatch(sendEphemeralPost('Your private key:\n' + privateKeyToString(key), args.channel_id));
            return Promise.resolve({});
        default:
            break;
        }
        return Promise.resolve({});
    };

    handlePost = async (commands, args) => {
        const users = await Client4.getProfilesInChannel(args.channel_id);
        // eslint-disable-next-line no-console
        console.log(users);

        const publicKeys = await Promise.all(users.map((user) => {
            return Client.retrievePublicKey(user.id).then((data) => {
                return Buffer.from(data.public_key, 'base64').toString();
            });
        }));
        // eslint-disable-next-line no-console
        console.log(publicKeys);

        const encrypted = await Promise.all(publicKeys.map((keyString) => {
            const key = keyFromString(keyString);
            // eslint-disable-next-line no-console
            console.log(publicKeyToString(key));
            const message = encrypt(key, commands[0]).toString('base64');
            return {
                message,
                public_key: keyString,
            };
        }));
        // eslint-disable-next-line no-console
        console.log(encrypted);

        const result = await Promise.all(encrypted.map((messageData) => {
            const {message} = messageData;
            // eslint-disable-next-line no-console
            return Client4.createPost({
                channel_id: args.channel_id,
                message,
                props: {public_key: messageData.public_key},
            });
        }));

        // eslint-disable-next-line no-console
        console.log(result);

        // eslint-disable-next-line no-console
        //this.store.dispatch(sendPost('Your private key:\n', args.channel_id));

        //await Client.sendPost(args.channel_id, commands[0], getPublicKeyFromPrivateKey(getPrivateKey()));
        return Promise.resolve({});
    };

    slashCommandWillBePostedHook = (message, contextArgs) => {
        const commands = message.split(/(\s+)/).filter((e) => e.trim().length > 0);

        if (commands[0] !== '/anonymous') {
            return Promise.resolve({});
        }
        if (commands.length < 2) {
            return Promise.resolve({error: {message: "Command isn't specified"}});
        }

        switch (commands[1]) {
        case 'keypair':
            return this.handleKeyPair(commands.splice(2), contextArgs);
        case 'a':
            return this.handlePost(commands.splice(2), contextArgs);
        default:
            break;
        }

        return Promise.resolve({message, args: contextArgs});
    }

    messageWillFormatHook = (post) => {
        // message text in database
        const {message} = post;

        // here is public key if needed
        // eslint-disable-next-line no-unused-vars
        const {props} = post;

        // eslint-disable-next-line no-console
        console.log(message);
        // eslint-disable-next-line no-console
        console.log(props);

        if (!props.public_key) {
            return message;
        }

        const key = loadKey();
        // eslint-disable-next-line no-console
        console.log(key);

        if (props.public_key !== publicKeyToString(key)) {
            return '';
        }

        const res = decrypt(key, Buffer.from(message, 'base64'));
        // eslint-disable-next-line no-console
        console.log('jjj');
        // eslint-disable-next-line no-console
        console.log(res);

        return res;
    }
}
