import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import Axios from 'axios';

import {id} from '../manifest';

const STATUS_OK = 200;

export default class Client {
    constructor() {
        this.url = `/plugins/${id}/api/v1`;
    }

    /**
     *  @param {string} publicKey
     *  @returns {Object} response from api call
     */
    storePublicKey = async (publicKey) => {
        return this.doPost(`${this.url}/pub_key`, {public_key: publicKey});
    };

    /**
     *  @param {string} userID
     *  @returns {Object} response from api call
     */
    retrievePublicKey = async (userID) => {
        return this.doGet(`${this.url}/pub_key?user_id=${userID}`);
    };

    /**
     *  @param {string} url, api endpoint
     *  @param {Object} headers, request headers
     *  @returns {Object} response from api call
     */
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

    /**
     *  @param {string} url, api endpoint
     *  @param {Object} body, request body
     *  @param {Object} headers, request headers
     *  @returns {Object} response from api call
     */
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
