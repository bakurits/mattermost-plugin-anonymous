import {
    generateAndStoreKeyPair,
    getPrivateKey,
    getPublicKeyFromPrivateKey,
    storePrivateKey,
} from '../encrypt/key_manager';
import {sendEphemeralPost} from '../actions/actions';

const base64Tester = RegExp('^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$');

export default class Hooks {
    constructor(store, settings) {
        this.store = store;
        this.settings = settings;
    }

    handleKeyPair = async (commands, args) => {
        let privateKey;
        switch (commands[0]) {
        case '--generate':
            // eslint-disable-next-line no-case-declarations
            const response = await generateAndStoreKeyPair();
            if (response.status !== 'OK') {
                return Promise.resolve({error: {message: 'Error occurred while trying to store key on a server'}});
            }
            this.store.dispatch(sendEphemeralPost('keys generated', args.channel_id));
            return Promise.resolve({});

        case '--overwrite':
            if (commands.length < 2) {
                return Promise.resolve({error: {message: "Private key isn't specified"}});
            }
            // eslint-disable-next-line no-magic-numbers
            if (commands.length > 2) {
                return Promise.resolve({error: {message: 'Too many arguments'}});
            }
            if (!base64Tester.test(commands[1])) {
                return Promise.resolve({error: {message: 'Invalid private key'}});
            }
            privateKey = Buffer.from(commands[1], 'base64');
            storePrivateKey(privateKey);
            // eslint-disable-next-line no-case-declarations
            const publicKey = getPublicKeyFromPrivateKey(privateKey).toString('base64');
            if (!publicKey) {
                return Promise.resolve({error: {message: 'Invalid private key'}});
            }
            return Promise.resolve({message: '/anonymous keypair --overwrite ' + publicKey, args});

        case '--export':
            privateKey = getPrivateKey();
            this.store.dispatch(sendEphemeralPost('your private key is :    ' + privateKey.toString('base64'), args.channel_id));
            return Promise.resolve({});
        default:
            break;
        }
        return Promise.resolve({});
    };

    slashCommandWillBePostedHook = (message, contextArgs) => {
        const commands = message.split(/(\s+)/).filter((e) => e.trim().length > 0);
        // eslint-disable-next-line no-console
        console.log(commands);
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
