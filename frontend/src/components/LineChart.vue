<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';

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

    const dataPoints = documentHistory.map(entry => ({
        label: entry.date,
        y: entry.count
    }));

    options.value.data[0].dataPoints = dataPoints;
}

async function fetchDocumentHistory() {
    // Mock implementation, replace with actual API call or logic
    return [
        { date: "Monday", count: 2 },
        { date: "Tuesday", count: 4 },
        { date: "Wednesday", count: 8 },
        { date: "Thursday", count: 4 },
        { date: "Friday", count: 10 },
        { date: "Saturday", count: 0 },
        { date: "Sunday", count: 6 }
    ];
}
</script>

<template>
    <CanvasJSChart :options="options" :style="styleOptions" @chart-ref="chartInstance"/>
</template> 
