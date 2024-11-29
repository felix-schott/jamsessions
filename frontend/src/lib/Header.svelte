<script lang="ts">
	import {
		selectedDateStr,
		filterMenuVisible,
		infoVisible,
		selectedDateRange,
		selectedSessions,
		loading
	} from '../stores';
	import { type SessionOptions, getSessions } from '../api';
	import type { Backline, Genre } from '../types';
	import SettingsIcon from '../lib/icons/SettingsIcon.svelte';
	import InfoIcon from '../lib/icons/InfoIcon.svelte';

	interface Props {
		positionRelative: boolean;
	}

	let { positionRelative }: Props = $props();

	const onDateChange = async () => {
		$loading = true;

		window.sessionStorage.setItem('selectedDateStr', $selectedDateStr);

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
	};
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
		<input
			type="date"
			min={new Date().toISOString().slice(0, 10)}
			bind:value={$selectedDateStr}
			onchange={onDateChange}
		/>
		{#if $selectedDateRange > 0}
			<span style="margin-left: 0.5em;">+{$selectedDateRange}</span>
		{/if}
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
		background-color: rgba(0, 0, 0, 0.4);
		transition: background-color ease-in-out 1s;
		padding: 1em;
		color: whitesmoke;
		justify-content: space-between;
		font-size: x-large;
	}

	.relative {
		position: relative;
		background-color: rgba(0, 0, 0, 0.6);
	}

	@media (max-width: 480px) {
		input[type='date'] {
			padding: 0.1em 0.2em;
			font-weight: 400;
			margin-left: 0.5em;
			width: 7em;
		}

		.top-bar {
			font-size: 1em;
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
