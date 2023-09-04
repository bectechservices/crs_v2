import Vue from "vue";
import axios from "axios";
import VueSweetalert2 from 'vue-sweetalert2';
import ExcelFilePreview from "./lib/excel";
import {formatMoney, parseStringToFloat} from "./lib/helpers";
import {endOfMonth, format, parse, subMonths} from "date-fns";
import {BillingQuarterlyReport, BillingSearchData, TransactionDetails} from "./lib/types";
import store from "../js/vuex"
import {BILLING_CONSTANTS} from "./vuex/constants";


Vue.use(VueSweetalert2);

interface Data {
    bpid: string;
    reportMonth: string;
    client: {
        name: string;
        safekeeping: string;
        bpid: string;
        sca: string;
        isNPRAClient: boolean
    },
    transactionDetails: Array<Array<string>>;
    hasUploadedFile: boolean;
    wantsToInputCurrency: boolean;
    currencyDetails: {
        currency: string;
        rate: number;
        date: string;
    },
    wantsToCalculateNAV: boolean;
    navInputMode: string;
    navValue: number;
    navGenerated: {
        position: number;
        cash_balance: number;
        liabilities: number;
        nav: number;
    },
    navStr: string;
    clientPositionChanged: boolean;
    clientReport: BillingQuarterlyReport
}

interface Methods {
    loadClientDetails: () => Promise<void>;
    generateInvoice: () => void;
    uploadTransactionDetails: () => Promise<void>;
    bpidEmpty: () => void;
    parseUploadedTransactionDetails: (event: any) => Promise<void>;
    makeTransactionDetailsFromStringArray: (data: Array<string>) => TransactionDetails;
    inputCurrencyDetails: () => void;
    uploadCurrencyDetails: () => Promise<void>;
    calculateClientNAV: () => void;
    saveClientNAV: () => Promise<void>;
    formatMoney: (amount: number) => string;
}

