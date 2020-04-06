import {
    generateAndStoreKeyPair,
    getPrivateKey,
    getPublicKeyFromPrivateKey,
    storePrivateKey,
} from '../encrypt/key_manager';
import {sendEphemeralPost} from '../actions/actions';

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
        default:
            break;
        }

        return Promise.resolve({message, args: contextArgs});
    }
}
