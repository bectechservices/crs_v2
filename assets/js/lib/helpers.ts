//@ts-ignore
import $ from "jquery";
import { endOfMonth, format, getQuarter, getYear, parse, subMonths } from "date-fns"

export function csrf_token(): string {
    return $('meta[name="csrf-token"]').attr("content");
}

export function worker_script(): string {
    return $('meta[name="worker_script"]').attr("content");
}

export function isExcelFile(file: File): boolean {
    return ["application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-excel"].includes(file.type)
}

export function isExcel97File(file: File): boolean {
    return file.type == "application/vnd.ms-excel"
}

export function isWordFile(file: File): boolean {
    return ["application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/msword"].includes(file.type)
}

export function isPDFFile(file: File): boolean {
    return "application/pdf" === file.type
}

export function isImage(file: File): boolean {
    return ["image/jpeg", "image/png"].includes(file.type)
}

export function convertFileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (event: ProgressEvent) => {
            let data = (event.target as any).result;
            try {
                resolve(btoa(data));
            } catch (e) {
                reject(e);
            }
        };
        reader.readAsBinaryString(file);
    });
}

export async function dataURLToBlob(dataURL: string): Promise<Blob> {
    return await (await fetch(dataURL)).blob()
}

export function lastDayOfQuarter(): string {
    const date = new Date();
    const month = date.getMonth() + 1;
    const year = date.getFullYear();

    if (month >= 1 && month <= 3) {
        return `31st December ${year - 1}`
    } else if (month >= 4 && month <= 6) {
        return `31st March ${year}`
    } else if (month >= 7 && month <= 9) {
        return `30th June ${year}`
    } else {
        return `30th September ${year}`
    }
}

export function lastSecLicenseRenewalDate(): string {
    const date = new Date();
    const month = date.getMonth() + 1;
    let year = date.getFullYear();
    if (month > 6) {
        return `1st July ${year + 1}`
    } else {
        return `1st July ${year}`
    }
}

export function currentQuarter(): string {
    const date = new Date();
    const month = date.getMonth() + 1;
    const year = date.getFullYear();

    if (month >= 1 && month <= 3) {
        return `Fourth Quarter ${year - 1}`
    } else if (month >= 4 && month <= 6) {
        return `First Quarter ${year}`
    } else if (month >= 7 && month <= 9) {
        return `Second Quarter ${year}`
    } else {
        return `Third Quarter ${year}`
    }
}

export function makeLastDayOfQuarter(quarter: string, year: string): string {
    if (Boolean(quarter) && Boolean(year)) {
        if (quarter == "1") {
            return `31st March ${year}`
        } else if (quarter == "2") {
            return `30th June ${year}`
        } else if (quarter == "3") {
            return `30th September ${year}`
        } else {
            return `31st December${year}`
        }
    }
    return lastDayOfQuarter();
}

export function makeLastSecLicenseRenewalDate(quarter: string, year: string): string {
    if (Boolean(quarter) && Boolean(year)) {
        const quarterInt = parseInt(quarter);
        const yearInt = parseInt(year);
        if (quarterInt > 2) {
            return `1st July ${yearInt + 1}`
        }
        return `1st July ${year}`
    }
    return lastSecLicenseRenewalDate();
}

export function makeCurrentQuarter(quarter: string, year: string): string {
    if (Boolean(quarter) && Boolean(year)) {
        if (quarter == "1") {
            return `First Quarter ${year}`
        } else if (quarter == "2") {
            return `Second Quarter ${year}`
        } else if (quarter == "3") {
            return `Third Quarter ${year}`
        } else {
            return `Fourth Quarter ${year}`
        }
    }
    return currentQuarter();
}

export function isNumeric(number: number | string) {
    return !isNaN(parseFloat(number as string)) && isFinite(number as number);
}

export function currentQuarterNumber(): string {
    const date = new Date()
    const month = date.getMonth() + 1;

    if (month >= 1 && month <= 3) {
        return "4"
    } else if (month >= 4 && month <= 6) {
        return "1"
    } else if (month >= 7 && month <= 9) {
        return "2"
    } else {
        return "3"
    }
}
export function currentMonthNumber(month: string): string {
    
    //const date = new Date()
    //const month = date.getMonth() + 1;

    if (month = "01") {
        return "01"
    } else if (month = "02") {
        return "02"
    } else if (month = "03") {
        return "03"
    } else if (month = "04") {
        return "04"
    } else if (month = "05") {
        return "05"
    }  else if (month = "06") {
        return "06"
    }else if (month = "07") {
        return "07"
    } else if (month = "08") {
        return "08"
    }else if (month = "09") {
        return "09"
    } else if (month = "10") {
        return "10"
    } else if (month = "11") {
        return "11"
    } else {
        return "12"
    }
}

export function currentMonthNumber(): number {
    const date = new Date();
    return date.getMonth() + 1;
}

export function formatMoney(amount: number, dp: number = 2) {
    return new Intl.NumberFormat('en-GH', { minimumFractionDigits: dp }).format(amount)
}

export function getShortDate(): string {
    const quarter = currentQuarterNumber();
    const year = currentYear();

    if (quarter == "1") {
        return `31st Mar ${year}`
    } else if (quarter == "2") {
        return `30th Jun ${year}`
    } else if (quarter == "3") {
        return `30th Sep ${year}`
    } else {
        return `31st Dec ${year}`
    }
}

//END OF DATE CHANGES
export function makeShortDate(quarter: string, year: string): string {
    if (quarter == "1") {
        return `31st Mar ${year}`
    } else if (quarter == "2") {
        return `30th Jun ${year}`
    } else if (quarter == "3") {
        return `30th Sep ${year}`
    } else {
        return `31st Dec ${year}`
    }
}

