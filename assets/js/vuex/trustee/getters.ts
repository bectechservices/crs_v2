import {GetterTree} from 'vuex';
import {IUnidentifiedPayment, RootState, TrusteeState} from '../../lib/types';

export const getters: GetterTree<TrusteeState, RootState> = {
    unidentifiedPaymentsData(state: TrusteeState) {
        return (key: number) => {
            return state.unidentifiedPayments.find((payment: IUnidentifiedPayment) => payment.key === key);
        }
    },
};