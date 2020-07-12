import {Client4} from 'mattermost-redux/client';

import {generateAndStoreKeyPair, loadKeyFromLocalStorage, storePrivateKey} from '../encrypt/key_manager';
import {decrypt as aesDecrypt, encrypt as aesEncrypt} from '../encrypt/aes.js';
import {sendEphemeralPost} from '../actions/actions';
import {newFromPrivateKey, newFromPublicKey} from '../encrypt/key';
import Client from '../api_client';
import Cache from '../cache';

export default class Hooks {
    constructor(store, settings) {
        this.store = store;
        this.settings = settings;
    }

    /**
     * @param {string[]} commands, slash command input
     * @param {object} args, contextArgs object
     * @returns {Promise<Object>} object with modified commands or an error message
     */
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
            privateKey = loadKeyFromLocalStorage();
            this.store.dispatch(sendEphemeralPost('Your private key:\n' + privateKey.PrivateKey, args.channel_id));
            return Promise.resolve({});
        default:
            break;
        }
        return Promise.resolve({message: '/anonymous keypair' + commands[0], args});
    };

    /**
     * @param {string[]} commands, slash command input
     * @param {object} args, contextArgs object
     * @returns {Promise<Object>} resolved promise after sending messages to all users in channel
     */
    handlePost = async (commands, args) => {
        await this.encryptMessage(args.channel_id, commands.join(' '));
        return Promise.resolve({});
    };

    encryptMessage = async (channelID, post) => {
        const users = await Client4.getProfilesInChannel(channelID);

        const userIDs = users.map((user) => {
            return user.id;
        });

        const publicKeysb64 = await Client.retrievePublicKey(userIDs);

        const publicKeys = publicKeysb64.public_keys.map((publicKey) => {
            return Buffer.from(publicKey, 'base64').toString();
        });

        const encrypted = await Promise.all(publicKeys.map((keyString) => {
            const encrypter = newFromPublicKey(keyString);
            const aesEncryptData = aesEncrypt(post);
            const encryptedAESKey = encrypter.encrypt(aesEncryptData.key).toString('base64');
            return {
                message: aesEncryptData.message,
                aes_key: encryptedAESKey,
                public_key: Buffer.from(keyString).toString('base64'),
            };
        }));

        const message = Buffer.from(JSON.stringify(encrypted)).toString('base64');
        await Client4.createPost({
            channel_id: channelID,
            message,
            props: {encrypted: true},
        });

        return Promise.resolve({});
    };

    /**
     * @param {Object} post, post that needs decryption
     * @returns {string} decrypted message
     */
    decryptMessage = (post) => {
        // message text in database
        const {message} = post;

        const decrypter = loadKeyFromLocalStorage();

        if (decrypter === null) {
            return "Message couldn't be decrypted!";
        }

        const messageObject = Array.from(JSON.parse(Buffer.from(message, 'base64').toString()));

        const myMessages = messageObject.filter((value) => {
            return (value.public_key === decrypter.PublicKey);
        });
        if (!myMessages || myMessages.length === 0) {
            return '';
        }
        const encryptedAESKey = myMessages[0].aes_key;
        if (!encryptedAESKey) {
            return "Message couldn't be decrypted!";
        }
        const aesKey = decrypter.decrypt(Buffer.from(encryptedAESKey, 'base64'));
        return aesDecrypt(myMessages[0].message, aesKey);
    }

    /**
     * @param {string} message, slash command
     * @param {object} contextArgs, contextArgs object
     * @returns {Promise<Object>} object with modified commands or an error message
     */
    slashCommandWillBePostedHook = (message, contextArgs) => {
        const commands = message.split(/(\s+)/).filter((e) => e.trim().length > 0);

        if (commands[0] !== '/anonymous') {
            return Promise.resolve({message, args: contextArgs});
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

    /**
     * @param {Object} post, post to be formatted
     * @returns {string} formatted message
     */
    messageWillFormatHook = (post) => {
        const {id} = post;
        const {props} = post;
        const {message} = post;

        if (!props || !props.encrypted) {
            return message;
        }

        const cachedMessage = Cache.get(id);
        if (cachedMessage) {
            return cachedMessage;
        }

        const decryptedMessage = this.decryptMessage(post);
        Cache.put(id, decryptedMessage);
        return decryptedMessage;
    }
}
