export default interface ShouldReadExcelFile {
    parseDataToJson(mergeSheets: boolean): Promise<Array<Array<string>>>
}