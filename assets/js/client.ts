import Vue from "vue"
import {IKeyable, RangedBasisPoint} from "./lib/types";
import {CLIENT_ACTIONS_CONSTANTS, CLIENT_CONSTANTS} from "./vuex/constants";
import store from "./vuex"
import VueSweetalert2 from 'vue-sweetalert2';
import {formatMoney} from "./lib/helpers";


Vue.use(VueSweetalert2);

interface Data {
    basis_point_type: string;
    bpid: string;
    minimum_charge: number;
    basis_point: number;
    ranged_basis_point: Array<RangedBasisPoint>;
    bps: Array<IKeyable>;
    charge_per_txn: number;
    third_party_transfer: number;
}

interface Methods {
    addNewBP: () => void;
    deleteBP: (bp: RangedBasisPoint) => void;
    submitClientBillingInfo: () => Promise<void>;
}

Vue.component(
    "ranged-basis-point",
    {
        template: `
                   <tr>
                        <td>
                            <input @input="onMinimumAmountChange" class="assetFormInput h-100 w-100 uk-input" v-model="minimum" required/>
                        </td>
                        <td>
                            <input @input="onMaximumAmountChange" class="assetFormInput h-100 w-100 uk-input" v-model="maximum" required/>
                        </td>
                        <td>
                            <input @input="onBasisPointChange" class="assetFormInput h-100 w-100 uk-input" v-model="basisPoint" required/>
                        </td>
                        <td>
                        <button class="w-100" v-if="canAdd" @click="$emit('add-bp-clicked')" type="button">Add</button>
                        <button class="w-100" v-if="canDelete" @click="$emit('delete-bp-clicked')" type="button">Delete</button>
                        </td>
                    </tr>
  `,
        props: ["canAdd", "canDelete", "bp"],
        data: function () {
            const basisPointData: RangedBasisPoint = this.$store.getters.rangedBasisPointData(
                this.$props.bp.key
            );
            return {
                minimum: basisPointData ? formatMoney(basisPointData.minimum_amount.toFixed(2) as any) : "0",
                maximum: basisPointData ? formatMoney(basisPointData.maximum_amount.toFixed(2) as any) : "0",
                basisPoint: basisPointData ? basisPointData.basis_point : 0,
            };
        },
        methods: {
            commitChanges: function () {
                this.$nextTick(() => {
                    this.$store.commit(
                        CLIENT_CONSTANTS.MODIFY_RANGED_BASIS_POINT_DATA,
                        {
                            key: this.$props.bp.key,
                            minimum_amount: parseFloat(this.minimum.replace(/,/g, '')),
                            maximum_amount: parseFloat(this.maximum.replace(/,/g, '')),
                            basis_point: this.basisPoint
                        } as RangedBasisPoint
                    );
                });
            },
            onMinimumAmountChange: function ({target}: { target: HTMLInputElement }) {
                if (target.value) {
                    this.minimum = formatMoney(parseFloat(target.value.replace(/,/g, '')));
                    this.commitChanges();
                }
            },
            onMaximumAmountChange: function ({target}: { target: HTMLInputElement }) {
                if (target.value) {
                    this.maximum = formatMoney(parseFloat(target.value.replace(/,/g, '')));
                    this.commitChanges();
                }
            },
            onBasisPointChange: function ({target}: { target: HTMLInputElement }) {
                if (target.value) {
                    this.basisPoint = parseFloat(target.value);
                    this.commitChanges();
                }
            }
        }
    });

export default new Vue<Data, Methods>({
    el: '.clientEditPage',
    store,
    data: {
        basis_point_type: "Flat",
        bpid: "",
        minimum_charge: 0,
        basis_point: 0,
        ranged_basis_point: [],
        bps: [],
        charge_per_txn: 0,
        third_party_transfer: 0
    },
    beforeMount() {
        this.$store.commit(CLIENT_CONSTANTS.CLEAR_ALL_RANGED_BASIS_POINT_INPUTS);
        const data = (window as any).ClientEditData;
        this.bpid = data.bpid;
        this.minimum_charge = data.billingInfo.minimum_charge;
        this.charge_per_txn = data.billingInfo.charge_per_transaction;
        this.third_party_transfer = data.billingInfo.third_party_transfer;
        if (data.basisPoints.length == 0) {
            this.addNewBP();
        } else if (data.basisPoints.length == 1 && (data.basisPoints[0].minimum == 0) && (data.basisPoints[0].maximum == 0)) {
            this.basis_point = data.basisPoints[0].basis_points;
            this.addNewBP();
        } else {
            data.basisPoints.forEach((each: any) => {
                this.bps.push({
                    key: each.id
                });
                this.$store.commit(CLIENT_CONSTANTS.ADD_RANGED_BASIS_POINT_INPUT, each.id);
                this.$store.commit(
                    CLIENT_CONSTANTS.MODIFY_RANGED_BASIS_POINT_DATA,
                    {
                        key: each.id,
                        minimum_amount: each.minimum,
                        maximum_amount: each.maximum,
                        basis_point: each.basis_points
                    } as RangedBasisPoint
                );
            });
            this.basis_point_type = "Tiered";
        }
    },
    methods: {
        addNewBP: function () {
            const key = Math.random() * 1000000;
            this.bps.push({
                key
            });
            this.$store.commit(CLIENT_CONSTANTS.ADD_RANGED_BASIS_POINT_INPUT, key);
        },
        deleteBP: function (bp: RangedBasisPoint) {
            this.bps = this.bps.filter((_bp: IKeyable) => _bp.key != bp.key);
            this.$store.commit(CLIENT_CONSTANTS.DELETE_RANGED_BASIS_POINT_INPUT, bp.key);
        },
        submitClientBillingInfo: async function () {
            try {
                const response = await this.$store.dispatch(
                    CLIENT_ACTIONS_CONSTANTS.UPLOAD_CLIENT_BILLING_INFORMATION, {
                        basis_point_type: this.basis_point_type,
                        bpid: this.bpid,
                        minimum_charge: parseFloat(this.minimum_charge as any),
                        basis_point: parseFloat(this.basis_point as any),
                        charge_per_txn: parseFloat(this.charge_per_txn as any),
                        third_party_transfer: parseFloat(this.third_party_transfer as any)
                    }
                );
                console.log(response);
                if (response.error) {
                    this.$swal({
                        toast: true,
                        position: 'top-end',
                        showConfirmButton: false,
                        timer: 3000,
                        type: 'error',
                        title: 'Something went wrong. Please try again'
                    })
                } else {
                    this.$swal({
                        toast: true,
                        position: 'top-end',
                        showConfirmButton: false,
                        timer: 3000,
                        type: 'success',
                        title: 'Billing details updated'
                    })
                }
            } catch (e) {
                console.log(e);
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'error',
                    title: 'Something went wrong. Please try again'
                })
            }
        },
    }
})