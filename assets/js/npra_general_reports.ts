import Vue from "vue";
import axios from "axios";
import {format} from "date-fns";
import collect, {Collection} from "collect.js";
import {
    FundAdministrator,
    GroupedOutstandingFDCertificate,
    IFund,
    IOutstandingFDCertificate,
    VarianceData
} from "./lib/types";
import {
    CURRENT_YEAR,
    currentQuarterNumber,
    formatMoney,
    getPreviousQuarterShortDate,
    getShortDate,
    makePreviousQuarterShortDate,
    makeShortDate
} from "./lib/helpers";
//@ts-ignore
import download from "downloadjs"
//@ts-ignore
import moment from "moment";


interface Data {
    quarter: string;
    year: string;
    outstandingFDCertificates: Array<GroupedOutstandingFDCertificate>;
    quarterlyReport: {
        pensions: Array<IFund>;
        provident: Array<IFund>;
        administrators: Array<FundAdministrator>
    },
    selectedQuarterShortDate: string;
    lastAUADate: string;
    localVariance: Array<VarianceData>;
    foreignVariance: Array<VarianceData>;
}

interface Methods {
    formatMoney: (amount: number) => string;
    loadDataForQuarter: () => Promise<void>;
    loadLocalVariance: (quarter: string, year: string) => Promise<void>;
    loadForeignVariance: (quarter: string, year: string) => Promise<void>;
    exportLocalVariance: () => Promise<void>;
    exportOutstandingFD: () => Promise<void>;
    exportQuarterlyReport: () => Promise<void>;
    fdIsMatured: (date: string) => boolean;
}

export default new Vue<Data, Methods>({
    el: ".npraGeneralReports",
    async beforeMount() {
        await this.loadDataForQuarter();
    },
    data: {
        quarter: currentQuarterNumber(),
        year: CURRENT_YEAR,
        outstandingFDCertificates: [],
        quarterlyReport: {
            pensions: [],
            provident: [],
            administrators: []
        },
        selectedQuarterShortDate: getShortDate(),
        lastAUADate: getPreviousQuarterShortDate(),
        localVariance: [],
        foreignVariance: [],
    },
    methods: {
        formatMoney(amount: number): string {
            return formatMoney(amount);
        }, 
        async loadDataForQuarter(): Promise<void> {
            this.selectedQuarterShortDate = makeShortDate(this.quarter, this.year);
            this.lastAUADate = makePreviousQuarterShortDate(this.quarter, this.year);
            const reportResponse = await axios.post("/load-npra-quarterly-report", {
                year: this.year,
                quarter: this.quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                }
            });
            this.quarterlyReport = reportResponse.data.data;

            const outstandingFDResponse = await axios.post("/load-outstanding-fd-certificates", {
                year: this.year,
                quarter: this.quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (outstandingFDResponse.data.data.certificates.length) {
                const certificates = outstandingFDResponse.data.data.certificates.map((cert: any) => {
                    return {
                        ...cert,
                        created_at: format(cert.effectiveDate, "MMMM"),
                        effectiveDate: format(cert.effectiveDate, "YYYY-MM-DD"),
                        maturity: format(cert.maturity, "YYYY-MM-DD"),
                    }
                });
                const grouped = collect(certificates).groupBy('created_at').all();
                for (let group in grouped) {
                    this.outstandingFDCertificates.push({
                        month: group,
                        certificates: (grouped[group] as Collection<IOutstandingFDCertificate>).all()
                    });
                }
            } else {
                this.outstandingFDCertificates = [];
            }
            await this.loadLocalVariance(this.quarter, this.year);
            await this.loadForeignVariance(this.quarter, this.year);
        },
        loadLocalVariance: async function (quarter: string, year: string) {
            const response = await axios.post("/load-npra-local-variance", {
                year,
                quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                }
            });
            this.localVariance = response.data.data.variance;
        },
        loadForeignVariance: async function (quarter: string, year: string) {
            const response = await axios.post("/load-npra-foreign-variance", {
                year,
                quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                    
                }
            });
            this.foreignVariance = response.data.data.variance;
        },
        exportLocalVariance: async function () {
            const quarter = this.quarter;
            const year = this.year;
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
            const quarter = this.quarter;
            const year = this.year;
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
            const quarter = this.quarter;
            const year = this.year;
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
        fdIsMatured: function (date: string): boolean {
            return moment().isAfter(moment(date))
        }
    }
})