import Vue from "vue";
import axios from "axios";
import {formatMoney, makeNPRAReportShortDate, parseStringToFloat, parseStringToInt} from "./lib/helpers";
import {IKeyable, INPRADeclaration, IOutstandingFDCertificate, IUnauthorizedTransaction} from "./lib/types";
import {NPRA_ACTION_CONSTANTS, NPRA_CONSTANTS} from "./vuex/constants";
import store from "./vuex";
import {format, parse} from "date-fns";
import VueSweetalert2 from 'vue-sweetalert2';
import ExcelFilePreview from "./lib/excel";


Vue.use(VueSweetalert2);

interface Data {
    reportShortDate: string;
    unauthorizedTransactions: Array<IKeyable>;
    outstandingFDCertificates: Array<IOutstandingFDCertificate>;
    declaration: INPRADeclaration
}

interface Methods {
    addNewTxn: () => void;
    deleteTxn: (txn: IUnauthorizedTransaction) => void;
    uploadTransactions: () => Promise<void>;
    uploadCertificates: () => Promise<void>;
    saveNPRADeclaration: () => Promise<void>;
    handleApiResponse: (response: { error: boolean }, successMessage: string) => void;
    submitChanges: () => Promise<void>;
    outstandingFDFileUpload: (event: any) => Promise<void>;
    formatMoney: (amount: number) => string;
    formatDate: (date: string, layout: string) => string;
    headerContainsMonth: (header: string) => boolean;
}

Vue.component(
    "unauthorized-transactions",
    {
        template: `
                   <tr>
                        <td style="width:10%;">
                            <input @input="onClientNameChange" type="text" class="assetFormInput h-100 w-100 uk-input" v-model="clientName" required/>
                        </td>
                        <td>
                            <input @input="onTxnDetailsChange" type="text" class="assetFormInput h-100 w-100 uk-input" v-model="txnDetails" required/>
                        </td>
                        <td>
                            <input @input="onDateChange" type="date" class="assetFormInput h-100 w-100 uk-input" v-model="date" required/>
                        </td>
                        <td>
                        <button class="w-100" v-if="canAdd" @click="$emit('add-txn-clicked')" type="button">Add</button>
                        <button class="w-100" v-if="canDelete" @click="$emit('delete-txn-clicked')" type="button">Delete</button>
                        </td>
                    </tr>
  `,
        props: ["canAdd", "canDelete", "txn"],
        data: function () {
            const transaction: IUnauthorizedTransaction = this.$store.getters.unauthorizedTransactionsData(
                this.$props.txn.key
            );
            return {
                clientName: transaction ? transaction.clientName : "",
                txnDetails: transaction ? transaction.txnDetails : "",
                date: transaction ? transaction.date : "",
            };
        },
        methods: {
            commitChanges: function () {
                this.$nextTick(() => {
                    this.$store.commit(
                        NPRA_CONSTANTS.MODIFY_UNAUTHORIZED_TRANSACTION_DATA,
                        {
                            key: this.$props.txn.key,
                            clientName: this.clientName,
                            txnDetails: this.txnDetails,
                            date: this.date
                        }
                    );
                });
            },
            onClientNameChange: function ({target}: { target: HTMLInputElement }) {
                this.clientName = target.value;
                this.commitChanges();
            },
            onTxnDetailsChange: function ({target}: { target: HTMLInputElement }) {
                this.txnDetails = target.value;
                this.commitChanges();
            },
            onDateChange: function ({target}: { target: HTMLInputElement }) {
                this.date = target.value;
                this.commitChanges();
            }
        }
    })

;

