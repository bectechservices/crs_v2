import Vue from "vue";
import axios from "axios";
import {FundManager, IKeyable, IUnidentifiedPayment} from "./lib/types";
import store from "./vuex";
import {TRUSTEE_ACTION_CONSTANTS, TRUSTEE_CONSTANTS} from "./vuex/constants";

interface Data {
    payments: Array<IKeyable>;
    fundManagers: Array<FundManager>;
    clientBPID: string;
    clientName: string;
}

interface Methods {
    addNewTxn: () => void;
    deleteTxn: (txn: IUnidentifiedPayment) => void;
    uploadPayments: () => Promise<void>;
}

Vue.component(
    "unidentified-payments",
    {
        template: `
                   <tr>
                        <td style="width:10%;">
                            <input @input="onTxnDateChange" type="date" class="assetFormInput h-100 w-100 uk-input" v-model="txn_date" required/>
                        </td>
                        <td>
                            <input @input="onValueDateChange" type="date" class="assetFormInput h-100 w-100 uk-input" v-model="value_date" required/>
                        </td>
                        <td>
                            <input @input="onNameOfCompanyChange" type="text" class="assetFormInput h-100 w-100 uk-input" v-model="name_of_company" required/>
                        </td>
                        <td>
                            <select v-model="txn_type" @change="onTxnTypeChange" required>
                                <option>Cash</option>
                                <option>Cheque</option>
                                <option>Transfer</option>
                                <option>Others</option>
                            </select>
                        </td>
                        <td>
                            <input @input="onAmountChange" type="number"  min="0" class="assetFormInput h-100 w-100 uk-input" v-model="amount" required/>
                        </td>
                        <td style="width:15%;">
                            <select v-model="fund_manager" @change="onFundManagerChange" class="uk-select" required>
                                <option v-for="(manager,key) in fundManagers" :key="key" :value="manager.id">ACC NO. {{manager.account_number.substring(manager.account_number.length - 3,manager.account_number.length)}}({{manager.name}})</option>
                            </select>
                        </td>
                        <td>
                            <input @input="onCollectionAccountChange" type="text" class="assetFormInput h-100 w-100 uk-input" v-model="collection_acc_num" required/>
                        </td>
                        <td>
                            <div class="uk-grid-small uk-child-width-auto uk-grid">
                                <label><input type="radio" value="pending" class="uk-radio" v-model="status" @change="onStatusChange">Pending</label>
                                <label><input type="radio" value="done" class="uk-radio" v-model="status" @change="onStatusChange">Done</label>
                            </div>
                        </td>
                        <td>
                        <button class="w-100" v-if="canAdd" @click="$emit('add-txn-clicked')" type="button">Add</button>
                        <button class="w-100" v-if="canDelete" @click="$emit('delete-txn-clicked')" type="button">Delete</button>
                        </td>
                    </tr>
  `,
        props: ["canAdd", "canDelete", "fundManagers", "txn"],
        data: function () {
            const payment: IUnidentifiedPayment = this.$store.getters.unidentifiedPaymentsData(
                this.$props.txn.key
            );
            return {
                txn_date: payment ? payment.txn_date : "",
                txn_type: payment ? payment.txn_type : "Cash",
                value_date: payment ? payment.value_date : "",
                name_of_company: payment ? payment.name_of_company : "",
                amount: payment ? payment.amount : "",
                fund_manager: payment ? (payment.fund_manager > 0 ? payment.fund_manager : this.$props.fundManagers[0].id) : this.$props.fundManagers[0].id,
                collection_acc_num: payment ? (payment.collection_acc_num ? payment.collection_acc_num : this.$props.fundManagers[0].account_number) : this.$props.fundManagers[0].account_number,
                status: payment ? payment.status : "pending"
            };
        },
        methods: {
            commitChanges: function () {
                this.$nextTick(() => {
                    this.$store.commit(
                        TRUSTEE_CONSTANTS.MODIFY_UNIDENTIFIED_PAYMENT_DATA,
                        {
                            key: this.$props.txn.key,
                            txn_date: this.txn_date,
                            txn_type: this.txn_type,
                            value_date: this.value_date,
                            name_of_company: this.name_of_company,
                            amount: this.amount,
                            fund_manager: this.fund_manager,
                            collection_acc_num: this.collection_acc_num,
                            status: this.status
                        }
                    );
                });
            },
            onTxnDateChange: function ({target}: { target: HTMLInputElement }) {
                this.txn_date = target.value;
                this.commitChanges();
            },
            onValueDateChange: function ({target}: { target: HTMLInputElement }) {
                this.value_date = target.value;
                this.commitChanges();
            },
            onNameOfCompanyChange: function ({target}: { target: HTMLInputElement }) {
                this.name_of_company = target.value;
                this.commitChanges();
            },
            onAmountChange: function ({target}: { target: HTMLInputElement }) {
                this.amount = parseFloat(target.value);
                this.commitChanges();
            },
            onFundManagerChange: function ({target}: { target: HTMLInputElement }) {
                this.fund_manager = parseInt(target.value);
                this.collection_acc_num = this.$props.fundManagers.find((manager: FundManager) => manager.id == this.fund_manager).account_number;
                this.commitChanges();
            },
            onCollectionAccountChange: function ({target}: { target: HTMLInputElement }) {
                this.collection_acc_num = target.value;
                this.commitChanges();
            },
            onStatusChange: function ({target}: { target: HTMLInputElement }) {
                this.status = target.value as "pending" | "done";
                this.commitChanges();
            },
            onTxnTypeChange: function ({target}: { target: HTMLInputElement }) {
                this.txn_type = target.value;
                this.commitChanges();
            },
        }
    })

