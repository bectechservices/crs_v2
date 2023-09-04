import Vue from 'vue';
import Vuex, {StoreOptions} from 'vuex';
import {RootState} from '../lib/types';
import {custodian} from './custodian';
import {ppt} from './ppt';
import {trustee} from './trustee';
import {npra} from './npra';
import {client} from './client';
import {billing} from './billing';
import {rootMutations} from './mutations';
import {getters} from "./getters";
import VuexPersistence from 'vuex-persist';

Vue.use(Vuex);

const store: StoreOptions<RootState> = {
    state: {
        loading: false,
        pvUploadErrors: []
    },
    modules: {
        custodian,
        ppt,
        trustee,
        npra,
        client,
        billing
    },
    getters,
    mutations: rootMutations,
    plugins: [new VuexPersistence().plugin as any]
};

export default new Vuex.Store<RootState>(store);