<template>
    <div class="w-full h-96">
        <canvas ref="lineChart"></canvas>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { Chart, registerables } from 'chart.js';

Chart.register(...registerables);

const props = defineProps<{
    data: {
        labels: string[];
        datasets: { label: string; data: number[]; borderColor: string; fill: boolean }[];
    };
}>();

const lineChart = ref<HTMLCanvasElement | null>(null);
let chartInstance: Chart | null = null;

onMounted(() => {
    if (lineChart.value) {
        chartInstance = new Chart(lineChart.value, {
            type: 'line',
            data: props.data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                layout: {
                    padding: {
                        left: 10,
                        right: 10,
                        top: 10,
                        bottom: 10,
                    },
                },
                scales: {
                    x: {
                        display: true,
                        title: {
                            display: true,
                            text: 'Date',
                        },
                    },
                    y: {
                        display: true,
                        title: {
                            display: true,
                            text: 'Count',
                        },
                    },
                },
            },
        });
    }
});

watch(
    () => props.data,
    (newData) => {
        if (chartInstance) {
            chartInstance.data = newData;
            chartInstance.update();
        }
    }
);
</script>

<style scoped>
.w-full {
    width: 100%;
}

.h-96 {
    height: 24rem; /* 96 * 0.25rem */
}

canvas {
    display: block;
    width: 100%;
    height: 100%;
}
</style>
