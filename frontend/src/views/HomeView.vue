<script setup lang="ts">
import type { Stock } from '@/domain/stock';
import { ref, onMounted } from 'vue';
import ModalComponent from '../components/modal.vue'
import { formatColombianDateTime, formatAsMoney, compareDecimals, getDelta, getSortChar } from '@/utils/format-utils';

const apiUrl = import.meta.env.VITE_API_URL;
const stocks = ref<Stock[]>([]);
const page = ref<number>(0);
const sortby = ref<string>("");
const sortorder = ref<boolean>(true);
const query = ref<string>("");
const canContinue = ref<boolean>(false);
const isModalOpen = ref<boolean>(false);
const recommendedStock = ref<Stock | null>(null);
const sortBtnClass = "ml-2 px-2 py-1 text-xs bg-white/20 hover:bg-white/30 rounded transition-colors";
const thClass = "px-6 py-4 text-left font-semibold whitespace-nowrap";
const topBtnClass = "px-4 py-2 bg-[#3B1CEA] text-white font-semibold rounded-lg hover:bg-[#2D15B8] transition-colors shadow-md disabled:bg-gray-400 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:bg-gray-400";
const columnTitles = [
  { display: 'Ticker', label: 'Ticker', sortable: true },
  { display: 'Delta', label: 'Delta', sortable: false },
  { display: 'Target From', label: 'TargetFrom', sortable: true },
  { display: 'Target To', label: 'TargetTo', sortable: true },
  { display: 'Company', label: 'Company', sortable: true },
  { display: 'Action', label: 'Action', sortable: true },
  { display: 'Brokerage', label: 'Brokerage', sortable: true },
  { display: 'Rating From', label: 'RatingFrom', sortable: true },
  { display: 'Rating To', label: 'RatingTo', sortable: true },
  { display: 'Time', label: 'Time', sortable: true }
] as const;

const onSearch = async () => {
  const queryVal = query.value || ""
  console.debug(`Searching for ${queryVal}`);
  page.value = 0;
  updateStocks(page.value, sortby.value, sortorder.value, queryVal)
}

const onClickSort = async (label: string) => {
  const column = columnTitles.find(col => col.label === label);
  if (!column) return;
  console.debug(`Toggle sort by ${column.label}`);
  stocks.value = [];
  page.value = 0;
  if (sortby.value === label) {
    sortorder.value = !sortorder.value
  } else {
    sortorder.value = true
  }
  sortby.value = label;
  await updateStocks(0, label, sortorder.value, query.value)
}

const onClickPrev = async () => {
  if (page.value > 0) {
    page.value = page.value - 1;
    updateStocks(page.value, sortby.value, sortorder.value, query.value);
  }
}

const onClickNext = async () => {
  await updateStocks(page.value + 1, sortby.value, sortorder.value, query.value);
  if (stocks.value.length > 0) {
    page.value = page.value + 1;
  }
}

async function updateStocks(offset: number = 0, sortBy: string | null = null, sortOrder: boolean = true, queryStr: string = "") {
  async function requestStocks(offset: number = 0): Promise<Stock[]> {
    const params = new URLSearchParams({
      offset: offset.toString(),
      sortby: sortBy || "",
      asc: sortOrder ? 'true' : 'false',
      query: queryStr
    });
    const response = await fetch(`${apiUrl}/stocks?${params}`);
    const data = await response.json();
    console.debug(data);
    return data;
  }
  try {
    const data = await Promise.all([requestStocks(offset), requestStocks(offset + 1)]);
    const [current, next] = data;
    canContinue.value = next.length > 0;
    stocks.value = current;
  } catch (error) {
    console.error('Error fetching stocks: ', error);
  }
}

async function getRecommendation() {
  console.debug("Getting recommendation data..");
  const response = await fetch(`${apiUrl}/recommendation`);
  const data = await response.json() as Stock;
  if (data) {
    recommendedStock.value = data;
  }
  console.debug("Obtained recommendation: ", data);
  isModalOpen.value = true
}

onMounted(async () => {
  await updateStocks();
}
)
</script>

<template>
  <main class="min-h-[calc(100vh-5rem)] bg-gradient-to-br from-white to-[#81A5F7] pt-12 px-8 pb-12">
    <div class="w-full max-w-7xl mx-auto">
      <div class="flex items-center justify-between mb-4">
        <div class="flex gap-2 items-center">
          <button @click="onClickPrev" :class="topBtnClass" :disabled="page === 0">
            Prev</button>
          <button @click="onClickNext" :class="topBtnClass" :disabled="!canContinue">Next
          </button>
          <span>Page {{ page }}</span>
        </div>
        <div class="flex gap-2 items-center">
          <button @click="getRecommendation" :class="topBtnClass"> ★ Get recommendation ★</button>
          <input v-model="query"
            class="px-4 py-2 border-2 border-[#3B1CEA] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3B1CEA] shadow-md"
            placeholder="Search by ticker...">
          <button @click="onSearch" :class="topBtnClass">Search</button>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full bg-white shadow-lg rounded-lg overflow-hidden mb-8">
          <thead class="bg-[#3B1CEA] text-white">
            <tr>
              <th v-for="column in columnTitles" :key="column.label" :class="thClass">
                <span>{{ column.display }}</span><button v-if="column.sortable" @click="() => onClickSort(column.label)"
                  :class="sortBtnClass">{{ getSortChar(column.label, sortby, sortorder) }}</button>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="stock in stocks" :key="stock.Id" class="border-b hover:bg-gray-50">
              <td class="px-6 py-4">{{ stock.Ticker }}</td>
              <td class="px-6 py-4">
                <span v-if="compareDecimals(stock.TargetFrom, stock.TargetTo) === 1"
                  class="text-green-500 whitespace-nowrap">▲ {{ getDelta(stock.TargetFrom, stock.TargetTo) }}</span>
                <span v-else-if="compareDecimals(stock.TargetFrom, stock.TargetTo) === -1"
                  class="text-red-500 whitespace-nowrap">▼ {{ getDelta(stock.TargetFrom, stock.TargetTo) }}</span>
                <span v-else class="text-blue-500 whitespace-nowrap">━ {{ getDelta(stock.TargetFrom, stock.TargetTo)
                  }}</span>
              </td>
              <td class="px-6 py-4">{{ formatAsMoney(stock.TargetFrom) }}</td>
              <td class="px-6 py-4">{{ formatAsMoney(stock.TargetTo) }}</td>
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
  <modal-component v-model="isModalOpen" v-bind:stock="recommendedStock!" />
</template>
