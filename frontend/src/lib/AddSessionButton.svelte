<script lang="ts">
	import { addSessionPopupVisible } from '../stores';
</script>

<button
	class="add-session-btn"
	title="Add session to the database"
	onclick={() => {
		$addSessionPopupVisible = true;
	}}>Session missing?</button
>

<style>
	.add-session-btn {
		position: absolute;
		bottom: 3em;
		right: 2em;
		display: flex;
		align-items: center;
		background-color: white;
		border-radius: 24px;
		border: 2px solid grey;
		z-index: 500001;
	}

	@media (max-width: 480px) {
		.add-session-btn {
			bottom: 6em;
			font-size: smaller;
			right: 0.5em;
		}
	}
</style>
