<!-- The "Map" component renders an OpenLayers map with two layers, one for sessions, one for venues. 
 Whenever the stores $selectedSessions and $venues change, the map gets re-rendered. -->

<script lang="ts">
	import 'ol/ol.css';
	import { Feature } from 'ol';
	import { Map, View } from 'ol';
	import type Point from 'ol/geom/Point';
	import Overlay from 'ol/Overlay';
	import { fromLonLat, useGeographic } from 'ol/proj';
	import { Tile as TileLayer, Vector as VectorLayer } from 'ol/layer';
	import GeoJSON from 'ol/format/GeoJSON';
	// import WebGLTile from "ol/layer/WebGLTile";
	import { OSM, Vector as VectorSource } from 'ol/source';
	// import { PMTilesRasterSource } from "ol-pmtiles";

	import { onMount, mount } from 'svelte';
	import {
		selectedSessions
		// activeVenueProperties,
	} from '../stores';
	import {
		type SessionProperties,
		type SessionPropertiesWithVenue,
		type VenueProperties
	} from '../types';
	import { flyTo, sessionStyle } from './mapUtils';
	import SessionPopup from './SessionPopup.svelte';
	import Message from './Message.svelte';
	import AddSessionButton from './AddSessionButton.svelte';

	interface Props {
		background: boolean;
		onClickBackground?: () => void;
	}

	let { background, onClickBackground }: Props = $props();

	type SessionsByVenue = { [key: number]: Feature[] };
	let sessionsById: { [key: number]: Feature } = {};
	let popupsById: { [key: number]: [Overlay, HTMLElement] } = {};

	// initialise variables
	let map: Map | undefined = $state(undefined);
	let view: View | undefined = $state(undefined);

	// const baseLayer = new WebGLTile({
	//     source: new PMTilesRasterSource({
	//         projection: 'EPSG:4326',
	//         url: process.env.TILES_ADDRESS.replace("/\/$/", "") + "/london.pmtiles",
	//         attributions: ['© <a href="https://openstreetmap.org">OpenStreetMap</a>'],
	//         tileSize: [512,512]
	//     })
	// });

	const sessionsSource = new VectorSource();
	const sessionsLayer = new VectorLayer({
		source: sessionsSource,
		style: sessionStyle
	});

	const venuesSource = new VectorSource();
	const venuesLayer = new VectorLayer({
		source: venuesSource
	});

	let activePopup: Overlay | null;

	let isMobile: boolean = $state(false);

	// reactivity - whenever the stores change (and the condition evaluates to true, the respective render function is called)
	$effect(() => {
		if (map !== undefined && $selectedSessions !== null) renderSessions();
	});

	$effect(() => {
		if (map) {
			if (background && window.matchMedia('(max-width: 480px)').matches) {
				// hide controls
				document.querySelectorAll('.ol-control').forEach((elem) => {
					(elem as HTMLElement).style.display = 'none';
				});
				// remove popups
				map.getOverlays().forEach((o) => map!.removeOverlay(o));
			} else {
				document.querySelectorAll('.ol-control').forEach((elem) => {
					(elem as HTMLElement).style.display = 'unset';
				});
			}
		}
	});

	// handler funcs - both follow the same structure
	// 1. overwrite features in source with values from store
	// 2. hide the other layer on map so that only one of venues/sessions is visible

	/** Clears the map before rendering the $selectedSessions svelte store as an OpenLayers
	 * GeoJSON layer on the map. Also calls `addPopup()` and zooms to the newly added features.
	 */
	const renderSessions = () => {
		// clear source before re-render
		sessionsSource.clear();
		if (map) map.getOverlays().forEach((i: Overlay) => map!.removeOverlay(i));

		// hide venues on map - only one layer can be visible at a time
		venuesLayer.setVisible(false);

		// add features from $selectedSessions store
		if ($selectedSessions!.features && $selectedSessions!.features.length != 0) {
			const features = new GeoJSON().readFeatures($selectedSessions, {
				dataProjection: 'EPSG:4326',
				featureProjection: 'EPSG:3857'
			});
			sessionsSource.addFeatures(features);

			let sessionsByVenue: SessionsByVenue = {};
			for (let ft of features) {
				let p = ft.getProperties();
				if (!(p.venue_id in sessionsByVenue)) {
					sessionsByVenue[p.venue_id] = [];
				}
				sessionsByVenue[ft.getProperties().venue_id].push(ft);

				sessionsById[p.session_id] = ft;
			}
			addPopups(sessionsByVenue);

			sessionsLayer.setVisible(true);

			if (view) {
				let activeSessionId = window.sessionStorage.getItem('activeSessionId');
				if (activeSessionId !== null) {
					let activeMarker = features.find(
						(i: Feature) => i.getProperties().session_id === parseInt(activeSessionId)
					);
					if (activeMarker) view.fit(activeMarker.getGeometry()?.getExtent()!, { maxZoom: 15 });
					// else console.log(`active marker ${activeSessionId} not in $selectedSessions`);
					// above is debug statement - normally it's fine if it happens when going back from e.g. venue view to start view
				} else {
					view.fit(sessionsSource.getExtent()!, {
						maxZoom: 13,
						duration: 1000
					});
				}
			}
		} else {
			mount(Message, {
				props: { message: 'No sessions found for this day!' },
				target: document.body
			});
		}
	};

	/** Display the popup belonging to a feature.
	 *
	 * @param feature - the Session Feature
	 * @param popup - the Overlay object
	 * @param elem - the target element of the Overlay object
	 */
	const showPopup = (feature: Feature, popup: Overlay, elem: HTMLElement) => {
		if (activePopup !== popup) {
			// remove existing popups
			map!.getOverlays().forEach((i) => map!.removeOverlay(i));

			// fly to the geometry and right zoom level first, then add the popup and center on that
			// has to be this sequence, as the center of the popup will be different depending on zoom level
			flyTo(view!, feature.getGeometry()?.getExtent()!, () => {
				map!.addOverlay(popup); // show popup when feature is clicked
				const rect = elem.getBoundingClientRect();
				flyTo(
					view!,
					map!.getCoordinateFromPixel([
						rect.right - (rect.right - rect.left) / 2,
						rect.top - (rect.top - rect.bottom) / 2
					])
				);
			});
			activePopup = popup;
		}
	};

	/**
	 * Adds a popup (ol.Overlay) with name and short description to a map feature, wraps svelte components SessionPopup
	 * and VenuePopup (depending on whether the property `session_id` is present in the feature properties).
	 *
	 * @param ft - an OpenLayers Feature object to attach the popup to
	 * @returns nothing
	 * */
	const addPopups = (sessionsByVenue: SessionsByVenue): void => {
		for (let [_, features] of Object.entries(sessionsByVenue)) {
			const propertiesList = features.map((i) => i.getProperties());
			const isSession = 'session_id' in propertiesList[0];

			let elem = document.createElement('div');
			let geom = features[0].getGeometry() as Point;
			let popup: Overlay;

			if (isSession) {
				mount(SessionPopup, {
					props: {
						propertiesList: propertiesList as SessionPropertiesWithVenue[],
						onclose: () => {
							map!.removeOverlay(popup);
							activePopup = null;
						}
					},
					target: elem
				});
			}
			document.getElementById('popups')!.appendChild(elem);

			popup = new Overlay({
				position: geom.getCoordinates(),
				positioning: 'bottom-left',
				element: elem,
				offset: [20, -5]
				// autoPan: true,
			});

			for (let p of propertiesList) {
				popupsById[p.session_id] = [popup, elem];
			}

			if (map && view)
				map.on('click', (ev) => {
					map!.forEachFeatureAtPixel(ev.pixel, (feature, _) => {
						// show popup
						if (feature == features[0]) {
							showPopup(feature, popup, elem);
						}
					});
				});
		}
	};

	const zoomToSession = (id: number) => showPopup(sessionsById[id], ...popupsById[id]);

	// when mounting the component, initialise the map object
	onMount(async () => {
		isMobile = window.matchMedia('(max-width: 480px)').matches;
		// useGeographic();

		view = new View({
			center: fromLonLat([-0.12574, 51.50853]),
			zoom: 8,
			extent: [
				// greater london
				...fromLonLat([-0.619354, 51.234407]),
				...fromLonLat([0.296254, 51.731281])
			]
		});

		map = new Map({
			target: 'map',
			layers: [
				new TileLayer({
					source: new OSM()
				}),
				sessionsLayer,
				venuesLayer
			],
			view: view
		});

		//
		map.on('pointermove', (ev) => {
			// change cursor to pointer on hover
			if (!ev.dragging) {
				map!.getTargetElement()!.style.cursor = map!.hasFeatureAtPixel(
					map!.getEventPixel(ev.originalEvent)
				)
					? 'pointer'
					: 'grab';
			}
		});
	});
