<script lang="ts">
	import ListIcon from './icons/ListIcon.svelte';
	import MapIcon from './icons/MapIcon.svelte';
    import type { TabOptions } from '../types';
    interface Props {
        activeTab: TabOptions
        class?: string
        onchange?: (activeTab: TabOptions) => void;
    }

    let { activeTab = $bindable(), onchange, class: _class = "" }: Props = $props();

	const changeTab = (tab: TabOptions) => {
		activeTab = tab
		if (onchange) onchange(activeTab);
	};
</script>

<div class="view-select {_class}">
	<button class:active={activeTab === "map"} onclick={() => changeTab("map")}><MapIcon class="clickable" title="Show map with sessions" /></button>
	<button class:active={activeTab === "list"} onclick={() => changeTab("list")}><ListIcon class="clickable" title="Show list of sessions" /></button>
</div>

<style>
	.view-select button {
		background-color: white;
		padding: 1em 1.2em; /* Some padding */
		cursor: pointer; /* Pointer/hand icon */
		float: left; /* Float the buttons side by side */
        border: 2px solid grey;
	}

    @media (max-width: 480px) {
        .view-select button {
            padding: 0.6em 0.7em;
        }
    }

	.view-select button.active {
		box-shadow: inset dimgrey 0px 0px 10px -2px;
		background: lightgrey;
		border-color: dimgrey;
		color:black;
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
