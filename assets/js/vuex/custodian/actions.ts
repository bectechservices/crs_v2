import {ActionTree} from 'vuex';
import {
    ApiResponse,
    CustodianState,
    IGovernance,
    IOfficialRemarks,
    IOtherSecInformations,
    ISchemeDetails,
    RootState,
    SchemeDetailsRequest
} from '../../lib/types';
import axios from 'axios';
import {ROOT_STATE_CONSTANTS} from '../constants';


export const actions: ActionTree<CustodianState, RootState> = {
    async uploadGovernanceData({commit}, data: IGovernance): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-governance", data, {
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
    async uploadSchemeDetails({commit}, data: ISchemeDetails): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-scheme-details", data, {
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
    async uploadOtherInformation({commit}, data: IOtherSecInformations): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-other-information", data, {
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
    async uploadOfficialReportRemarks({commit}, data: IOfficialRemarks): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/upload-official-remarks", data, {
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
    async fetchSchemeDetails({commit}, data: SchemeDetailsRequest): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/fetch-scheme-details", data, {
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
    async recalculateFormulas({commit}, data: SchemeDetailsRequest): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/recalculate-scheme-details", data, {
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
    async sendChangesToChecker({commit}): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: true, message: ""}
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/sec-send-progress-to-checker", {}, {
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
    async deleteSchemeDetails({commit}, data: { bid: string }): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            const response = await axios.post("/delete-scheme-details", data, {
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