<script lang="ts">
    import SessionDetails from "../../lib/SessionDetails.svelte";
    import type { SessionFeature, VenueFeature } from "../../types";
    import { getSessionById, getVenueById } from "../../api";
    import type { SessionSlugData } from "./+page.js";

    export let data: SessionSlugData;

    const getSessionByIdWithErrorHandling = async (): Promise<[SessionFeature, VenueFeature]> => {
        try {
            let sessionProperties = await getSessionById(data.sessionId);
            let venueProperties = await getVenueById(sessionProperties.properties.venue!);
            return [sessionProperties, venueProperties]
        } catch (e) {
            alert("An error occured when waiting for data from the server: " + (e as Error).message)
            throw e
        }
    }
</script>

{#await getSessionByIdWithErrorHandling() then p}
    <SessionDetails venueProperties={p[1].properties} sessionProperties={p[0].properties} />
{/await}