import { GetterTree } from 'vuex';
import { CustodianState, RootState, IKeyableShareError, IKeyableCustodianTransactionError, IKeyableSchemesUnderCustodyError, IMultipleGovernanceFieldsData, IKeyableShare, IKeyableCustodianTransaction, IKeyableSchemesUnderCustody } from '../../lib/types';

export const getters: GetterTree<CustodianState, RootState> = {
    custodianMultipleInputFieldHasError(state: CustodianState): boolean {
        return state.fieldHasError;
    },
    ordinaryShareInputHasError(state: CustodianState) {
        return (key: number) => {
            return state.ordinarySharesInputError.find((share: IKeyableShareError) => share.key === key)
        }
    },
    preferenceShareInputHasError(state: CustodianState) {
        return (key: number) => {
            return state.preferenceSharesInputError.find((share: IKeyableShareError) => share.key === key)
        }
    },
    custodianTransactionsInputHasError(state: CustodianState) {
        return (key: number) => {
            return state.custodianTransactionsInputError.find((transaction: IKeyableCustodianTransactionError) => transaction.key === key)
        }
    },
    schemesUnderCustodyInputHasError(state: CustodianState) {
        return (key: number) => {
            return state.schemesUnderCustodyInputError.find((scheme: IKeyableSchemesUnderCustodyError) => scheme.key === key)
        }
    },
    multipleFieldsData(state: CustodianState): IMultipleGovernanceFieldsData {
        return {
            ordinarySharesInputData: state.ordinarySharesInputData,
            preferenceSharesInputData: state.preferenceSharesInputData,
            custodianTransactionsInputData: state.custodianTransactionsInputData,
            schemesUnderCustodyInputData: state.schemesUnderCustodyInputData
        }
    },
    ordinaryShareInputData(state: CustodianState) {
        return (key: number) => {
            return state.ordinarySharesInputData.find((share: IKeyableShare) => share.key === key);
        }
    },
    preferenceShareInputData(state: CustodianState) {
        return (key: number) => {
            return state.preferenceSharesInputData.find((share: IKeyableShare) => share.key === key);
        }
    },
    custodianTransactionsInputData(state: CustodianState) {
        return (key: number) => {
            return state.custodianTransactionsInputData.find((transaction: IKeyableCustodianTransaction) => transaction.key === key);
        }
    },
    schemesUnderCustodyInputData(state: CustodianState) {
        return (key: number) => {
            return state.schemesUnderCustodyInputData.find((scheme: IKeyableSchemesUnderCustody) => scheme.key === key);
        }
    }
};