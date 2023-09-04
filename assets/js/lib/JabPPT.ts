import * as PptxGenJS from "pptxgenjs/dist/pptxgen";
import {PPTTemplate} from "./types";
import {
    arrayValuesToPercentage,
    getShortMonthNamesForPastMonths,
    getShortMonthNamesForPastMonthsRelativeQuarter,
    makePastQuarters
} from "./helpers";

export default class JabPPT {
    private pptx = new PptxGenJS();
    private COMPANY_NAME = "BECTECH LTD";
    private CUST_NAME = "Standard Chartered Bank";
    private CLIENT: any[];
    private MISC: any[];

    private BRANDCOLOR_BLUE = `${this.rgbToHex(31, 73, 125)}`;
    private BRANDCOLOR_BLUELGHT = `${this.rgbToHex(74, 126, 187)}`;
    private BRANDCOLOR_GREEN = `${this.rgbToHex(119, 147, 60)}`;
    private BRANDCOLOR_RED = `${this.rgbToHex(175, 0, 0)}`;
    private TEXTCOLOR = "999999";
    private HEADERCOLOR = "525355";
    private BLUEBLACK = "0B2135";
    private HUELGHT = "ffffff";
    private HUEDARK = "000000";
    private SCBBLUELIGHTER = "CDE3FB";
    private DEFFONTFACE = "Arial";
    private BASEFONTSIZE = 16;
    private HEADERTITLESIZE = 30;

    private xAxisTextPosition = 8;
    private yAxisTextPosition = 6.5;
    private baseTextHeight =1.65;

    private assetPaths = {
        clientLogo: {path: "./media/logo/twifo.png"},
    };

    private assets = {
        stanChartBgLogo: "./media/sclogo.png",
        sclogoStandAlone: "./media/logo/standard-chartered-mobile.png",
        stanChartBgImage: "./media/Picture1.jpg",
        scbGradientA: "./media/gradientA.png",
        scbGradientB: "./media/gradientB.png",
        scbTextLogo: "./media/sctext.png",
        clientLogo: "",
        pvReportIcon: "./media/logo/twifo.png"
    }

    private dataChartPieLocs = [
        {
            name: "AUC By Assets Class",
            labels: [],
            values: []
        }
    ];
    private aucTrendChart = [
        {
            // Should be the year
            name: "2019",
            values: [],
            labels: []
        }
    ];
    private trnsVolAssetChart = [
        {},
        {}
    ];

    private arrDataRegions = [
        {
            name: "Monthly Trade Volume",
            labels: [],
            values: []
        }
    ];

    private fdMaturitiesData = [
        {
            name: "FD Maturities",
            labels: [],
            values: []
        }
    ];
    private pptTemplateOptions: PPTTemplate;

    constructor(client: any[], misc: any[], options: PPTTemplate) {
        this.CLIENT = client;
        this.MISC = misc;
        this.pptTemplateOptions = options;
        this.pptx.setAuthor(this.CUST_NAME);
        this.pptx.setCompany(this.COMPANY_NAME);
        this.pptx.setRevision("20");
        this.pptx.setSubject("Trustee Report");
        this.pptx.setTitle(`${misc[0].report_name} Trustee Report Powerpoint`);
        this.pptx.setLayout("LAYOUT_WIDE");
        this.assetPaths.clientLogo = {path: `/file-manager/${client[0].image}`};
        this.assets.clientLogo = `/file-manager/${client[0].image}`;
        for (let j = 0; j < client.length; j++) {
            this.dataChartPieLocs[j] = {
                name: "AUC By Assets Class",
                labels: misc[j].asset_by_class.map((each: any) => each.dataName),
                values: misc[j].asset_by_class.map((each: any) => each.dataValue)
            };
            this.aucTrendChart[j] = {
                name: "AUC",
                values: misc[j].chart_trend_data,
                labels: getShortMonthNamesForPastMonths(misc[j].current_quarter_number)
            } as any;
            this.arrDataRegions[j] = {
                name: "Trade Volume",
                labels: getShortMonthNamesForPastMonthsRelativeQuarter(misc[j].current_quarter_number),
                values: misc[j].txn_trend_vols
            } as any;
            this.fdMaturitiesData[j] = {
                name: "FD Maturities",
                labels: Object.keys(misc[j].fdMaturities),
                values: Object.values(misc[j].fdMaturities)
            } as any;
            this.trnsVolAssetChart[j] = misc[j].tradeByAssetClassData;
        }
    }

