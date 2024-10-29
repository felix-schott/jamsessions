<script lang="ts">
	import InfoPopup from '$lib/InfoPopup.svelte';
	import LoadingScreen from '$lib/LoadingScreen.svelte';
	import Header from '$lib/Header.svelte';
	import Map from '$lib/Map.svelte';
	import { addSessionPopupVisible } from '../stores';
	import AddSessionPopup from '$lib/AddSessionPopup.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<Header />
<main>
	{@render children?.()}
	<LoadingScreen />
	<InfoPopup />
	<div style="height: 100%;">
		<Map />
		<button
			class="add-session-btn"
			title="Add session to the database"
			onclick={() => {
				$addSessionPopupVisible = true;
			}}>Session missing?</button
		>
	</div>
	<AddSessionPopup />
</main>

<style>
	main {
		height: 100%;
		width: 100%;
		background-color: grey;
		/* display: flex;
      flex-direction: column; */
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
	}

	/* @media (max-width: 320px) {
		.add-session-btn {
			font-size: smaller;
		}
	} */
</style>
