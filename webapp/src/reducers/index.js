import {combineReducers} from 'redux';

import ActionTypes from '../action_types';

const encryptionState = (state = true, action) => {
    switch (action.type) {
    case ActionTypes.DISABLE_ENCRYPTION:
        return false;
    case ActionTypes.ENABLE_ENCRYPTION:
        return true;
    default:
        return state;
    }
};

export default combineReducers({
    encryptionState,
});
