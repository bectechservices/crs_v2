let Quartyear = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb", "Mar",]

// STYLE VARIABLES
var BRANDCOLOR_GREEN = "#77933c";
var BRANDCOLOR_BLUELGHT = "#4a7ebb";
var PRIMARYFONT_COLOR = "#0B2135";
var SHADOWCOLOR = "rgba(200,200,208,0.65)";
var SHADOWBLUR = 12;
var SHADOWOFFSETX = 0;
var SHADOWOFFSETY = 5;

function makeTrusteeDates(quarter) {
    let dates = [];
    if (quarter === 1) {
        let num = 9;
        for (let i = 0; i < 6; i++) {
            if (num > 11) {
                num = 0;
            }
            dates.push(Quartyear[num++]);
        }
    } else {
        for (let i = 0; i < quarter * 3; i++) {
            dates.push(Quartyear[i]);
        }
    }
    return dates;
}

function makeTrusteeDatesRelativeToCurrentQuarter(quarter) {
    let dates = [];
    for (let i = 0; i < quarter * 3; i++) {
        dates.push(Quartyear[i]);
    }
    return dates;
}

// CANVAS SETUP
// SETUP
let setup = canvas => {
    let context, dpr;
    context = canvas.getContext("2d");
    canvas.style.width = '100%';
    canvas.style.height = '100%';
    dpr = window.devicePixelRatio || 1.4;
    canvas.width = canvas.offsetWidth * dpr;
    canvas.height = canvas.offsetHeight * dpr;
    context.scale(dpr, dpr);
    return context;
};


