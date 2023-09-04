import {MutationTree} from 'vuex';
import {ClientState, RangedBasisPoint} from '../../lib/types';

export const mutations: MutationTree<ClientState> = {
    addRangedBasisPointInput(state: ClientState, key: number): void {
        state.rangedBasisPoint.push({
            minimum_amount: 0,
            key,
            maximum_amount: 0,
            basis_point: 0,
        });
    },
    deleteRangedBasisPointInput(state: ClientState, key: number): void {
        state.rangedBasisPoint = state.rangedBasisPoint.filter(bp => bp.key !== key);
    },
    modifyRangedBasisPointData(state: ClientState, bp: RangedBasisPoint): void {
        state.rangedBasisPoint = state.rangedBasisPoint.map((_bp: RangedBasisPoint) => {
            if (_bp.key === bp.key) {
                return bp;
            }
            return _bp;
        })
    },
    clearAllRangedBasisPointInputs(state: ClientState) {
        state.rangedBasisPoint = [];
    }
};