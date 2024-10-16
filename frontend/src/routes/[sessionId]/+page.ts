import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
    return {
        sessionId: params.sessionId
    }
}

export interface SessionSlugData {
    sessionId: number
}