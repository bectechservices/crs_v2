import {Module} from 'vuex';
import {getters} from './getters';
import {actions} from './actions';
import {mutations} from './mutations';
import {PPTState, RootState} from '../../lib/types';

export const state: PPTState = {
    templateOptions: {
        total_summary_of_auc: false,
        auc_trend: false,
        trade_volumes: false,
        pv_report: false,
        total_contribution: false,
        corporate_action: false,
        gog_and_fd_maturities: false,
        appendix_i: false,
        appendix_ii: false,
        unidentified_payments: false
    },
    bpid: "",
    year: 0,
    quarter: 0
};

const namespaced: boolean = false;

export const ppt: Module<PPTState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};