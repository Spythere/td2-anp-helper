<template>
  <div class="h-full text-center">
    <!-- <input type="file" accept=".anp" multiple="false" ref="fileUpload" @change="onFileUpload" /> -->

    <div class="flex gap-2 text-center justify-center">
      <button class="bg-slate-600 hover:bg-slate-500 p-1" @click="handleFindUpload">Wybierz plik ANP</button>
    </div>
    <div class="font-bold mt-3">{{ path }}</div>

    <div v-for="activeRun in sortedRuns" :key="activeRun.Content">
      {{ activeRun.Content }} (od {{ new Date(activeRun.TimestampFrom).toLocaleTimeString('pl-PL') }})
    </div>

    <div>{{ activeRuns.length }}</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { GetActiveRuns, GetFilePath } from '../wailsjs/go/main/App';
import { ActiveRun } from './types/dataTypes';

export default defineComponent({
  data() {
    return {
      path: '',
      activeRuns: [] as ActiveRun[],
    };
  },

  mounted() {
    setInterval(() => {
      this.getActiveRuns();
    }, 5000);
  },

  computed: {
    sortedRuns() {
      return this.activeRuns.sort((a, b) => {
        return a.TimestampFrom - b.TimestampFrom;
      });
    },
  },

  methods: {
    async handleFindUpload() {
      await GetFilePath().then((result) => {
        this.path = result;
        this.getActiveRuns();
      });
    },

    async getActiveRuns() {
      if (this.path == '') return;

      const dataStr = await GetActiveRuns();
      this.activeRuns = JSON.parse(dataStr);
    },
  },
});
</script>

<style scoped></style>
