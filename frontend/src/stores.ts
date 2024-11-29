import { writable, type Writable } from "svelte/store";
import type { SessionWithVenueFeatureCollection, VenueFeature, TabOptions } from './types';

// stored in sessionStorage for cross-page persistence: activeSessionId, selectedGenre, selectedDate, selectedBackline

// containers for data returned by REST API
export const selectedSessions: Writable<SessionWithVenueFeatureCollection | null> = writable(null);
export const venuesById: Writable<{ [key: number]: VenueFeature } | null> = writable(null);

// state controlling variables
export const addSessionPopupVisible: Writable<boolean> = writable(false);
export const filterMenuVisible: Writable<boolean> = writable(false);
export const infoVisible: Writable<boolean> = writable(false);
export const loading: Writable<boolean> = writable(false);
export const message: Writable<string> = writable("");
export const activeTab: Writable<TabOptions> = writable("map");
export const editingSession: Writable<boolean> = writable(false);
export const selectedDateRange: Writable<number> = writable(0);
export const selectedDateStr: Writable<string> = writable("");