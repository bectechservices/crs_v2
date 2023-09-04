import {GetterTree} from 'vuex';
import {RootState} from "../lib/types";

export const getters: GetterTree<RootState, RootState> = {
    uploadedPvErrors(state: RootState) {
        return state.pvUploadErrors;
    }
};