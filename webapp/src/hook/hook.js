import {Client4} from 'mattermost-redux/client';

import {
    generateAndStoreKeyPair,
    getPrivateKey,
    getPublicKeyFromPrivateKey,
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
        let privateKey;
        let publicKey;
        let response;
        let pubKeyString;
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
            privateKey = Buffer.from(commands[1], 'base64');
            publicKey = getPublicKeyFromPrivateKey(privateKey);
            if (!publicKey) {
                return Promise.resolve({error: {message: 'Invalid private key'}});
            }

            pubKeyString = publicKey.toString('base64');
            storePrivateKey(privateKey);
            return Promise.resolve({message: '/anonymous keypair --overwrite ' + pubKeyString, args});

        case '--export':
            privateKey = getPrivateKey();
            this.store.dispatch(sendEphemeralPost('Your private key:\n' + privateKey.toString('base64'), args.channel_id));
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
            return Client.retrievePublicKey(user.id);
        }));
        // eslint-disable-next-line no-console
        console.log(publicKeys);

        const encrypted = await Promise.all(publicKeys.map((publicKey) => {
            // eslint-disable-next-line no-console
            console.log(Buffer.from(publicKey.public_key, 'base64'));
            return encrypt(Buffer.from(publicKey.public_key, 'base64'), commands[0]).then((data) => {
                return {
                    data,
                    public_key: publicKey.public_key,
                };
            });
        }));
        // eslint-disable-next-line no-console
        console.log(encrypted);

        const messages = encrypted.map((cypherObjectWithPublicKey) => {
            const {data} = cypherObjectWithPublicKey;
            const publicKey = cypherObjectWithPublicKey.public_key;
            const message = JSON.stringify(
                {
                    ciphertext: data.ciphertext.toString('base64'),
                    ephemPublicKey: data.ephemPublicKey.toString('base64'),
                    iv: data.iv.toString('base64'),
                    mac: data.mac.toString('base64'),
                }
            );

            return {
                message,
                public_key: publicKey,
            };
        });
        // eslint-disable-next-line no-console
        console.log(messages);

        const result = await Promise.all(messages.map((messageData) => {
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

        const privateKey = getPrivateKey();
        const publicKey = getPublicKeyFromPrivateKey(privateKey);

        // eslint-disable-next-line no-console
        console.log(privateKey);
        // eslint-disable-next-line no-console
        console.log(publicKey.toString('base64'));
        if (props.public_key !== publicKey.toString('base64')) {
            return '';
        }

        const messageJson = JSON.parse(message);

        const encrypted = {};
        encrypted.ciphertext = Buffer.from(messageJson.ciphertext, 'base64');

        // eslint-disable-next-line no-console
        console.log(Buffer.from(messageJson.ciphertext, 'base64'));

        encrypted.ephemPublicKey = Buffer.from(messageJson.ephemPublicKey, 'base64');
        encrypted.iv = Buffer.from(messageJson.iv, 'base64');
        // eslint-disable-next-line no-console
        console.log(Buffer.from(messageJson.iv, 'base64'));
        encrypted.mac = Buffer.from(messageJson.mac, 'base64');

        // eslint-disable-next-line no-console
        console.log(encrypted);
        return decrypt(privateKey, encrypted).then((plaintext) => {
            // eslint-disable-next-line no-console
            console.log(plaintext.toString());
            return plaintext.toString();
        });
    }
}