export default new Vue<Data, Methods>({
    el: '.npraPage',
    store,
    data: {
        reportShortDate: makeNPRAReportShortDate(),
        unauthorizedTransactions: [],
        outstandingFDCertificates: [],
        declaration: {
            nameOfOfficer: "",
            designation: "",
            headOfCustodyServices: "",
            date: ""
        }
    },
    async beforeMount() {
        const transactions = this.$store.state.npra.unauthorizedTransactions;
        if (transactions.length) {
            transactions.forEach((txn: IUnauthorizedTransaction) => {
                this.unauthorizedTransactions.push({key: txn.key})
            })
        } else {
            this.addNewTxn();
        }

        const declarationResponse = await axios.post("/load-current-npra-declaration", {}, {
            headers: {
                "Content-Type": "application/json"
            }
        });
        this.declaration = declarationResponse.data.data.declaration;
    },
    methods: {
        addNewTxn: function () {
            const key = Math.random() * 1000000;
            this.unauthorizedTransactions.push({
                key
            });
            this.$store.commit(NPRA_CONSTANTS.ADD_UNAUTHORIZED_TRANSACTION_INPUT, key);
        },
        deleteTxn: function (txn: IUnauthorizedTransaction) {
            this.unauthorizedTransactions = this.unauthorizedTransactions.filter((_txn: IKeyable) => _txn.key != txn.key);
            this.$store.commit(NPRA_CONSTANTS.DELETE_UNAUTHORIZED_TRANSACTION_INPUT, txn.key);
        },
        uploadTransactions: async function () {
            try {
                const response = await this.$store.dispatch(
                    NPRA_ACTION_CONSTANTS.UPLOAD_UNAUTHORIZED_TRANSACTIONS
                );
                this.handleApiResponse(response, "Transactions uploaded");
                this.$store.commit(NPRA_CONSTANTS.CLEAR_ALL_UNAUTHORIZED_TRANSACTION);
                this.unauthorizedTransactions = [];
                this.addNewTxn()
            } catch (e) {
                console.log(e);
                this.handleApiResponse({error: true}, "")
            }
        },
        uploadCertificates: async function () {
            try {
                const response = await this.$store.dispatch(
                    NPRA_ACTION_CONSTANTS.UPLOAD_OUTSTANDING_FD_CERTIFICATES,
                    this.outstandingFDCertificates
                );
                this.handleApiResponse(response, "Outstanding certificates uploaded");
                this.outstandingFDCertificates = [];
            } catch (e) {
                console.log(e);
                this.handleApiResponse({error: true}, "")
            }
        },
        saveNPRADeclaration: async function () {
            if (this.declaration.nameOfOfficer &&
                this.declaration.designation &&
                this.declaration.headOfCustodyServices &&
                this.declaration.date) {
                const response = await axios.post("/upload-npra-declaration", {
                    nameOfOfficer: this.declaration.nameOfOfficer,
                    designation: this.declaration.designation,
                    headOfCustodyServices: this.declaration.headOfCustodyServices,
                    date: parse(this.declaration.date)
                }, {
                    headers: {
                        "Content-Type": "application/json"
                    }
                });
                this.handleApiResponse(response.data, "Declaration updated")
            } else {
                this.$swal({
                    type: 'error',
                    title: "All inputs are required.",
                    animation: false
                }).then(console.log);
            }
        },
        handleApiResponse: function (response: { error: boolean }, successMessage: string) {
            if (response.error) {
                console.log(response);
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'error',
                    title: "Something went wrong. Please try again."
                }).then(console.log)
            } else {
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'success',
                    title: successMessage
                }).then(console.log)
            }
        },
        submitChanges: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "You won't be able to revert this!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Yes'
            });
            if (result.value) {
                const response = await store.dispatch(NPRA_ACTION_CONSTANTS.SEND_CHANGES_TO_CHECKERS);
                this.handleApiResponse(response, "Changes Submitted");
                console.log(response)
            }
        },
        outstandingFDFileUpload: async function (event: any) {
            this.outstandingFDCertificates = [];
            const parser = new ExcelFilePreview(event.target.files[0]);
            const data = await parser.parseDataToJson();
            const months = ['JANUARY', 'FEBRUARY', 'MARCH', 'APRIL', 'MAY', 'JUNE', 'JULY', 'AUGUST', 'SEPTEMBER', 'OCTOBER', 'NOVEMBER', 'DECEMBER'];
            data.forEach((datum: Array<any>) => {
                if (datum.length > 0 && !this.headerContainsMonth(datum[0]) && (datum[0] !== "Fund Manager" && datum[1] !== "Client Name")) {
                    datum[2] = parseStringToFloat(datum[2]);
                    datum[4] = parseStringToFloat(datum[4]);
                    datum[5] = parseStringToInt(datum[5]);
                    datum[7] = parse(datum[7]);
                    datum[8] = parse(datum[8]);
                    this.outstandingFDCertificates.push({
                        fundManager: datum[0],
                        clientName: datum[1],
                        amount: datum[2],
                        issuer: datum[3],
                        rate: datum[4],
                        tenor: datum[5],
                        term: datum[6],
                        effectiveDate: datum[7],
                        maturity: datum[8],
                    });
                }
            });
        },
        formatMoney(amount: number): string {
            return formatMoney(amount)
        },
        formatDate(date: string, layout: string): string {
            return format(date, layout);
        },
        headerContainsMonth(header: string): boolean {
            const months = ['JANUARY', 'FEBRUARY', 'MARCH', 'APRIL', 'MAY', 'JUNE', 'JULY', 'AUGUST', 'SEPTEMBER', 'OCTOBER', 'NOVEMBER', 'DECEMBER'];
            let contains = false;
            months.forEach((each: string) => {
                if (header.toUpperCase() === each || header.toUpperCase().startsWith(each)) {
                    contains = true;
                    return;
                }
            });
            return contains;
        }
    }
})