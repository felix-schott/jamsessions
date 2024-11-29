<script lang="ts">
	import ListIcon from './icons/ListIcon.svelte';
	import MapIcon from './icons/MapIcon.svelte';
	import type { TabOptions } from '../types';
	import { activeTab } from '../stores';
	import { onMount } from 'svelte';
	interface Props {
		class?: string;
	}

	let { class: _class = '' }: Props = $props();

	let inSessionView = $state(false); // whether or not we're in the session view (SessionDetails)
	onMount(() => {
		inSessionView =
			window.location.pathname.split('/').filter((i) => i).length === 2 &&
			window.location.pathname.endsWith(window.sessionStorage.getItem('activeSessionId')!);
	});
</script>

<div class="{_class} view-select">
	<button
		class:active={$activeTab === 'map'}
		onclick={() => {
			if (!window.matchMedia('(max-width: 480px)').matches && inSessionView)
				window.location.assign('/');
			$activeTab = 'map';
		}}><MapIcon class="clickable" title="Show map with sessions" /></button
	>
	<button
		class:active={$activeTab === 'list' || $activeTab === 'session'}
		onclick={() => {
			if (!window.matchMedia('(max-width: 480px)').matches && inSessionView)
				window.location.assign('/');

			if (inSessionView) {
				$activeTab = 'session';
			} else {
				$activeTab = 'list';
			}
		}}><ListIcon class="clickable" title="Show list of sessions" /></button
	>
</div>

<style>
	.hidden {
		display: none;
	}

	.view-select button {
		background-color: white;
		padding: 1em 1.2em; /* Some padding */
		cursor: pointer; /* Pointer/hand icon */
		float: left; /* Float the buttons side by side */
		border: 2px solid grey;
	}

	.view-select button.active {
		box-shadow: inset dimgrey 0px 0px 10px -2px;
		background: lightgrey;
		border-color: dimgrey;
		color: black;
	}

	.view-select button:not(:last-child) {
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
		border-right: none; /* Prevent double borders */
	}

	.view-select button:last-child {
		border-top-left-radius: 0;
		border-bottom-left-radius: 0;
	}

	@media (max-width: 480px) {
		.view-select {
			display: flex;
			flex-direction: column;
		}

		.view-select button {
			padding: 0.6em 0.7em;
		}

		.view-select button:not(:last-child) {
			border-top-right-radius: 8px;
			border-bottom-right-radius: 0;
			border-bottom-left-radius: 0;
			border-top-left-radius: 8px;
			border-right: 2px solid grey;
			border-bottom: none; /* Prevent double borders */
		}

		.view-select button:last-child {
			border-top-left-radius: 0;
			border-bottom-left-radius: 8px;
			border-bottom-right-radius: 8px;
			border-top-right-radius: 0;
		}
	}

	/* Clear floats (clearfix hack) */
	.view-select:after {
		content: '';
		clear: both;
		display: table;
	}

	/* Add a background color on hover */
	.view-select button:hover {
		background-color: grey;
	}
</style>
