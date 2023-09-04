import {ActionTree} from 'vuex';
import {ApiResponse, ClientState, RangedBasisPoint, RootState} from '../../lib/types';
import {ROOT_STATE_CONSTANTS} from "../constants";
import axios from "axios";


export const actions: ActionTree<ClientState, RootState> = {
    async uploadClientBillingInfo({commit, state}, data: any): Promise<ApiResponse> {
        let apiResponse: ApiResponse = {error: false, message: ""};
        commit(ROOT_STATE_CONSTANTS.SET_LOADING_STATUS, true);
        try {
            let basisPoints: Array<RangedBasisPoint> = [];
            if (data.basis_point_type == "Flat") {
                basisPoints.push({key: 0, maximum_amount: 0, minimum_amount: 0, basis_point: data.basis_point})
            } else {
                basisPoints = state.rangedBasisPoint;
            }
            const response = await axios.post("/update-client-billing-info", {
                bpid: data.bpid,
                minimum_charge: data.minimum_charge,
                charge_per_txn: data.charge_per_txn,
                third_party_transfer: data.third_party_transfer,
                basis_points: basisPoints
            }, {
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