    private rgbToHex(r: number, g: number, b: number) {
        return ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
    }

    private getTimestamp() {
        let dateNow = new Date();
        let dateMM = dateNow.getMonth() + 1;
        let dateDD = dateNow.getDate();
        let dateYY = dateNow.getFullYear();
        let h = dateNow.getHours();
        let m = dateNow.getMinutes();

        return (
            dateNow.getFullYear() +
            "" +
            (dateMM <= 9 ? "0" + dateMM : dateMM) +
            "" +
            (dateDD <= 9 ? "0" + dateDD : dateDD) +
            (h <= 9 ? "0" + h : h) +
            (m <= 9 ? "0" + m : m)
        );
    }


    private genSlides_Content() {
        {
            let cover_page_slide = this.pptx.addNewSlide("COVER_SLIDE");
            cover_page_slide.addText(`${this.MISC[0].report_name} Trustee Meeting`, {
                x: "25%",
                y: "15%",
                w: "70%",
                h: "7%",
                fontFace: this.DEFFONTFACE,
                color: this.HUELGHT,
                bold: true,
                fontSize: `${this.HEADERTITLESIZE - 5 }`,
                align: "l",
                valign: "m",
                margin: 0.4
            } as any);
            cover_page_slide.addText(`${this.MISC[0].quarter} Custodian Report`, {
                x: "25%",
                y: "25%",
                w: "70%",
                h: "6%",
                fontFace: this.DEFFONTFACE,
                color: `${this.rgbToHex(255, 255, 255)}`,
                fontSize: `${this.BASEFONTSIZE + 4}`,
                align: "l",
                valign: "m",
                margin: 0.4
            } as any);            
            cover_page_slide.addImage({
                x: "2%",
                y: 0.0,
                w: "11%",
                h: "13%",
                path: this.assets.clientLogo,
                sizing: {type: "contain", w: "11%", h: "13%"}
            } as any);
        }

        // AGENDA SLIDE
        {
            let agendaSlide = this.pptx.addNewSlide("PAGE_SLIDES");
            agendaSlide.addText("Agenda", {
                x: "4%",
                y: "6%",
                fontSize: `${this.HEADERTITLESIZE - 6}`,
                fontFace: this.DEFFONTFACE,
                color: this.HEADERCOLOR,
                fontWeight:"bold",
                bold: true,
                align: "l",
                valign: "middle"
            } as any);

            // TOC
            let tocHeadings = [];
            if (this.pptTemplateOptions.pv_report) {
                tocHeadings.push("Summary of Valuation Report")
            }
            if (this.pptTemplateOptions.total_summary_of_auc) {
                tocHeadings.push("Summary of Assets under custody")
            }
            if (this.pptTemplateOptions.auc_trend) {
                tocHeadings.push("AUC Trends")
            }
            if (this.pptTemplateOptions.trade_volumes) {
                tocHeadings.push("Trade Volumes")
            }
            if (this.pptTemplateOptions.total_contribution) {
                tocHeadings.push("Total Contributions")
            }
            if (this.pptTemplateOptions.corporate_action) {
                tocHeadings.push("Corporate Actions")
            }
            if (this.pptTemplateOptions.gog_and_fd_maturities) {
                tocHeadings.push("GOG and FD Maturities")
            }
            if (this.pptTemplateOptions.appendix_i) {
                tocHeadings.push("Appendix I- AUC per Fund Manager")
            }
            if (this.pptTemplateOptions.appendix_ii) {
                tocHeadings.push("Appendix II- Contributions")
            }
            if (this.pptTemplateOptions.unidentified_payments) {
                tocHeadings.push("Unidentified Payments")
            }

            const shouldUseNewStyle = tocHeadings.length <= 6;
            let yValTxt;
            tocHeadings.forEach((value: string, key: number) => {
                if (shouldUseNewStyle) {
                    yValTxt = (key + 1) + (0.43 * key) + 1.64;
                } else {
                    yValTxt = (key + 1) + (0.047 * key) + 0.95;
                }  

                agendaSlide.addText([
                    {
                        text:value,
                        options:{bullet:{code:'274F'}}
                    }],
                   {
                        x: `${this.xAxisTextPosition}%`,
                        y: `${(this.yAxisTextPosition * yValTxt) - 0.26}%`,
                        h: `${this.baseTextHeight}%`,
                        fontFace: this.DEFFONTFACE,
                        fontSize: `${this.BASEFONTSIZE}`,
                        color: this.HUEDARK,
                        valign: "middle"
                } as any);
            });
        }

        // TIERED CLIENT
        for (let i = 0; i < this.CLIENT.length; i++) {
            if (this.CLIENT.length == 2) {
                this.tieredSlide(this.CLIENT[i].name, i);
            }
            if (this.pptTemplateOptions.pv_report) {
                this.aucSummaryTblSlide(i);
            }
            if (this.pptTemplateOptions.total_summary_of_auc) {
                this.aucSummarySlide2(i);
            }
            if (this.pptTemplateOptions.auc_trend) {
                this.aucTrendSlide(i);
            }
            if (this.pptTemplateOptions.trade_volumes) {
                this.transactionVolSlide(i);
            }
            this.pvReportSlide();
            if (this.pptTemplateOptions.total_contribution) {
                this.monthlyContributionTableSlide(i);
            }
            if (this.pptTemplateOptions.corporate_action) {
                this.corporateActionActivitiesSlide(i);
            }
            if (this.pptTemplateOptions.gog_and_fd_maturities) {
                this.maturitiesForGOGSlide(i);
            }
            if (this.pptTemplateOptions.appendix_i) {
                this.MISC[i].summaryForFDs.forEach((fd: string) => {
                    this.pvSummaryForFundManagersSlide(fd, i);
                });
            }
            if (this.pptTemplateOptions.appendix_ii) {
                this.MISC[i].contributionsForFDs.forEach((fd: string) => {
                    this.contributionsForFundManagersSlide(fd, i);
                });
            }
            if (this.pptTemplateOptions.unidentified_payments) {
                this.unidentifiedPaymentsSlide(i);
            }
        }

        // THANK YOU SLIDE
        this.pptx.defineSlideMaster({
            title: "THANK_YOU_SLIDE",
            bkgd:  this.HUELGHT,
            margin: [1.25, 0.25, 0.25, 0.25],
            objects: [
                // BGIMAGE
                {
                    image: {
                        x: 0.0,
                        y: "0%",
                        w: "100%",
                        h: "100%",
                        path: this.assets.scbGradientB,
                        sizing: {type: "contain", w: "100%", h: "100%"}
                    }
                },
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
            ]
        });

        let thankYouSlide = this.pptx.addNewSlide("THANK_YOU_SLIDE");
        thankYouSlide.addText(`Thank You`, {
            x: "22.1%",
            y: "45%",
            w: "50%",
            h: "7%",
            fontWeight:"bold",
            bold: true,
            fontSize: `${this.HEADERTITLESIZE + 12}`,
            fontFace: this.DEFFONTFACE,
            color: this.HUELGHT,
            align: "center",
            valign: "middle"
            } as any);
        // THANK YOU SLIDE
        }

