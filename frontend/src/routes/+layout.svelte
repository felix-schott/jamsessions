<script lang="ts">
	import InfoPopup from '$lib/InfoPopup.svelte';
	import LoadingScreen from '$lib/LoadingScreen.svelte';
	import Header from '$lib/Header.svelte';
	import Map from '$lib/Map.svelte';
	import { addSessionPopupVisible } from '../stores';
	import AddSessionPopup from '$lib/AddSessionPopup.svelte';
	import ViewSelect from '$lib/ViewSelect.svelte';
	import ListView from '$lib/ListView.svelte';
	import type { TabOptions } from '../types';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	let activeTab: TabOptions = $state('map');
	let headerRelative = $derived(activeTab !== 'map');
</script>

<div id="app" class="flex-column">
	<LoadingScreen />
	<Header positionRelative={headerRelative} />
	<main>
		{@render children?.()}
		<InfoPopup />
		<ViewSelect class="view-select-btns" bind:activeTab />
		<div style="height: 100%; display: flex; flex-direction: row;">
			<Map background={activeTab !== 'map'} />
			<button
				class="add-session-btn"
				title="Add session to the database"
				onclick={() => {
					$addSessionPopupVisible = true;
				}}>Session missing?</button
			>
			{#if activeTab === 'list'}
				<ListView />
			{/if}
		</div>
		<AddSessionPopup />
	</main>
</div>

<style>
	#app {
		margin: 0;
		display: flex;
		height: 100vh;
		width: 100vw;
	}

	.flex-column {
		flex-direction: column;
	}

	@media (max-width: 480px) {
		.flex-column {
			flex-direction: column-reverse;
		}
	}

	main {
		flex-grow: 1;
		min-height: 0;
	}

	:global(.view-select-btns) {
		position: absolute;
		bottom: 2em;
		left: 2em;
		z-index: 100;
	}

	.add-session-btn {
		position: absolute;
		bottom: 3em;
		right: 2em;
		display: flex;
		align-items: center;
		background-color: white;
		border-radius: 24px;
		border: 2px solid grey;
	}

	@media (max-width: 480px) {
		.add-session-btn {
			bottom: 6em;
			font-size: smaller;
			right: 0.5em;
		}

		:global(.view-select-btns) {
			position: absolute;
			bottom: unset;
			top: 1em;
			left: 1em;
			font-size: smaller;
			z-index: 100;
		}
	}

	/* @media (max-width: 320px) {
		.add-session-btn {
			font-size: smaller;
		}
	} */
</style>
