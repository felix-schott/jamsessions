<script lang="ts">
	import type { SessionWithVenueFeature } from '../types';
	import { selectedSessions } from '../stores';
	import { untrack } from 'svelte';

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
					window.location.assign('/' + session.properties.session_id);
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
		</p>
	{/each}
{/each}
