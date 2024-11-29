export const sanitisePathElement = (s: string) => s.replaceAll(" ", "-").replaceAll(/[^a-zA-Z\-]/g, "").toLowerCase();

export const extractDomain = (url: string) => {
    let matches = url.matchAll(/(?<=(https:\/\/|http:\/\/))(.*?)(\/|$)/g);
    if (matches) {
        return [...matches][0][2]
    }

    return url // if matching fails, return input as fallback
}

export const processVenueAndSessionUrl = (venueUrl: string, sessionUrl: string): { venueWebsite: string, sessionWebsite: string | null } => {
    let venueDomain = extractDomain(venueUrl);
    let sessionDomain = extractDomain(sessionUrl);
    if (venueDomain === sessionDomain) {
        if (venueUrl === sessionUrl) { // if they're exactly the same, just display one
            return { venueWebsite: venueDomain, sessionWebsite: null }
        }
        // we're assuming:
        //  venueUrl = example.com
        //  sessionUrl = example.com/events/...
        let sessionPathComponents = new URL(sessionUrl).pathname.slice(1).split("/");
        // leave venue as is, include first path param in session url
        return { venueWebsite: venueDomain, sessionWebsite: sessionDomain + "/" + sessionPathComponents[0] }

    }
    return { venueWebsite: venueDomain, sessionWebsite: sessionDomain }
}