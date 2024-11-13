<script lang="ts">
	import Popup from './Popup.svelte';
	import type { SessionPropertiesWithVenue } from '../types';
	import TimeIcon from './icons/TimeIcon.svelte';
	import FileTrayIcon from './icons/FileTrayIcon.svelte';
	import { constructTimeString } from './timeUtils';

	interface Props {
		propertiesList: SessionPropertiesWithVenue[];
		onclick?: () => any;
		onclose: () => any;
	}

	let { propertiesList, onclick, onclose }: Props = $props();

	let id = 'popup-session-' + propertiesList.map((i) => i.session_id).join('-'); // : "venue-" + properties.venue_id}`;

	const onClick = (activeSessionId: number) => {
		window.sessionStorage.setItem('activeSessionId', activeSessionId.toString());
		window.location.assign('/' + activeSessionId);
	};
</script>

<Popup {id} title="Click to expand" {onclick} {onclose}>
	{#snippet heading()}
		{#if propertiesList.length > 1}
			<div class="venue-name-multi">
				<b>{propertiesList[0].venue_name}</b>
			</div>
		{:else}
			<span onclick={() => onClick(propertiesList[0].session_id!)}>
				<b>{propertiesList[0].session_name}</b> at <b>{propertiesList[0].venue_name}</b>
			</span>
		{/if}
	{/snippet}
	{#snippet content()}
		<span>
			<div style="display: flex; flex-direction: column;">
				{#each propertiesList as properties}
					<span style="margin-top: 0.4em;" onclick={() => onClick(properties.session_id!)}>
						{#if propertiesList.length > 1}
							<b>{properties.session_name}</b><br />
						{/if}
						<table>
							<tbody>
								<tr
									><td><TimeIcon title="Time of event" /></td>
									<td>{@html constructTimeString(properties)}</td>
								</tr>
								<tr
									><td><FileTrayIcon title="Genre" /></td>
									<td>{properties.genres.map((i) => i.replace('_', ' ')).join(', ')}</td>
								</tr>
							</tbody>
						</table>
						<div style="text-align: right; padding: 0.3em;">
							<i>View more ...</i>
						</div>
					</span>
				{/each}
			</div>
		</span>
	{/snippet}
</Popup>

<style>
	td {
		vertical-align: top;
	}

	td:first-child {
		padding-top: 0.2em;
	}

	.venue-name-multi {
		font-size: larger;
		margin-bottom: 0.2em;
	}
</style>
