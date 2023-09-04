import { Module } from 'vuex';
import { getters } from './getters';
import { actions } from './actions';
import { mutations } from './mutations';
import { CustodianState,RootState } from '../../lib/types';

export const state: CustodianState = {
    ordinarySharesInputData: [],
    preferenceSharesInputData: [],
    custodianTransactionsInputData: [],
    schemesUnderCustodyInputData: [],
    ordinarySharesInputError: [],
    preferenceSharesInputError: [],
    custodianTransactionsInputError: [],
    schemesUnderCustodyInputError: [],
    fieldHasError: false
};

const namespaced: boolean = false;

export const custodian: Module<CustodianState, RootState> = {
    namespaced,
    state,
    getters,
    actions,
    mutations
};