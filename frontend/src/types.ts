import type { FeatureCollection, Feature, Point } from "geojson";

export enum MapLayer {
	NONE,
    SESSIONS,
    VENUES
}

export enum Backline {
	PA = 'PA', 
	GUITAR_AMP = 'Guitar_Amp', 
	BASS_AMP = 'Bass_Amp', 
	KEYS = 'Keys',
	DRUMS = 'Drums', 
	MIC = 'Microphone', 
	MISC_PERCUSSION = 'MiscPercussion',
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

export enum Genre {
	ANY = 'ANY',
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
	WORLD_MUSIC = 'World_Music',
}

export enum Interval {
	ONCE = 'Once',
	DAILY = 'Daily',
	WEEKLY = 'Weekly',
	FIRSTOFMONTH = 'FirstOfMonth',
	SECONDOFMONTH = 'SecondOfMonth',
	THIRDOFMONTH = 'ThirdOfMonth' ,
	FOURTHOFMONTH = 'FourthOfMonth', 
	LASTOFMONTH = 'LastOfMonth'
}

export interface SessionProperties {
    session_id?: number
	venue?: number
	session_name: string
	genres: Genre[]
	description: string
	start_time_utc: Date | string
	interval: Interval
	duration_minutes: number
	session_comments: string[]
	session_website: string
	dt_updated_utc?: Date
}

export interface SessionPropertiesWithVenue extends SessionProperties, VenueProperties {};

export type SessionFeatureCollection = FeatureCollection<Point, SessionProperties>;
export type SessionWithVenueFeatureCollection = FeatureCollection<Point, SessionPropertiesWithVenue>;
export type VenuesFeatureCollection = FeatureCollection<Point, VenueProperties>;
export type SessionFeature = Feature<Point, SessionProperties>;
export type SessionWithVenueFeature = Feature<Point, SessionPropertiesWithVenue>;
export type VenueFeature = Feature<Point, VenueProperties>;