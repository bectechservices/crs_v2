import * as XLSX from "xlsx";
import ShouldReadExcelFile from "./IExcel";

export default class ExcelFilePreview implements ShouldReadExcelFile {
    private readonly file: File;
    private readAsBinary: boolean = Boolean(
        typeof FileReader !== "undefined" &&
        (FileReader.prototype && FileReader.prototype.readAsBinaryString)
    );

    constructor(file: File) {
        this.file = file;
    }

    private fixData(data: any): string {
        var o = "",
            l = 0,
            w = 10240;
        for (; l < data.byteLength / w; ++l)
            o += String.fromCharCode.apply(null, new Uint8Array(
                data.slice(l * w, l * w + w)
            ) as any);
        o += String.fromCharCode.apply(null, new Uint8Array(
            data.slice(o.length)
        ) as any);
        return o;
    }

    private processWorkbooks(wb: XLSX.WorkBook, mergeSheets: boolean) {
        const results = this.toJson(wb, mergeSheets);
        if (results.PATCH && mergeSheets) {
            let data = Object.values(results).flat(1);
            data.unshift(["___PATCHED_PV__"]);
            return data;
        }
        return Object.values(results).flat(1);
    }

    private toJson(
        workbook: XLSX.WorkBook,
        mergeSheets: boolean
    ): {
        [key: string]: Array<Array<string>>;
    } {
        const result: any = {};
        if (mergeSheets) {
            workbook.SheetNames.forEach(function (sheetName: string) {
                const roa = XLSX.utils.sheet_to_json(workbook.Sheets[sheetName], {
                    raw: false,
                    header: 1
                });
                if (roa.length > 0) result[sheetName] = roa;
            });
            if (workbook.SheetNames.length > 1) {
                result['PATCH'] = '___PATCHED_PV___';
            }
        } else {
            const sheetName = workbook.SheetNames[0];
            const roa = XLSX.utils.sheet_to_json(workbook.Sheets[sheetName], {
                raw: false,
                header: 1
            });
            if (roa.length > 0) result[sheetName] = roa;
        }
        return result;
    }

    parseDataToJson(mergeSheets: boolean = false): Promise<Array<Array<string>>> {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = (event: ProgressEvent) => {
                let data = (event.target as any).result;
                const readtype: XLSX.ParsingOptions = {
                    type: this.readAsBinary ? "binary" : "base64"
                };
                if (!this.readAsBinary) {
                    const arr = this.fixData(data);
                    data = btoa(arr);
                }
                try {
                    const workbook = XLSX.read(data, readtype);
                    resolve(this.processWorkbooks(workbook, mergeSheets));
                } catch (e) {
                    reject(e);
                }
            };
            if (this.readAsBinary) {
                reader.readAsBinaryString(this.file);
            } else {
                reader.readAsArrayBuffer(this.file);
            }
        });
    }
}
