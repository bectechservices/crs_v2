google.charts.load('current', {'packages': ['bar', 'corechart']});
google.charts.setOnLoadCallback(drawChart);


let dynamoColor = () => {
    let r = Math.floor((Math.random()* 5) * 225);
    let g = Math.floor((Math.random()* 3)* 220);
    let b = Math.floor((Math.random()* 5)* 200);
    let rgb = "rgba(" + r + "," + g + "," + b + ", 0.7)";

    function rgbToHex(swatch) {
        return ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
    }

    return "#" + rgbToHex(rgb)
}

let generateColors = function (number) {
    let result = [];
    for (let i = 0; i < number; i++) {
        result.push(dynamoColor())
    }
    return result;
}

function drawChart() {
    for (let i = 0; i < window.NumberOfRecords; i++) {
        var data = google.visualization.arrayToDataTable([
            ["Quarters", ...window[`SummaryForPast3MonthsLabels_${i}`]],
            ["Q4", ...window[`DataForMonth1_${i}`]],
            ["Q1", ...window[`DataForMonth2_${i}`]],
            ["Q2", ...window[`DataForMonth3_${i}`]],
        ]);
        var stackerData = google.visualization.arrayToDataTable([
            ["Quarters", ...window[`SummaryForPast3MonthsLabels_${i}`]],
            ["Q4", ...window[`DataForMonth1_${i}`]],
            ["Q1", ...window[`DataForMonth2_${i}`]],
            ["Q2", ...window[`DataForMonth3_${i}`]],
        ]);

        var options = {
            chart: {
                title: '',
                subtitle: '',
                //   legend: { position: 'bottom' }
            },
            colors: generateColors(window[`DataForMonth3_${i}`].length)
        };

        var options_fullStacked = {
            isStacked: 'percent',
            height: 300,
            legend: {position: 'top', maxLines: 3},
            hAxis: {
                // title:'',
                minValue: 0,
                ticks: [0, .3, .6, .9, 1]
            },
            vAxis: {
                // title:'',
            },
            is3D: true,
            colors: generateColors(window[`DataForMonth3_${i}`].length)
        };

        var chart = new google.charts.Bar(document.getElementById(`aucBarVolumes_${i}`));
        chart.draw(data, google.charts.Bar.convertOptions(options));

        var aucLines = new google.visualization.LineChart(document.getElementById(`aucLines_${i}`));
        aucLines.draw(data, options);

        var aucStacker = new google.visualization.BarChart(document.getElementById(`aucStacked_${i}`));
        aucStacker.draw(stackerData, options_fullStacked);
    }
}