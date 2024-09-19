<template>
    <div class="chart-container">
        <canvasjs-chart :options="chartOptions"></canvasjs-chart>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { CanvasJSChart } from '@canvasjs/vue-charts';

const chartOptions = ref({
    animationEnabled: true,
    title: {
        text: "Sample Chart"
    },
    axisX: {
        title: "X Axis Title"
    },
    axisY: {
        title: "Y Axis Title"
    },
    data: [{
        type: "line",
        name: "Sample Data",
        dataPoints: [
            { x: 1, y: 10 },
            { x: 2, y: 15 },
            { x: 3, y: 20 }
        ]
    }]
});
</script>

<style scoped>
.chart-container {
    width: 100%;
    height: 400px;
}
</style>
