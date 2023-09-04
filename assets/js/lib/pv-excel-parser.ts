import {BondData, ParsedPVReport} from './types';

export default class PVReportExcelParser {
    static reportDelimeter: string = "TOTAL";
    static securityDataDelimeter: string = "SECURITY NAME";
    static summaryDelimeter: string = "summary";
    static summaryDataDelimeter: string = "description";

    static parsePVReport(data: Array<Array<string>>): Array<ParsedPVReport> {
        let headers: Array<string> = [];
        let bonds: Array<BondData> = [];
        let summary: Array<Array<string>> = [];
        let readItems = 0;
        let isReadingBond = false;
        let isReadingSummary = false;
        let currentBond: BondData = {
            bond: "",
            values: []
        };
        const reports: Array<ParsedPVReport> = [];

        data.forEach((datum: Array<string>, index: number) => {
            if (datum.length > 0) {
                if ((readItems < 3)) {
                    if (readItems == 2) {
                        if (data.length > 3 && data[3].length > 0) {
                            if (["safekeeping account:", "client:"].includes((data[3][0]).trim().toLowerCase())) {
                                headers[1] = `${headers[1]}${datum[0]}`;
                                headers[2] = `${data[3][0]} ${data[4][0]}`;
                                readItems += 2;
                                return;
                            }
                        }
                    }
                    if (datum[0] && !datum[0].includes("THE PRICES QUOTED ARE INTENDED FOR INTERNAL ADMINISTRATIVE")) {
                        headers.push(datum[0]);
                        readItems++;
                    }
                    return;
                }
                if (!isReadingBond && !isReadingSummary) {
                    if ((datum.length == 1) && (data.length > index + 1)) {
                        const nextRow = data[index + 1];
                        if (nextRow.length > 0 && nextRow[0] == "SECURITY NAME") {
                            console.log("RowNum: ",index+1)
                            console.log("NextRow: ",nextRow)
                            isReadingBond = true;
                            currentBond.bond = datum[0];
                            
                            return;
                        }
                    }
                } else if (!isReadingSummary) {
                    if ((datum[0] || "").trim().toUpperCase() == PVReportExcelParser.reportDelimeter) {
                        isReadingBond = false;
                        currentBond.values.push(datum);
                        bonds.push(currentBond);
                        currentBond = {
                            bond: "",
                            values: []
                        };
                        return;
                    }
                    if ((datum[0] || "") != PVReportExcelParser.securityDataDelimeter) {
                        currentBond.values.push(datum);
                    }
                }
                if (!isReadingSummary && !isReadingBond) {
                    if ((datum[0] || "").trim().toLowerCase() == PVReportExcelParser.summaryDelimeter) {
                        isReadingSummary = true;
                        return;
                    }
                } else if (!isReadingBond) {
                    if ((datum[0] || "").trim().toUpperCase() == PVReportExcelParser.reportDelimeter) {
                        summary.push(datum);
                        isReadingSummary = false;
                        readItems = 0;
                        reports.push({
                            headers,
                            bonds,
                            summary
                        });
                        headers = [];
                        bonds = [];
                        summary = [];
                        return;
                    }
                    if ((datum[0] || "").trim().toLowerCase() != PVReportExcelParser.summaryDataDelimeter) {
                        summary.push(datum);
                    }
                }
            }
        });

        return reports;
    }

}