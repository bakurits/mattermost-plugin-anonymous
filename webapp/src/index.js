import manifest from './manifest';
import { generateKeyPair, storeKeyPair } from './encrypt/key_manager'

async function handleKeyPair(commands, args){
    switch (commands[0]){
        case "--generate":
            const keys = await generateKeyPair();
            const privateKey = keys[0];
            const publicKey = keys[1];  
            const response = await storeKeyPair(privateKey, publicKey)
            console.log(response)
            console.log('code ', response.code)
            break;
        case "--overwrite":
            break;
        case "--export":
            break;
        default:
            break;
    }
}

function hook(message, args){
    const commands = message.split(/(\s+)/).filter( e => e.trim().length > 0)
    console.log(commands)
    if (commands[0] !== "/anonymous"){
        return Promise.resolve({})
    }
    if (commands.length < 2){
        return Promise.resolve({error: {message: "Command isn't specified"}})
    }

    switch(commands[1]){
        case "keypair":
            handleKeyPair(commands.splice(2), args)
            break;

        default:
            break
    }
}

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerSlashCommandWillBePostedHook(hook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
