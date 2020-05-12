import manifest from './manifest';
import Hooks from './hook/hook';

export default class Plugin {
    initialize(registry, store) {
        const hook = new Hooks(store, null);
        registry.registerSlashCommandWillBePostedHook(hook.slashCommandWillBePostedHook);
        registry.registerMessageWillFormatHook(hook.messageWillFormatHook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
