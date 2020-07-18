import Constants from '../src/constants';

const Store = {
    getState: jest.fn(() => {
        const reducerId = Constants.REDUCER_ID;
        const state = {};
        state[reducerId] = {};
        state[reducerId].encryptionState = true;
        return state;
    }),
};

export default Store;
