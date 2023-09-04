import {ActionTree} from 'vuex';
import {ApiResponse, IUnidentifiedPayment, RootState, TrusteeState} from '../../lib/types';
import {ROOT_STATE_CONSTANTS} from "../constants";
import axios from "axios";


export const actions: ActionTree<TrusteeState, RootState> = {
    async uploadUnidentifiedPayments({commit}, data: Array<IUnidentifiedPayment>): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-unidentified-payments", data, {
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                }
            });
            apiResponse = response.data;
        } catch (error) {
            apiResponse = error.response.data;
        } finally {
            commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, false);
        }
        return apiResponse;
    },
};