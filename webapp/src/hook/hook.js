import {generateAndStoreKeyPair, getKeyPair} from '../encrypt/key_manager';
import {sendEphemeralPost} from '../actions/actions';

export default class Hooks {
    constructor(store, settings) {
        this.store = store;
        this.settings = settings;
    }

    // eslint-disable-next-line consistent-return
    handleKeyPair = async (commands, args) => {
        switch (commands[0]) {
        case '--generate':
            // eslint-disable-next-line no-case-declarations
            await generateAndStoreKeyPair();
            this.store.dispatch(sendEphemeralPost('keys generated', args.channel_id));
            return Promise.resolve({});
        case '--overwrite':
            return Promise.resolve({});
        case '--export':
            // eslint-disable-next-line no-case-declarations
            const keys = await getKeyPair();
            // eslint-disable-next-line no-case-declarations,no-console
            const privateKey = keys[0];
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
