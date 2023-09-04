// RGBTOHEX
function rgbToHex(r, g, b) {
  return ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}
//
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

// DEFINE VARIABLES
let COMPANY_NAME = "BECTECH";
let CUST_NAME = "Standard Chartered";
let CLIENT_NAME = "Twifo Oil Palm Plantation";
let BRANDCOLOR_BLUE = `${rgbToHex(31, 73, 125)}`;
let BRANDCOLOR_GREEN = `${rgbToHex(119, 147, 60)}`;
let TEXTCOLOR = "999999";
let BLUEBLACK = '0B2135';
let BGCOLORLGHT = "ffffff";
let CHARSPERLINE = 130;

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
let slideTitleText = { fontSize: 14, color: "0088CC", bold: true };
let slideTitleOpts = { fontSize: 13, color: "9F9F9F" };

// ASSETS PATHS
let assetPaths = {
  stanChartBgLogo: { path: "./media/logo/standardchartered@2x.png" },
  stanChartBgImage: { path: "./media/Picture1.jpg" },
  slideBgPatternImageA: { path: "./media/Picture2.jpg" },
  slideBgPatternImageB: { path: "./media/Picture3.jpg" },
  slideBrandLineImage: { path: "./media/Picture4.png" },
  clientLogo: { path: "./media/logo/npra.png" }
};

let assets = {
  stanChartBgLogo: "./media/logo/sc.jpeg",
  stanChartBgImage: "./media/Picture1.jpg",
  slideBgPatternImageA: "./media/Picture2.jpg",
  slideBgPatternImageB: "./media/Picture3.jpg",
  scBrandLineImage: "./media/Picture4.png",
  clientLogo: "./media/logo/twifo.png"
};

// TEST COLOR SUITE
let colors = {
  yellow: "7FFF00",
  blue: "0b2135",
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
        // { rect: { x: 0.0, y: 0, w: "100%", h: "13%", fill: colors.darkRed } },
        // { rect: { x: 0.0, y: "13%", w: "100%", h: "87%", fill: colors.blue } },
        // { rect: { x: 0.0, y: 0, w: "10%", h: "13%", fill: colors.greenYellow } },
        // IMAGE POSITIONING
        // { rect: { x: "83%", y: 0, w: "15%", h: "13%", fill: colors.greenYellow } },
        // { rect: { x: "50%", y: "68%", w: "50%", h: "30%", fill: colors.darkRed } },
        // TextA
        // {
        //   rect: {
        //     x: "50%",
        //     y: "68%",
        //     w: "50%",
        //     h: "9%",
        //     fill: colors.greenYellow
        //   }
        // },
        // TextB
        // {
        //   rect: {
        //     x: "50%",
        //     y: "77.5%",
        //     w: "50%",
        //     h: "7.5%",
        //     fill: "ffffff"
        //   }
        // },
        //TextC
        // {
        //   rect: {
        //     x: "50%",
        //     y: "85.5%",
        //     w: "50%",
        //     h: "6%",
        //     fill: "ffffff"
        //   }
        // },
        //TextD
        // {
        //   rect: {
        //     x: "50%",
        //     y: "91.8%",
        //     w: "50%",
        //     h: "6%",
        //     fill: "ffffff"
        //   }
        // },

        // { rect: { x: '55%', y: '68%', w: "50%", h: "30%", fill: colors.yellow } },
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
              fontFace: "Calibri",
              color: `${rgbToHex(109, 110, 113)}`,
              fontWeight: "bold",
              fontSize: 25,
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
              fontFace: "Calibri",
              color: `${rgbToHex(109, 110, 113)}`,
              fontSize: 18,
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
              fontFace: "Calibri",
              color: `${rgbToHex(109, 110, 113)}`,
              fontSize: 18,
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
              fontFace: "Calibri",
              color: `${rgbToHex(155, 187, 89)}`,
              fontSize: 17,
              fontWeight: "bolder",
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
      margin: [0.25, 0.20, 0.25, 0.20],
      objects: [
        // Body
        {
          rect: { x: 0.0, y: "17%", w: "100%", h: "11%", fill: BRANDCOLOR_GREEN }
        },
        {
          rect: { x: 0.0, y: "29%", w: "100%", h: "11%", fill: colors.yellow }
        },
        {
          rect: { x: 0.0, y: "41%", w: "100%", h: "11%", fill: BRANDCOLOR_GREEN }
        },
        {
          rect: { x: 0.0, y: "53%", w: "100%", h: "11%", fill: BRANDCOLOR_BLUE }
        },
        {
          rect: { x: 0.0, y: "65%", w: "100%", h: "11%", fill: '000000' }
        },
        {
          rect: { x: 0.0, y: "77%", w: "100%", h: "11%", fill: BRANDCOLOR_BLUE }
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
            x: 0.4,
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
              fontSize: 12,
              fontFace: "Calibri"
            },
            text: "Trustee Report"
          }
        }
      ],
      slideNumber: {
        x: "93%",
        y: "94%",
        color: TEXTCOLOR,
        fontFace: "Calibri",
        fontSize: 12
      }
    });
    // MASTER PAGE TEMPLATE
  }

  let arrTypes = typeof type === "string" ? [type] : type;
  arrTypes.forEach(function(type, idx) {
    eval("genSlides_" + type + "(pptx)");
  });

  pptx.save("SampleReportFile" + type + "_" + getTimestamp());
}

function genSlides_Content(pptx) {
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
      fontSize: 36,
      fontFace: "Calibri",
      color: "0b2135",
      align: "c",
      valign: "m",
      bold: true
    });
    // OVAL
    agendaSlide.addText("01", {
      shape: pptx.shapes.OVAL,
      x: '1%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
    agendaSlide.addText("02", {
      shape: pptx.shapes.OVAL,
      x: '10%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
    agendaSlide.addText("03", {
      shape: pptx.shapes.OVAL,
      x: '18%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_GREEN, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
    agendaSlide.addText("04", {
      shape: pptx.shapes.OVAL,
      x: '22%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_GREEN, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
    agendaSlide.addText("05", {
      shape: pptx.shapes.OVAL,
      x: '26%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
    agendaSlide.addText("06", {
      shape: pptx.shapes.OVAL,
      x: '26%',
      y: "16%",
      w: '6%',
      h: '9.9%',
      fill: { type: "solid", color: BRANDCOLOR_BLUE, alpha: 1 },
      align: "c",
      fontSize: 24,
      bold: true,
      fontFace: "Calibri",
      color: BGCOLORLGHT
    });
  }
}
