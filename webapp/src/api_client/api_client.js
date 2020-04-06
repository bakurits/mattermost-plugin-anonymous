import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import Axios from 'axios';

import {id} from '../manifest';

const STATUS_OK = 200;

export default class Client {
    constructor() {
        this.url = `/plugins/${id}/api/v1`;
    }

    storePublicKey = async (publicKey) => {
        return this.doPost(`${this.url}/pub_key`, {public_key: publicKey.toString('base64')});
    };

    retrievePublicKey = async () => {
        return this.doGet(`${this.url}/pub_key`);
    };

    doGet = async (url, headers = {}) => {
        const opts = Client4.getOptions(headers);
        const options = {
            headers: opts.headers,
            withCredentials: opts.credentials === 'include',
        };
        const response = await Axios.get(url, options);
        if (response.status === STATUS_OK) {
            return response.data;
        }

        throw new ClientError(Client4.url, {
            message: response.statusText || '',
            status_code: response.status,
            url,
        });
    };

    doPost = async (url, body, headers = {}) => {
        const opts = Client4.getOptions(headers);
        const options = {
            headers: opts.headers,
            withCredentials: opts.credentials === 'include',
        };
        const response = await Axios.post(url, body, options);
        if (response.status === STATUS_OK) {
            return response.data;
        }

        throw new ClientError(Client4.url, {
            message: response.statusText || '',
            status_code: response.status,
            url,
        });
    };
}
