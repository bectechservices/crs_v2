import {MutationTree} from 'vuex';
import {PPTState} from '../../lib/types';

export const mutations: MutationTree<PPTState> = {
    storePPTOptions(state: PPTState, data: PPTState): void {
        state.bpid = data.bpid;
        state.templateOptions = data.templateOptions;
        state.year = data.year;
        state.quarter = data.quarter;
    }
};