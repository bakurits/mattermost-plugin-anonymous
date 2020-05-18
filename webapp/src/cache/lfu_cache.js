
export default class LFUCache {
    constructor(size) {
        this.size = size;
        this.data = {};
    }
    get = (key) => {
        if (this.data[key]) {
            this.data[key].frequency++;
            return this.data[key].value;
        }
        return null;
    };

    delete = (key) => {
        const oldValue = this.get(key);
        delete this.data[key];
        return oldValue;
    };

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