</script>

<div
	id="map"
	class:relative={!isMobile || (isMobile && background)}
	onzoom={() => zoomToSession(parseInt(window.sessionStorage.getItem('activeSessionId')!))}
	class:map-background={background}
	class:map-foreground={!background}
>
	{#if background}
		<div id="map-background-blur" onclick={onClickBackground}></div>
	{/if}
	{#if !background || !isMobile}
		{#await new Promise((resolve) => setTimeout(resolve, 1)) then}
			<AddSessionButton />
		{/await}
	{/if}
	<div id="popups"></div>
</div>

<style>
	#map {
		height: 100%;
		transition: width ease-in-out 1s;
	}

	.relative {
		position: relative;
	}

	.map-foreground {
		width: 100%;
		pointer-events: unset;
	}

	.map-background {
		flex-shrink: 0;
		width: 50%;
	}

	#map-background-blur {
		position: absolute;
		background: transparent;
		height: 100%;
		width: 100%;
		z-index: 800000;
		pointer-events: none;
	}

	:global(.ol-control button) {
		height: 2em;
		width: 2em;
		font-size: large;
	}

	:global(.ol-zoom) {
		bottom: 7em;
		top: unset;
		right: unset;
		left: 2em;
		border: 2px grey solid;
		border-radius: 8px;
	}

	@media (max-width: 480px) {
		.map-background {
			width: 20%;
			pointer-events: none;
		}

		#map-background-blur {
			pointer-events: all;
		}

		:global(.ol-zoom) {
			top: 1.5em;
			right: 0.5em;
			bottom: unset;
			left: unset;
		}

		:global(.ol-control button) {
			height: 1.375em;
			width: 1.375em;
		}

		:global(.ol-touch .ol-control button) {
			font-size: 1.3em;
		}

		:global(.ol-attribution.ol-uncollapsible) {
			top: 0;
			left: unset;
			right: 0;
			bottom: unset;
			border-top-left-radius: 0;
			border-bottom-left-radius: 4px;
		}
	}

	#popups {
		z-index: 500000;
	}
</style>
