import Vue from "vue";
import axios from "axios";
import collect from "collect.js";
import {AuditablePV, FormattedPVData, IMaturable, PVReport, PVSummary} from "./lib/types";
import VueSweetalert2 from 'vue-sweetalert2';
import {formatMoney} from "./lib/helpers";
import {isBefore, parse} from "date-fns";


Vue.use(VueSweetalert2);

interface Data {
    reports: Array<AuditablePV>;
    pvDate: string;
    pvType: string;
    bpid: string;
}

interface Methods {
    deletePV: () => Promise<void>;
    formatMoney: (amount: number) => string;
    sumOf: (summary: Array<PVSummary>, key: string) => string;
}

export default new Vue<Data, Methods>({
    el: '.auditPVPage',
    data: {
        reports: [],
        pvDate: '',
        pvType: "",
        bpid: ""
    },
    async beforeMount() {
        const url = new URL(window.location.href);
        const bpid = url.searchParams.get("bpid") as string;
        const quarterDate = url.searchParams.get("quarter") as string;
        const pvType = url.searchParams.get("type") as string;

        const response = await axios.post("/pv-details-for-eyeballing", {
            bpid,
            quarterDate,
            pv_type: pvType
        });
        this.pvDate = quarterDate.replace(/-/gi, "/");
        this.pvType = pvType;
        this.bpid = bpid;

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
    methods: {
        deletePV: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "You won't be able to revert this!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Delete'
            });
            if (result.value) {
                let response: any;
                if (this.pvType == "npra") {
                    response = await axios.delete("/delete-npra-uploaded-pv", {
                        data: {
                            bpid: this.bpid,
                            date: this.pvDate
                        }
                    });
                } else if (this.pvType == "sec") {
                    response = await axios.delete("/delete-sec-uploaded-pv", {
                        data: {
                            bpid: this.bpid,
                            date: this.pvDate
                        }
                    });
                } else if (this.pvType == "billing") {
                    response = await axios.delete("/delete-billing-uploaded-pv", {
                        data: {
                            bpid: this.bpid,
                            date: this.pvDate
                        }
                    });
                }
                if (!response.data.error) {
                    window.history.back();
                }
            }
        },
        formatMoney: function (amount: number) {
            return formatMoney(amount);
        },
        sumOf: function (summary: Array<PVSummary>, key: string): string {
            return (collect(summary).sum(key) as number).toFixed(2)
        }
    }
})