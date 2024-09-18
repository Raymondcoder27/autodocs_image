
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import LineChart from '@/components/LineChart.vue';

const templateStore = useTemplateStore();
const documentStore = useDocumentStore();

const totalTemplates = ref(0);
const totalDocuments = ref(0);
const generationRate = ref(0);
const failureRate = ref(0);
const chartData = ref({ labels: [], datasets: [] });

onMounted(async () => {
    await fetchMetrics();
    await fetchChartData();
});

async function fetchMetrics() {
    await templateStore.fetchTemplates();
    await documentStore.fetchDocuments();

    totalTemplates.value = templateStore.templates.length;
    totalDocuments.value = documentStore.documents.length;

    // Calculate generation rate and failure rate
    const successfulGenerations = documentStore.documents.filter(doc => doc.status === 'success').length;
    const failedGenerations = documentStore.documents.filter(doc => doc.status === 'failure').length;
    const totalGenerations = successfulGenerations + failedGenerations;

    generationRate.value = totalGenerations / 7; // Assuming weekly data
    failureRate.value = (failedGenerations / totalGenerations) * 100;
}

async function fetchChartData() {
    // Fetch historical data for the line chart
    const documentHistory = await documentStore.fetchDocumentHistory(); // Assuming this method exists
    const templateHistory = await templateStore.fetchTemplateHistory(); // Assuming this method exists

    const labels = documentHistory.map(entry => entry.date);
    const documentData = documentHistory.map(entry => entry.count);
    const templateData = templateHistory.map(entry => entry.count);

    chartData.value = {
        labels,
        datasets: [
            {
                label: 'Documents Generated',
                data: documentData,
                borderColor: 'blue',
                fill: false,
            },
            {
                label: 'Templates Uploaded',
                data: templateData,
                borderColor: 'green',
                fill: false,
            },
        ],
    };
}
</script>

<template>
    <div class="p-5">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-5">
            <div class="bg-white p-5 rounded-lg shadow text-center">
                <p class="text-2xl">{{ totalTemplates }}</p>
                <h3 class="text-lg font-semibold">Total Templates</h3>
            </div>
            <div class="bg-white p-5 rounded-lg shadow text-center">
                <p class="text-2xl">{{ totalDocuments }}</p>
                <h3 class="text-lg font-semibold">Total Documents</h3>
            </div>
            <div class="bg-white p-5 rounded-lg shadow text-center">
                <p class="text-2xl">{{ generationRate }} per day</p>
                <h3 class="text-lg font-semibold">Generation Rate</h3>
            </div>
            <div class="bg-white p-5 rounded-lg shadow text-center">
                <p class="text-2xl">{{ failureRate }}%</p>
                <h3 class="text-lg font-semibold">Failure Rate</h3>
            </div>
        </div>
        <div class="bg-white p-5 rounded-lg shadow">
            <line-chart :data="chartData" />
        </div>
    </div>
</template>
