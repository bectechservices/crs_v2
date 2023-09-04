import Vue from "vue";
import axios from "axios";
// import {format} from "date-fns";
// import collect, {Collection} from "collect.js";
import {

} from "./lib/types";
import {
    CURRENT_YEAR
} from "./lib/helpers";
//@ts-ignore
import download from "downloadjs"
//@ts-ignore
//import moment from "moment";


interface Data {
    month: string;
    year: string;
    bpOrSca: string;
}

interface Methods {
}

export default new Vue<Data, Methods>({
    el: ".npra030Reports",

    data: {
        month: '06',
        year: CURRENT_YEAR,
        bpOrSca: '',
    },
    methods: {

        export301MonthlyReport: async function () {
            // this.month = makeShortDate(this.month, this.year);
             const month = this.month
             const year = this.year;
             
             console.log("##########301 month ##########", month)
             console.log("##########301 year ##########", year)
             const response = await axios.post("/npra-301-report-excel", {
                 year,
                 month
             }, {
                 headers: {
                     'Content-Type': 'application/json'
                 },
                 responseType: 'blob'
             });
             download(response.data, `NPRA301-${month}-${year} REPORT.xlsx`);
         },
         export302MonthlyReport: async function () {
            //this.date = makeMonthFormalDate(this.month, this.year);
             const month = this.month
            const bpOrSca = this.bpOrSca
             const year = this.year;
             console.log("##########302 month ##########", month)
             console.log("##########302 year ##########", year)
             console.log("##########302 bpOrSca ##########", bpOrSca)             
             const response = await axios.post("/npra-302-report-excel", {
                 year,
                 month,
                 bpOrSca
             }, {
                 headers: {
                     'Content-Type': 'application/json'
                 },
                 responseType: 'blob'
             });
             download(response.data, `NPRA302-${month}-${year} REPORT.xlsx`);
         } 
    }
})