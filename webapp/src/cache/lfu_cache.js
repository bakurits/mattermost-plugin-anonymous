
export default class LFUCache {
    constructor(size) {
        this.size = size;
        this.data = {};
    }

    /**
     * @param {string} key, key for which we need a value
     * @returns {string | null} value if present, null if not
     */
    get = (key) => {
        if (this.data[key]) {
            this.data[key].frequency++;
            return this.data[key].value;
        }
        return null;
    };

    /**
     * @param {string} key for which we need to delete an entry
     * @returns {string | null} value which we are deleting
     */
    delete = (key) => {
        const oldValue = this.get(key);
        delete this.data[key];
        return oldValue;
    };

    /**
     * @param {string} key for which we want to store a value
     * @param {string | null} val, value to store
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
        let keys = Object.keys(this.data);
        keys.sort((key1, key2) => {
            if (this.data[key1].frequency < this.data[key2].frequency) {
                return -1;
            }
            if (this.data[key1].frequency > this.data[key2].frequency) {
                return 1;
            }
            return 0;
        });
        return keys[0];
    }
}

