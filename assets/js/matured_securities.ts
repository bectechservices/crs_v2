import Vue from "vue";
import VueSweetalert2 from "vue-sweetalert2";
import ExcelFilePreview from "./lib/excel";
import axios from "axios";
import {getCurrentQuarterFormalDate, parseStringToFloat} from "./lib/helpers";

Vue.use(VueSweetalert2);

interface Data {
    loading: boolean;
    hasLoadedExcelFile: boolean;
    maturedSecurities: Array<Array<string>>;
    header: Array<string>
}

interface Methods {
    uploadDocument: () => Promise<void>;
    onFileChange: (event: any) => Promise<void>;
}

export default new Vue<Data, Methods>({
    el: '.maturedSecuritiesApp',
    data: {
        loading: false,
        hasLoadedExcelFile: false,
        maturedSecurities: [],
        header: []
    },
    methods: {
        uploadDocument: async function () {
            if (this.hasLoadedExcelFile) {
                this.loading = true;
                try {
                    const date = getCurrentQuarterFormalDate();
                    const response = await axios.post("/matured-securities", {
                        maturities: this.maturedSecurities.map((each: Array<string>) => {
                            return {
                                client: each[0],
                                issuer: each[1],
                                amount_invested: parseStringToFloat(each[2]),
                                value: parseStringToFloat(each[3]),
                                date
                            }
                        }),
                    }, {
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                    if (response.data.error) {
                        this.$swal({
                            toast: true,
                            position: 'top-end',
                            showConfirmButton: false,
                            timer: 5000,
                            type: 'error',
                            title: response.data.message
                        }).then(console.log)
                    } else {
                        this.$swal({
                            toast: true,
                            position: 'top-end',
                            showConfirmButton: false,
                            timer: 3000,
                            type: 'success',
                            title: 'Matured Securities Uploaded'
                        });
                        window.location.reload();
                    }
                } catch (error) {
                    this.$swal({
                        toast: true,
                        position: 'top-end',
                        showConfirmButton: false,
                        timer: 5000,
                        type: 'error',
                        title: error.message
                    }).then(console.log)
                } finally {
                    this.loading = false;
                }
            }
        },
        onFileChange: async function (event: any) {
            try {
                const parser = new ExcelFilePreview(event.target.files[0]);
                const data = (await parser.parseDataToJson()).filter((each: Array<string>) => each.length >= 4);
                this.header = data.shift() as Array<string>;
                this.maturedSecurities = data;
                this.$nextTick(() => {
                    this.hasLoadedExcelFile = true;
                })
            } catch (e) {
                this.hasLoadedExcelFile = false;
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 5000,
                    type: 'error',
                    title: e.message
                }).then(console.log)
            }
        }
    }
})