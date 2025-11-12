<script setup lang="ts">
import type { Stock } from '@/domain/stock';
import {ref, onMounted} from 'vue';
const apiUrl = import.meta.env.VITE_API_URL;
const stocks = ref<Stock[]>([])

const formatColombianDateTime = (isoString: string): string => {
  const date = new Date(isoString);
  return date.toLocaleString('es-CO', {
    timeZone: 'America/Bogota',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

onMounted(async () => {
  try {
    const response = await fetch(apiUrl);
    const data = await response.json();
    console.debug(data);
    stocks.value = data;
  } catch (error) {
    console.error('Error fetching stocks: ', error);
  }
})
</script>

<template>
  <main class="min-h-[calc(100vh-5rem)] bg-gradient-to-br from-white to-[#81A5F7] flex justify-center pt-12 px-8">
    <div class="w-full max-w-7xl">
      <div class="flex justify-end mb-4">
        <input class="px-4 py-2 border-2 border-[#3B1CEA] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3B1CEA] shadow-md" placeholder="Search...">
      </div>
      <table class="w-full bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-[#3B1CEA] text-white">
        <tr>
          <th class="px-6 py-4 text-left font-semibold">Ticker</th>
          <th class="px-6 py-4 text-left font-semibold">Target From</th>
          <th class="px-6 py-4 text-left font-semibold">Target To</th>
          <th class="px-6 py-4 text-left font-semibold">Company</th>
          <th class="px-6 py-4 text-left font-semibold">Action</th>
          <th class="px-6 py-4 text-left font-semibold">Brokerage</th>
          <th class="px-6 py-4 text-left font-semibold">Rating From</th>
          <th class="px-6 py-4 text-left font-semibold">Rating To</th>
          <th class="px-6 py-4 text-left font-semibold">Time</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="stock in stocks" :key="stock.Id" class="border-b hover:bg-gray-50">
          <td class="px-6 py-4">{{ stock.Ticker }}</td>
          <td class="px-6 py-4">{{ stock.TargetFrom }}</td>
          <td class="px-6 py-4">{{ stock.TargetTo }}</td>
          <td class="px-6 py-4">{{ stock.Company }}</td>
          <td class="px-6 py-4">{{ stock.Action }}</td>
          <td class="px-6 py-4">{{ stock.Brokerage }}</td>
          <td class="px-6 py-4">{{ stock.RatingFrom }}</td>
          <td class="px-6 py-4">{{ stock.RatingTo }}</td>
          <td class="px-6 py-4">{{ formatColombianDateTime(stock.Time) }}</td>
        </tr>
      </tbody>
    </table>
    </div>
  </main>
</template>
