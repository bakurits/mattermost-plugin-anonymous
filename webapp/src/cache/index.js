import LFUCache from './lfu_cache';

const CACHE_SIZE = 50;
const Cache = new LFUCache(CACHE_SIZE);

export default Cache;
