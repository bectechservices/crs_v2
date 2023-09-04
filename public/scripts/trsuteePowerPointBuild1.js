// RGBTOHEX
function rgbToHex(r, g, b) {
  return ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

// DEFINE VARIABLES
let COMPANY_NAME = "BECTECH";
let CUST_NAME = "Standard Chartered";
let CLIENT_NAME = "Twifo Oil Palm Plantation";
let BRANDCOLOR_BLUE = `${rgbToHex(31, 73, 125)}`;
let BRANDCOLOR_GREEN = `${rgbToHex(119, 147, 60)}`;
let TEXTCOLOR = "999999";
let BGCOLORLGHT = "ffffff";
let CHARSPERLINE = 130;
// console.log(BRANDCOLOR_BLUE);

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
  slideBrandLineImage: "./media/Picture4.png",
  clientLogo: "./media/logo/twifo.png"
};
console.log(assetPaths.stanChartBgImage);
// INSTANTIATE PPTX
let pptx = new PptxGenJS();

// PRESENTATIONAL & LAYOUT PROPS
pptx.setAuthor(CUST_NAME);
pptx.setCompany(COMPANY_NAME);
pptx.setRevision("15");
pptx.setSubject("Trustee Report");
pptx.setTitle(`${CLIENT_NAME} Trustee Report Powerpoint`);

pptx.setLayout("LAYOUT_WIDE");

// TEST COLOR SUITE
let colors = {
  yellow: "7FFF00",
  blue: "0b2135",
  darkRed: "8B0000",
  greenYellow: "ADFF2F"
};
// TEST COLOR SUITE END

//BUILDING SLIDES

// MASTER SLIDE TEMPLATE
console.log(`${68 / 4.6 + 68.4}%`);
pptx.defineSlideMaster({
  title: "COVER_SLIDE",
  bkgd: BGCOLORLGHT,
  objects: [
    // { rect: { x: 0.0, y: 0, w: "100%", h: "13%", fill: colors.darkRed } },
    // { rect: { x: 0.0, y: "13%", w: "100%", h: "87%", fill: colors.blue } },
    // { rect: { x: 0.0, y: 0, w: "10%", h: "13%", fill: colors.greenYellow } },
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
          fontFace: "Helvetica",
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
          fontFace: "Arial",
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
          fontFace: "Arial",
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
          fontFace: "Arial",
          color: `${rgbToHex(155, 187, 89)}`,
          fontSize: 17,
          fontWeight: 'bolder',
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
        path: "./media/Picture4.png"
      }
    }
  ]
});
var slide = pptx.addNewSlide("COVER_SLIDE");
// slide.addText("How To Create PowerPoint Presentations", {
//   x: 0.5,
//   y: 0.7,
//   fontSize: 18
// });
