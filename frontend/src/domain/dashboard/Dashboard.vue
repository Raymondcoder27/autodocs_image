<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useTemplateStore } from "@/domain/templates/stores";
import { useDocumentStore } from "@/domain/documents/stores";
import DatePicker from "@/components/DatePicker.vue";
import LineChart from "@/components/LineChart.vue";
import axios from "axios";
import api from "@/config/api";
import { useLogStore } from "@/domain/requests/stores";

const templateStore = useTemplateStore();
const documentStore = useDocumentStore();
const logStore = useLogStore();

const startDate = ref(new Date().toISOString().split("T")[0]);
const endDate = ref(new Date().toISOString().split("T")[0]);

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
  // await fetchLogs();
});

async function fetchMetrics() {
  await templateStore.fetchTemplates();
  await documentStore.fetchDocuments();
  // await fetchLogs();

  totalTemplates.value = templateStore.templates.length;
  totalDocuments.value = documentStore.documents.length;

  successfulGenerations.value = documentStore.documents.length;
//   failedGenerations.value =
//     logStore.logs?.filter((l) => l.requestStatus === "FAILED").length || 0;

failedGenerations.value = documentStore.failedDocuments.length;
  const totalGenerations =
    successfulGenerations.value + failedGenerations.value;

  // Fetch document history
  const documentHistory = await fetchDocumentHistory();

  if (documentHistory && Array.isArray(documentHistory)) {
    const numberOfDays = documentHistory.length;

    generationRate.value = totalGenerations / numberOfDays;
    failureRate.value = (failedGenerations.value / totalGenerations) * 100;

    // Log document history for debugging
    console.log("Document History:", documentHistory);
  } else {
    console.error("Document history is empty or not in expected format.");
  }
}

async function fetchRange() {
  try {
    const response = await api.get("/daterange-metrics", {
      params: {
        startDate: startDate.value,
        endDate: endDate.value,
      },
    });
    if (response.status !== 200) {
      throw new Error("Failed to fetch metrics");
    }
    const responseData = response.data;
    if (responseData.code !== 200) {
      throw new Error("Failed to fetch metrics");
    }
    const data = responseData.data;

    totalTemplates.value = data.totalTemplates;
    totalDocuments.value = data.totalDocuments;
    failedGenerations.value = data.failedGenerations;
    generationRate.value = data.generationRate;
    failureRate.value = data.failureRate;

    console.log("Metrics:", data);
  } catch (error) {
    console.error("Error fetching metrics:", error);
  }
}

async function fetchDocumentHistory() {
  try {
    const response = await api.get("/document-history");
    if (response.status !== 200) {
      throw new Error("Failed to fetch document history");
    }
    const responseData = response.data;
    if (responseData.code !== 200) {
      throw new Error("Failed to fetch document history");
    }
    const data: { date: string; count: number }[] = responseData.data;
    return data.map((entry: { date: string; count: number }) => ({
      date: entry.date.trim(),
      count: entry.count,
    }));
  } catch (error) {
    console.error("Error fetching document history:", error);
    return [];
  }
}

async function fetchChartData() {
  try {
    // const response = await axios.get('http://localhost:8080/document-history');
    const response = await api.get("/document-history", {
      params: {
        startDate: startDate.value,
        endDate: endDate.value,
      },
    });
    if (response.status !== 200) {
      throw new Error("Failed to fetch document history");
    }
    const responseData = response.data;
    if (responseData.code !== 200) {
      throw new Error("Failed to fetch document history");
    }
    const documentHistory = responseData.data;

    chartData.value = documentHistory.map((entry) => ({
      label: entry.date.trim(),
      y: entry.count,
    }));

    console.log("Chart data:", chartData.value);
  } catch (error) {
    console.error("Error fetching document history:", error);
  }
}
</script>

<template>
  <div class="p-0">
    <div class="mb-2 bg-white font-semibold text-gray-700 rounded-lg max-w-sm mx-auto shadow-lg p-2">
  <div class="text-center text-gray-700 text-sm mb-0">
    Choose Dates to View Report
  </div>
  <div class="flex items-center space-x-0.5">
    <DatePicker
      v-model="startDate"
      label="Start Date:"
      id="start-date"
      class="flex-1 text-xs"
    />
    <DatePicker
      v-model="endDate"
      label="End Date:"
      id="end-date"
      class="flex-1 text-xs"
    />
    <button
      @click="fetchRange"
      class="py-1 mt-4 px-2 bg-blue-500 text-white text-xs rounded-md hover:bg-blue-600 transition duration-200"
    >
      Send
    </button>
  </div>
</div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-2 mb-2">
      <div
        class="bg-white border border-blue-100 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-blue-700">{{ totalTemplates }}</p>
        <h3 class="text-sm font-semibold text-blue-500">Total Templates</h3>
      </div>
      <div
        class="bg-white border border-warning-500 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-warning-600">{{ totalDocuments }}</p>
        <h3 class="text-sm font-semibold text-warning-600">Total Documents</h3>
      </div>
      <div
        class="bg-white border border-green-300 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-green-700">
          {{ successfulGenerations }}
        </p>
        <h3 class="text-sm font-semibold text-green-500">
          Successful Generations
        </h3>
      </div>
      <div
        class="bg-white border border-red-100 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-red-500">{{ failedGenerations }}</p>
        <h3 class="text-sm font-semibold text-red-300">Failed Generations</h3>
      </div>
      <div
        class="bg-white border border-gray-400 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-gray-600">
          {{ generationRate.toFixed(2) }}%
        </p>
        <h3 class="text-sm font-semibold text-gray-500">
          Generation Rate (daily)
        </h3>
      </div>
      <div
        class="bg-white border border-blue-100 p-1 rounded-lg shadow text-center"
      >
        <p class="text-2xl font-bold text-blue-300">
          {{ failureRate.toFixed(2) }}%
        </p>
        <h3 class="text-sm font-semibold text-blue-200">Failure Rate</h3>
      </div>
    </div>
    <div class="bg-white p-4 rounded-lg shadow w-full  max-h-[385px]">
      <line-chart :chartData="chartData" />
    </div>
  </div>
</template>
