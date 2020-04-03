
import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import Axios from 'axios';

const STATUS_OK = 200;

export default class Client {
    constructor() {
        this.url = '/plugins/anonymous/api/v1';
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
        console.log('dopost!!!  ')


        Axios.get(url, Client4.getOptions(options)).then(async (response) => {
            console.log('axios  ', response)
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

