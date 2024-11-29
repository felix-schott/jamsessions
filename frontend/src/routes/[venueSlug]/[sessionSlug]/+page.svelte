<script lang="ts">
	import SessionDetails from '../../../lib/SessionDetails.svelte';
	import type { SessionComment, SessionWithVenueFeature } from '../../../types';
	import { getCommentsBySessionId, getSessionById } from '../../../api';
	import type { SessionSlugData } from './+page.js';
	import { selectedSessions } from '../../../stores';

	interface Props {
		data: SessionSlugData;
	}

	let { data }: Props = $props();

	const getSessionByIdWithErrorHandling = async (): Promise<
		[SessionWithVenueFeature, SessionComment[]]
	> => {
		try {
			let sessionFeature = await getSessionById(data.sessionId);
			$selectedSessions = { type: 'FeatureCollection', features: [sessionFeature] }; // propagate to map
			let comments = await getCommentsBySessionId(data.sessionId);
			return [sessionFeature, comments];
		} catch (e) {
			alert('An error occured when waiting for data from the server: ' + (e as Error).message);
			throw e;
		}
	};
</script>

{#await getSessionByIdWithErrorHandling() then p}
	<SessionDetails sessionComments={p[1]} sessionProperties={p[0].properties} />
{/await}
