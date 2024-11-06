<script lang="ts">
	import type { SessionWithVenueFeature } from '../types';
	import { getSessions } from '../api';
	import { selectedSessions } from '../stores';

	type SessionsByDate = { [key: string]: SessionWithVenueFeature[] };
	// const getSessionsForNext7Days = async (today: Date): Promise<SessionsByDate> => {
	//     let sessionsByDate: SessionsByDate = {}
	//     for (let i of Array(7).keys()) {
	//         let d = today
	//         d.setDate(d.getDate() + i);
	//         let sessions = await getSessions({
	//             date: d
	//         })
	//         sessionsByDate[d.toDateString()] = sessions.features
	//     }
	//             // debug

	//     for (let i of Array(7).keys()) {
	//         let d = today
	//         d.setDate(d.getDate() + 7 + i);
	//         let sessions = await getSessions({
	//             date: d
	//         })
	//         sessionsByDate[d.toDateString()] = sessions.features
	//     }
	//     return sessionsByDate
	// };
</script>

<div class="session-list-wrapper">
	<div>
		<!-- {#await getSessionsForNext7Days(new Date()) then sessions}
		{#each Object.entries(sessions) as dateAndSessions}
            <h2>{dateAndSessions[0]}</h2>
            {#each dateAndSessions[1] as session}
            <b>{new Date(session.properties.start_time_utc).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</b>: {session.properties.session_name} at {session.properties.venue_name} 
            {/each}
        {/each}
	{/await} -->
		{#each $selectedSessions!.features as session}
			<h2>{new Date(window.sessionStorage.getItem('selectedDateStr')!).toLocaleDateString()}</h2>
			<b
				>{new Date(session.properties.start_time_utc).toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit'
				})}</b
			>: {session.properties.session_name} at {session.properties.venue_name}
		{/each}
	</div>
</div>

<style>
	.session-list-wrapper {
		background: white;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1em;
	}

	.session-list-wrapper > div {
		overflow: auto;
		height: 80%;
	}
</style>
