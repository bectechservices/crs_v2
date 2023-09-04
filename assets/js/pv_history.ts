import Vue from "vue";
import {currentMonthNumber, currentYear, formatMoney, makeMonthFormalDate} from "./lib/helpers";
import {AuditablePV, FormattedPVData, IMaturable, PVReport, PVSummary} from "./lib/types";
import axios from "axios";
import {isBefore, parse} from "date-fns";
import collect from "collect.js";

interface Data {
    bpOrSca: string;
    year: string;
    month: string;
    type: string;
    reports: Array<AuditablePV>;
    pvDate: string;
}

interface Methods {
    loadPVData: () => Promise<void>;
    formatMoney: (amount: number) => string;
    sumOf: (summary: Array<PVSummary>, key: string) => string;
}

export default new Vue<Data, Methods>({
    el: '.pvHistory',
    data: {
        bpOrSca: '',
        type: 'sec',
        year: `${currentYear()}`,
        month: `${currentMonthNumber()}`,
        reports: [],
        pvDate: ''
    },
    methods: {
        loadPVData: async function () {
            const date = makeMonthFormalDate(this.month, this.year);
            const response = await axios.post("/pv-details-for-eyeballing", {
                bpid: this.bpOrSca,
                quarterDate: date,
                pv_type: this.type
            });
            this.pvDate = date;
            this.reports = response.data.data.reportData.map((each: any) => {
                const report = collect(each.reports).groupBy("security_type").items;
                let reports: Array<FormattedPVData> = [];
                for (let item in report) {
                    reports.push({
                        //Im so sorry.. i did a terrible job
                        bond: (report as any)[item].toArray()[0].security_type,
                        data: (report as any)[item].toArray().map((each: PVReport & IMaturable) => {
                            const date = parse(each.date_to);
                            let is_matured = false;
                            if (date.getFullYear() > 0) {
                                is_matured = isBefore(date, new Date());
                            }
                            return {...each, is_matured};
                        })
                    })
                }
                return {
                    client: each.client,
                    reports,
                    summary: each.summary.map((_summary: any) => {
                        return {
                            description: _summary.security_type,
                            norminal_value: _summary.nominal_value,
                            cummulative_cost: _summary.cumlative_cost,
                            value: _summary.lcy_amount,
                            percentage_of_total: _summary.percentage_of_total,
                        }
                    })
                }
            });
        },
        formatMoney: function (amount: number) {
            return formatMoney(amount);
        },
        sumOf: function (summary: Array<PVSummary>, key: string): string {
            return (collect(summary).sum(key) as number).toFixed(2)
        }
    }
});