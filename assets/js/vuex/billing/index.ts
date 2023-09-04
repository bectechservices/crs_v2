import {Module} from 'vuex';
import {getters} from './getters';
import {actions} from './actions';
import {mutations} from './mutations';
import {BillingState, RootState} from '../../lib/types';

export const state: BillingState = {
    searchData: {
        bpOrSca: "",
        date: ""
    }
};

const namespaced: boolean = false;

export const billing: Module<BillingState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};