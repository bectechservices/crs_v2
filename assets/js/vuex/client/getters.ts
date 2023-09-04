import {GetterTree} from 'vuex';
import {ClientState, RangedBasisPoint, RootState} from '../../lib/types';

export const getters: GetterTree<ClientState, RootState> = {
    rangedBasisPointData(state: ClientState) {
        return (key: number) => {
            return state.rangedBasisPoint.find((bp: RangedBasisPoint) => bp.key === key);
        }
    }
};