for (let ll = 0; ll < window.NumberOfRecords; ll++) {
    let draw = Chart.controllers.line.prototype.draw;
    let volumeTrendOptions = {
        legend: {
            display: false,
            labels: {
                fontColor: "#0b2135",
                fontSize: 15,
                fontFamily: "'Helvetica', 'Arial'"
            }
        },
        scales: {
            xAxes: [{
                ticks: {
                    beginAtZero: true
                }
            }]
        },
        tooltips: {
            callbacks: {
                label: function (tooltipItem, data) {
                    var dataset = data.datasets[tooltipItem.datasetIndex];
                    var index = tooltipItem.index;
                    return dataset.labels[index] + ": " + dataset.data[index];
                }
            }
        },
        layout: {
            padding: {
                left: 35,
                right: 35,
                top: 0,
                bottom: 0
            }
        }
    };
    const aucTrendLabels = makeTrusteeDates(window.CRS_PPT_MISC_0.current_quarter_number);
    const aucLabels = makeTrusteeDatesRelativeToCurrentQuarter(window.CRS_PPT_MISC_0.current_quarter_number);

    const aucChartTrendEl = document.querySelector(`#aucChartTrend_${ll}`);
    if (aucChartTrendEl) {
        let ctx = setup(aucChartTrendEl);
        Chart.controllers.line = Chart.controllers.line.extend({
            draw: function () {
                let ctx = this.chart.chart.ctx;
                ctx.save();
                ctx.shadowColor = SHADOWCOLOR;
                ctx.shadowBlur = SHADOWBLUR;
                ctx.shadowOffsetX = SHADOWOFFSETX;
                ctx.shadowOffsetY = SHADOWOFFSETY;
                ctx.stroke();
                draw.apply(this, arguments);
                ctx.restore();
            }
        });

        // AUC TREND
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: aucTrendLabels,
                datasets: [{
                    data: window[`AUCChartTrendData_${ll}`],
                    labels: aucTrendLabels,
                    borderColor: "#3e95cd",
                    fill: false,
                    lineTension: 0,
                    pointBackgroundColor: "#fff",
                    pointBorderColor: "#ffb88c",
                    pointHoverBackgroundColor: "#ffb88c",
                    pointHoverBorderColor: "#fff",
                    pointRadius: 7,
                    pointHoverRadius: 7,
                }
                ]
            },
            options: {
                legend: {
                    display: false,
                    labels: {
                        fontColor: '#0b2135',
                        fontSize: 15,
                        fontFamily: "'Helvetica', 'Arial'",
                    }
                },
                scales: {
                    yAxes: [{
                        ticks: {
                            callback: function (value) {
                                if (value >= 1000 && value < 1000000) {
                                    return `${value / 1000}K`
                                } else if (value >= 1000000 && value < 1000000000) {
                                    return `${value / 1000000}M`
                                } else if (value >= 1000000000) {
                                    return `${value / 1000000000}B`
                                } else if (value < 2) {
                                    return parseFloat(value).toFixed(2);
                                }
                            }
                        }
                    }],
                    xAxes: [{
                        ticks: {
                            beginAtZero: true
                        }
                    }]
                },
                tooltips: {
                    callbacks: {
                        label: function (tooltipItem, data) {
                            var dataset = data.datasets[tooltipItem.datasetIndex];
                            var index = tooltipItem.index;
                            return dataset.labels[index] + ': ' + dataset.data[index];
                        }
                    }
                },
                layout: {
                    padding: {
                        left: 35,
                        right: 35,
                        top: 0,
                        bottom: 0
                    }
                }
            }
        });
    }
    const aucMonthTrendEl = document.querySelector(`#aucMonthTrend_${ll}`);
    if (aucMonthTrendEl) {
        let monthTrend = setup(aucMonthTrendEl);
        // Monthly Trend
        new Chart(monthTrend, {
            type: 'bar',
            data: {
                labels: aucLabels,
                datasets: [{
                    data: window[`TxnVolsTrendData_${ll}`],
                    labels: Quartyear,
                    borderColor: "#3e95cd",
                    fill: true,
                    backgroundColor: '#4a7ebb'
                }
                ]
            },
            options: volumeTrendOptions
        });
    }
    const volAssetsEl = document.querySelector(`#aucVolAsset_${ll}`);
    if (volAssetsEl) {
        let volAssets = setup(volAssetsEl);
        // Volume Assets
        new Chart(volAssets, {
            type: "line",
            data: {
                labels: window[`volAssetsMonths_${ll}`],
                datasets: window[`volAssetsGraphData_${ll}`],
            },
            options: {
                legend: {
                    display: false,
                    labels: {
                        fontColor: "#0b2135",
                        fontSize: 15,
                        fontFamily: "'Helvetica', 'Arial'"
                    }
                },
                scales: {
                    yAxes: [
                        {
                            ticks: {
                                beginAtZero: true
                            }
                        }
                    ]
                },
                tooltips: {
                    callbacks: {
                        label: function (tooltipItem, data) {
                            var dataset = data.datasets[tooltipItem.datasetIndex];
                            var index = tooltipItem.index;
                            return dataset.labels[index] + ": " + dataset.data[index];
                        }
                    }
                },
                layout: {
                    padding: {
                        left: 30,
                        right: 30,
                        top: 0,
                        bottom: 0
                    }
                }
            }
        });
    }
    const maturitiesFDEl = document.querySelector(`#maturitiesFD_${ll}`);
    if (maturitiesFDEl) {
        let maturitiesFD = setup(maturitiesFDEl);
        const fdMaturitiesData = window[`CRS_PPT_MISC_${ll}`].fdMaturities;
        // Maturities Graph
        new Chart(maturitiesFD, {
            type: "bar",
            data: {
                labels: Object.keys(fdMaturitiesData),
                datasets: [
                    {
                        data: Object.values(fdMaturitiesData),
                        labels: Object.keys(fdMaturitiesData),
                        borderColor: "#3e95cd",
                        fill: true,
                        backgroundColor: "#4a7ebb"
                    }
                ]
            },
            options: {
                legend: {
                    display: false,
                    labels: {
                        fontColor: "#0b2135",
                        fontSize: 15,
                        fontFamily: "'Helvetica', 'Arial'"
                    }
                },
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero: true,
                            stepSize: 1
                        }
                    }]
                },
                tooltips: {
                    callbacks: {
                        label: function (tooltipItem, data) {
                            var dataset = data.datasets[tooltipItem.datasetIndex];
                            var index = tooltipItem.index;
                            return dataset.labels[index] + ": " + dataset.data[index];
                        }
                    }
                },
                layout: {
                    padding: {
                        left: 35,
                        right: 35,
                        top: 0,
                        bottom: 0
                    }
                }
            }
        });
    }
}