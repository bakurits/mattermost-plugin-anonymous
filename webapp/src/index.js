import manifest from './manifest';
import Hooks from './hook/hook';
export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        const hook = new Hooks(store, null);
        registry.registerSlashCommandWillBePostedHook(hook.slashCommandWillBePostedHook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
