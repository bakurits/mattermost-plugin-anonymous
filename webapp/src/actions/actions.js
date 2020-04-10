import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import {PostTypes} from 'mattermost-redux/action_types';

import manifest from '../manifest';

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
