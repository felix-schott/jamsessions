import type { PageLoad } from './$types';
import { activeTab } from '../../stores';

export const load: PageLoad = ({ params }) => {
    activeTab.set("venue");
    return {
        venueSlug: params.venueSlug,
        venueId: parseInt(params.venueSlug.split("-")[params.venueSlug.split("-").length - 1])
    }
}

export interface VenueSlugData {
    venueSlug: string
    venueId: number
}