<script lang="ts">
	import SessionList from '$lib/SessionList.svelte';
	import { onMount } from 'svelte';
	import { type SessionOptions, getSessions } from '../api';
	import type { Backline, Genre } from '../types';
	import { loading, selectedDateStr, selectedSessions } from '../stores';

	onMount(async () => {
		$loading = true;
		// retrieve data from session storage
		let storedDateStr = window.sessionStorage.getItem('selectedDateStr');
		if (storedDateStr === null) {
			$selectedDateStr = new Date().toISOString().slice(0, 10); // today's date
		} else {
			$selectedDateStr = storedDateStr;
		}

		// load data from backend
		try {
			let params: SessionOptions = {
				date: new Date($selectedDateStr),
				backline: window.sessionStorage.getItem('selectedBackline')?.split(',') as Backline[],
				genre: window.sessionStorage.getItem('selectedGenre') as Genre
			};
			if (
				window.sessionStorage.getItem('selectedDateRange') &&
				window.sessionStorage.getItem('selectedDateRange') !== '0'
			) {
				let endDate = new Date($selectedDateStr);
				endDate!.setDate(
					endDate!.getDate() + parseInt(window.sessionStorage.getItem('selectedDateRange')!)
				);
				params['endDate'] = endDate;
			}
			$selectedSessions = await getSessions(params);
		} catch (e) {
			alert('An error occured when waiting for data from the server: ' + (e as Error).message);
			$loading = false;
			throw e;
		}
		$loading = false;
	});
</script>

<SessionList />
