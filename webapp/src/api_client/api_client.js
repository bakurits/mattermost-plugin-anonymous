
import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import Axios from 'axios';

const STATUS_OK = 200;

export default class Client {
    constructor() {
        this.url = '/plugins/anonymous/api';
    }

    storePublicKey = async (publicKey) => {
        return this.doPost(`${this.url}/pub_key`, publicKey);
    };

    retrievePublicKey = async () => {
        return this.doGet(`${this.url}/pub_key`);
    };

    doGet = async (url, headers = {}) => {
        const options = {
            method: 'get',
            headers,
        };
        Axios.get(url, Client4.getOptions(options)).then(async (response) => {
            if (response.status === STATUS_OK) {
                return response.data;
            }

            throw new ClientError(Client4.url, {
                message: response.statusText || '',
                status_code: response.status,
                url,
            });
        });
    };

    doPost = async (url, body, headers = {}) => {
        const options = {
            method: 'post',
            body,
            headers,
        };

        Axios.get(url, Client4.getOptions(options)).then(async (response) => {
            if (response.status === STATUS_OK) {
                return response.data;
            }

            throw new ClientError(Client4.url, {
                message: response.statusText || '',
                status_code: response.status,
                url,
            });
        });
    };
}

//store public key on server
export function storePublicKey(publicKey, callback) {
    // eslint-disable-next-line no-warning-comments
    // Todo: post a request to server to store public key
    callback(0);
}

//retrieve public key from server
export function retrievePublicKey(callback) {
    // eslint-disable-next-line no-warning-comments
    // Todo: get public key from server
    callback([1, 1]);
}

