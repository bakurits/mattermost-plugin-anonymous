
export default class LFUCache {
    constructor(size) {
        this.size = size;
        this.data = {};
    }

    /**
     * @param {string} key, postID for which we need a message
     * @returns {string | null} decrypted message if present, null if not
     */
    get = (key) => {
        if (this.data[key]) {
            this.data[key].frequency++;
            return this.data[key].value;
        }
        return null;
    };

    /**
     * @param {string} key, postID for which we need to delete an entry
     * @returns {string | null} value which we are deleting
     */
    delete = (key) => {
        const oldValue = this.get(key);
        delete this.data[key];
        return oldValue;
    };

    /**
     * @param {string} key, postID for which we need a message
     * @param {string | null} val, decrypted message
     * @returns {string | null} entry which we replaced
     */
    put = (key, val) => {
        const oldVal = this.delete(key);
        if (val !== null) {
            if (Object.keys(this.data).length === this.size) {
                this.delete(this.getLFU());
            }
            this.data[key] = {frequency: 0, value: val};
        }
        return oldVal;
    };

    /**
     * @returns {string} key for the least frequently used entry in the dictionary
     */
    getLFU = () => {
        return Object.keys(this.data).sort((key1, key2) => {
            if (this.data[key1].frequency < this.data[key2].frequency) {
                return -1;
            }
            if (this.data[key1].frequency > this.data[key2].frequency) {
                return 1;
            }
            return 0;
        })[0];
    }
}

