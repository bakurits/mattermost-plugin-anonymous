import manifest from './manifest';
import {generateAndStoreKeyPair, getKeyPair} from './encrypt/key_manager';

// eslint-disable-next-line no-unused-vars,consistent-return
async function handleKeyPair(commands, args) {
    let response;
    switch (commands[0]) {
    case '--generate':
        response = await generateAndStoreKeyPair();
        // eslint-disable-next-line no-console
        console.log(response);
        return Promise.resolve(response);
    case '--overwrite':
        return Promise.resolve({});
    case '--export':
        // eslint-disable-next-line no-case-declarations
        const keys = await getKeyPair();
        // eslint-disable-next-line no-case-declarations,no-console
        console.log(keys);
        return Promise.resolve({message: keys.privateKey.toString('base64'), args});
    default:
        break;
    }
}

function hook(message, args) {
    const commands = message.split(/(\s+)/).filter((e) => e.trim().length > 0);
    // eslint-disable-next-line no-console
    console.log(commands);
    if (commands[0] !== '/anonymous') {
        return Promise.resolve({message, args});
    }
    if (commands.length < 2) {
        return Promise.resolve({error: {message: "Command isn't specified"}});
    }

    switch (commands[1]) {
    case 'keypair':
        return handleKeyPair(commands.splice(2), args);
    default:
        break;
    }

    //TODO: finish this
    return Promise.resolve({message: {}, args});
}

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerSlashCommandWillBePostedHook(hook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
