import {generateAndStoreKeyPair, loadKeyFromLocalStorage, storePrivateKey} from '../encrypt/key_manager';
import {decrypt as aesDecrypt, encrypt as aesEncrypt} from '../encrypt/aes.js';
import {sendEphemeralPost} from '../actions/actions';
import {newFromPrivateKey, newFromPublicKey} from '../encrypt/key';
import Cache from '../cache';
import Constants from '../constants';

export default class Hooks {
    constructor(store, settings, client) {
        this.store = store;
        this.settings = settings;
        this.Client = client;
    }

    /**
     * @param {string[]} commands, slash command input
     * @param {object} args, contextArgs object
     * @returns {Promise<Object>} object with modified commands or an error message
     */
    handleKeyPair = async (commands, args) => {
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
            privateKey = newFromPrivateKey(atob(privateKeyString));
            if (!privateKey) {
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
    handlePostCommand = async (commands, args) => {
        const encryptedData = await this.encryptMessage(args.channel_id, commands.join(' '));
        if (encryptedData.success !== true) {
            return Promise.resolve({error: 'could not encrypt message properly'});
        }
        const message = encryptedData.message;
        await this.Client.createPost({
            channel_id: args.channel_id,
            message,
            props: {encrypted: true},
        });

        return Promise.resolve({});
    };

    /**
     * @param {string} channelID, channel id
     * @param {Object} post, message to be encrypted
     * @returns {Object} success status of the operation with encrypted message
     */
    encryptMessage = async (channelID, post) => {
        const users = await this.Client.getProfilesInChannel(channelID);

        const userIDs = users.map((user) => {
            return user.id;
        });

        const publicKeysb64 = await this.Client.retrievePublicKeys(userIDs);
        if (publicKeysb64 === null) {
            return {success: false, message: ''};
        }

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

        return {success: true, message: JSON.stringify(encrypted)};
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
        let messageObject;
        try {
            messageObject = Array.from(JSON.parse(message));
        } catch (e) {
            return "Message couldn't be decrypted!";
        }
        const myMessages = messageObject.filter((value) => {
            return (value.public_key === decrypter.PublicKey);
        });
        if (!myMessages || myMessages.length === 0) {
            return 'Message could not be decrypted!';
        }
        const encryptedAESKey = myMessages[0].aes_key;
        const encryptedMessage = myMessages[0].message;
        if (!encryptedAESKey) {
            return "Message couldn't be decrypted!";
        }
        const aesKey = decrypter.decrypt(Buffer.from(encryptedAESKey, 'base64'));
        return aesDecrypt(encryptedMessage, aesKey);
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
            return this.handlePostCommand(commands.splice(2), contextArgs);
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

    /**
     * @param {Object} post, post to be processed
     * @returns {object} processed post
     */
    messageWillBePostedHook = async (post) => {
        if (this.store.getState()[Constants.REDUCER_ID].encryptionState !== true) {
            return Promise.resolve({post});
        }
        const newPost = post;
        newPost.props = {encrypted: true};
        const encryptedData = await this.encryptMessage(post.channel_id, post.message);
        if (encryptedData.success !== true) {
            return Promise.resolve({error: {message: 'could not encrypt properly'}});
        }
        newPost.message = encryptedData.message;
        return Promise.resolve({post: newPost});
    }
}
