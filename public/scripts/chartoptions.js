let barVerticalOptions_stacked = {
  tooltips: {
    enabled: false
  },
  hover: {
    animationDuration: 0
  },
  scales: {
    xAxes: [
      {
        ticks: {
          beginAtZero: true,
          fontFamily: "'Open Sans Bold', sans-serif",
          fontSize: 11
        },
        scaleLabel: {
          display: false
        },
        gridLines: {},
        stacked: true
      }
    ],
    yAxes: [
      {
        gridLines: {
          display: false,
          color: "#fff",
          zeroLineColor: "#fff",
          zeroLineWidth: 0
        },
        ticks: {
          fontFamily: "'Open Sans Bold', sans-serif",
          fontSize: 11,
          callback: function(value){
            if(value >= 1000 && value < 1000000){
              return `${value/1000}K`
            }else if(value >= 1000000 && value < 1000000000){
              return `${value/1000000 }M`
            }else if(value >=1000000000  ){
              return `${value/1000000000}B`
            }else if(value < 2){
              return parseFloat(value).toFixed(2);
            }
          }
        },
        stacked: true
      }
    ]
  },
  legend: {
    display: false
  },
  layout: {
    padding: {
      left: 35,
      right: 35,
      top: 0,
      bottom: 0
    }
  },
  animation: {
    onComplete: function() {
      var chartInstance = this.chart;
      var ctx = chartInstance.ctx;
      ctx.textAlign = "center";
      ctx.font = "12px Open Sans";
      ctx.fillStyle = "#fff";

      Chart.helpers.each(
        this.data.datasets.forEach(function(dataset, i) {
          var meta = chartInstance.controller.getDatasetMeta(i);
          Chart.helpers.each(
            meta.data.forEach(function(bar, index) {
              data = dataset.data[index];
              if (i == 0) {
                ctx.fillText(data, bar._model.x - 2, bar._model.y + 50);
              } else {
                ctx.fillText(data, bar._model.x - 2, bar._model.y + 50);
              }
            }),
            this
          );
        }),
        this
      );
    }
  },
  pointLabelFontFamily: "Quadon Extra Bold",
  scaleFontFamily: "Quadon Extra Bold"
};
let barHorizontalOptions_stacked = {
  tooltips: {
    enabled: false
  },
  hover: {
    animationDuration: 0
  },
  scales: {
    xAxes: [
      {
        ticks: {
          beginAtZero: true,
          fontFamily: "'Open Sans Bold', sans-serif",
          fontSize: 11,
          display: false,
        },
        scaleLabel: {
          display: false
        },
        gridLines: {
          zeroLineColor: "transparent",
          display: false,
          drawTicks: false
        },
        stacked: true
      }
    ],
    yAxes: [
      {
        gridLines: {
          display: false,
          color: "#fff",
          zeroLineColor: "#fff",
          zeroLineWidth: 0
        },
        ticks: {
          fontFamily: "'Open Sans Bold', sans-serif",
          fontSize: 11
        },
        stacked: true
      }
    ]
  },
  legend: {
    display: false
  },
  layout: {
    padding: {
      left: 0,
      right: 0,
      top: 35,
      bottom: 35
    }
  },
  animation: {
    onComplete: function() {
      var chartInstance = this.chart;
      var ctx = chartInstance.ctx;
      ctx.textAlign = "left";
      ctx.font = "9px Open Sans";
      ctx.fillStyle = "#fff";

      Chart.helpers.each(
        this.data.datasets.forEach(function(dataset, i) {
          var meta = chartInstance.controller.getDatasetMeta(i);
          Chart.helpers.each(
            meta.data.forEach(function(bar, index) {
              data = dataset.data[index];
              ctx.fillText(data, bar._model.x - 20, bar._model.y + 4);
            }),
            this
          );
        }),
        this
      );
    }
  },
  pointLabelFontFamily: "Arial",
  scaleFontFamily: "Arial"
};

