<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import {
		loading,
		selectedSessions,
		filterMenuVisible,
		visibleLayer,
		infoVisible
	} from '../stores';
	import { getSessions, type SessionOptions } from '../api';
	import { MapLayer, type SessionFeatureCollection, Genre, Backline } from '../types';
	import SettingsIcon from '../lib/icons/SettingsIcon.svelte';
	import InfoIcon from '../lib/icons/InfoIcon.svelte';

	interface Props {
		positionRelative: boolean
	}

	let { positionRelative }: Props = $props();

	const today = new Date();
	const todayString = today.toISOString().slice(0, 10);
	let selectedDateStr: string = $state("");

	const onDateChange = async () => {
		$loading = true;
		window.sessionStorage.setItem('selectedDateStr', selectedDateStr);
		try {
			let params: SessionOptions = {
				date: new Date(selectedDateStr),
				backline: window.sessionStorage.getItem('selectedBackline')?.split(',') as Backline[],
				genre: window.sessionStorage.getItem('selectedGenre') as Genre
			}
			if (window.sessionStorage.getItem('selectedTimeRange') && window.sessionStorage.getItem('selectedTimeRange') !== '0') {
				let endDate = new Date(selectedDateStr);
				endDate!.setDate(endDate!.getDate() + parseInt(window.sessionStorage.getItem('selectedTimeRange')!))
				params["endDate"] = endDate
			}
			$selectedSessions = await getSessions(params);
		} catch (e) {
			alert('An error occured when waiting for data from the server: ' + (e as Error).message);
			throw e;
		}

		$visibleLayer = MapLayer.SESSIONS;
		$loading = false;
	};

	onMount(async () => {
		let storedDateStr = window.sessionStorage.getItem("selectedDateStr");
		if (storedDateStr === null) {
			selectedDateStr = todayString;
		} else {
			selectedDateStr = storedDateStr
		}
		onDateChange();
	});
</script>

<div class="top-bar vertically-centered" class:relative={positionRelative}>
	<span
		class="title clickable"
		title="About this website"
		onclick={() => {
			$infoVisible = true;
		}}
	>
		<b>Jam Sessions</b>
in <b>London</b><InfoIcon
			style="margin-left: 0.2em; cursor: pointer;"
			colour="whitesmoke"
			title="About this website"
		/></span
	>
	<div class="vertically-centered">
		<SettingsIcon
			title="Adjust filters"
			class="clickable"
			height="2em"
			width="2em"
			onclick={() => {
				$filterMenuVisible = true;
			}}
		/>
		<input type="date" min={todayString} bind:value={selectedDateStr} onchange={onDateChange} />
	</div>
</div>

<style>
	:global(.clickable) {
		cursor: pointer;
	}

	.top-bar {
		position: absolute;
		z-index: 500000;
		top: 0;
		left: 0;
		right: 0;
		/* height: 8%;
		max-height: 5em; */
		padding: 1em 1.5em;
		background: rgba(0, 0, 0, 0.4);
		padding: 1em;
		color: whitesmoke;
		justify-content: space-between;
		font-size: x-large;
	}

	.relative {
		position: relative;
	}

	@media (max-width: 480px) {
		input[type='date'] {
			padding: 0.1em 0.2em;
			font-weight: 400;
			margin-left: 0.5em;
			width: 8.5em;
		}

		.top-bar {
			font-size: 1em;
			max-height: 6rem;
			top: unset;
			bottom: 0;
			padding: 0.7em;
			justify-content: space-around;
		}
	}

	@media (max-width: 320px) {
		.top-bar {
			font-size: 0.9em;
		}

		input[type='date'] {
			padding: 0.1em 0.2em;
			font-weight: 400;
			margin-left: 0.5em;
			width: 7.5em;
		}
	}

	.title {
		white-space: pre-wrap;
	}

	.vertically-centered {
		display: flex;
		align-items: center;
	}

	input[type='date'] {
		background: rgba(0, 0, 0, 0.5);
		color: white;
		margin-left: 1em;
		border-radius: 8px;
		border: 1px solid transparent;
		padding: 0.6em 1.2em;
		font-size: 1em;
		font-weight: 500;
		font-family: inherit;
		cursor: pointer;
		transition: border-color 0.25s;
	}
</style>
