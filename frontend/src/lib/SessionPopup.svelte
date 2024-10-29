<script lang="ts">
	import Popup from './Popup.svelte';
	import type { SessionPropertiesWithVenue } from '../types';
	import TimeIcon from './icons/TimeIcon.svelte';
	import FileTrayIcon from './icons/FileTrayIcon.svelte';

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

{#if propertiesList.length == 1}
	<Popup {id} title="Click to expand" blur={true} {onclick} {onclose}>
		{#snippet heading()}
			<span onclick={() => onClick(propertiesList[0].session_id!)}>
				<b>{propertiesList[0].session_name}</b> at <b>{propertiesList[0].venue_name}</b>
			</span>
		{/snippet}
		{#snippet content()}
			<span>
				<div
					style="display: flex; flex-direction: column;"
					onclick={() => onClick(propertiesList[0].session_id!)}
				>
					<div>
						<TimeIcon title="Time of event" /><span style="margin-left: 0.5em;">
							{new Date(propertiesList[0].start_time_utc).toLocaleTimeString()} -
							{new Date(
								new Date(propertiesList[0].start_time_utc).getTime() +
									propertiesList[0].duration_minutes * 60000
							).toLocaleTimeString()}
						</span>
					</div>
					<div>
						<FileTrayIcon title="Genre" /><span style="margin-left: 0.5em;">
							{propertiesList[0].genres.map((i) => i.replace('_', ' ')).join(', ')}
						</span>
					</div>
				</div>
			</span>
		{/snippet}
	</Popup>
{:else}
	<Popup {id} title="Click to expand" {onclick} {onclose}>
		{#snippet heading()}
			<div style="font-size: larger; margin-bottom: 0.2em;">
				<b>{propertiesList[0].venue_name}</b>
			</div>
		{/snippet}
		{#snippet content()}
			<span>
				<div style="display: flex; flex-direction: column;">
					{#each propertiesList as properties, idx}
						<span style="margin-top: 0.4em;" onclick={() => onClick(properties.session_id!)}>
							<div>
								<b>{properties.session_name}</b><br />
								<TimeIcon title="Time of event" /><span style="margin-left: 0.5em;">
									{new Date(properties.start_time_utc).toLocaleTimeString()} -
									{new Date(
										new Date(properties.start_time_utc).getTime() +
											properties.duration_minutes * 60000
									).toLocaleTimeString()}
								</span>
							</div>
							<div>
								<FileTrayIcon title="Genre" /><span style="margin-left: 0.5em;">
									{properties.genres.map((i) => i.replace('_', ' ')).join(', ')}
								</span>
							</div>
							<div style="text-align: right; padding: 0.3em;">
								<i>View more ...</i>
							</div>
							<!-- {#if idx !== propertiesList.length - 1}
								<hr />
							{/if} -->
						</span>
					{/each}
				</div>
			</span>
		{/snippet}
	</Popup>
{/if}
