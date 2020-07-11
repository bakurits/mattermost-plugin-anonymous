import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import {PostTypes} from 'mattermost-redux/action_types';

import ActionTypes from '../action_types';

import Constants from '../constants';

import manifest from '../manifest';
import Client from '../api_client';

/**
 * @param {string} message, message to be posted
 * @param {number} channelId, channel id in which message should be posted
 */
export function sendEphemeralPost(message, channelId) {
    return (dispatch, getState) => {
        const timestamp = Date.now();
        const post = {
            id: manifest.id + Date.now(),
            user_id: getState().entities.users.currentUserId,
            channel_id: channelId || getCurrentChannelId(getState()),
            message,
            type: 'system_ephemeral',
            create_at: timestamp,
            update_at: timestamp,
            root_id: '',
            parent_id: '',
            props: {},
        };

        dispatch({
            type: PostTypes.RECEIVED_NEW_POST,
            data: post,
            channelId,
        });
    };
}

export function enableEncryption(channelId) {
    return async (dispatch, getState) => {
        if (getState()[Constants.REDUCER_ID].encryptionState !== true) {
            const response = await Client.setEncryptionStatus(channelId, true);
            if (response.status !== 'OK') {
                return;
            }
            dispatch({
                type: ActionTypes.ENABLE_ENCRYPTION,
            });
        }
    };
}

export function disableEncryption(channelId) {
    return async (dispatch, getState) => {
        if (getState()[Constants.REDUCER_ID].encryptionState !== false) {
            const response = await Client.setEncryptionStatus(channelId, false);
            if (response.status !== 'OK') {
                return;
            }
            dispatch({
                type: ActionTypes.DISABLE_ENCRYPTION,
            });
        }
    };
}

export function toggleEncryption(channelId) {
    return async (dispatch, getState) => {
        if (getState()[Constants.REDUCER_ID].encryptionState === false) {
            const response = await Client.setEncryptionStatus(channelId, true);
            if (response.status !== 'OK') {
                return;
            }
            dispatch({
                type: ActionTypes.ENABLE_ENCRYPTION,
            });
        } else {
            const response = await Client.setEncryptionStatus(channelId, false);
            if (response.status !== 'OK') {
                return;
            }
            dispatch({
                type: ActionTypes.DISABLE_ENCRYPTION,
            });
        }
    };
}

export function initializeEncryptionStatus(status) {
    return async (dispatch) => {
        if (status === true) {
            dispatch({
                type: ActionTypes.ENABLE_ENCRYPTION,
            });
        } else {
            dispatch({
                type: ActionTypes.DISABLE_ENCRYPTION,
            });
        }
    };
}

// eslint-disable-next-line no-unused-vars
export function handleEncryptionStatusChange(store) {
    return (msg) => {
        store.dispatch(initializeEncryptionStatus(msg.data.status));
    };
}
