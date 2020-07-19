
const Client = {
    retrievePublicKey: jest.fn(async () => {
        return Promise.resolve({public_keys: []});
    }),
    getProfilesInChannel: jest.fn(async () => {
        return Promise.resolve([{id: 'user1'}]);
    }),
    createPost: jest.fn(async () => {
        return Promise.resolve({});
    }),
};

export default Client;
