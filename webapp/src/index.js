import React from 'react';

import manifest from './manifest';

import Hooks from './hook/hook';
import Icon from './components/icon';

export default class Plugin {
    initialize(registry, store) {
        const hook = new Hooks(store, null);
        registry.registerChannelHeaderButtonAction(
            // eslint-disable-next-line react/jsx-filename-extension
            <Icon/>,
            (channel) => {
                // eslint-disable-next-line no-console
                console.log(channel);
            },
            'toggle encryption',
            'toggle encryption'
        );
        registry.registerSlashCommandWillBePostedHook(hook.slashCommandWillBePostedHook);
        registry.registerMessageWillFormatHook(hook.messageWillFormatHook);
    }
}

window.registerPlugin(manifest.id, new Plugin());