export function getPreviousQuarterShortDate(): string {
    const quarter = currentQuarterNumber();
    const year = (new Date()).getFullYear();
    if (quarter == "1") {
        return `31st Dec ${year - 1}`
    } else if (quarter == "2") {
        return `31st Mar ${year}`
    } else if (quarter == "3") {
        return `30th Jun ${year}`
    } else {
        return `30th Sep ${year}`
    }
}

export function makePreviousQuarterShortDate(quarter: string, year: string): string {
    if (quarter == "1") {
        return `31st Dec ${parseInt(year) - 1}`
    } else if (quarter == "2") {
        return `31st Mar ${year}`
    } else if (quarter == "3") {
        return `30th Jun ${year}`
    } else {
        return `30th Sep ${year}`
    }
}

export const CURRENT_YEAR = `${currentYear()}`;

export function makeNPRAReportShortDate(): string {
    const date = subMonths(parse(new Date()), 1);
    const month = date.getMonth() + 1;
    let monthStr = `${month}`;
    if (month < 10) {
        monthStr = `0${month}`
    }
    return `${`${date.getFullYear()}`.substr(2, 4)}${monthStr}`
}

export function makeQuarterBeginning(): string {
    const quarter = currentQuarterNumber();
    if (quarter == "1") {
        return `1st January`
    } else if (quarter == "2") {
        return `1st April`
    } else if (quarter == "3") {
        return `1st July`
    } else {
        return `1st October`
    }
}

export function parseStringToFloat(value: string): number {
    const regex = /[A-Za-z,/]/gi;
    if (!value)
        return 0;
    value = value.replace(regex, '');
    if (!value)
        return 0;
    return parseFloat(value);
}

export function parseStringToInt(value: string): number {
    const regex = /[A-Za-z,/]/gi;
    if (!value)
        return 0;
    value = value.replace(regex, '');
    if (!value)
        return 0;
    return parseInt(value);
}

export function getShortMonthNamesForPastMonths(quarter: number): Array<string> {
    const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    let dates = [];
    if (quarter === 1) {
        let num = 9;
        for (let i = 0; i < 6; i++) {
            if (num > 11) {
                num = 0;
            }
            dates.push(months[num++]);
        }
    } else {
        for (let i = 0; i < quarter * 3; i++) {
            dates.push(months[i]);
        }
    }
    return dates;
}

export function getShortMonthNamesForPastMonthsRelativeQuarter(quarter: number): Array<string> {
    const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    let dates = [];
    for (let i = 0; i < quarter * 3; i++) {
        dates.push(months[i]);
    }
    return dates;
}
export function searchMonthNumber(month: number): Array<string> {
    const months = ["01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"];
    let dates = [];
    for (let i = 0; i < month * 1; i++) {
        dates.push(months[i]);
    }
    return dates;
}

export function currentYear(): number {
    let quarter = currentQuarterNumber();
    let year = (new Date()).getFullYear();
    if (quarter == "4") {
        return year - 1
    }
    return year;
}

export function getCurrentQuarterFormalDate() {
    const quarter = currentQuarterNumber();
    const year = (new Date()).getFullYear();
    if (quarter == "1") {
        return `${year}-03-31`
    } else if (quarter == "2") {
        return `${year}-06-30`
    } else if (quarter == "3") {
        return `${year}-09-30`
    } else {
        return `${year - 1}-12-31`
    }
}

export function makeQuarterFormalDate(quarter: string, year: string) {
    if (quarter == "1") {
        return `${year}-03-31`
    } else if (quarter == "2") {
        return `${year}-06-30`
    } else if (quarter == "3") {
        return `${year}-09-30`
    } else {
        return `${year}-12-31`
    }
}

export function makeMonthFormalDate(month: string, year: string): string {
    const date = parse(new Date(parseInt(year), parseInt(month) - 1));
    return format(endOfMonth(date), 'YYYY-MM-DD');
}

export function makePastQuarters(label: string, quarter: number, make: number): Array<string> {
    let output: Array<string> = [];
    let currentQuarter = quarter + 1; //to include current quarter
    for (let i = 0; i < make; i++) {
        currentQuarter--;
        if (currentQuarter == 0) {
            currentQuarter = 4;
        }
        output.push(`${label}${currentQuarter}`)
    }
    return output;
}

export function arrayValuesToPercentage(data: Array<number>): Array<number> {
    let totalValues = 0;
    if (data.length) {
        totalValues = data.reduce((accumulator: number, currentValue: number) => accumulator + currentValue);
    }
    return data.map((each: number) => ((each / totalValues) * 100) / 100)
}

export function sumMemberVariableInObjectArray(data: Array<any>, field: string): number {
    let total: number = 0;
    data.forEach((each: any) => {
        total += each[field];
    });
    return total;
}

export function benchmark(name: string) {
    const start = new Date();
    return {
        stop: function () {
            const end = new Date();
            const time = end.getTime() - start.getTime();
            console.log('Timer:', name, 'finished in', time, 'ms');
        }
    }
}

export function makeClientInfoForPVList(headers: Array<string>): string {
    const dateMatches = (headers[1] as any).match(/-.+/g);
    let date = "";
    if (dateMatches && dateMatches.length > 0) {
        date = dateMatches[0];
    }
    const nameMatches = (headers[2] as any).match(/-.+/g);
    let name = "";
    if (nameMatches && nameMatches.length > 0) {
        name = nameMatches[0].replace(/-/g, '');
    }
    return `${name} ${date}`
}

export function makeTrusteeQuarterFancyQuarterDateFromDate(date: string) {
    const parsedDate = parse(date);
    return `Q${getQuarter(parsedDate)} ${getYear(parsedDate)}`;
}