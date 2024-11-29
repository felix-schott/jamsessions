import { Backline, Genre, type CommentBody, type SessionWithVenueFeatureCollection, type SessionComment, type VenuesFeatureCollection, type SessionWithVenueFeature, type VenueFeature, type VenueProperties, type SessionProperties, type SessionPropertiesWithVenue, type SessionFeatureCollection } from "./types";

const API_ROOT = process.env.API_ADDRESS?.replace("/\/$/", "");
const API_VERSION = "v1";

const API_ADDRESS = API_ROOT + "/" + API_VERSION;

interface ErrorResponse {
    detail: string
}

const headers = new Headers();
headers.append("Content-Type", "application/json")

export const getVenues = async (): Promise<VenuesFeatureCollection> => {
    let response = await fetch(API_ADDRESS + "/venues")
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        return await response.json() as VenuesFeatureCollection
    }
}

export interface SessionOptions {
    date?: Date
    endDate?: Date
    genre?: Genre
    backline?: Backline[]
}

export const getSessions = async (opts: SessionOptions = {}): Promise<SessionWithVenueFeatureCollection> => {
    // parse opts and build URL
    let queryParams: string[] = []
    if (opts.date && !opts.endDate) {
        queryParams.push("date=" + opts.date.toISOString().slice(0, 10))
    }
    if (opts.date && opts.endDate) {
        queryParams.push("date=" + opts.date.toISOString().slice(0, 10) + "%2F" + opts.endDate.toISOString().slice(0, 10))
    }
    if (opts.genre && opts.genre != Genre.ANY) {
        queryParams.push("genre=" + opts.genre)
    }
    opts.backline = opts.backline?.filter(i => i);
    if (opts.backline && opts.backline.length != 0) {
        queryParams.push("backline=" + opts.backline.join(","))
    }
    const url = API_ADDRESS + "/jamsessions" + (queryParams.length > 0 ? "?" : "") + queryParams.join("&")
    let response = await fetch(url)
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        let body = await response.json() as SessionWithVenueFeatureCollection;
        return body
    }
}

export const getSessionById = async (id: number): Promise<SessionWithVenueFeature> => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id)
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        return await response.json() as SessionWithVenueFeature;
    }
}

export const getSessionsByVenueId = async (venueId: number): Promise<SessionWithVenueFeatureCollection> => {
    let response = await fetch(API_ADDRESS + "/venues/" + venueId + "/jamsessions")
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        return await response.json() as SessionWithVenueFeatureCollection;
    }
}

export const getVenueById = async (id: number): Promise<VenueFeature> => {
    let response = await fetch(API_ADDRESS + "/venues/" + id)
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        return await response.json() as VenueFeature;
    }
}

export const getCommentsBySessionId = async (id: number): Promise<SessionComment[]> => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id + "/comments")
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    } else {
        return await response.json() as SessionComment[];
    }
}

export const postCommentForSessionById = async (id: number, payload: CommentBody) => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id + "/comments", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: headers
    })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}

export const postSuggestionForSessionById = async (id: number, payload: CommentBody) => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id + "/suggestions", {
        method: "POST",
        headers: headers,
        body: JSON.stringify(payload)
    })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}

// export const postVenue = async (payload: VenueProperties) => {
//     let response = await fetch(API_ADDRESS + "/venues", {
//         method: "POST",
//         body: JSON.stringify(payload)
//     })
//     if (!response.ok) {
//         throw new Error((await response.json() as ErrorResponse)["detail"])
//     }
// }

export const postSession = async (payload: SessionProperties | SessionPropertiesWithVenue) => {
    let response = await fetch(API_ADDRESS + "/jamsessions", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: headers,
    })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}


export const patchSessionById = async (id: number, payload: SessionProperties | {}): Promise<void> => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id, {
        method: "PATCH",
        body: JSON.stringify(payload)
    })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}

export const patchVenueById = async (id: number, payload: VenueProperties | {}): Promise<void> => {
    let response = await fetch(API_ADDRESS + "/venues/" + id, {
        method: "PATCH",
        body: JSON.stringify(payload)
    })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}

export const deleteSessionById = async (id: number): Promise<void> => {
    let response = await fetch(API_ADDRESS + "/jamsessions/" + id, { method: "DELETE" });
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}

export const deleteVenueById = async (id: number): Promise<void> => {
    let response = await fetch(API_ADDRESS + "/venues/" + id, { method: "DELETE" })
    if (!response.ok) {
        throw new Error((await response.json() as ErrorResponse)["detail"])
    }
}