;

export default new Vue<Data, Methods>({
    el: '.unidentifiedPaymentsPage',
    store,
    beforeMount: async function () {
        const response = await axios.post("/load-fund-managers", {}, {
            headers: {
                'Content-Type': 'application/json'
            }
        });
        this.fundManagers = response.data.data.managers;
        const payments = this.$store.state.trustee.unidentifiedPayments;
        if (payments.length) {
            payments.forEach((payment: IUnidentifiedPayment) => {
                this.payments.push({key: payment.key})
            })
        } else {
            this.addNewTxn();
        }
    },
    data: {
        payments: [],
        fundManagers: [],
        clientBPID: "",
        clientName: ""
    },
    methods: {
        addNewTxn: function () {
            const key = Math.random() * 1000000;
            this.payments.push({
                key
            });
            this.$store.commit(TRUSTEE_CONSTANTS.ADD_UNIDENTIFIED_PAYMENT_INPUT, key);
        },
        deleteTxn: function (txn: IUnidentifiedPayment) {
            this.payments = this.payments.filter((payment: IKeyable) => payment.key != txn.key);
            this.$store.commit(TRUSTEE_CONSTANTS.DELETE_UNIDENTIFIED_PAYMENT_INPUT, txn.key);
        },
        uploadPayments: async function () {
            if (this.clientBPID) {
                try {
                    const response = await this.$store.dispatch(
                        TRUSTEE_ACTION_CONSTANTS.UPLOAD_UNIDENTIFIED_PAYMENTS,
                        this.$store.state.trustee.unidentifiedPayments.map((payment: IUnidentifiedPayment) => {
                            return {...payment, client_bpid: this.clientBPID}
                        })
                    );
                    console.log(response);
                    this.$store.commit(TRUSTEE_CONSTANTS.CLEAR_ALL_UNIDENTIFIED_PAYMENTS);
                    this.payments = [];
                    this.addNewTxn()
                } catch (e) {
                    console.log(e);
                }
            } else {
                alert("Select a client to proceed");
            }
        }
    },
    watch: {
        clientBPID: function (current: string) {
            this.clientName = (this.$refs[current] as any).innerText;
        }
    }
})