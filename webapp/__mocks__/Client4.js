
const Client4 = {
    getProfilesInChannel: jest.fn(async () => {
        return Promise.resolve([{id: 'user1'}]);
    }),
    createPost: jest.fn(async () => {
        return Promise.resolve({});
    }),
};

export default Client4;
