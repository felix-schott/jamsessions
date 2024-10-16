import { writable, type Writable } from "svelte/store";
import type { VenuesFeatureCollection, SessionProperties, SessionFeatureCollection, VenueFeature} from './types';
import { MapLayer } from "./types";

// stored in sessionStorage for cross-page persistence: activeSessionId, selectedGenre, selectedDate, selectedBackline

// containers for data returned by REST API
export const selectedSessions: Writable<SessionFeatureCollection | null> = writable(null);
export const venuesById: Writable<{[key: number]: VenueFeature} | null> = writable(null);

// state controlling variables
export const visibleLayer: Writable<MapLayer> = writable(MapLayer.NONE);
export const addSessionPopupVisible: Writable<boolean> = writable(false);
export const filterMenuVisible: Writable<boolean> = writable(false);
export const infoVisible: Writable<boolean> = writable(false);
export const loading: Writable<boolean> = writable(false);
export const message: Writable<string> = writable("");