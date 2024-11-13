<script lang="ts">
	import {
		loading,
		filterMenuVisible,
		selectedSessions,
		visibleLayer,
		selectedDateRange
	} from '../stores';
	import Modal from './Modal.svelte';
	import { Backline, Genre, type SessionFeatureCollection, MapLayer } from '../types';
	import { getSessions, type SessionOptions } from '../api';
	import MultiSelect from './MultiSelect.svelte';
	import MicrophoneIcon from './icons/MicrophoneIcon.svelte';
	import FileTrayIcon from './icons/FileTrayIcon.svelte';
	import ResetIcon from './icons/ResetIcon.svelte';
	import { onMount } from 'svelte';

	let selectedGenre: Genre = $state(Genre.ANY);
	let selectedBackline: Backline[] = $state([]);
	let selectedTimeRange: number = $state(0);

	let onChangeBackline = () => {
		selectedBackline = Array.from(document.querySelectorAll('#backline-select > option:checked'))
			.filter((i) => i)
			.map((i) => (i as HTMLOptionElement).value as Backline);
		if (selectedBackline)
			window.sessionStorage.setItem('selectedBackline', selectedBackline.join(','));
	};

	let onSubmit = async () => {
		$filterMenuVisible = false;

		// update session storage
		if (selectedTimeRange) {
			window.sessionStorage.setItem('selectedDateRange', selectedTimeRange.toString());
			$selectedDateRange = selectedTimeRange;
		}
		if (selectedGenre) window.sessionStorage.setItem('selectedGenre', selectedGenre);
		onChangeBackline();

		// request data
		$loading = true;
		try {
			let params: SessionOptions = {
				date: new Date(window.sessionStorage.getItem('selectedDateStr')!),
				backline: selectedBackline,
				genre: selectedGenre
			};
			if (selectedTimeRange != 0) {
				let endDate = new Date(params.date!);
				endDate!.setDate(endDate!.getDate() + selectedTimeRange);
				params['endDate'] = endDate;
			}
			$selectedSessions = await getSessions(params);
		} catch (e) {
			alert('An error occured when waiting for data from the server: ' + (e as Error).message);
			throw e;
		}
		$visibleLayer = MapLayer.SESSIONS;
		$loading = false;
	};

	let onReset = () => {
		selectedGenre = Genre.ANY;
		selectedBackline = [];
		selectedTimeRange = 0;
		$selectedDateRange = 0;
		window.sessionStorage.setItem('selectedBackline', '');
		window.sessionStorage.setItem('selectedGenre', Genre.ANY);
		window.sessionStorage.setItem('selectedDateRange', '0');
	};

	onMount(() => {
		let storedDateRange = window.sessionStorage.getItem('selectedDateRange');
		if (storedDateRange !== null) {
			$selectedDateRange = parseInt(storedDateRange);
			selectedTimeRange = $selectedDateRange;
		}

		let storedBackline = window.sessionStorage.getItem('selectedBackline');
		if (storedBackline !== null) {
			selectedBackline = storedBackline.split(',') as Backline[];
		}

		let storedGenre = window.sessionStorage.getItem('selectedGenre');
		if (storedGenre !== null) {
			selectedGenre = storedGenre as Genre;
		}
	});
</script>

<Modal
	isVisible={() => $filterMenuVisible}
	hide={() => {
		$filterMenuVisible = false;
	}}
>
	<h2 style="font-size: large;">Select date range</h2>
	<p style="font-size: larger;">
		Display sessions for the next<input
			type="number"
			title="Select number of days to display"
			bind:value={selectedTimeRange}
			style="margin-left: 0.5em;"
			min="0"
			max="30"
		/>
		days
	</p>
	<h2 style="font-size: large;">Filter by genre/backline</h2>
	<table>
		<tbody>
			<tr>
				<td><FileTrayIcon title="Select genre" class="icon-auto" /></td>
				<td
					><select title="Select genre" name="genre" bind:value={selectedGenre}>
						{#each Object.values(Genre) as opt}
							{#if window.sessionStorage.getItem('selectedGenre') === opt}
								<option value={opt} selected>{opt.replace('_', ' ')}</option>
							{:else}
								<option value={opt}>{opt.replace('_', ' ')}</option>
							{/if}
						{/each}
					</select></td
				>
			</tr>
			<tr>
				<td><MicrophoneIcon title="Select backline provided by venue" /></td>
				<td>
					<MultiSelect
						title="Select backline provided by genre"
						id="backline-select"
						name="backline"
					>
						{#each Object.values(Backline) as opt}
							{#if selectedBackline.includes(opt)}
								<option value={opt} selected>{opt.replace('_', ' ')}</option>
							{:else}
								<option value={opt}>{opt.replace('_', ' ')}</option>
							{/if}
						{/each}
					</MultiSelect>
				</td>
			</tr>
		</tbody>
	</table>
	<div
		style="display: flex; width: 100%; height: 100%; justify-content: space-between; align-items: center; margin-top: 1em;"
	>
		<ResetIcon title="Unset filters" onclick={onReset} style="cursor: pointer;" />
		<button title="Apply filters and load sessions" onclick={onSubmit}>Apply</button>
		<ResetIcon title="" style="visibility: hidden;" />
	</div>
</Modal>

<style>
	td {
		vertical-align: top;
		padding: 0.5em;
	}

	td:first-child {
		font-weight: bold;
	}
</style>
