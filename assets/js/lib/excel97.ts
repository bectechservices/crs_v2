import ShouldReadExcelFile from "./IExcel";
import * as parse5 from "parse5";

export default class Excel97FilePreview implements ShouldReadExcelFile {
    private readonly file: File;
    private static isReadingSummary: boolean = false;
    private static shouldSkipNTimes: number = 0;

    constructor(file: File) {
        this.file = file;
    }

    parseDataToJson(mergeSheets: boolean = false): Promise<Array<Array<string>>> {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = (event: ProgressEvent) => {
                let data = (event.target as any).result;
                const content = this.getHTMLBody(data);
                if (!content) {
                    reject("File content empty")
                }
                const parsedHTML = parse5.parse(content);
                resolve(this.readExcelDataFromHTML((parsedHTML as any).childNodes[0].childNodes[1].childNodes.filter((each: any) => each.attrs && each.attrs.length > 0)))
            };
            reader.readAsText(this.file);
        });
    }

    private getHTMLBody(html: string): string {
        const regex = /(?:<body[^>]*>)(.*)(?:<\/body>)/s;
        const matches = regex.exec(html);
        if (matches !== null && matches.length > 1) {
            return matches[1];
        }
        return "";
    }

    private readExcelDataFromHTML(content: Array<any>): Array<Array<string>> {
        const results: Array<Array<string>> = [];
        content.forEach((data: any) => {
            const result = this.readDomElementRecursively(data);
            if (result.length) {
                results.push(...this.__parseReturnedDomArrayInAHackyWay(result));
            }
        });
        return results;
    }

    private readDomElementRecursively(element: any): Array<string> {
        let result: Array<string> = [];
        if (element.value) {
            result.push(element.value);
            if ((element.value.trim().toLowerCase() === "total") && element.parentNode && element.parentNode.nodeName === "b") {
                if (Excel97FilePreview.isReadingSummary) {
                    Excel97FilePreview.shouldSkipNTimes = 2;
                    Excel97FilePreview.isReadingSummary = false;
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                } else {
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                    result.push("\n");
                    result.push("\n");
                    result.push("");
                }
            } else if (element.value.trim().toLowerCase() === "summary") {
                Excel97FilePreview.isReadingSummary = true;
            }
        }
        if (element.nodeName == "td" && !(element.childNodes as any[]).length) {
            result.push("");
            result.push("\n");
        }
        if (element.childNodes !== undefined) {
            element.childNodes.forEach((datum: any) => {
                if (Excel97FilePreview.shouldSkipNTimes > 0) {
                    Excel97FilePreview.shouldSkipNTimes--;
                } else {
                    result.push(...this.readDomElementRecursively(datum));
                }

            });
        }
        return result;
    }

    private __parseReturnedDomArrayInAHackyWay(data: Array<string>): Array<Array<string>> {
        const results: Array<Array<string>> = [];
        let numberOfNLs = 0;
        let isReadingNLs = false;
        let tempBuffer: Array<string> = [];
        data.forEach((datum: string) => {
            if (datum !== " ") {
                if (datum === "↵" || datum.charAt(0) === "\n") {
                    if (!isReadingNLs) {
                        isReadingNLs = true;
                        numberOfNLs++;
                    } else {
                        numberOfNLs++;
                    }
                } else {
                    if (numberOfNLs == 2 || numberOfNLs == 3) {
                        tempBuffer.push(datum)
                    } else {
                        if (tempBuffer.length) {
                            results.push(tempBuffer);
                            tempBuffer = [];
                        }
                        tempBuffer.push(datum);
                    }
                    numberOfNLs = 0;
                    isReadingNLs = false;
                }
            }
        });
        if (tempBuffer.length) {
            results.push(tempBuffer);
        }
        return results;
    }

    static hasOnlyNullsTillEndOfArray(array: Array<string>, currentPosition: number): boolean {
        let result = true;
        for (let i = currentPosition + 1; i < array.length; i++) {
            result = array[i] == undefined;
        }
        return result;
    }

    static parseConvertedPVToAKnownFormat(data: Array<Array<string>>): Array<Array<string>> {
        const results: Array<Array<string>> = [];
        let numberOfNLs = 0;
        let isReadingNLs = false;
        let tempBuffer: Array<string> = [];
        data.forEach((datum: Array<string>) => {
            if (datum.length && datum[0] !== "") {
                for (let i = 0; i < datum.length; i++) {
                    if ((datum[0] || "").startsWith("E.O. & E") || (datum[0] || "").startsWith("Date: "))
                        continue;
                    if (datum[i] == undefined) {
                        if (!isReadingNLs) {
                            isReadingNLs = true;
                            numberOfNLs++;
                        } else {
                            numberOfNLs++;
                        }
                    } else {
                        if (numberOfNLs == 3) {
                            tempBuffer.push("");
                            tempBuffer.push(datum[i]);
                        } else if (numberOfNLs == 5) {
                            tempBuffer.push("");
                            tempBuffer.push("");
                            tempBuffer.push(datum[i]);
                        } else if (numberOfNLs == 14) {
                            tempBuffer.push("");
                            tempBuffer.push("");
                            tempBuffer.push("");
                            tempBuffer.push("");
                            tempBuffer.push("");
                            tempBuffer.push(datum[i]);
                        } else {
                            tempBuffer.push(datum[i]);
                        }
                        numberOfNLs = 0;
                        isReadingNLs = false;
                    }
                }
            }
            results.push(tempBuffer);
            tempBuffer = [];
        });
        if (tempBuffer.length) {
            results.push(tempBuffer);
        }
        return results;
    }
}