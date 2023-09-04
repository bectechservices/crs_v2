function rgbToHex(r, g, b) {
    return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

am4core.ready(function () {
    for (let i = 0; i < window.NumberOfRecords; i++) {
        am4core.useTheme(am4themes_animated);
        const el = document.querySelector(`#aucChart_${i}`);
        if (el) {
            let chart = am4core.create(
                `aucChart_${i}`,
                am4charts.PieChart3D
            );

            chart.hiddenState.properties.opacity = 0;

            chart.data = window[`AUCAssetByClassData_${i}`];

            let series = chart.series.push(new am4charts.PieSeries());
            
            chart.innerRadius = am4core.percent(45);
            series.slices.template.stroke = am4core.color("#f8f8f8");
            series.slices.template.strokeWidth = 1;
            series.slices.template.strokeOpacity = 1;
            
            chart.legend = new am4charts.Legend();
            chart.legend.position = "right";
            
            // Disable ticks and labels
            series.labels.template.disabled = true;
            series.ticks.template.disabled = true;
            
            // Disable tooltips
            series.slices.template.tooltipText = "";

            series.dataFields.value = "dataValue";
            series.dataFields.category = "dataName";
            series.colors.list = [
                am4core.color(`${rgbToHex(79, 129, 189)}`),//Corporate Bond
                am4core.color(`${rgbToHex(192, 80, 77)}`),//Fixed Deposits
                am4core.color(`${rgbToHex(56, 210, 0)}`),//Govt bond
                am4core.color(`${rgbToHex(212, 212, 212)}`),//Equities
                am4core.color(`${rgbToHex(75, 172, 198)}`),//Treasury Bills
                am4core.color(`${rgbToHex(0, 139, 179)}`),//Receivables
                am4core.color(`${rgbToHex(32,172, 255)}`),//Cash Balance
            ];
        }
    }
});
