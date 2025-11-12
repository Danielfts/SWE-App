<script setup lang="ts">
import type { Stock } from '@/domain/stock';
import { ref, onMounted } from 'vue';
const apiUrl = import.meta.env.VITE_API_URL;
const stocks = ref<Stock[]>([]);
const page = ref<number>(0);
const sortBtnClass = "ml-2 px-2 py-1 text-xs bg-white/20 hover:bg-white/30 rounded transition-colors";
const thClass = "px-6 py-4 text-left font-semibold";

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

const onClickPrev = async () => {
  if (page.value > 0) {
    page.value = page.value - 1;
    queryStocks(page.value);
  }
}

const onClickNext = async () => {
  page.value = page.value + 1;
  queryStocks(page.value);
}

async function queryStocks(offset: number = 0) {
  try {
    const params = new URLSearchParams({
      offset : offset.toString()
    });
    const response = await fetch(`${apiUrl}?${params}`);
    const data = await response.json();
    console.debug(data);
    stocks.value = data;
  } catch (error) {
    console.error('Error fetching stocks: ', error);
  }
}

onMounted(async () => {
  await queryStocks();
}
)
</script>

<template>
  <main class="min-h-[calc(100vh-5rem)] bg-gradient-to-br from-white to-[#81A5F7] pt-12 px-8 pb-12">
    <div class="w-full max-w-7xl mx-auto">
      <div class="flex items-center justify-between mb-4">
        <div class="flex gap-2">
          <button @click="onClickPrev"
            class="px-4 py-2 bg-[#3B1CEA] text-white font-semibold rounded-lg hover:bg-[#2D15B8] transition-colors shadow-md">
            < Prev</button>
              <button @click="onClickNext"
                class="px-4 py-2 bg-[#3B1CEA] text-white font-semibold rounded-lg hover:bg-[#2D15B8] transition-colors shadow-md">Next
                > </button>
        </div>
        <input
          class="px-4 py-2 border-2 border-[#3B1CEA] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3B1CEA] shadow-md"
          placeholder="Search...">
      </div>
      <div class="overflow-x-auto">
        <table class="w-full bg-white shadow-lg rounded-lg overflow-hidden mb-8">
        <thead class="bg-[#3B1CEA] text-white">
          <tr>
            <th :class="thClass">Ticker<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Target From<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Target To<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Company<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Action<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Brokerage<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Rating From<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Rating To<button :class="sortBtnClass">-</button> </th>
            <th :class="thClass">Time<button :class="sortBtnClass">-</button> </th>
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
    </div>
  </main>
</template>
