<script setup lang="ts">
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
        <tr class="border-b hover:bg-gray-50">
          <td class="px-6 py-4">AKBA</td>
          <td class="px-6 py-4">$8.00</td>
          <td class="px-6 py-4">$6.00</td>
          <td class="px-6 py-4">Akebia Therapeutics</td>
          <td class="px-6 py-4">target lowered by</td>
          <td class="px-6 py-4">Goldman Sachs</td>
          <td class="px-6 py-4">Buy</td>
          <td class="px-6 py-4">Buy</td>
          <td class="px-6 py-4">2025-10-31</td>
        </tr>
        <tr class="border-b hover:bg-gray-50">
          <td class="px-6 py-4">TSLA</td>
          <td class="px-6 py-4">$250.00</td>
          <td class="px-6 py-4">$300.00</td>
          <td class="px-6 py-4">Tesla Inc</td>
          <td class="px-6 py-4">target raised by</td>
          <td class="px-6 py-4">Morgan Stanley</td>
          <td class="px-6 py-4">Hold</td>
          <td class="px-6 py-4">Buy</td>
          <td class="px-6 py-4">2025-11-10</td>
        </tr>
      </tbody>
    </table>
    </div>
  </main>
</template>
