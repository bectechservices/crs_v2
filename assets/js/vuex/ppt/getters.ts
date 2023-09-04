import {GetterTree} from 'vuex';
import {PPTState, RootState} from '../../lib/types';

export const getters: GetterTree<PPTState, RootState> = {
    pptTemplateOptions(state: PPTState) {
        return state.templateOptions;
    }
};