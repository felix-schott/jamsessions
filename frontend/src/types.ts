import type { FeatureCollection, Feature, Point } from "geojson";

export enum MapLayer {
	NONE,
	SESSIONS,
	VENUES
}

export type TabOptions = "map" | "list" | "session";

export enum Backline {
	PA = 'PA',
	GUITAR_AMP = 'Guitar_Amp',
	BASS_AMP = 'Bass_Amp',
	KEYS = 'Keys',
	DRUMS = 'Drums',
	MIC = 'Microphone',
	MISC_PERCUSSION = 'MiscPercussion',
}

export enum Genre {
	ANY = 'Any',
	STRAIGHT_AHEAD = 'Straight-Ahead_Jazz',
	MODERN_JAZZ = 'Modern_Jazz',
	TRAD_JAZZ = 'Trad_Jazz',
	JAZZ_FUNK = 'Jazz-Funk',
	FUSION = 'Fusion',
	LATIN_JAZZ = 'Latin_Jazz',
	FUNK = 'Funk',
	BLUES = 'Blues',
	FOLK = 'Folk',
	ROCK = 'Rock',
	POP = 'Pop',
	WORLD_MUSIC = 'World_Music',
}

export enum Interval {
	ONCE = 'Once',
	DAILY = 'Daily',
	WEEKLY = 'Weekly',
	FORTNIGHTLY = 'Fortnightly',
	FIRSTOFMONTH = 'FirstOfMonth',
	SECONDOFMONTH = 'SecondOfMonth',
	THIRDOFMONTH = 'ThirdOfMonth',
	FOURTHOFMONTH = 'FourthOfMonth',
	LASTOFMONTH = 'LastOfMonth',
	IRREGULARWEEKLY = 'IrregularWeekly'
}

export interface VenueProperties {
	venue_id?: number
	venue_name: string
	address_first_line: string
	address_second_line?: string
	city: string
	postcode: string
	venue_website: string
	backline: Backline[]
	venue_comments?: string[]
	venue_dt_updated_utc?: Date
}

export interface SessionProperties {
	session_id?: number
	venue?: number
	session_name: string
	genres: Genre[]
	description: string
	start_time_utc: Date | string
	dates?: string[]
	interval: Interval
	duration_minutes: number
	session_website: string
	rating?: number
	dt_updated_utc?: Date
}

export interface SessionComment {
	comment_id: number
	session: number
	author: string
	content: string
	dt_posted: string
	rating: number // between 1 and 5
}

export interface CommentBody {
	session?: number
	author?: string
	content: string
	rating?: number
}

export interface SessionPropertiesWithVenue extends SessionProperties, VenueProperties { };

export type SessionFeatureCollection = FeatureCollection<Point, SessionProperties>;
export type SessionWithVenueFeatureCollection = FeatureCollection<Point, SessionPropertiesWithVenue>;
export type VenuesFeatureCollection = FeatureCollection<Point, VenueProperties>;
export type SessionFeature = Feature<Point, SessionProperties>;
export type SessionWithVenueFeature = Feature<Point, SessionPropertiesWithVenue>;
export type VenueFeature = Feature<Point, VenueProperties>;