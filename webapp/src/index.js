import React from 'react';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import reducer from 'reducers';

import manifest from './manifest';

import Hooks from './hook/hook';
import Icon from './components/iconContainer';
import {toggleEncryption} from './actions/actions';
import ChannelChangeListener, {initializeEncryptionStatusForChannel} from './hook/channelChangeListener';

export default class Plugin {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        this.channelChangeListener = new ChannelChangeListener(store);
        initializeEncryptionStatusForChannel(store.dispatch, getCurrentChannelId(store.getState()));

        registry.registerChannelHeaderButtonAction(
            // eslint-disable-next-line react/jsx-filename-extension
            <Icon/>,
            (channel) => {
                store.dispatch(toggleEncryption(channel.id));
            },
            'toggle encryption',
            'toggle encryption'
        );

        const hook = new Hooks(store, null);

        registry.registerSlashCommandWillBePostedHook(hook.slashCommandWillBePostedHook);
        registry.registerMessageWillFormatHook(hook.messageWillFormatHook);
    }

    uninitialize() {
        this.channelChangeListener.unsubscribe();
    }
}

window.registerPlugin(manifest.id, new Plugin());
