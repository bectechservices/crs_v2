// RGBTOHEX
function rgbToHex(r, g, b) {
  return ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}
//TIMESTAMP
function getTimestamp() {
  let dateNow = new Date();
  let dateMM = dateNow.getMonth() + 1;
  dateDD = dateNow.getDate();
  (dateYY = dateNow.getFullYear()), (h = dateNow.getHours());
  m = dateNow.getMinutes();
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
function toUpper() {
  var TblContent = document.querySelectorAll(".aucSummaryTbl tr td");
  for (var i = 0; i < TblContent.length; ++i) {
    TblContent[i].setAttribute(
      "style",
      "text-transform:uppercase;padding:9px;"
    );
  }
}

// DEFINE VARIABLES
let COMPANY_NAME = "BECTECH";
let CUST_NAME = "Standard Chartered";
let CLIENT_NAME = "Twifo Oil Palm Plantation";

// CLIENT PARAMETERS
let summaryDateRange = "March 31, 2019";
let pvReportpath = "https://www.google.com";
let PVMONTH = "Month";
// CLIENT PARAMETERS

// PAGE SETTING
let slideTabConfig = { x: 0.4, y: 0.13, w: 12.5, colW: [9, 3.5] };
let slideTabConfigL = {
  color: "9F9F9F",
  margin: 3,
  border: [0, 0, { pt: "1", color: "CFCFCF" }, 0]
};
let slideTabConfigR = {
  color: "9F9F9F",
  margin: 3,
  border: [0, 0, { pt: "1", color: "CFCFCF" }, 0],
  align: "right"
};
let slideSubTitle = {
  x: 0.5,
  y: 0.7,
  w: 4,
  h: 0.3,
  fontSize: 18,
  fontFace: "Arial",
  color: "0088CC",
  fill: "FFFFFF"
};

// STYLING VARIABLES
let BRANDCOLOR_BLUE = `${rgbToHex(31, 73, 125)}`;
let BRANDCOLOR_BLUELGHT = `${rgbToHex(74, 126, 187)}`;
let BRANDCOLOR_GREEN = `${rgbToHex(119, 147, 60)}`;
let TEXTCOLOR = "999999";
let BLUEBLACK = "0B2135";
let BGCOLORLGHT = "ffffff";
let CHARSPERLINE = 130;
let DEFFONTFACE = "Calibri";
let BASEFONTSIZE = 16;
let HEADERTITLESIZE = `${16 + 18}`;
// STYLING VARIABLES

// ASSETS PATHS
let assetPaths = {
  stanChartBgLogo: { path: "./media/logo/standardchartered@2x.png" },
  stanChartBgImage: { path: "./media/Picture1.jpg" },
  slideBgPatternImageA: { path: "./media/Picture2.jpg" },
  slideBgPatternImageB: { path: "./media/Picture3.jpg" },
  slideBrandLineImage: { path: "./media/Picture4.png" },
  clientLogo: { path: "./media/logo/npra.png" },
  pvReportIcon: { path: "./media/excel.svg" }
};

let assets = {
  stanChartBgLogo: "./media/logo/sc.jpeg",
  sclogoStandAlone: "./media/logo/sclogoAlone.jpg",
  stanChartBgImage: "./media/Picture1.jpg",
  slideBgPatternImageA: "./media/Picture2.jpg",
  slideBgPatternImageB: "./media/Picture3.jpg",
  scBrandLineImage: "./media/Picture4.png",
  clientLogo: "./media/logo/twifo.png",
  pvReportIcon: "./media/logo/twifo.png"
};

// TEST COLOR SUITE
let colors = {
  yellow: "7FFF00",
  blue: BLUEBLACK,
  darkRed: "8B0000",
  greenYellow: "ADFF2F"
};
// TEST COLOR SUITE END

//BUILDING SLIDES
function execGenSlidesFuncs(type) {
  let pptx = new PptxGenJS();

  // PRESENTATIONAL & LAYOUT PROPS
  pptx.setAuthor(CUST_NAME);
  pptx.setCompany(COMPANY_NAME);
  pptx.setRevision("15");
  pptx.setSubject("Trustee Report");
  pptx.setTitle(`${CLIENT_NAME} Trustee Report Powerpoint`);

  pptx.setLayout("LAYOUT_WIDE");
  {
    // COVER SLIDE
    pptx.defineSlideMaster({
      title: "COVER_SLIDE",
      bkgd: BGCOLORLGHT,
      margin: [1.25, 0.25, 0.25, 0.25],
      objects: [
        // BGIMAGE
        {
          image: {
            x: 0.0,
            y: "13%",
            w: "100%",
            h: "87%",
            path: assets.stanChartBgImage,
            sizing: { type: "contain", w: "100%", h: "87%" }
          }
        },
        // CLIENTLOGO
        {
          image: {
            x: "2%",
            y: 0.0,
            w: "9%",
            h: "13%",
            path: assets.clientLogo,
            sizing: { type: "contain", w: "9%", h: "13%" }
          }
        },
        // COMPANY LOGO
        {
          image: {
            x: "83%",
            y: 0.0,
            w: "15%",
            h: "13%",
            path: assets.stanChartBgLogo,
            sizing: { type: "contain", w: "15%", h: "13%" }
          }
        },
        {
          text: {
            text: `${CLIENT_NAME} Trustee Meeting`,
            options: {
              x: "50%",
              y: "73%",
              w: "50%",
              h: "7%",
              fontFace: DEFFONTFACE,
              color: `${rgbToHex(109, 110, 113)}`,
              fontWeight: "bold",
              fontSize: `${BASEFONTSIZE + 9}`,
              align: "r",
              valign: "m",
              margin: 0.4
            }
          }
        },
        {
          text: {
            text: `${"Q1"} ${"2019"} Custodian Report`,
            options: {
              x: "50%",
              y: "79%",
              w: "50%",
              h: "6%",
              fontFace: DEFFONTFACE,
              color: `${rgbToHex(109, 110, 113)}`,
              fontSize: `${BASEFONTSIZE + 2}`,
              align: "r",
              valign: "m",
              margin: 0.4
            }
          }
        },
        {
          text: {
            text: `${"May 2019"}`,
            options: {
              x: "50%",
              y: "84%",
              w: "50%",
              h: "4%",
              fontFace: DEFFONTFACE,
              color: `${rgbToHex(109, 110, 113)}`,
              fontSize: `${BASEFONTSIZE + 2}`,
              align: "r",
              valign: "m",
              margin: 0.4
            }
          }
        },
        {
          text: {
            text: "sc.com | Here for good",
            options: {
              x: "50%",
              y: "92.8%",
              w: "50%",
              h: "4%",
              fontFace: DEFFONTFACE,
              color: `${rgbToHex(155, 187, 89)}`,
              fontSize: `${BASEFONTSIZE + 1}`,
              bold: true,
              align: "r",
              valign: "m",
              margin: 0.4
            }
          }
        },
        {
          image: {
            x: 0.0,
            y: "13%",
            w: "100%",
            h: 0.05,
            path: assets.scBrandLineImage
          }
        }
      ]
    });
    // COVER SLIDE
    // MASTER PAGE TEMPLATE
    pptx.defineSlideMaster({
      title: "PAGE_SLIDES",
      bkgd: BGCOLORLGHT,
      margin: [0.25, 0.5, 0.25, 0.5],
      objects: [
        // SCB BRAND LINE
        {
          image: {
            x: 0.0,
            y: "14%",
            w: "100%",
            h: 0.06,
            path: assets.scBrandLineImage
          }
        },
        // Footer Content
        {
          image: {
            x: "2.4%",
            y: "92%",
            w: "12%",
            h: "8%",
            path: assets.stanChartBgLogo,
            sizing: { type: "contain", w: "11%", h: "8%" }
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
              color: TEXTCOLOR,
              fontSize: `${BASEFONTSIZE - 4}`,
              fontFace: DEFFONTFACE
            },
            text: "Trustee Report"
          }
        }
      ],
      slideNumber: {
        x: "93%",
        y: "94%",
        color: TEXTCOLOR,
        fontSize: `${BASEFONTSIZE - 4}`,
        fontFace: DEFFONTFACE
      }
    });
    // MASTER PAGE TEMPLATE

    // SLIDES WITH TABLE
    //AUCSUMMARY TABLE
    pptx.defineSlideMaster({
      title: "AUCSUMMARY_SLIDE",
      margin: [0.5, 0.5, 0.5, 0.5],
      objects: [
        {
          text: {
            text: `AUC Summary (%) – ${summaryDateRange}`,
            options: {
              x: "0%",
              y: "0%",
              w: "100%",
              align: "c",
              h: "17%",
              valign: "m",
              fontSize: HEADERTITLESIZE,
              color: BLUEBLACK,
              fontFace: DEFFONTFACE,
              margin: 0.4
            }
          }
        },
        // SCB BRAND LINE
        {
          image: {
            x: 0.0,
            y: "14%",
            w: "100%",
            h: 0.06,
            path: assets.scBrandLineImage
          }
        },
        // Footer Content
        {
          image: {
            x: "2.4%",
            y: "92%",
            w: "12%",
            h: "8%",
            path: assets.stanChartBgLogo,
            sizing: { type: "contain", w: "11%", h: "8%" }
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
              color: TEXTCOLOR,
              fontSize: `${BASEFONTSIZE - 4}`,
              fontFace: DEFFONTFACE
            },
            text: "Trustee Report"
          }
        }
      ],
      slideNumber: {
        x: "93%",
        y: "94.5%",
        color: TEXTCOLOR,
        fontSize: `${BASEFONTSIZE - 4}`,
        fontFace: DEFFONTFACE
      }
    });
    // MONTHLY CONTRIBUTIONS TABLE
    pptx.defineSlideMaster({
      title: "MONTHLY_SLIDE",
      margin: [0.5, 2.2, 0.5, 2.2],
      objects: [
        {
          text: {
            text: "Monthly Contributions",
            options: {
              x: "0%",
              y: "0%",
              w: "100%",
              align: "c",
              h: "17%",
              valign: "m",
              fontSize: HEADERTITLESIZE,
              color: BLUEBLACK,
              fontFace: DEFFONTFACE,
              margin: 0.4
            }
          }
        },
        // SCB BRAND LINE
        {
          image: {
            x: 0.0,
            y: "14%",
            w: "100%",
            h: 0.06,
            path: assets.scBrandLineImage
          }
        },
        // Footer Content
        {
          image: {
            x: "2.4%",
            y: "92%",
            w: "12%",
            h: "8%",
            path: assets.stanChartBgLogo,
            sizing: { type: "contain", w: "11%", h: "8%" }
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
              color: TEXTCOLOR,
              fontSize: `${BASEFONTSIZE - 4}`,
              fontFace: DEFFONTFACE
            },
            text: "Trustee Report"
          }
        }
      ],
      slideNumber: {
        x: "93%",
        y: "94.5%",
        color: TEXTCOLOR,
        fontSize: `${BASEFONTSIZE - 4}`,
        fontFace: DEFFONTFACE
      }
    });
    // THANK YOU SLIDE
    pptx.defineSlideMaster({
      title: "THANK_YOU_SLIDE",
      bkgd: BGCOLORLGHT,
      margin: [1.25, 0.25, 0.25, 0.25],
      objects: [
        // BGIMAGE
        {
          image: {
            x: 0.0,
            y: "0%",
            w: "100%",
            h: "87%",
            path: assets.slideBgPatternImageB,
            sizing: { type: "contain", w: "100%", h: "87%" }
          }
        },
        {
          text: {
            text: "THANK YOU",
            options: {
              x: "22.1%",
              y: "57%",
              w: "50%",
              h: "7%",
              fontFace: DEFFONTFACE,
              color: BLUEBLACK,
              fontWeight: "bold",
              fontSize: `${BASEFONTSIZE + 15}`,
              align: "c",
              valign: "m",
              margin: 0.4
            }
          }
        },
        {
          image: {
            x: "40%",
            y: "32%",
            w: "12%",
            h: "22%",
            path: assets.sclogoStandAlone,
            sizing: { type: "contain", w: "14%", h: "21%" }
          }
        }
      ]
    });
    // THANK YOU SLIDE
  }
  

  let arrTypes = typeof type === "string" ? [type] : type;
  arrTypes.forEach(function(type, idx) {
    eval("genSlides_" + type + "(pptx)");
  });
  pptx.save(CLIENT_NAME +'Trustee Report Powerpoint' + type + "_" + getTimestamp());
}

