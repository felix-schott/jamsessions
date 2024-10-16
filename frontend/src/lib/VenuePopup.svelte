<script lang="ts">
    import Popup from "./Popup.svelte";
    import type { VenueProperties } from "../types";

    export let properties: VenueProperties;

    let id = "popup-venue-" + properties.venue_id;
</script>

<Popup {id} on:click on:close>
    <span slot="title">
        <b>${properties.venue_name}</b>
    </span>
</Popup>