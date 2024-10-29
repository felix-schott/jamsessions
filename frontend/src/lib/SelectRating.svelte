<script lang="ts">
	import StarEmptyIcon from './icons/StarEmptyIcon.svelte';
	import StarFilledIcon from './icons/StarFilledIcon.svelte';

	interface Props {
		style?: string;
		size?: string;
		onchange?: (rating: number) => void;
	}

	let { style = '', size = '1em', onchange }: Props = $props();

	let currentRating = $state(0);
	let currentRatingHover = $state(0);

	const onClick = (idx: number) => {
		currentRating = idx + 1;
        if (onchange) onchange(currentRating);
	};
	const onHover = (idx: number) => {
        currentRatingHover = idx + 1;
    };
</script>

<!-- set currentRatingHover back to 0 when the mouse leaves the div -->
<div class="rating-select" {style} onmouseleave={() => {currentRatingHover = 0}}> 
	{#each Array.from(Array(5).keys()) as idx}
		{#if currentRating >= idx + 1 || currentRatingHover >= idx + 1}
			<StarFilledIcon
				title="Rate {idx + 1} stars"
				height={size}
				width={size}
				onmouseover={() => onHover(idx)}
				onclick={() => onClick(idx)}
			/>
		{:else}
			<StarEmptyIcon
				title="Rate {idx + 1} stars"
				height={size}
				width={size}
				onmouseover={() => onHover(idx)}
				onclick={() => onClick(idx)}
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
