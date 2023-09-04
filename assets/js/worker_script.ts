import ShouldReadExcelFile from "./lib/IExcel";
import {isExcel97File, isExcelFile} from "./lib/helpers";
import Excel97FilePreview from "./lib/excel97";
import ExcelFilePreview from "./lib/excel";
import PVReportExcelParser from "./lib/pv-excel-parser";

onmessage = function (e) {
    e.data.forEach(async (file: File) => {
        let preview: ShouldReadExcelFile;
        if (isExcelFile(file)) {
            if (isExcel97File(file)) {
                preview = new Excel97FilePreview(file);
            } else {
                preview = new ExcelFilePreview(file);
            }
            let results: string[][] = [];
            try {
                results = await preview.parseDataToJson(true);
            } catch (e) {
                if (isExcel97File(file)) {
                    preview = new ExcelFilePreview(file);
                    results = await preview.parseDataToJson(true);
                }
            }
            if (results.length && (results[0][0] === '___PATCHED_PV__')) {
                results.shift();
                results.pop();
                results = Excel97FilePreview.parseConvertedPVToAKnownFormat(results);
            }
            postMessage({
                type: "data", data: PVReportExcelParser.parsePVReport(
                    results
                )
            })
        }
    })
    postMessage({
        type: "terminate", data: []
    })
};