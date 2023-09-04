import Vue from "vue";
import {format} from "date-fns";
import {
    CURRENT_YEAR,
    currentQuarterNumber,
    formatMoney,
    getCurrentQuarterFormalDate,
    getShortDate,
    lastDayOfQuarter,
    makeNPRAReportShortDate,
    makeQuarterBeginning
} from "./lib/helpers";
import axios from "axios";
import {
    FundAdministrator,
    GroupedOutstandingFDCertificate,
    IFund,
    INPRADeclaration,
    IOutstandingFDCertificate,
    IUnauthorizedTransaction
} from "./lib/types";
import collect, {Collection} from "collect.js";
import VueSweetalert2 from 'vue-sweetalert2';
//@ts-ignore
import download from "downloadjs"
//@ts-ignore
import moment from "moment";


Vue.use(VueSweetalert2);

interface Data {
    unauthorisedLetter: {
        date: string;
        incidents: {
            from: string;
            to: string;
        },
        reporter: {
            name: string;
            position: string;
        }
    },
    declaration: INPRADeclaration,
    monthlyReport: {
        shortDate: string;
    }
    unauthorizedTransactions: Array<IUnauthorizedTransaction>;
    outstandingFDCertificates: Array<GroupedOutstandingFDCertificate>;
    quarterlyReport: {
        pensions: Array<IFund>;
        provident: Array<IFund>;
        administrators: Array<FundAdministrator>
    },
    currentQuarterShortDate: string;
    currentQuarterDate: string;
}

interface Methods {
    formatMoney: (amount: number) => string;
    formatDate: (date: string, layout: string) => string;
    saveOutstandingFDChanges: () => Promise<void>;
    handleApiResponse: (response: { error: boolean }, successMessage: string) => void;
    exportLocalVariance: () => Promise<void>;
    exportOutstandingFD: () => Promise<void>;
    exportQuarterlyReport: () => Promise<void>;
    exportUnauthorizedTransactions: () => Promise<void>;
    exportNPRALetter: () => Promise<void>;
    exportUnAuthorizedLetter: () => Promise<void>;
    fdIsMatured: (date: string) => boolean;
}


export default new Vue<Data, Methods>({
    el: ".npraReportPage",
    data: {
        unauthorisedLetter: {
            date: format(new Date(), "MMMM Do, YYYY"),
            incidents: {
                from: makeQuarterBeginning(),
                to: lastDayOfQuarter()
            }
            ,
            reporter: {
                name: "",
                position: ""
            }
        },
        declaration: {
            nameOfOfficer: "",
            designation: "",
            headOfCustodyServices: "",
            date: ""
        },
        monthlyReport: {
            shortDate: makeNPRAReportShortDate()
        },
        unauthorizedTransactions: [],
        outstandingFDCertificates: [],
        quarterlyReport: {
            pensions: [],
            provident: [],
            administrators: []
        },
        currentQuarterShortDate: getShortDate(),
        currentQuarterDate: getCurrentQuarterFormalDate()
    },
    async beforeMount() {
        const response = await axios.post("/load-current-npra-declaration", {}, {
            headers: {
                "Content-Type": "application/json"
            }
        });
        this.declaration = response.data.data.declaration;

        const txnResponse = await axios.post("/load-unauthorized-transactions", {
            year: CURRENT_YEAR,
            quarter: currentQuarterNumber()
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        });
        this.unauthorizedTransactions = txnResponse.data.data.transactions.map((txn: IUnauthorizedTransaction) => {
            return {...txn, date: format(txn.date, "YYYY-MM-DD")}
        });

        const certResponse = await axios.post("/load-outstanding-fd-certificates", {
            year: CURRENT_YEAR,
            quarter: currentQuarterNumber()
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        });
        const certificates = certResponse.data.data.certificates.map((cert: any, index: number) => {
            return {
                ...cert,
                created_at: format(cert.effectiveDate, "MMMM"),
                effectiveDate: format(cert.effectiveDate, "YYYY-MM-DD"),
                maturity: format(cert.maturity, "YYYY-MM-DD")
            }
        });
        const grouped = collect(certificates).groupBy('created_at').all();
        for (let group in grouped) {
            this.outstandingFDCertificates.push({
                month: group,
                certificates: (grouped[group] as Collection<IOutstandingFDCertificate>).all()
            });
        }

        const reportResponse = await axios.post("/load-npra-quarterly-report", {
            year: CURRENT_YEAR,
            quarter: currentQuarterNumber()
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        });
        this.quarterlyReport = reportResponse.data.data;

        const mailSenderInfoResponse = await axios.post("/get-mail-sender-info", {}, {headers: {"Content-Type": "application/json"}});
        this.unauthorisedLetter.reporter = mailSenderInfoResponse.data.data.sender;
    },
    methods: {
        formatMoney(amount: number): string {
            return formatMoney(amount)
        },
        formatDate(date: string, layout: string): string {
            return format(date, layout);
        },
        saveOutstandingFDChanges: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "This will update the values of the fd certificates!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Yes'
            });
            if (result.value) {
                let data: any = [];
                this.outstandingFDCertificates.forEach((each: any) => {
                    each.certificates.forEach((one: any) => {
                        data.push({id: one.id, value: one.receiptReceived == "1"})
                    })
                });
                const response = await axios.post("/update-outstanding-fds", {
                    data
                });
                this.handleApiResponse(response.data, "Outstanding receipt changes uploaded")
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
        exportLocalVariance: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/npra-local-variance-excel", {
                year,
                quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Local Variance Q${quarter}-${year}.xlsx`);
        },
        exportOutstandingFD: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/npra-outstanding-fd-excel", {
                year,
                quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Outstanding Receipt Q${quarter}-${year}.xlsx`);
        },
        exportQuarterlyReport: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/npra-monthly-report-excel", {
                year,
                quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `NPRA Q${quarter}-${year} REPORT.xlsx`);
        },

        exportNPRALetter: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/npra-monthly-report-word", {}, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `NPRA Q${quarter}-${year} REPORT.docx`);
        },
        exportUnAuthorizedLetter: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/npra-unauthorized-report-word", {}, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `NPRA Form Q${quarter}-${year}.docx`);
        },
        fdIsMatured: function (date: string): boolean {
            return moment().isAfter(moment(date))
        },
        exportUnauthorizedTransactions: async function () {
            const quarter = currentQuarterNumber();
            const year = CURRENT_YEAR;
            const response = await axios.post("/export-unauthorized-transactions", {
                year,
                quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `NPRA Q${quarter}-${year} UNAUTHORIZED TRANSACTIONS.xlsx`);
        }
    }
});