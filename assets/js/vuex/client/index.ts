import {Module} from 'vuex';
import {getters} from './getters';
import {actions} from './actions';
import {mutations} from './mutations';
import {ClientState, RootState} from '../../lib/types';

export const state: ClientState = {
    rangedBasisPoint: [],
};

const namespaced: boolean = false;

export const client: Module<ClientState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};