export default new Vue<Data, Methods>({
    el: '.billingDashboard',
    store,
    data: {
        bpid: "",
        reportMonth: format(subMonths(new Date(), 1), "YYYY-MM"),
        client: {
            name: "",
            safekeeping: "",
            bpid: "",
            isNPRAClient: true,
            sca: ""
        },
        transactionDetails: [],
        hasUploadedFile: false,
        wantsToInputCurrency: false,
        currencyDetails: {
            currency: "",
            rate: 0,
            date: ""
        },
        wantsToCalculateNAV: false,
        navInputMode: 'Generated',
        navValue: 0,
        navGenerated: {
            position: 0,
            cash_balance: 0,
            liabilities: 0,
            nav: 0
        },
        navStr: "",
        clientPositionChanged: false,
        clientReport: {
            id: 0,
            client_id: "",
            month: "",
            approved: true,
            approved_by: 0
        }
    },
    mounted() {
        const searchData = this.$store.state.billing.searchData;
        if (searchData.bpOrSca && searchData.date) {
            this.bpid = searchData.bpOrSca;
            this.reportMonth = searchData.date;
            this.$nextTick(() => {
                this.loadClientDetails().catch(console.log);
            })
        }
    },
    methods: {
        loadClientDetails: async function () {
            this.hasUploadedFile = false;
            this.clientPositionChanged = false;
            this.wantsToInputCurrency = false;
            this.wantsToCalculateNAV = false;
            const searchData: BillingSearchData = {bpOrSca: this.bpid, date: this.reportMonth}
            this.$store.commit(BILLING_CONSTANTS.STORE_BILLING_SEARCH_DATA, searchData);
            const response = await axios.post("/load-client-details", {bpid: this.bpid});
            const client = response.data.data.client;
            this.client.name = client.name;
            this.client.bpid = this.bpid;
            this.client.safekeeping = client.account_number;
            // this.client.isNPRAClient = client.type === "NPRA";
            this.client.isNPRAClient = true;
            if (client.id === 0) {
                //search using SCA
                const response2 = await axios.post("/load-client-details-with-code", {code: this.bpid});
                this.client.name = response2.data.data.client.client_name;
                this.client.bpid = response2.data.data.client.bpid;
                this.client.sca = response2.data.data.client.code;
                this.client.safekeeping = "";
                // this.client.isNPRAClient = false; //hack until we add sec/npra to unique codes
            }
            const positionRequest = await axios.post("/load-client-position", {
                bp_or_sca: this.bpid,
                date: this.reportMonth
            });
            this.navGenerated.position = positionRequest.data.data.position;
            this.navGenerated.nav = positionRequest.data.data.position;
            this.navStr = this.formatMoney(positionRequest.data.data.position);
            //Look out for wrong values
            const navRequest = await axios.post("/load-client-nav", {
                bp_or_sca: this.bpid,
                date: this.reportMonth
            });
            const navData = navRequest.data.data.nav;
            this.clientReport = navRequest.data.data.report;

            if (navData.id != 0 && navData.position == 0) {
                this.navInputMode = "Manual";
                this.navValue = navData.nav;
            } else {
                if (navData.position !== this.navGenerated.position) {
                    this.clientPositionChanged = true;
                }
                this.navGenerated.liabilities = navData.liabilities;
                this.navGenerated.nav = this.navGenerated.position - navData.liabilities;
            }
        },
        generateInvoice: function () {
            if (!this.client.name) {
                this.bpidEmpty();
            } else {
                window.location.assign(`/billing?bpOrSca=${this.bpid}&period=${format(endOfMonth(this.reportMonth + "-01"), "YYYY-MM-DD")}`);
            }
        },
        uploadTransactionDetails: async function () {
            if (!this.client.name) {
                this.bpidEmpty();
            } else {
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
                    const parsedData: Array<TransactionDetails> = this.transactionDetails.map((each: Array<string>): TransactionDetails => {
                        if (each[0] && each[0].trim()) {
                            each[0] = parse(each[0]) as any;
                        }
                        if (each[6] && each[6].trim()) {
                            each[6] = parseStringToFloat(each[6]) as any;
                        } else {
                            each[6] = 0 as any;
                        }
                        if (each[7] && each[7].trim()) {
                            each[7] = parseStringToFloat(each[7]) as any;
                        } else {
                            each[7] = 0 as any;
                        }
                        if (each[8] && each[8].trim()) {
                            each[8] = parseStringToFloat(each[8]) as any;
                        } else {
                            each[8] = 0 as any;
                        }
                        if (each[9] && each[9].trim()) {
                            each[9] = parseStringToFloat(each[9]) as any;
                        } else {
                            each[9] = 0 as any;
                        }
                        return this.makeTransactionDetailsFromStringArray(each)
                    });
                    const response = await axios.post("/upload-billing-transaction-details", {
                        transactions: parsedData
                    });
                    if (response.data.error) {
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
                            title: 'Transaction details uploaded'
                        });
                        this.transactionDetails = [];
                    }
                }
            }
        },
        bpidEmpty: function () {
            this.$swal({
                type: 'error',
                title: 'Oops...',
                text: 'Please select a client to proceed'
            });
        },
        parseUploadedTransactionDetails: async function (event: any) {
            this.hasUploadedFile = true;
            this.wantsToInputCurrency = false;
            this.wantsToCalculateNAV = false;
            const parser = new ExcelFilePreview(event.target.files[0]);
            const data = await parser.parseDataToJson();
            data.forEach((details: Array<string>) => {
                if (details.length > 0 && details[0].toLowerCase() !== "date") {
                    this.transactionDetails.push(details)
                }
            });
        },
        makeTransactionDetailsFromStringArray: function (data: Array<string>): TransactionDetails {
            return {
                bpid: this.client.bpid,
                reporting_date: endOfMonth(this.reportMonth + "-01") as any,
                sca: this.client.sca,
                date: data[0],
                reference: data[1],
                security_name: data[2],
                security_category: data[3],
                charge_type: data[4],
                charge_item: data[5],
                number_of_units: data[6] as any,
                market_value: data[7] as any,
                charge_amount_with_tax: data[8] as any,
                invoice_amount_with_tax: data[9] as any,
            }
        },
        inputCurrencyDetails: function () {
            this.wantsToInputCurrency = true;
            this.wantsToCalculateNAV = false;
            this.hasUploadedFile = false;
        },
        uploadCurrencyDetails: async function () {
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
                const response = await axios.post("/upload-billing-currency-details", {
                    currency: this.currencyDetails.currency,
                    rate: parseFloat(this.currencyDetails.rate as any),
                    date: parse(this.currencyDetails.date),
                    bpid: this.client.bpid,
                    sca: this.client.sca,
                    reporting_date: endOfMonth(this.reportMonth + "-01")
                });
                if (response.data.error) {
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
                        title: 'Currency details uploaded'
                    });
                    this.currencyDetails = {
                        currency: "",
                        rate: 0,
                        date: ""
                    };
                }
            }
        },
        calculateClientNAV: function () {
            this.wantsToInputCurrency = false;
            this.wantsToCalculateNAV = true;
            this.hasUploadedFile = false;
        },
        saveClientNAV: async function () {
            let hasError = false;
            if (this.navInputMode == 'Manual') {
                if (!this.navValue) {
                    hasError = true;
                    this.$swal({
                        type: 'error',
                        title: "The NAV input is required.",
                        animation: false
                    }).then(console.log);
                }
            } else {
                if (this.client.isNPRAClient) {
                    if (!this.navGenerated.liabilities) {
                        hasError = true;
                        this.$swal({
                            type: 'error',
                            title: "The liabilities input is required.",
                            animation: false
                        }).then(console.log);
                    }
                } else {
                    if (!this.navGenerated.cash_balance) {
                        hasError = true;
                        this.$swal({
                            type: 'error',
                            title: "The cash balance input is required.",
                            animation: false
                        }).then(console.log);
                    }
                }
            }
            if (!hasError) {
                const result = await this.$swal({
                    title: 'Are you sure?',
                    text: "This will update the client's NAV!",
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#3085d6',
                    cancelButtonColor: '#d33',
                    confirmButtonText: 'Yes'
                });
                if (result.value) {
                    let NAV = 0;
                    if (this.navInputMode == 'Manual') {
                        NAV = this.navValue;
                    } else {
                        NAV = this.navGenerated.nav;
                    }
                    const response = await axios.post("/update-clients-nav", {
                        bpid: this.client.bpid,
                        sca: this.client.sca,
                        nav: parseFloat(NAV as any),
                        position: parseFloat(this.navGenerated.position as any),
                        cash_balance: parseFloat(this.navGenerated.cash_balance as any),
                        liabilities: parseFloat(this.navGenerated.liabilities as any),
                        date: this.reportMonth
                    });
                    if (response.data.error) {
                        this.$swal({
                            toast: true,
                            position: 'top-end',
                            showConfirmButton: false,
                            timer: 3000,
                            type: 'error',
                            title: 'Something went wrong. Please try again'
                        })
                    } else {
                        this.clientPositionChanged = false;
                        this.$swal({
                            toast: true,
                            position: 'top-end',
                            showConfirmButton: false,
                            timer: 3000,
                            type: 'success',
                            title: 'NAV value saved!'
                        });
                    }

                }
            }
        },
        formatMoney(amount: number): string {
            return formatMoney(amount)
        },
    },
    watch:
        {
            'navGenerated.cash_balance':

                function (newVal) {
                    this.navGenerated.nav = this.navGenerated.position - newVal;
                    this.navStr = this.formatMoney(this.navGenerated.position - newVal);
                }

            ,
            'navGenerated.liabilities':

                function (newVal) {
                    this.navGenerated.nav = this.navGenerated.position - newVal;
                    this.navStr = this.formatMoney(this.navGenerated.position - newVal);
                }
        }
})