    private tieredSlide(client: string, index: number) {
        this.pptx.defineSlideMaster({
            title: `TIERED_SLIDE_${index}`,
            bkgd: this.HUELGHT,
            margin: [1.25, 0.25, 0.25, 0.25],
            objects: [
                {
                    text: {
                        //TIER TEXT
                        text: client,
                        options: {
                            x: "23.1%",
                            y: "47.5%",
                            w: "50%",
                            h: "7%",
                            fontFace: this.DEFFONTFACE,
                            color: this.BRANDCOLOR_BLUE,
                            bold: true,
                            italic: true,
                            fontSize: `${this.BASEFONTSIZE + 12}`,
                            align: "c",
                            margin: 0.4
                        }
                    }
                },
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
            ]
        });
        this.pptx.addNewSlide(`TIERED_SLIDE_${index}`)
    }

    private fdMaturitiesSlide(index: number) {
        var slide = this.pptx.addNewSlide("PAGE_SLIDES");
        slide.addText("Maturities For Fixed Deposits", {
            x: "4%",
            y: "6%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle"
        } as any);

        // TOP-L
        var optsChartBar2 = {
            x: 2.65,
            y: 1.3,
            w: "55%",
            h: "74%",
            barDir: "col",
            barGrouping: "stacked",
            dataLabelColor: this.HUELGHT,
            dataLabelFontFace: this.DEFFONTFACE,
            dataLabelFontSize: 12,
            dataLabelFontBold: true,
            showValue: false,
            catAxisLabelColor: this.TEXTCOLOR,
            catAxisLabelFontFace: this.DEFFONTFACE,
            catAxisLabelFontSize: 12,
            catAxisOrientation: "minMax",
            showLegend: false,
            chartColors: [this.BRANDCOLOR_BLUELGHT],
            shadow: {
                offset: 2,
                blur: 2,
                type: "outer"
            },
            showTitle: true,
            title: "FD Maturities",
            titleFontFace: this.DEFFONTFACE,
            titleFontSize: `${this.BASEFONTSIZE + 1}`,
            legendColor: this.BLUEBLACK,
            showDataTable: true,
            showDataTableKeys: false,
        };
        slide.addChart((this.pptx as any).charts.BAR, [this.fdMaturitiesData[index] as any], optsChartBar2 as any);
    }

