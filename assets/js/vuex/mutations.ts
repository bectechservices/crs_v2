import {MutationTree} from 'vuex';
import {RootState} from '../lib/types';

export const rootMutations: MutationTree<RootState> = {
    setLoadingStatus(state: RootState, status: boolean): void {
        state.loading = status;
    },
    storePvUploadErrors(state: RootState, errors: Array<string>): void {
        state.pvUploadErrors = errors;
    }
};