<script lang="ts">
	import '../app.css';
	import InfoPopup from '$lib/InfoPopup.svelte';
	import LoadingScreen from '$lib/LoadingScreen.svelte';
	import Header from '$lib/Header.svelte';
	import Map from '$lib/Map.svelte';
	import type { SessionOptions } from '../api';
	import { loading, activeTab, editingSession } from '../stores';
	import AddSessionPopup from '$lib/AddSessionPopup.svelte';
	import ViewSelect from '$lib/ViewSelect.svelte';
	import SidePanel from '$lib/SidePanel.svelte';
	import FilterMenu from '$lib/FilterMenu.svelte';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	let headerRelative = $derived($activeTab !== 'map');

	let viewSelect: ViewSelect | undefined = $state();

	// session view: $selectedSessions must be active session, previous params must be stored in session storage
	// venue view: $selectedSessions must be session associated with venue, previous params must be stored in session storage
	// start view: $selectedSessions must be all session matching the filters

	// loaddata in page.ts/page.svelte??
</script>

<div id="app" class="flex-column">
	<LoadingScreen />
	<FilterMenu />
	<Header positionRelative={headerRelative} />
	<main>
		<InfoPopup />
		{#await new Promise((resolve) => setTimeout(resolve, 1)) then}
			<ViewSelect class="view-select-btns" bind:this={viewSelect} />
		{/await}
		<div class="content-wrapper">
			<Map
				background={$activeTab !== 'map'}
				onClickBackground={() => {
					if (window.matchMedia('(max-width: 480px)').matches) $activeTab = 'map';
				}}
			/>

			<SidePanel
				background={$activeTab === 'map'}
				hide={() => {
					if ($activeTab === 'session') {
						if ($editingSession) {
							$editingSession = false;
						} else window.location.assign('/');
					} else if ($activeTab === 'list') {
						$activeTab = 'map';
					} else if ($activeTab === 'venue') {
						window.location.assign('/');
					} else {
						throw 'not implemented for ' + $activeTab;
					}
				}}>{@render children?.()}</SidePanel
			>
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

	.content-wrapper {
		height: 100%;
		display: flex;
		flex-direction: row;
	}

	:global(.content-wrapper > div) {
		flex-grow: 1;
		flex-shrink: 1;
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

	@media (max-width: 480px) {
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
