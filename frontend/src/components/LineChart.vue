<!-- <script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import axios from 'axios';

const chart = ref(null);
const templateStore = useTemplateStore();
const documentStore = useDocumentStore();

const options = ref({
    animationEnabled: true,
    exportEnabled: true,
    theme: "light2",
    title: {
        text: "Weekly Document Generation Report"
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
        dataPoints: []
    }]
});

const styleOptions = {
    width: "100%",
    height: "360px"
};

const chartInstance = (chartInstance: any) => {
    chart.value = chartInstance;
};

onMounted(async () => {
    await fetchMetrics();
});

async function fetchMetrics() {
    await templateStore.fetchTemplates();
    await documentStore.fetchDocuments();

    const documentHistory = await fetchDocumentHistory();
    console.log('Fetched document history:', documentHistory);

    const dataPoints = documentHistory.map(entry => ({
        label: entry.date,
        y: entry.count
    }));
    console.log('Data points for chart:', dataPoints);

    options.value.data[0].dataPoints = dataPoints;
}

async function fetchDocumentHistory() {
    try {
        const response = await axios.get('http://localhost:8080/document-history');
        if (response.status !== 200) {
            throw new Error('Failed to fetch document history');
        }
        const responseData = response.data;
        if (responseData.code !== 200) {
            throw new Error('Failed to fetch document history');
        }
        const data: { date: string, count: number }[] = responseData.data;
        return data.map((entry: { date: string, count: number }) => ({
            date: entry.date,   // Use proper date formatting if necessary
            count: entry.count
        }));
    } catch (error) {
        console.error('Error fetching document history:', error);
        return [];
    }
}
</script>

<template>
    <CanvasJSChart :options="options" :style="styleOptions" @chart-ref="chartInstance"/>
</template> -->


<script setup lang="ts">
import { defineProps, watch, onMounted,ref } from 'vue';
import { Chart, registerables } from 'chart.js';

Chart.register(...registerables);

const props = defineProps({
    data: {
        type: Object,
        required: true
    }
});

const chartRef = ref(null);

watch(() => props.data, (newData) => {
    if (chartRef.value) {
        chartRef.value.data = newData;
        chartRef.value.update();
    }
}, { deep: true });

onMounted(() => {
    const ctx = document.getElementById('line-chart').getContext('2d');
    chartRef.value = new Chart(ctx, {
        type: 'line',
        data: props.data,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: 'Chart.js Line Chart'
                }
            }
        }
    });
});
</script>

<template>
    <canvas id="line-chart"></canvas>
</template>