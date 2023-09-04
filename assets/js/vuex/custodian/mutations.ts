import {MutationTree} from 'vuex';
import {
    CustodianState,
    IKeyableCustodianTransaction,
    IKeyableCustodianTransactionError,
    IKeyableSchemesUnderCustody,
    IKeyableSchemesUnderCustodyError,
    IKeyableShare,
    IKeyableShareError
} from '../../lib/types';
import {isNumeric} from '../../lib/helpers';

export const mutations: MutationTree<CustodianState> = {
    addOrdinaryShareInput(state: CustodianState, key: number): void {
        state.ordinarySharesInputData.push({key, name: "", percentage: 0, shareholding: 0});
        state.ordinarySharesInputError.push({key, name: false, percentage: false, shareholding: false})
    },
    deleteOrdinaryShareInput(state: CustodianState, key: number): void {
        state.ordinarySharesInputData = state.ordinarySharesInputData.filter(share => share.key !== key);
        state.ordinarySharesInputError = state.ordinarySharesInputError.filter(share => share.key !== key)
    },
    modifyOrdinaryShareInputData(state: CustodianState, share: IKeyableShare): void {
        state.ordinarySharesInputData = state.ordinarySharesInputData.map((_share: IKeyableShare) => {
            if (_share.key === share.key) {
                return share;
            }
            return _share;
        })
    },
    addPreferenceShareInput(state: CustodianState, key: number): void {
        state.preferenceSharesInputData.push({key, name: "", percentage: 0, shareholding: 0});
        state.preferenceSharesInputError.push({key, name: false, percentage: false, shareholding: false});
    },
    deletePreferenceShareInput(state: CustodianState, key: number): void {
        state.preferenceSharesInputData = state.preferenceSharesInputData.filter(share => share.key !== key);
        state.preferenceSharesInputError = state.preferenceSharesInputError.filter(share => share.key !== key)
    },
    modifyPreferenceShareInputData(state: CustodianState, share: IKeyableShare): void {
        state.preferenceSharesInputData = state.preferenceSharesInputData.map((_share: IKeyableShare) => {
            if (_share.key === share.key) {
                return share;
            }
            return _share;
        })
    },
    addCustodianTransactionInput(state: CustodianState, key: number): void {
        state.custodianTransactionsInputData.push({
            key,
            nameOfTrustee: "",
            relationshipWithTrustee: "",
            typeOfTransaction: "",
            amount: 0
        });
        state.custodianTransactionsInputError.push({
            key,
            nameOfTrustee: false,
            relationshipWithTrustee: false,
            typeOfTransaction: false,
            amount: false
        })
    },
    deleteCustodianTransactionInput(state: CustodianState, key: number): void {
        state.custodianTransactionsInputData = state.custodianTransactionsInputData.filter(share => share.key !== key);
        state.custodianTransactionsInputError = state.custodianTransactionsInputError.filter(share => share.key !== key)
    },
    modifyCustodianTransactionInputData(state: CustodianState, transaction: IKeyableCustodianTransaction): void {
        state.custodianTransactionsInputData = state.custodianTransactionsInputData.map((txn: IKeyableCustodianTransaction) => {
            if (txn.key === transaction.key) {
                return transaction;
            }
            return txn;
        })
    },
    addSchemeUnderCustodyInput(state: CustodianState, key: number): void {
        state.schemesUnderCustodyInputData.push({
            key,
            nameOfFirm: "",
            nameOfScheme: "",
            relationshipWithTrustee: "",
            volume: 0,
            markedToMarketValue: 0
        });
        state.schemesUnderCustodyInputError.push({
            key,
            nameOfFirm: false,
            nameOfScheme: false,
            relationshipWithTrustee: false,
            volume: false,
            markedToMarketValue: false
        });
    },
    deleteSchemeUnderCustodyInput(state: CustodianState, key: number): void {
        state.schemesUnderCustodyInputData = state.schemesUnderCustodyInputData.filter(share => share.key !== key);
        state.schemesUnderCustodyInputError = state.schemesUnderCustodyInputError.filter(share => share.key !== key);
    },
    modifySchemeUnderCustodyInputData(state: CustodianState, scheme: IKeyableSchemesUnderCustody): void {
        state.schemesUnderCustodyInputData = state.schemesUnderCustodyInputData.map((_scheme: IKeyableSchemesUnderCustody) => {
            if (_scheme.key === scheme.key) {
                return scheme;
            }
            return _scheme;
        })
    },
    validateGovernanceMultipleFieldData(state: CustodianState) {
        state.fieldHasError = false;
        state.ordinarySharesInputData.forEach((share: IKeyableShare) => {
            const error: IKeyableShareError = {key: share.key, name: false, percentage: false, shareholding: false};

            if (!share.name) {
                error.name = true;
            }
            if (!isNumeric(share.percentage)) {
                error.percentage = true;
            }
            if (!isNumeric(share.shareholding)) {
                error.shareholding = true;
            }

            if (error.name || error.percentage || error.shareholding) {
                state.fieldHasError = true;
            }
            state.ordinarySharesInputError = state.ordinarySharesInputError.map((err: IKeyableShareError) => {
                if (error.key === err.key) {
                    return error;
                }
                return err;
            });
        });
        state.preferenceSharesInputData.forEach((share: IKeyableShare) => {
            const error: IKeyableShareError = {key: share.key, name: false, percentage: false, shareholding: false};
            if (!share.name) {
                error.name = true;
            }
            if (!isNumeric(share.percentage)) {
                error.percentage = true;
            }
            if (!isNumeric(share.shareholding)) {
                error.shareholding = true;
            }

            if (error.name || error.percentage || error.shareholding) {
                state.fieldHasError = true;
            }
            state.preferenceSharesInputError = state.preferenceSharesInputError.map((err: IKeyableShareError) => {
                if (error.key === err.key) {
                    return error;
                }
                return err;
            });
        });
        state.custodianTransactionsInputData.forEach((transaction: IKeyableCustodianTransaction) => {
            const error: IKeyableCustodianTransactionError = {
                key: transaction.key,
                nameOfTrustee: false,
                relationshipWithTrustee: false,
                typeOfTransaction: false,
                amount: false
            };

            if (!transaction.nameOfTrustee) {
                error.nameOfTrustee = true;
            }
            if (!transaction.relationshipWithTrustee) {
                error.relationshipWithTrustee = true;
            }
            if (!transaction.typeOfTransaction) {
                error.typeOfTransaction = true;
            }
            if (!isNumeric(transaction.amount)) {
                error.amount = true;
            }

            if (error.nameOfTrustee || error.relationshipWithTrustee || error.typeOfTransaction || error.amount) {
                state.fieldHasError = true;
            }
            state.custodianTransactionsInputError = state.custodianTransactionsInputError.map((err: IKeyableCustodianTransactionError) => {
                if (error.key === err.key) {
                    return error;
                }
                return err;
            });
        });
        state.schemesUnderCustodyInputData.forEach((scheme: IKeyableSchemesUnderCustody) => {
            const error: IKeyableSchemesUnderCustodyError = {
                key: scheme.key,
                nameOfFirm: false,
                nameOfScheme: false,
                relationshipWithTrustee: false,
                volume: false,
                markedToMarketValue: false
            };
            if (!scheme.nameOfFirm) {
                error.nameOfFirm = true;
            }
            if (!scheme.nameOfScheme) {
                error.nameOfScheme = true;
            }
            if (!scheme.relationshipWithTrustee) {
                error.relationshipWithTrustee = true;
            }
            if (!scheme.relationshipWithTrustee) {
                error.relationshipWithTrustee = true;
            }
            if (!isNumeric(scheme.volume)) {
                error.volume = true;
            }
            if (!isNumeric(scheme.markedToMarketValue)) {
                error.markedToMarketValue = true;
            }
            if (error.nameOfFirm || error.relationshipWithTrustee || error.nameOfScheme || error.volume || error.markedToMarketValue) {
                state.fieldHasError = true;
            }
            state.schemesUnderCustodyInputError = state.schemesUnderCustodyInputError.map((err: IKeyableSchemesUnderCustodyError) => {
                if (error.key === err.key) {
                    return error;
                }
                return err;
            });
        });
    },
    clearAllData(state: CustodianState): void {
        state.ordinarySharesInputData = [];
        state.ordinarySharesInputError = [];
        state.preferenceSharesInputData = [];
        state.preferenceSharesInputError = [];
        state.custodianTransactionsInputData = [];
        state.custodianTransactionsInputError = [];
        state.schemesUnderCustodyInputData = [];
        state.schemesUnderCustodyInputError = [];
        state.fieldHasError = false;
    }
};