<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import LineChart from '@/components/LineChart.vue';
import axios from 'axios';

const templateStore = useTemplateStore();
const documentStore = useDocumentStore();

const totalTemplates = ref(0);
const totalDocuments = ref(0);
const generationRate = ref(0);
const failureRate = ref(0);
const successfulGenerations = ref(0);
const failedGenerations = ref(0);
const chartData = ref([]);

onMounted(async () => {
    await fetchMetrics();
    await fetchChartData();
});

async function fetchMetrics() {
    await templateStore.fetchTemplates();
    await documentStore.fetchDocuments();

    totalTemplates.value = templateStore.templates.length;
    totalDocuments.value = documentStore.documents.length;

    successfulGenerations.value = documentStore.documents.length;
    failedGenerations.value = documentStore.documents.filter(doc => doc.status === 'failure').length;
    const totalGenerations = successfulGenerations.value + failedGenerations.value;

    // Fetch document history to determine the number of days
    const documentHistory = await fetchDocumentHistory();
    const numberOfDays = documentHistory.length;

    generationRate.value = totalGenerations / numberOfDays;
    failureRate.value = (failedGenerations.value / totalGenerations) * 100;
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
            date: entry.date.trim(),
            count: entry.count
        }));
    } catch (error) {
        console.error('Error fetching document history:', error);
        return [];
    }
}

async function fetchChartData() {
    try {
        const response = await axios.get('http://localhost:8080/document-history');
        if (response.status !== 200) {
            throw new Error('Failed to fetch document history');
        }
        const responseData = response.data;
        if (responseData.code !== 200) {
            throw new Error('Failed to fetch document history');
        }
        const documentHistory = responseData.data;

        chartData.value = documentHistory.map(entry => ({
            label: entry.date.trim(),
            y: entry.count
        }));

        console.log('Chart data:', chartData.value);
    } catch (error) {
        console.error('Error fetching document history:', error);
    }
}
</script>

<template>
    <div class="p-0">
        <div class="mb-3">
            <DatePicker v-model="startDate" label="Start Date" id="start-date" />
            <DatePicker v-model="endDate" label="End Date" id="end-date" />
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-3 mb-3">
            <div class="bg-white border border-blue-100 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-blue-700">{{ totalTemplates }}</p>
                <h3 class="text-lg font-semibold text-blue-500">Total Templates</h3>
            </div>
            <div class="bg-white border border-warning-500 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-warning-600">{{ totalDocuments }}</p>
                <h3 class="text-lg font-semibold text-warning-600">Total Documents</h3>
            </div>
            <div class="bg-white border border-green-300 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-green-700">{{ successfulGenerations }}</p>
                <h3 class="text-lg font-semibold text-green-500">Successful Generations</h3>
            </div>
            <div class="bg-white border border-red-100 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-red-500">{{ failedGenerations }}</p>
                <h3 class="text-lg font-semibold text-red-300">Failed Generations</h3>
            </div>
            <div class="bg-white border border-gray-400 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-gray-600">{{ generationRate.toFixed(2) }}%</p>
                <h3 class="text-lg font-semibold text-gray-500">Generation Rate (daily)</h3>
            </div>
            <div class="bg-white border border-blue-100 p-5 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-blue-300">{{ failureRate.toFixed(2) }}%</p>
                <h3 class="text-lg font-semibold text-blue-200">Failure Rate</h3>
            </div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow w-full h-full">
            <line-chart :chartData="chartData" />
        </div>
    </div>
</template>