function genSlides_Content(pptx) {
  // CHARTS DATA
  let yearlyData = [
    "Jan",
    "Feb",
    "Mar",
    "April",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
    "Jan",
    "Feb",
    "Mar"
  ];
  // PIE CHART
  let dataChartPieLocs = [
    {
      name: "Location",
      labels: [
        "RECEIVABLES",
        "EQUITIES",
        "TREASURY BILLS",
        "CASH BALANCE",
        "CORPORATE BOND",
        "FIXED DEPOSITS",
        "GOVT DEBT-BOND/LOANS"
      ],
      values: [16, 11, 7, 3, 10, 20, 33]
    }
  ];
  // TABLE CHART
  var aucTrendChart = [
    {
      name: "2018/2019",
      values: [2.79, 2.73, 3.44, 4.24, 4.64, 4.67],
      labels: ["Oct", "Nov", "Dec", "Jan", "Feb", "Mar"]
    }
  ];
  let trnsVolAssetChart = [
    {
      name: "Corporate Bond",
      values: [0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0],
      labels: yearlyData
    },
    {
      name: "T'Bill",
      values: [0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0],
      labels: yearlyData
    },
    {
      name: "Fixed Deposit",
      values: [0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
      labels: yearlyData
    },
    {
      name: "Govt. Debt-Fixed",
      values: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0],
      labels: yearlyData
    }
  ];
  let arrDataRegions = [
    {
      name: "Monthly Trade Volume",
      labels: yearlyData,
      values: [8, 4, 0, 1, 4, 2, 2, 0, 1, 1, 4, 1, 0, 2]
    }
  ];
  // CHARTS DATA END

  // COVER SLIDE
  {
    let slide = pptx.addNewSlide("COVER_SLIDE");
  }

  // AGENDA SLIDE
  {
    let agendaSlide = pptx.addNewSlide("PAGE_SLIDES");
    agendaSlide.addText("Agenda", {
      x: "11%",
      y: "7%",
      fontSize: HEADERTITLESIZE,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      align: "c",
      valign: "m"
    });
    // TOC
    // EACH NUMBER AND AGENDA HEADING
    agendaSlide.addText("01", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "20%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("Summary of Valuation Report", {
      x: "8%",
      y: "20%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_BLUE,
      valign: "m"
    });

    agendaSlide.addText("02", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "30.8%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("Summary of Assets under Custody", {
      x: "8%",
      y: "30.8%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_BLUE,
      valign: "m"
    });

    agendaSlide.addText("03", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "41.5%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_GREEN, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("AUC Trend", {
      x: "8%",
      y: "41.5%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_GREEN,
      valign: "m"
    });
    agendaSlide.addText("04", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "51.6%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_GREEN, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("Trade Volumes", {
      x: "8%",
      y: "51.6%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_GREEN,
      valign: "m"
    });
    agendaSlide.addText("05", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "61.8%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("Contributions", {
      x: "8%",
      y: "61.6%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_BLUE,
      valign: "m"
    });
    agendaSlide.addText("06", {
      shape: pptx.shapes.OVAL,
      x: "2%",
      y: "71.9%",
      w: "5%",
      h: "9.2%",
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: `${BASEFONTSIZE + 3}`,
      bold: true,
      fontFace: DEFFONTFACE,
      color: BGCOLORLGHT
    });
    agendaSlide.addText("Portfolio Valuation Report", {
      x: "8%",
      y: "71.9%",
      h: "10%",
      fontFace: DEFFONTFACE,
      fontSize: `${BASEFONTSIZE + 4}`,
      color: BRANDCOLOR_BLUE,
      valign: "m"
    });
  }

  // AUC SUMMARY SLIDES
  {
    pptx.addSlidesForTable("aucSummaryTbl", {
      master: "AUCSUMMARY_SLIDE",
      x: "4%",
      y: "20%"
    });
  }
  {
    let aucSummarySlide2 = pptx.addNewSlide("AUCSUMMARY_SLIDE");
    // aucSummarySlide2.addText(".", {
    //   x: 9.8,
    //   y: 4.0,
    //   w: 3.2,
    //   h: 3.2,
    //   fill: "F1F1F1",
    //   color: "F1F1F1"
    // });
    aucSummarySlide2.addChart(pptx.charts.PIE, dataChartPieLocs, {
      title: "AUC BY Asset Class",
      x: 0,
      y: 1.2,
      w: "100%",
      h: "74%",
      dataBorder: { pt: "1", color: "F1F1F1" },
      showLegend: true,
      legendPos: "b",
      showTitle: true,
      dataLabelColor: BRANDCOLOR_BLUE,
      showLabel: true,
      showValue: false,
      showPercent: true,
      shadow: {
        offset: 4,
        blur: 18,
        type: "outer"
      }
    });
  }

  // AUC TREND SLIDE
  {
    let aucTrendSlide = pptx.addNewSlide("PAGE_SLIDES");
    aucTrendSlide.addText(`AUC TREND as at ${"Q1 2019"}`, {
      x: "11%",
      y: "7%",
      fontSize: HEADERTITLESIZE,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      align: "c",
      valign: "m"
    });
    var aucTrendChartOpts = {
      x: 0.5,
      y: 1.3,
      w: "80%",
      h: "74%",
      valAxisTitle: "AUC (GHS mn)",
      valAxisTitleFontSize: 10,
      valAxisLabelFontFace: DEFFONTFACE,
      showValAxisTitle: true,
      valAxisTitleColor: BLUEBLACK,
      showDataTable: true,
      showDataTableKeys: false,
      chartColors: [BRANDCOLOR_BLUELGHT],
      shadow: {
        offset: 2,
        blur: 2,
        type: "outer"
      }
    };
    aucTrendSlide.addChart(pptx.charts.LINE, aucTrendChart, aucTrendChartOpts);
  }

  // Transaction Volume Slide
  {
    let trnsVolSlide = pptx.addNewSlide("PAGE_SLIDES");
    trnsVolSlide.addText(`Transaction Volumes for 2018 and ${"Q1 2019"}`, {
      x: "11%",
      y: "7%",
      fontSize: HEADERTITLESIZE,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      align: "c",
      valign: "m"
    });

    // TOP-L
    var optsChartBar2 = {
      x: 0.5,
      y: 1.3,
      w: 6.0,
      h: "50%",
      barDir: "col",
      barGrouping: "stacked",
      dataLabelColor: BGCOLORLGHT,
      dataLabelFontFace: DEFFONTFACE,
      dataLabelFontSize: 12,
      dataLabelFontBold: true,
      showValue: false,
      catAxisLabelColor: TEXTCOLOR,
      catAxisLabelFontFace: DEFFONTFACE,
      catAxisLabelFontSize: 12,
      catAxisOrientation: "minMax",
      showLegend: false,
      showTitle: false,
      chartColors: [BRANDCOLOR_BLUELGHT],
      shadow: {
        offset: 2,
        blur: 2,
        type: "outer"
      },
      showTitle: true,
      title: "Monthly Trade Volume",
      // showLegend: true,
      // legendPos: "b",
      legendColor: BLUEBLACK
    };
    trnsVolSlide.addChart(pptx.charts.BAR, arrDataRegions, optsChartBar2);

    // TOP-R
    var trnsVolAssetChartOpts = {
      x: 7,
      y: 1.3,
      w: 6.0,
      h: "50%",
      valAxisTitle: "No. of trades",
      valAxisTitleFontSize: 10,
      valAxisLabelFontFace: DEFFONTFACE,
      showValAxisTitle: true,
      valAxisTitleColor: BLUEBLACK,
      showDataTable: true,
      showDataTableKeys: false,
      chartColors: [BRANDCOLOR_BLUELGHT, BRANDCOLOR_GREEN, BRANDCOLOR_BLUE],
      shadow: {
        offset: 2,
        blur: 2,
        type: "outer"
      },
      showTitle: true,
      title: "Trade Volume by Asset Class",
      showLegend: true,
      legendPos: "b",
      legendColor: BLUEBLACK,
      lineDataSymbol: "none"
    };
    trnsVolSlide.addChart(
      pptx.charts.LINE,
      trnsVolAssetChart,
      trnsVolAssetChartOpts
    );

    trnsVolSlide.addText(` Total volumes of trades as at Q1– 11`, {
      x: "10%",
      y: "63%",
      fontSize: BASEFONTSIZE,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      align: "c",
      valign: "m",
      width: "100%",
      h: "20%",
      bold: true
    });
    trnsVolSlide.addText(`Govt Bonds/Notes – ${"27"}%`, {
      x: "5%",
      y: 5.55,
      fontSize: `${BASEFONTSIZE - 2}`,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      bold: true,
      h: "6%"
    });
    trnsVolSlide.addText(`Corporate Bonds – ${"36.36"}%`, {
      x: "5%",
      y: 5.75,
      fontSize: `${BASEFONTSIZE - 2}`,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      bold: true,
      h: "6%"
    });
    trnsVolSlide.addText(`Treasury Bills – ${"18.18"}%`, {
      x: "5%",
      y: 5.95,
      fontSize: `${BASEFONTSIZE - 2}`,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      bold: true,
      h: "6%"
    });
    trnsVolSlide.addText(`Fixed Deposit  – ${"18.18"}%`, {
      x: "5%",
      y: 6.15,
      fontSize: `${BASEFONTSIZE - 2}`,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      bold: true,
      h: "6%"
    });
  }

  // Monthly Contributions Slide
  {
    pptx.addSlidesForTable("monthlyContributionTbl", {
      master: "MONTHLY_SLIDE",
      x: "15%",
      y: "20%"
    });
  }

  // PV Report Slide
  {
    let pvReportSlide = pptx.addNewSlide("PAGE_SLIDES");
    pvReportSlide.addText("Portfolio Valuation Report", {
      x: "11%",
      y: "7%",
      fontSize: HEADERTITLESIZE,
      fontFace: DEFFONTFACE,
      color: BLUEBLACK,
      align: "c",
      valign: "m"
    });
    pvReportSlide.addImage({
      x: "42%",
      y: "42%",
      w: "14%",
      h: "14%",
      hyperlink: {
        url: pvReportpath,
        tooltip: `${PVMONTH} Report`
      },
      path: assets.pvReportIcon,
      sizing: { type: "contain", w: "14%", h: "14%" }
    });
  }

  {
    let thankYouSlide = pptx.addNewSlide("THANK_YOU_SLIDE");
  }
}
