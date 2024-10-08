<!-- <script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import axios from 'axios';
import api from '@/config/api';

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
    // console.log('Fetched document history:', documentHistory);

    const dataPoints = documentHistory.map(entry => ({
        label: entry.date,
        y: entry.count
    }));
    // console.log('Data points for chart:', dataPoints);

    options.value.data[0].dataPoints = dataPoints;
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
</template>
 -->






 <script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTemplateStore } from '@/domain/templates/stores';
import { useDocumentStore } from '@/domain/documents/stores';
import axios from 'axios';
import api from '@/config/api';

const chart = ref(null);
const templateStore = useTemplateStore();
const documentStore = useDocumentStore();

const options = ref({
    animationEnabled: true,
    exportEnabled: true,
    theme: "light",
    title: {
        text: "This Week's Report"
    },
    axisX: {
        title: "Days of the Week",
        labelTextAlign: "center",
        labelAngle: 0,
        interval: 1 
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

// async function fetchMetrics() {
//     await templateStore.fetchTemplates();
//     await documentStore.fetchDocuments();

//     const documentHistory = await fetchDocumentHistory();
    
//     if (documentHistory && Array.isArray(documentHistory)) {
//         // Map data to fit the chart's data format
//         options.value.data[0].dataPoints = documentHistory.map(entry => ({
//             label: entry.date,
//             y: entry.count
//         }));

//         // Log chart data for debugging
//         console.log('Chart Data:', options.value.data[0].dataPoints);
//     } else {
//         console.error('Document history is empty or not in expected format.');
//     }
// }


// Function to fetch metrics and reorganize data
async function fetchMetrics() {
    try {
        const response = await api.get('/document-history');
        const responseData = response.data;

        // Check for successful response
        if (responseData.code !== 200 || !Array.isArray(responseData.data)) {
            throw new Error('Invalid response format');
        }

        // Get the current day index (0=Sunday, 1=Monday, ..., 6=Saturday)
        const currentDate = new Date();
        const currentDayIndex = currentDate.getDay(); // 0-6 for Sun-Sat

        // Define the ordered array of days
        const allDays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

        // Reorder the days so the current day is last
        const orderedDays = [
            ...allDays.slice(currentDayIndex + 1), // Days after the current day
            ...allDays.slice(0, currentDayIndex + 1) // Days before and including the current day
        ];

        // Create a map for easy data access
        const dayDataMap = {};
        responseData.data.forEach(entry => {
            dayDataMap[entry.date] = entry.count;
        });

        // Create the dataPoints in the desired order
        const dataPoints = orderedDays.map(day => ({
            label: day,
            y: dayDataMap[day] || 0 // Default to 0 if no data available
        }));

        // Assign the organized dataPoints to the chart options
        options.value.data[0].dataPoints = dataPoints;

        // Log for debugging
        console.log('Data Points:', dataPoints);
    } catch (error) {
        console.error('Error fetching document history:', error);
    }
}




// async function fetchDocumentHistory() {
//     try {
//         const response = await api.get('/document-history');
//         if (response.status !== 200) {
//             throw new Error('Failed to fetch document history');
//         }
//         const responseData = response.data;
//         if (responseData.code !== 200) {
//             throw new Error('Failed to fetch document history');
//         }
//         const data: { date: string, count: number }[] = responseData.data;
//         return data.map((entry: { date: string, count: number }) => ({
//             date: entry.date,   
//             count: entry.count
//         }));
//     } catch (error) {
//         console.error('Error fetching document history:', error);
//         return [];
//     }
// }

async function fetchDocumentHistory() {
    try {
        const response = await api.get('/document-history');

        // Check for successful response
        if (response.status !== 200) {
            throw new Error('Failed to fetch document history');
        }

        const responseData = response.data;

        // Validate the response structure
        if (responseData.code !== 200 || !Array.isArray(responseData.data)) {
            throw new Error('Invalid response format');
        }

        // Return the data directly
        return responseData.data;
    } catch (error) {
        console.error('Error fetching document history:', error);
        return []; // Return an empty array in case of error
    }
}

</script>

<template>
    <CanvasJSChart :options="options" :style="styleOptions" @chart-ref="chartInstance"/>
</template>
