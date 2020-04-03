import manifest from './manifest';
import {generateKeyPair} from './encrypt/key_manager';

// eslint-disable-next-line no-unused-vars
async function handleKeyPair(commands, args) {
    let generateReturn;
    let response;
    switch (commands[0]) {
    case '--generate':
        generateReturn = await generateKeyPair();
        response = generateReturn[0];
        // eslint-disable-next-line no-console
        console.log(response);
        // eslint-disable-next-line no-console
        console.log('code ', response.code);
        break;
    case '--overwrite':
        break;
    case '--export':
        break;
    default:
        break;
    }
}

function hook(message, args) {
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
        handleKeyPair(commands.splice(2), args);
        break;
    default:
        break;
    }

    //TODO: finish this
    return Promise.resolve({});
}

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerSlashCommandWillBePostedHook(hook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
