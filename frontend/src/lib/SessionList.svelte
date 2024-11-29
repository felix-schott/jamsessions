<script lang="ts">
	import { Interval, type SessionWithVenueFeature } from '../types';
	import { selectedSessions } from '../stores';
	import { untrack } from 'svelte';
	import { sanitisePathElement } from './uriUtils';
	import InfoIcon from './icons/InfoIcon.svelte';

	type SessionsByDate = { [key: string]: SessionWithVenueFeature[] };
	let sessionsByDate: SessionsByDate = $state({});

	$effect(() => {
		if ($selectedSessions && $selectedSessions.features) {
			sessionsByDate = {};
			untrack(() => {
				for (let s of $selectedSessions?.features!) {
					for (let d of s.properties.dates!) {
						if (!(d in sessionsByDate)) {
							sessionsByDate[d] = [];
						}
						sessionsByDate[d].push(s);
					}
				}
			});
		}
	});
</script>

{#each Object.entries(sessionsByDate) as [d, sessions]}
	<h2>{new Date(d).toLocaleDateString()}</h2>
	{#each sessions as session}
		<p
			class="clickable"
			onclick={() => {
				window.sessionStorage.setItem('activeSessionId', session.properties.session_id!.toString());
				if (window.matchMedia('(max-width: 480px)').matches) {
					window.location.assign(
						`/${sanitisePathElement(session.properties.venue_name)}-${session.properties.venue}/${sanitisePathElement(session.properties.session_name)}-${session.properties.session_id}`
					);
				} else {
					document.getElementById('map')?.dispatchEvent(new Event('zoom'));
				}
			}}
		>
			<b
				>{new Date(session.properties.start_time_utc).toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit'
				})}</b
			>: {session.properties.session_name} at {session.properties.venue_name}
			{#if session.properties.interval === Interval.IRREGULARWEEKLY}
				<span
					class="irregular-tag"
					title="Irregular session - may or may not take place on this day."
					>irregular<InfoIcon
						style="margin-left: 0.2em;"
						colour="white"
						title="Irregular session - may or may not take place on this day."
					/></span
				>
			{/if}
		</p>
	{/each}
{/each}

<style>
	.irregular-tag {
		margin-left: 0.3em;
		padding: 0.1em 0.6em;
		background: var(--accent-color);
		color: white;
		font-size: smaller;
		border: none;
		border-radius: 24px;
		display: inline-flex;
		justify-content: center;
		align-items: center;
		vertical-align: text-top;
	}
</style>
