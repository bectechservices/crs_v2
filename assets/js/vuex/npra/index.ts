import {Module} from 'vuex';
import {getters} from './getters';
import {actions} from './actions';
import {mutations} from './mutations';
import {NPRAState, RootState} from '../../lib/types';

export const state: NPRAState = {
    unauthorizedTransactions: []
};

const namespaced: boolean = false;

export const npra: Module<NPRAState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};