import {MutationTree} from 'vuex';
import {BillingSearchData, BillingState} from '../../lib/types';

export const mutations: MutationTree<BillingState> = {
    storeBillingSearchData(state: BillingState, data: BillingSearchData): void {
        state.searchData = data;
    }
};