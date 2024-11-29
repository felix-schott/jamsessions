import type { PageLoad } from './$types';
import { activeTab } from '../../../stores';

export const load: PageLoad = ({ params }) => {
    activeTab.set("session");
    return {
        sessionId: parseInt(params.sessionSlug.split("-")[params.sessionSlug.split("-").length - 1])
    }
}

export interface SessionSlugData {
    sessionId: number
}