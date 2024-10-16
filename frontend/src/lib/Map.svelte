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

	import { onMount } from 'svelte';
	import {
		visibleLayer,
		selectedSessions,
		// activeVenueProperties,
	} from '../stores';
	import { MapLayer, type SessionProperties, type SessionPropertiesWithVenue, type VenueProperties } from '../types';
	import { flyTo, sessionStyle } from './mapUtils';
	import SessionPopup from './SessionPopup.svelte';
	import VenuePopup from './VenuePopup.svelte';
	import Message from './Message.svelte';

	type SessionsByVenue = { [key: number]: Feature[] };

	// initialise variables
	let map: Map;
	let view: View;

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

	// reactivity - whenever the stores change (and the condition evaluates to true, the respective render function is called)
	$: if ($visibleLayer == MapLayer.SESSIONS && $selectedSessions !== null) renderSessions();
	// $: if ($visibleLayer == MapLayer.VENUES && $venues != null) renderVenues();

	// handler funcs - both follow the same structure
	// 1. overwrite features in source with values from store
	// 2. hide the other layer on map so that only one of venues/sessions is visible

	/** Clears the map before rendering the $selectedSessions svelte store as an OpenLayers
	 * GeoJSON layer on the map. Also calls `addPopup()` and zooms to the newly added features.
	 */
	const renderSessions = () => {
		// clear source before re-render
		sessionsSource.clear();
		if (map) map.getOverlays().forEach((i: Overlay) => map.removeOverlay(i));

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
				if (!(ft.getProperties().venue_id in sessionsByVenue)) {
					sessionsByVenue[ft.getProperties().venue_id] = [];
				}
				sessionsByVenue[ft.getProperties().venue_id].push(ft);
			}
			addPopups(sessionsByVenue);

			sessionsLayer.setVisible(true);

			if (view) {
				let activeSessionId = window.sessionStorage.getItem('activeSessionId');
				if (activeSessionId !== null) {
					view.fit(
						features
							.find((i: Feature) => i.getProperties().session_id === parseInt(activeSessionId))
							?.getGeometry()
							?.getExtent()!,
						{ maxZoom: 15 }
					);
				} else {
					view.fit(sessionsSource.getExtent()!, {
						maxZoom: 13,
						duration: 1000
					});
				}
			}
		} else {
			new Message({
				props: { message: 'No sessions found for this day!' },
				target: document.body
			});
		}
	};

	/** Clears the map before rendering the $venues svelte store as an OpenLayers
	 * GeoJSON layer on the map. Also calls `addPopup()` and zooms to the newly added features.
	 */
	const renderVenues = () => {
		// clear source before re-render
		venuesSource.clear();
		if (map) map.getOverlays().forEach((i) => map.removeOverlay(i));

		// hide sessions on map - only one layer can be visible at a time
		sessionsLayer.setVisible(false);

		// add features
		// if ($venues!.features && $venues!.features.length != 0) {
		// 	const features = new GeoJSON().readFeatures($venues, {
		// 		dataProjection: 'EPSG:4326',
		// 		featureProjection: 'EPSG:3857'
		// 	});
		// 	venuesSource.addFeatures(features);
		// 	for (let ft of features) {
		// 		addPopup(ft);
		// 	}

		// 	venuesLayer.setVisible(true);

		// 	if (view) {
		// 		// if ($activeVenueId) {
		// 		// 	view.fit(
		// 		// 		features
		// 		// 			.find((i) => i.getProperties().venue_id == $activeVenueId)
		// 		// 			?.getGeometry()
		// 		// 			?.getExtent()!
		// 		// 	);
		// 		// } else {
		// 		view.fit(venuesSource.getExtent()!, {
		// 			maxZoom: 13,
		// 			duration: 1000
		// 		});
		// 		// }
		// 	}
		// }
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


			// TODO one popup for multiple sessions, might have to change SessionPopup

			const propertiesList = features.map(i => i.getProperties());
			const isSession = 'session_id' in propertiesList[0];

			let elem = document.createElement('div');
			let geom = features[0].getGeometry() as Point;
			let popup: Overlay;

			if (isSession) {
				const sessionPopup = new SessionPopup({
					props: {
						propertiesList: propertiesList as SessionPropertiesWithVenue[]
					},
					target: elem
				});
				// sessionPopup.$on('click', () => {
				// 	window.location.assign('/' + properties.session_id); // navigate to page /[slug] where [slug] is a session id
				// });
				sessionPopup.$on('close', () => map.removeOverlay(popup));
			} else {
				// new VenuePopup({
				// 	props: {
				// 		// onClick: () => {
				// 		//     console.log("venue popup clicked");
				// 		// },
				// 		// onClose: () => map.removeOverlay(popup),
				// 		properties: properties as VenueProperties
				// 	},
				// 	target: elem
				// });
			}
			document.getElementById('popups')!.appendChild(elem);

			popup = new Overlay({
				position: geom.getCoordinates(),
				positioning: 'bottom-left',
				element: elem,
				offset: [20, -5]
				// autoPan: true,
			});

			// let options = {
			//     root: document.querySelector(".ol-viewport"),
			//     rootMargin: "0px",
			//     threshold: 0.1,
			// };

			// //
			// const movePopupOnScreen = (
			//     entries: IntersectionObserverEntry[],
			//     popup: Overlay,
			// ) => {
			//     console.log("yeah outside")
			//     const rightIsOutside = entries[0].intersectionRect.left < entries[0].boundingClientRect.left;
			//     const leftIsOutside = entries[0].intersectionRect.right < entries[0].boundingClientRect.right;
			//     const topIsOutside = entries[0].intersectionRect.bottom < entries[0].boundingClientRect.bottom;
			//     const bottomIsOutside = entries[0].intersectionRect.top < entries[0].boundingClientRect.top;

			//     if (rightIsOutside) {
			//         popup.setPositioning("center-left")
			//     } else if (leftIsOutside) {
			//         popup.setPositioning("center-right")
			//     } else if (topIsOutside) {
			//         popup.setPositioning("bottom-center")
			//     } else if (bottomIsOutside) {
			//         popup.setPositioning("top-center")
			//     }
			// };

			// const observer = new IntersectionObserver(
			//     (entries, _) => movePopupOnScreen(entries, popup),
			//     options,
			// );
			// observer.observe(elem)

			map.on('click', (ev) => {
				map.forEachFeatureAtPixel(ev.pixel, (feature, _) => {
					// remove existing overlays
					map.getOverlays().forEach((i) => map.removeOverlay(i));
					// show popup
					if (feature == features[0]) {
						map.addOverlay(popup); // show popup when feature is clicked

						// zoom to center of popup
						const rect = elem.getBoundingClientRect();
						flyTo(view, map.getCoordinateFromPixel([rect.right - (rect.right - rect.left) / 2, rect.top - (rect.top - rect.bottom) / 2]))
					}
				});
			});
		}
	};

	// when mounting the component, initialise the map object
	onMount(async () => {
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
				map.getTargetElement()!.style.cursor = map.hasFeatureAtPixel(
					map.getEventPixel(ev.originalEvent)
				)
					? 'pointer'
					: 'grab';
			}
		});
	});
</script>

<div id="map" class="map"></div>
<div id="popups" />

<style>
	#map {
		height: 100%;
		width: 100%;
	}

	#popups {
		z-index: 500000;
	}

	:global(.ol-zoom) {
		top: unset;
		bottom: 0.5em;
		left: 0.5em;
	}
	@media (max-width: 480px) {
		:global(.ol-zoom) {
			top: 0.5em;
			right: 0.5em;
			bottom: unset;
			left: unset;
		}

		:global(.ol-attribution.ol-uncollapsible) {
			top: 0;
			left: 0;
			right: unset;
			bottom: unset;
			border-top-left-radius: 0;
			border-bottom-right-radius: 4px;
		}
	}
</style>
