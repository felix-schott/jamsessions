<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import StarEmptyIcon from './icons/StarEmptyIcon.svelte';
	import StarFilledIcon from './icons/StarFilledIcon.svelte';

	const dispatch = createEventDispatcher();
	export let style: string = '';
	export let size: string = '1em';

	let currentRating = 0;
	let currentRatingHover = 0;

	const onClick = (idx: number) => {
		currentRating = idx + 1;
        dispatch('change', { rating: currentRating });
	};
	const onHover = (idx: number) => {
        currentRatingHover = idx + 1;
    };
</script>

<!-- set currentRatingHover back to 0 when the mouse leaves the div -->
<div class="rating-select" {style} on:mouseleave={() => {currentRatingHover = 0}}> 
	{#each Array.from(Array(5).keys()) as idx}
		{#if currentRating >= idx + 1 || currentRatingHover >= idx + 1}
			<StarFilledIcon
				title="Rate {idx + 1} stars"
				height={size}
				width={size}
				on:mouseover={() => onHover(idx)}
				on:click={() => onClick(idx)}
			/>
		{:else}
			<StarEmptyIcon
				title="Rate {idx + 1} stars"
				height={size}
				width={size}
				on:mouseover={() => onHover(idx)}
				on:click={() => onClick(idx)}
			/>
		{/if}
	{/each}
</div>

<style>
	.rating-select {
		display: flex;
		cursor: pointer;
	}
</style>
