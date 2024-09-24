<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import DatePicker from '@/components/DatePicker.vue';
import LineChart from '@/components/LineChart.vue';
import axios from 'axios';
import api from '@/config/api';

const templateStore = useTemplateStore();
const documentStore = useDocumentStore();


const startDate = ref(new Date().toISOString().split('T')[0]);
const endDate = ref(new Date().toISOString().split('T')[0]);


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
        const response = await api.get('/document-history');
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
        // const response = await axios.get('http://localhost:8080/document-history');
        const response = await api.get('/document-history', {
            params: {
                startDate: startDate.value,
                endDate: endDate.value
            }
        });
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
        <div class="mb-2 bg-white font-semibold text-gray-500 rounded-md max-w-[350px] mx-auto shadow shadow-gray-400">
            <div class="text-semibold text-gray-500 text-xs text-center pt-1">CHOOSE DATES TO VIEW REPORT.</div>
            <div class="flex text-xs ml-[35px]">
                <DatePicker v-model="startDate" label="START DATE:  " id="start-date" class="pt-1 pb-2" />
            <DatePicker v-model="endDate" label="END DATE:  " id="end-date" class="pt-1 pb-2" />
            </div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-2 mb-2">
            <div class="bg-white border border-blue-100 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-blue-700">{{ totalTemplates }}</p>
                <h3 class="text-lg font-semibold text-blue-500">Total Templates</h3>
            </div>
            <div class="bg-white border border-warning-500 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-warning-600">{{ totalDocuments }}</p>
                <h3 class="text-lg font-semibold text-warning-600">Total Documents</h3>
            </div>
            <div class="bg-white border border-green-300 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-green-700">{{ successfulGenerations }}</p>
                <h3 class="text-lg font-semibold text-green-500">Successful Generations</h3>
            </div>
            <div class="bg-white border border-red-100 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-red-500">{{ failedGenerations }}</p>
                <h3 class="text-lg font-semibold text-red-300">Failed Generations</h3>
            </div>
            <div class="bg-white border border-gray-400 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-gray-600">{{ generationRate.toFixed(2) }}%</p>
                <h3 class="text-lg font-semibold text-gray-500">Generation Rate (daily)</h3>
            </div>
            <div class="bg-white border border-blue-100 p-1 rounded-lg shadow text-center">
                <p class="text-3xl font-bold text-blue-300">{{ failureRate.toFixed(2) }}%</p>
                <h3 class="text-lg font-semibold text-blue-200">Failure Rate</h3>
            </div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow w-full overflow-y-auto">
            <line-chart :chartData="chartData" />
        </div>
    </div>
</template>
