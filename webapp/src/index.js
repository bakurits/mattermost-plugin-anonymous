import React from 'react';
import {Client4} from 'mattermost-redux/client';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import Client from 'api_client';

import reducer from 'reducers';

import manifest from './manifest';

import Hooks from './hook/hook';
import Icon from './components/iconContainer';
import {handleEncryptionStatusChange, toggleEncryption} from './actions/actions';
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

        const hook = new Hooks(store, null, Client4, Client);

        registry.registerSlashCommandWillBePostedHook(hook.slashCommandWillBePostedHook);
        registry.registerMessageWillFormatHook(hook.messageWillFormatHook);
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_encryption_status_change`, handleEncryptionStatusChange(store));
    }

    uninitialize() {
        this.channelChangeListener.unsubscribe();
    }
}

window.registerPlugin(manifest.id, new Plugin());
