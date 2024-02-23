<script setup>
import axios from 'axios'
import moment from 'moment'
import { computed, nextTick, onBeforeMount, onMounted, onUnmounted, reactive, ref, watch } from "vue";

let newEntries = reactive([]);

let hourrate = ref(15)

// placeholder
// let newEntries = reactive([{ "date": "2024-02-03", "start": "08:00", "end": "12:33", "totalHours": 4.55 }, { "date": "2024-02-04", "start": "08:00", "end": "12:04", "totalHours": 4.066666666666666 }, { "date": "2024-02-10", "start": "08:00", "end": "12:52", "totalHours": 4.866666666666666 }, { "date": "2024-02-11", "start": "08:00", "end": "12:47", "totalHours": 4.783333333333333 }, { "date": "2024-02-17", "start": "08:00", "end": "13:08", "totalHours": 5.133333333333334 }, { "date": "2024-02-18", "start": "08:00", "end": "13:45", "totalHours": 5.75 }, { "date": "2024-02-24", "start": "08:00", "end": "10:28", "totalHours": 2.466666666666667 }, { "date": "2024-02-25", "start": "08:00", "end": "11:44", "totalHours": 3.7333333333333334 }, { "date": "2024-03-02", "start": "08:00", "end": "12:08", "totalHours": 4.133333333333334 }, { "date": "2024-03-03", "start": "08:00", "end": "12:35", "totalHours": 4.583333333333333 }, { "date": "2024-03-09", "start": "08:00", "end": "12:08", "totalHours": 4.133333333333334 }, { "date": "2024-03-10", "start": "08:00", "end": "12:55", "totalHours": 4.916666666666667 }, { "date": "2024-03-16", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-17", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-18", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-19", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-20", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-21", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-22", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-23", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-24", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-25", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-26", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-27", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-28", "start": "08:00", "end": "12:00", "totalHours": 4 }, { "date": "2024-03-30", "start": "08:00", "end": "12:00", "totalHours": 4 }]);


onMounted(() => {
  if (getWithExpiry('workedhours')) {
    const storedEntries = JSON.parse(getWithExpiry('workedhours'));
    storedEntries.forEach(entry => newEntries.push(entry));
  } else {
    newEntries.push({
      'date': '',
      'start': '',
      'end': '',
      'totalHours': '',
      'earned': ''
    })
  }
})
  

watch(newEntries, (currentState, prevState) => {
  currentState.forEach(entry => {
    entry.totalHours = moment(entry.end, 'HH:mm').diff(moment(entry.start, 'HH:mm'), 'hours', true);
    entry.earned = Number((hourrate.value * entry.totalHours).toFixed(2))
  });
  setWithExpiry('workedhours', JSON.stringify(newEntries), 900000)
});

watch(hourrate, (currentState, prevState) => {
  newEntries.forEach(entry => {
    entry.totalHours = moment(entry.end, 'HH:mm').diff(moment(entry.start, 'HH:mm'), 'hours', true);
    entry.earned = Number((hourrate.value * entry.totalHours).toFixed(2))
  });
  setWithExpiry('workedhours', JSON.stringify(newEntries), 900000)
})

function setWithExpiry(key, value, ttl) {
	const now = new Date()

	// `item` is an object which contains the original value
	// as well as the time when it's supposed to expire
	const item = {
		value: value,
		expiry: now.getTime() + ttl,
	}
	localStorage.setItem(key, JSON.stringify(item))
}

function getWithExpiry(key) {
	const itemStr = localStorage.getItem(key)
	// if the item doesn't exist, return null
	if (!itemStr) {
		return null
	}
	const item = JSON.parse(itemStr)
	const now = new Date()
	// compare the expiry time of the item with the current time
	if (now.getTime() > item.expiry) {
		localStorage.removeItem(key)
		return null
	}
	return item.value
}


