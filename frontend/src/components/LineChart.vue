<script setup lang="ts">
import { ref } from 'vue';

const chart = ref(null);

const options = {
    animationEnabled: true,
    exportEnabled: true,
    theme: "light2",
    title: {
        text: "Weekly Document generation report"
    },
    axisX: {
        title: "Day of Week",
        valueFormatString: "YYYY",
        labelTextAlign: "center",
        labelAngle: 0
    },
    axisY: {
        title: "No. of PDFs Generated",
        valueFormatString: "#"
    },
    data: [{
        type: "line",
        yValueFormatString: "# PDFs",
        dataPoints: [
        { label: "MON", y: 2 },
            { label: "TUE", y: 4 },
            { label: "WED", y: 8 },
            { label: "THUR", y: 4 },
            { label: "FRI", y: 10 },
            { label: "SAT", y: 0 },
            { label: "SUN", y: 6 }
            // { label: "MON", y: 800093.07 },
            // { label: "TUE", y: 5702000.87 },
            // { label: "WED", y: 585710.45 },
            // { label: "THUR", y: 10836638.18 },
            // { label: "FRI", y: 3387744.72 },
            // { label: "SAT", y: 7159285.80 },
            // { label: "SUN", y: 8000010.33 }
        ]
    }]
};

const styleOptions = {
    width: "100%",
    height: "360px"
};

const chartInstance = (chartInstance: any) => {
    chart.value = chartInstance;
};
</script>

<template>
    <CanvasJSChart :options="options" :style="styleOptions" @chart-ref="chartInstance"/>
</template>