    public generatePPT() {
        // COVER SLIDE
        this.pptx.defineSlideMaster({
            title: "COVER_SLIDE",
            bkgd: this.HUELGHT,
            margin: [1.25, 0.25, 0.25, 0.25],
            objects: [
                {
                    image: {
                        x: 0.0,
                        y: 0.0,
                        w: "100%",
                        h: "100%",
                        path: this.assets.scbGradientA,
                        sizing: {type: "contain", w: "100%", h: "100%"}
                    }
                },
                {
                    image: {
                        x: "1.45%",
                        y: "0.65%",
                        w: "86%",
                        h: "100%",
                        path: this.assets.scbTextLogo,
                        sizing: {type: "cover", w: "86%", h: "100%"}
                    }
                }
            ],

        });
        // END COVER SLIDE

        // MASTER PAGE TEMPLATE
        this.pptx.defineSlideMaster({
            title: "PAGE_SLIDES",
            bkgd: this.HUELGHT,
            margin: [0.25, 0.5, 0.25, 0.5],
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                },
                {
                    text: {
                        options: {
                            x: "0.5%",
                            y: "1%",
                            w: "30%",
                            h: "6%",
                            align: "l",
                            valign: "m",
                            color: this.BRANDCOLOR_RED,
                            fontSize: `${this.BASEFONTSIZE - 8}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "CONFIDENTIAL"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        // END MASTER PAGE TEMPLATE

        //AUCSUMMARY TABLE
        this.pptx.defineSlideMaster({
            title: "AUCSUMMARY_SLIDE",
            margin: [0.5, 0.5, 0.5, 0.5],
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                },
                {
                    text: {
                        options: {
                            x: "0.5%",
                            y: "1%",
                            w: "30%",
                            h: "6%",
                            align: "l",
                            valign: "m",
                            color: this.BRANDCOLOR_RED,
                            fontSize: `${this.BASEFONTSIZE - 8}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "CONFIDENTIAL"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        
        // END SLIDES WITH TABLE
        this.genSlides_Content();

        this.pptx.save(this.MISC[0].report_name + `Q${this.MISC[0].current_quarter_number}` + " Trustee Report_" + this.getTimestamp());
    }

    private aucSummaryTblSlide(index: number) {
        var slide = (this.pptx as any).addSlidesForTable(`aucSummaryTbl_${index}`, {
            master: "AUCSUMMARY_SLIDE",
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "90%",
            addText: {
                text: `AUC Summary (%) – ${this.MISC[0].current_quarter_month_year}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    h: "17%",
                    valign: "m",
                    fontWeight:"bold",
                    bold: true,
                    color: this.BLUEBLACK,
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.4
                },
            }
        });
    }

    private aucSummarySlide2(index: number) {
        this.pptx.defineSlideMaster({
            title: `CHARTPIE_SLIDE_${index}`,
            margin: [0.9, 0.9, 0.9, 0.9],
            bkgd: this.SCBBLUELIGHTER,
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                },
                {
                    text: {
                        options: {
                            x: "0.5%",
                            y: "1%",
                            w: "30%",
                            h: "6%",
                            align: "l",
                            valign: "m",
                            color: this.BRANDCOLOR_RED,
                            fontSize: `${this.BASEFONTSIZE - 8}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "CONFIDENTIAL"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        var slide = this.pptx.addNewSlide(`CHARTPIE_SLIDE_${index}`);
        slide.addText(`AUC Summary (%) – ${this.MISC[index].current_quarter_month_year}`, {
            x: "4%",
            y: "6%",
            w: "100%",
            h: "9%",
            margin: 0.15,
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle"
        } as any);
        slide.addText("AUC BY ASSETS CLASS", {
            x: "9%",
            y: "16%",
            fontSize: `${this.BASEFONTSIZE - 2}`,
            fontFace: this.DEFFONTFACE,
            color: this.BRANDCOLOR_BLUE,
            bold: true,
            align: "c",
            valign: "middle"
        } as any);
        slide.addChart((this.pptx as any).charts.PIE, [this.dataChartPieLocs[index] as any], {
            title: "AUC BY ASSETS CLASS",
            legendFontSize: 13,
            legendFontFace: this.DEFFONTFACE,
            legendColor: this.BRANDCOLOR_BLUE,
            titleColor: this.BRANDCOLOR_BLUE,
            titleFontFace: this.DEFFONTFACE,
            titleFontSize: 15,
            titleAlign: 'center',
            titlePos: {x: 0.5, y: 0},
            x: 1.80,
            y: 1.35,
            w: "80%",
            h: "75%",
            dataBorder: {pt: "1", color: "F8F8F8"},
            showLegend: true,
            legendPos: "r",
            showTitle: false,
            dataLabelColor: this.BRANDCOLOR_BLUE,
            showLabel: false,
            showValue: false,
            showPercent: true
        })
    }

    private aucTrendSlide(index: number) {
        var slide = this.pptx.addNewSlide("PAGE_SLIDES");
        slide.addText(`AUC TREND as at ${this.MISC[0].quarter}`, {
            x: "4%",
            y: "6%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle"
        } as any);
        let isBillion = false;
        const billVal = this.aucTrendChart[index].values.find((each) => each >= 1000000000);
        if (billVal) {
            isBillion = true;
        }
        var aucTrendChartOpts = {
            x: 0.5,
            y: 1.3,
            w: "80%",
            h: "74%",
            valAxisTitle: `AUC (GHS ${isBillion ? 'Bn' : 'M'})`,
            valAxisTitleFontSize: 10,
            valAxisLabelFontFace: this.DEFFONTFACE,
            showValAxisTitle: true,
            valAxisTitleColor: this.BLUEBLACK,
            showDataTable: true,
            showDataTableKeys: false,
            chartColors: [this.BRANDCOLOR_BLUELGHT],
            shadow: {
                offset: 2,
                blur: 2,
                type: "outer"
            }
        };
        let year = "";
        const date = this.MISC[0].quarter.split(" ");
        if (date[0] == "Q1") {
            year = `${parseInt(date[1]) - 1} / ${date[1]}`
        } else {
            year = date[1];
        }

        slide.addChart((this.pptx as any).charts.LINE, [{
            name: this.aucTrendChart[index].name + ` for ${year}`,
            labels: this.aucTrendChart[index].labels,
            values: this.aucTrendChart[index].values.map((each) => (isBillion ? (each / 1000000000) : (each / 1000000)).toFixed(2))
        } as any], aucTrendChartOpts);
    }

    private transactionVolSlide(index: number) {
        var slide = this.pptx.addNewSlide("PAGE_SLIDES");
        slide.addText(this.MISC[0].txn_vols_heading, {
            x: "11%",
            y: "2%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle"
        } as any);

        // TOP-L
        var optsChartBar2 = {
            x: 0.5,
            y: 0.88,
            w: 6.0,
            h: "55%",
            barDir: "col",
            barGrouping: "stacked",
            dataLabelColor: this.HUELGHT,
            dataLabelFontFace: this.DEFFONTFACE,
            dataLabelFontSize: 12,
            dataLabelFontBold: true,
            showValue: false,
            catAxisLabelColor: this.TEXTCOLOR,
            catAxisLabelFontFace: this.DEFFONTFACE,
            catAxisLabelFontSize: 12,
            catAxisOrientation: "minMax",
            showLegend: false,
            chartColors: [this.BRANDCOLOR_BLUELGHT],
            shadow: {
                offset: 2,
                blur: 2,
                type: "outer"
            },
            showTitle: true,
            title: "Monthly Trade Volume",
            titleFontFace: this.DEFFONTFACE,
            titleFontSize: `${this.BASEFONTSIZE + 1}`,
            legendColor: this.BLUEBLACK,
            showDataTable: true,
            showDataTableKeys: false,
        };
        slide.addChart((this.pptx as any).charts.BAR, [this.arrDataRegions[index] as any], optsChartBar2 as any);

        // TOP-R
        var trnsVolAssetChartOpts = {
            x: 7,
            y: 0.88,
            w: 6.0,
            h: "55%",
            valAxisTitle: "No. of trades",
            valAxisTitleFontSize: 10,
            valAxisLabelFontFace: this.DEFFONTFACE,
            showValAxisTitle: true,
            valAxisTitleColor: this.BLUEBLACK,
            showDataTable: true,
            showDataTableKeys: false,
            chartColors: this.MISC[index].tradeByAssetClassData.map((each: any) => each.color.replace(/#/, '')),
            shadow: {
                offset: 2,
                blur: 2,
                type: "outer"
            },
            showTitle: true,
            title: "Trade Volume by Assets Class",
            titleFontFace: this.DEFFONTFACE,
            titleFontSize: `${this.BASEFONTSIZE + 1}`,
            showLegend: true,
            legendPos: "b",
            legendColor: this.BLUEBLACK,
            lineDataSymbol: "none"
        };
        slide.addChart(
            (this.pptx as any).charts.LINE,
            this.trnsVolAssetChart[index] as any,
            trnsVolAssetChartOpts
        );

        slide.addText(`Total volumes of trades as at Q${this.MISC[index].current_quarter_number} - ${this.MISC[index].total_number_of_txn_vols}`, {
            x: "10%",
            y: "62%",
            fontSize: `${this.BASEFONTSIZE}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            align: "l",
            valign: "middle",
            width: "100%",
            h: "20%",
        } as any);
        let count = 0;
        for (let summary in this.MISC[index].tradeByAssetClassSummary) {
            slide.addText(`${summary} – ${this.MISC[index].tradeByAssetClassSummary[summary]}`, {
                x: "5%",
                y: 5.55 + (count++ * 0.25),
                fontSize: `${this.BASEFONTSIZE}`,
                fontFace: this.DEFFONTFACE,
                color: this.HEADERCOLOR,
                align: "l",
                valign: "middle",
                h: "6%"
            } as any);
        }
    }

    private createMasterSlide(name: string): string {
        const slideName = `SLIDE_${name.replace(/\s/, '')}`;
        this.pptx.defineSlideMaster({
            title: slideName,
            margin: [0.5, 2.2, 0.5, 2.2],
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                },
                {
                    text: {
                        options: {
                            x: "0.5%",
                            y: "1%",
                            w: "30%",
                            h: "6%",
                            align: "l",
                            valign: "m",
                            color: this.BRANDCOLOR_RED,
                            fontSize: `${this.BASEFONTSIZE - 8}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "CONFIDENTIAL"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        return slideName;
    }

    private monthlyContributionTableSlide(index: number) {
        var slide = (this.pptx as any).addSlidesForTable(`quarterMonthContribution_${index}`, {
            // QUARTER SEARCHED
            master: this.createMasterSlide(`Contributions as at Q${this.MISC[index].current_quarter_number}`),
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "90%",
            addText: {
                text: `Contributions as at Q${this.MISC[index].current_quarter_number}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    bold:true,
                    h: "17%",
                    valign: "m",
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    color: this.HEADERCOLOR,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.4
                },
            }
        });
    }

    // QUARTER SEARCHED
    private corporateActionActivitiesSlide(index: number) {
        var slide = (this.pptx as any).addSlidesForTable(`corporateActionsTbl_${index}`, {
            master: this.createMasterSlide(`Corporate Actions for Q${this.MISC[index].current_quarter_number}`),
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "90%",
            addText: {
                text: `Corporate Actions for Q${this.MISC[index].current_quarter_number}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    bold:true,
                    h: "17%",
                    valign: "m",
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    color: this.HEADERCOLOR,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.4
                },
            }
        });
    }

    private maturitiesForGOGSlide(index: number) {
        var slide = (this.pptx as any).addSlidesForTable(`quarterGOGSecuritiesTbl_${index}`, {
            master: this.createMasterSlide(`GOG Securities for Q${this.MISC[index].current_quarter_number}`),
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "90%",
            addText: {
                text: "Maturities for GOG Securities",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    fontSize: `${this.BASEFONTSIZE}`,
                    fontFace: this.DEFFONTFACE,
                    color: '000000',
                    bold: true,
                    align: "l",
                    valign: "middle"
                }
            },
            //@ts-ignore
            addText: {
                text: `GOG Securities for Q${this.MISC[index].current_quarter_number}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    bold: true,
                    h: "17%",
                    valign: "m",
                    fontSize:   `${this.HEADERTITLESIZE - 6}`,
                    color: this.HEADERCOLOR,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.4
                },
            }
        });
        this.fdMaturitiesSlide(index);
    }

    private unidentifiedPaymentsSlide(index: number) {
        var slide = (this.pptx as any).addSlidesForTable(`unidentifiedPaymentsTbl_${index}`, {
            master: this.createMasterSlide(`Unidentified payments for Q${this.MISC[index].current_quarter_number}`),
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "90%",
            addText: {
                text: `Unidentified payments for Q${this.MISC[index].current_quarter_number}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    h: "17%",
                    valign: "m",
                    fontWeight:"bold",
                    bold: true,
                    color: this.HEADERCOLOR,
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.4
                },
            }
        });
    }

    private pvSummaryForFundManagersSlide(manager: string, index: number) {
        // APPENDIX SLIDE 
        this.pptx.defineSlideMaster({
            title: `APPENDIX_SLIDE ${manager}`,
            margin: [0.85, 0.85, 0.85, 0.85],
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        var slide = (this.pptx as any).addSlidesForTable(`summaryFor_${manager}_${index}`, {
            master: `APPENDIX_SLIDE ${manager}`,
            x: "10%",
            y: "18%",
            addText: {
                text: `Appendix I - Summary for ${manager.replace(/\d{1,}_/, '')}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    h: "9%",
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    fontFace: this.DEFFONTFACE,
                    color: this.HEADERCOLOR,
                    fontWeight:"bold",
                    bold: true,
                    align: "l",
                    valign: "middle",
                    margin: 0.15,
                },
            }
        });
    }

    private contributionsForFundManagersSlide(manager: string, index: number) {
        // CONTRIBUTIONS TABLE 
        this.pptx.defineSlideMaster({
            title: `CONTRIBUTION_SLIDE${manager}`,
            margin: [2.5, 2.5, 2.5, 2.5],
            objects: [
                // Footer Content
                {
                    image: {
                        x: "15%",
                        y: "20%",
                        w: "85%",
                        h: "85%",
                        path: this.assets.stanChartBgLogo,
                        sizing: {type: "contain", w: "85%", h: "85%"}
                    }
                },
                {
                    text: {
                        options: {
                            x: "70%",
                            y: "92%",
                            w: "30%",
                            h: "8%",
                            align: "c",
                            valign: "m",
                            color: this.TEXTCOLOR,
                            fontSize: `${this.BASEFONTSIZE - 7}`,
                            fontFace: this.DEFFONTFACE
                        },
                        text: "Trustee Report"
                    }
                }
            ],
            slideNumber: {
                x: "2.4%",
                y: "92%",
                color: this.TEXTCOLOR,
                fontSize: `${this.BASEFONTSIZE - 4}`,
                fontFace: this.DEFFONTFACE
            } as any
        });
        var slide = (this.pptx as any).addSlidesForTable(`contributionsFor_${manager}_${index}`, {
            master: `CONTRIBUTION_SLIDE${manager}`,
            x: "10%",
            y: "18%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle",
            width: "80%",
            addText: {
                text: `Appendix II - Contributions for ${manager.replace(/\d{1,}_/, '')}`,
                placeholder: "body",
                options: {
                    x: "4%",
                    y: "0%",
                    w: "100%",
                    align: "l",
                    h: "9%",
                    valign: "m",
                    bold: true,
                    fontSize: `${this.HEADERTITLESIZE - 6}`,
                    color: this.HEADERCOLOR,
                    fontFace: this.DEFFONTFACE,
                    margin: 0.15,
                },
            }
        });
    }

    private pvReportSlide() {
        var slide = this.pptx.addNewSlide("PAGE_SLIDES");
        slide.addText("Portfolio Valuation Report", {
            x: "4%",
            y: "6%",
            fontSize: `${this.HEADERTITLESIZE - 6}`,
            fontFace: this.DEFFONTFACE,
            color: this.HEADERCOLOR,
            fontWeight:"bold",
            bold: true,
            align: "l",
            valign: "middle"
        } as any);
    }
}