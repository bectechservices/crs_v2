import {Module} from 'vuex';
import {getters} from './getters';
import {actions} from './actions';
import {mutations} from './mutations';
import {RootState, TrusteeState} from '../../lib/types';

export const state: TrusteeState = {
    unidentifiedPayments: [],
    searchInputData: {
        bpid: "",
        quarter: "",
        year: ""
    }
};

const namespaced: boolean = false;

export const trustee: Module<TrusteeState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};