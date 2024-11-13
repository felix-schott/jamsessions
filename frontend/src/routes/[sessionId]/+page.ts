import type { PageLoad } from './$types';
import { activeTab } from '../../stores';

export const load: PageLoad = ({ params }) => {
    activeTab.set("session");
    return {
        sessionId: params.sessionId
    }
}

export interface SessionSlugData {
    sessionId: number
}