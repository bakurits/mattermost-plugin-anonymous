import {Client4} from 'mattermost-redux/client';

import {
    generateAndStoreKeyPair, loadFromLocalStorage,
    storePrivateKey,
} from '../encrypt/key_manager';
import {sendEphemeralPost} from '../actions/actions';
import Client from '../api_client';
import {newFromPublicKey, newFromPrivateKey} from '../encrypt/key';

export default class Hooks {
    constructor(store, settings) {
        this.store = store;
        this.settings = settings;
    }

    handleKeyPair = async (commands, args) => {
        let key;
        let response;
        let privateKey;
        let privateKeyString;
        if (commands.length === 0) {
            return Promise.resolve({message: '/anonymous keypair', args});
        }

        switch (commands[0]) {
        case '--generate':
            response = await generateAndStoreKeyPair();
            if (response && response.status !== 'OK') {
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

            privateKeyString = commands[1];
            privateKey = newFromPrivateKey(privateKeyString);
            if (!key) {
                return Promise.resolve({error: {message: 'Invalid private key'}});
            }

            storePrivateKey(privateKey);
            return Promise.resolve({message: '/anonymous keypair --overwrite ' + privateKey.PublicKey, args});

        case '--export':
            privateKey = loadFromLocalStorage();
            this.store.dispatch(sendEphemeralPost('Your private key:\n' + privateKey.PrivateKey, args.channel_id));
            return Promise.resolve({});
        default:
            break;
        }
        return Promise.resolve({message: '/anonymous keypair' + commands[0], args});
    };

    handlePost = async (commands, args) => {
        const users = await Client4.getProfilesInChannel(args.channel_id);

        const publicKeys = await Promise.all(users.map((user) => {
            return Client.retrievePublicKey(user.id).then((data) => {
                return Buffer.from(data.public_key, 'base64').toString();
            });
        }));

        const encrypted = await Promise.all(publicKeys.map((keyString) => {
            const encrypter = newFromPublicKey(keyString);
            const messageText = commands.join(' ');
            const message = encrypter.encrypt(messageText).toString('base64');
            return {
                message,
                public_key: Buffer.from(keyString).toString('base64'),
            };
        }));

        const message = Buffer.from(JSON.stringify(encrypted)).toString('base64');
        await Client4.createPost({
            channel_id: args.channel_id,
            message,
            props: {encrypted: true},
        });

        return Promise.resolve({});
    };

    slashCommandWillBePostedHook = (message, contextArgs) => {
        const commands = message.split(/(\s+)/).filter((e) => e.trim().length > 0);

        if (commands[0] !== '/anonymous') {
            return Promise.resolve({});
        }
        if (commands.length < 2) {
            return Promise.resolve({message, args: contextArgs});
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

        const {props} = post;

        if (!props || !props.encrypted) {
            return message;
        }

        const decrypter = loadFromLocalStorage();

        if (decrypter === null) {
            return "Message couldn't be decrypted!";
        }

        const messageObject = Array.from(JSON.parse(Buffer.from(message, 'base64').toString()));

        const myMessages = messageObject.filter((value) => {
            return (value.public_key === decrypter.PublicKey);
        });
        if (myMessages.length === 0) {
            return '';
        }

        return decrypter.decrypt(Buffer.from(myMessages[0].message, 'base64'));
    }
}
