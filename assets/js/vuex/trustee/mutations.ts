import {MutationTree} from 'vuex';
import {IUnidentifiedPayment, TrusteeDashboardSearchData, TrusteeState, npra321SearchData, npra321State} from '../../lib/types';
import {format} from 'date-fns';

export const mutations: MutationTree<TrusteeState> = {
    addUnidentifiedPaymentInput(state: TrusteeState, key: number): void {
        state.unidentifiedPayments.push({
            client_bpid: "",
            key, txn_type: "Cash", txn_date: format(new Date(), "YYYY-MM-DD"),
            value_date: format(new Date(), "YYYY-MM-DD"),
            name_of_company: "",
            amount: 0,
            fund_manager: 0,
            collection_acc_num: "",
            status: "pending"
        });
    },
    deleteUnidentifiedPaymentInput(state: TrusteeState, key: number): void {
        state.unidentifiedPayments = state.unidentifiedPayments.filter(payment => payment.key !== key);
    },
    modifyUnidentifiedPaymentData(state: TrusteeState, payment: IUnidentifiedPayment): void {
        state.unidentifiedPayments = state.unidentifiedPayments.map((_payment: IUnidentifiedPayment) => {
            if (_payment.key === payment.key) {
                return payment;
            }
            return _payment;
        })
    },
    clearAllUnidentifiedPayments(state: TrusteeState) {
        state.unidentifiedPayments = [];
    },
    storeDashboardSearchInput(state: TrusteeState, data: TrusteeDashboardSearchData) {
        state.searchInputData = data;
    },
    // store321SearchInput(state: npra321State, data: npra321SearchData) {
    //     state.search321InputData = data;
    // }
};