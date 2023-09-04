import {ActionTree} from 'vuex';
import {ApiResponse, IOutstandingFDCertificate, NPRAState, RootState} from '../../lib/types';
import {ROOT_STATE_CONSTANTS} from "../constants";
import axios from "axios";


export const actions: ActionTree<NPRAState, RootState> = {
    async uploadUnauthorizedTransactions({commit, state}): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-unauthorized-transactions", {transactions: state.unauthorizedTransactions}, {
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
    async uploadOutstandingFDCertificates({commit}, outstandingFDCertificates: Array<IOutstandingFDCertificate>): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-outstanding-fd-certificates", {certificates: outstandingFDCertificates}, {
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
    async npraSendChangesToChecker({commit}): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: true, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/npra-send-progress-to-checker", {}, {
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
    }
};