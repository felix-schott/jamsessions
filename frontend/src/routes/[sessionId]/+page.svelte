<script lang="ts">
    import SessionDetails from "../../lib/SessionDetails.svelte";
    import type { SessionComment, SessionFeature, VenueFeature } from "../../types";
    import { getCommentsBySessionId, getSessionById, getVenueById } from "../../api";
    import type { SessionSlugData } from "./+page.js";

    interface Props {
        data: SessionSlugData;
    }

    let { data }: Props = $props();

    const getSessionByIdWithErrorHandling = async (): Promise<[SessionFeature, VenueFeature, SessionComment[]]> => {
        try {
            let sessionProperties = await getSessionById(data.sessionId);
            let venueProperties = await getVenueById(sessionProperties.properties.venue!);
            let comments = await getCommentsBySessionId(data.sessionId);
            return [sessionProperties, venueProperties, comments] 
        } catch (e) {
            alert("An error occured when waiting for data from the server: " + (e as Error).message)
            throw e
        }
    }
</script>

{#await getSessionByIdWithErrorHandling() then p}
    <SessionDetails venueProperties={p[1].properties} sessionComments={p[2]} sessionProperties={p[0].properties} />
{/await}