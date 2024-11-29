<script lang="ts">
	import { getSessionsByVenueId } from '../../api';
	import type { VenueSlugData } from './+page.js';
	import type { SessionWithVenueFeatureCollection, SessionPropertiesWithVenue } from '../../types';
	import TimeIcon from '$lib/icons/TimeIcon.svelte';
	import FileTrayIcon from '$lib/icons/FileTrayIcon.svelte';
	import { sanitisePathElement, extractDomain } from '$lib/uriUtils';
	import { constructTimeString } from '$lib/timeUtils';
	import { selectedSessions } from '../../stores';

	interface Props {
		data: VenueSlugData;
	}

	let { data }: Props = $props();

	const onClick = (properties: SessionPropertiesWithVenue) => {
		window.sessionStorage.setItem('activeSessionId', properties.session_id!.toString());
		window.location.assign(
			`/${sanitisePathElement(properties.venue_name)}-${properties.venue}/${sanitisePathElement(properties.session_name)}-${properties.session_id}`
		);
	};

	const getSessionsByVenueIdWithErrorHandling =
		async (): Promise<SessionWithVenueFeatureCollection> => {
			try {
				let response = await getSessionsByVenueId(data.venueId);
				$selectedSessions = response;
				console.log('length', $selectedSessions.features.length);
				return response;
			} catch (e) {
				alert('An error occured when waiting for data from the server: ' + (e as Error).message);
				throw e;
			}
		};
</script>

{#await getSessionsByVenueIdWithErrorHandling() then fc}
	<div class="venue-overview">
		<div class="venue-info">
			<h2>{fc.features[0].properties.venue_name}</h2>
			<div>
				<a href={fc.features[0].properties.venue_website} target="_blank"
					>{extractDomain(fc.features[0].properties.venue_website)}</a
				><br />
				{fc.features[0].properties.address_first_line}<br />
				{#if fc.features[0].properties.address_second_line}
					{fc.features[0].properties.address_second_line}<br />
				{/if}
				{fc.features[0].properties.city}
				{fc.features[0].properties.postcode}<br />
				<a
					target="_blank"
					href="https://www.google.com/maps/place/{fc.features[0].properties.address_first_line.replaceAll(
						' ',
						'+'
					)},+{fc.features[0].properties.city.replaceAll(
						' ',
						'+'
					)}+{fc.features[0].properties.postcode.replaceAll(' ', '+')}/">View on Google Maps</a
				>
			</div>
		</div>
		<div class="venue-jamsessions" style="margin-top: 0.4em;">
			<h3>Jam sessions</h3>
			{#each fc.features as f}
				<div class="venue-session" onclick={() => onClick(f.properties)}>
					<b>{f.properties.session_name}</b><br />
					<table>
						<tbody>
							<tr
								><td><TimeIcon title="Time of event" /></td>
								<td>{@html constructTimeString(f.properties)}</td>
							</tr>
							{#if f.properties.genres && f.properties.genres.length !== 0}
								<tr
									><td><FileTrayIcon title="Genre" /></td>
									<td>{f.properties.genres.map((i) => i.replace('_', ' ')).join(', ')}</td>
								</tr>
							{/if}
						</tbody>
					</table>
				</div>
			{/each}
		</div>
	</div>
{/await}

<style>
	.venue-overview {
		display: flex;
		flex-direction: column;
	}

	.venue-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		font-size: smaller;
	}

	.venue-jamsessions {
		min-height: 0;
	}

	.venue-session {
		background-color: whitesmoke;
		border-radius: 24px;
		padding: 1em;
		margin-bottom: 1em;
		cursor: pointer;
	}

	/* 
	table {
		background: whitesmoke;
		padding: 0.5em;
		border-radius: 12px;
	} */

	td {
		vertical-align: top;
		padding: 0 0.5em;
	}

	td:first-child {
		padding-top: 0.1em;
	}
</style>