function generatePDF() {
  let entriesPrepared = reactive([])

  newEntries.forEach(entry => {
    entriesPrepared.push({
      date: entry.date,
      start: entry.start,
      end: entry.end,
      totalHours: entry.totalHours,
      earned: entry.earned
    })
  })

  axios.post('http://localhost:1997/api/generate', newEntries, {
    responseType: 'blob', // Important for handling binary data like PDF
  })
    .then((response) => {
      // Create a new Blob object using the response data
      const file = new Blob([response.data], { type: 'application/pdf' });

      // Generate a URL for the Blob object
      const fileURL = URL.createObjectURL(file);

      // Create an anchor (<a>) element to facilitate downloading
      const downloadLink = document.createElement('a');
      downloadLink.href = fileURL;
      downloadLink.setAttribute('download', 'uren_registratie.pdf'); // Set the file name for the download
      document.body.appendChild(downloadLink);
      downloadLink.click();

      // Clean up by revoking the Blob URL and removing the temporary link
      URL.revokeObjectURL(fileURL);
      document.body.removeChild(downloadLink);
    })
    .catch((error) => {
      console.error('Error fetching the PDF:', error);
    });
}

const addRow = (index) => {
  if (index === Object.keys(newEntries).length - 1) {
    newEntries.push({
      date: '',
      start: '',
      end: '',
      totalHours: '',
      earned: ''
    });
  }
}

const removeRow = (index) => {
  newEntries.splice(index, 1);
}

</script>

<template>
  <div class="divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow">
    <div class="px-4 py-5 sm:px-6">
      <h1 class="text-lg leading-6 font-medium text-gray-900">Gewerkte uren</h1>
    </div>
    <div class="px-4 py-5 sm:p-6">
      <div class="flex items-center space-x-4 mb-4">
        <div class="w-1/4">Datum</div>
        <div class="w-1/6">Start tijd</div>
        <div class="w-1/6">Eind tijd</div>
        <div class="w-1/6">Verdienst</div>
      </div>
      <div v-for="(input, index) in newEntries" :key="index" class="flex items-center space-x-4 mb-4">
        <input v-model="input.date" type="date" name="date" id="date"
          class="w-1/4 rounded-md border-gray-300 py-2 px-3 text-gray-900 shadow-sm focus:ring-indigo-500 focus:border-indigo-500" />
        <input v-model="input.start" type="time" name="start" id="start"
          class="w-1/6 rounded-md border-gray-300 py-2 px-3 text-gray-900 shadow-sm focus:ring-indigo-500 focus:border-indigo-500" />
        <input v-model="input.end" type="time" name="end" id="end"
          class="w-1/6 rounded-md border-gray-300 py-2 px-3 text-gray-900 shadow-sm focus:ring-indigo-500 focus:border-indigo-500" />
        <!-- <input :value="input.cost" type="number" name="cost" id="cost" readonly
          class="w-1/6 rounded-md border-gray-300 py-2 px-3 text-gray-900 shadow-sm" /> -->
          <div class="flex items-center w-1/6 rounded-md border-gray-300 shadow-sm cursor-not-allowed bg-gray-100">
            <span class="pl-2 text-gray-900 cursor-not-allowed bg-gray-100">€</span>
            <input :value="input.earned" type="text" name="earned" id="earned" readonly
                  class="w-full rounded-md py-2 px-3 text-gray-900 bg-gray-100 shadow-sm outline-none border-none cursor-not-allowed"
                  style="padding-left: 10px;" />
          </div>
        <span v-if="newEntries.length >= 2" @click="removeRow(newEntries.indexOf(input))" class="cursor-pointer">
          <svg class="h-6 w-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
        </span>
        <div v-if="index === newEntries.length - 1" @click="addRow(index)" class="cursor-pointer text-blue-500">
          Rij toevoegen
        </div>
      </div>
      <div class="flex items-center space-x-4 mb-4">
        <div @click="generatePDF()" class="cursor-pointer text-white bg-blue-500 hover:bg-blue-700 rounded-md px-4 py-2 inline-block">
          PDF maken
        </div>
        <div class="flex items-center w-1/6 rounded-md border-gray-300 shadow-sm">
            <span class="pl-2 text-gray-900">€</span>
            <input v-model="hourrate" type="number" name="hourrate" id="hourrate" class="w-full roounded-md py-2 px-3 text-gray-900 shadow-sm focus:ring-indigo-500 focus:border-indigo-500" style="padding-left: 10px;" />
          </div>
      </div>
    </div>
  </div>
</template>
