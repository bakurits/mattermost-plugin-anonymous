
const Client = {
    retrievePublicKey: jest.fn(async () => {
        return Promise.resolve({public_keys: []});
    }),
};

export default Client;
