import Vue from "vue";
import {currentMonthNumber, currentYear, formatMoney, makeMonthFormalDate} from "./lib/helpers";
import {Auditable0301, Formatted301Data, IMaturable, Npra0301} from "./lib/types";
import axios from "axios";
import {isBefore, parse} from "date-fns";
import collect from "collect.js";

interface Data {
    year: string;
    month: string;
    reports: Array<Auditable0301>;
    pvDate: string;
}

interface Methods {
    load0301Data: () => Promise<void>;
    formatMoney: (amount: number) => string;
}

export default new Vue<Data, Methods>({
    el: '.loadNpra0301',
    data: {
        year: `${currentYear()}`,
        month: `${currentMonthNumber()}`,
        reports: [],
        pvDate: ''
    },
    methods: {
        load0301Data: async function () {
            const date = makeMonthFormalDate(this.month, this.year);
            const response = await axios.post("/details-0301-for-eyeballing", {
                quarterDate: date
            });
            this.pvDate = date;
            console.log("*******Level 2*********", response)
            this.reports = response.data.data.reportData.map((each: any) => {
                console.log("*******Level 2.1 *********", this.reports)
                 const report = collect(each.reports).items;
                 console.log("*******Level 3 *********", report)
                let reports: Array<Formatted301Data> = [];

                for (let item in report) {
                    reports.push({
                        data: (report as any)[item].toArray().map((each: Npra0301 & IMaturable) => {
                            console.log("*******Level 4*********", reports)
                            const date = parse(each.reporting_date);
                            let is_matured = false;
                            if (date.getFullYear() > 0) {
                                console.log("*******Level 5*********", is_matured)
                                is_matured = isBefore(date, new Date());
                            }
                            return {...each, is_matured};
                        })
                    })
                }
                return {
                    reports

                }
            });
        },
        formatMoney: function (amount: number) {
            return formatMoney(amount);
        }
    }
});