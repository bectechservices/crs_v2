import {MutationTree} from 'vuex';
import {IOutstandingFDCertificate, IUnauthorizedTransaction, NPRAState} from '../../lib/types';
import {format} from 'date-fns';

export const mutations: MutationTree<NPRAState> = {
    addUnauthorizedTransactionInput(state: NPRAState, key: number): void {
        state.unauthorizedTransactions.push({
            clientName: "",
            key,
            date: format(new Date(), "YYYY-MM-DD"),
            txnDetails: "",
        });
    },
    deleteUnauthorizedTransactionInput(state: NPRAState, key: number): void {
        state.unauthorizedTransactions = state.unauthorizedTransactions.filter(txn => txn.key !== key);
    },
    modifyUnauthorizedTransactionData(state: NPRAState, txn: IUnauthorizedTransaction): void {
        state.unauthorizedTransactions = state.unauthorizedTransactions.map((_txn: IUnauthorizedTransaction) => {
            if (_txn.key === txn.key) {
                return txn;
            }
            return _txn;
        })
    },
    clearAllUnauthorizedTransactions(state: NPRAState) {
        state.unauthorizedTransactions = [];
    }
};