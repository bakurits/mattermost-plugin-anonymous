import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import Client from '../api_client';
import {initializeEncryptionStatus} from '../actions/actions';

export default class ChannelChangeListener {
    constructor(store) {
        this.store = store;
        this.previousValue = getCurrentChannelId(this.store.getState());
        this.unsubscribe = this.store.subscribe(this.handleChange);
    }

    handleChange = async () => {
        const newValue = getCurrentChannelId(this.store.getState());
        if (newValue !== this.previousValue) {
            this.previousValue = newValue;
            await initializeEncryptionStatusForChannel(newValue);
        }
    }

    unsubscribe = () => {
        this.unsubscribe();
    }
}

export const initializeEncryptionStatusForChannel = async (channelID) => {
    const encryptionStatus = await Client.getEncryptionStatus(channelID);
    this.store.dispatch(initializeEncryptionStatus(encryptionStatus));
};
