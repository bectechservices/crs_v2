import {GetterTree} from 'vuex';
import {IUnauthorizedTransaction, NPRAState, RootState} from '../../lib/types';

export const getters: GetterTree<NPRAState, RootState> = {
    unauthorizedTransactionsData(state: NPRAState) {
        return (key: number) => {
            return state.unauthorizedTransactions.find((txn: IUnauthorizedTransaction) => txn.key === key);
        }
    }
};