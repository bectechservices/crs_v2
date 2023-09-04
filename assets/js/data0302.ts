import Vue from "vue";
import {currentMonthNumber, currentYear, formatMoney, makeMonthFormalDate} from "./lib/helpers";
import {Auditable0302, Formatted302Data, IMaturable, Npra0302} from "./lib/types";
import axios from "axios";
import {isBefore, parse} from "date-fns";
import collect from "collect.js";

interface Data {
    bpOrSca: string;
    year: string;
    month: string;
    report0302: Array<Auditable0302>;
    pvDate: string;
}

interface Methods {
    load0302Data: () => Promise<void>;
    formatMoney: (amount: number) => string;
}

export default new Vue<Data, Methods>({
    el: '.data0302ReportPage',
    data: {
        bpOrSca: 'BP0568916',
        year: `${currentYear()}`,
        month: `6`,  //${currentMonthNumber()}
        report0302: [],
        pvDate: ''
    },
    methods: {
        load0302Data: async function () {
            console.log("##################Level 1 ###################")
            const date = makeMonthFormalDate(this.month, this.year);
            const response = await axios.post("/pv-details-for-npra0302", {
                bpid: this.bpOrSca,
                quarterDate: date
            });

            console.log("################## Level 2 ###################",response)
            this.pvDate = date;
            let finalData:any = [];
            this.report0302 = response.data.data.reportData.map((each: any) => {
                const report = collect(each.report0302).groupBy("client_code").items;
                // let report0302: Array<Formatted302Data> = [];
                for (let item in report) {
                    // report0302.push({
                    //     //Im so sorry.. i did a terrible job
                    //     bpid: (report as any)[item].toArray()[0].client_code,
                        finalData = finalData.concat((report as any)[item].toArray().map((each: Npra0302 & IMaturable) => {
                            const date = parse(each.reporting_date);
                            let is_matured = false;
                            if (date.getFullYear() > 0) {
                                
                                is_matured = isBefore(date, new Date());
                            }
                            return {...each, is_matured};
                        }))
                }
               // console.log("&&&&&&&&& Level 2 &&&&&&&&&&", finalData)
                return {
                 //   client: each.client,
                  //  report0302: finalData
                }                

            });
               this.report0302 = {client:'',data:finalData} as any
               console.log("########Final Level 3 ########")
        },
        formatMoney: function (amount: number) {
            return formatMoney(amount);
        },
    }
});