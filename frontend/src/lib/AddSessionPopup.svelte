<script lang="ts">
	import { addSessionPopupVisible } from '../stores';
	import Modal from './Modal.svelte';

	import { getVenues, postSession } from '../api';
	import {
		Backline,
		Genre,
		Interval,
		type SessionProperties,
		type SessionPropertiesWithVenue,
		type VenueFeature,
		type VenueProperties,
		type VenuesFeatureCollection
	} from '../types';
	import InfoIcon from './icons/InfoIcon.svelte';
	import { constructIntervalString, minutesBetweenTimestamps } from './timeUtils';

	let venueId: string = $state('');
	let form: HTMLFormElement | undefined = $state();

	let sessionName: string = $state('');
	let sessionTimeStart: string = $state('');
	let sessionTimeFinish: string = $state('');
	let sessionDescription: string = $state('');
	let sessionWebsite: string = $state('');
	let sessionInterval: Interval | undefined = $state();
	let sessionDate: string = $state('');
	let venueName: string = $state('');
	let venueAddress1: string = $state('');
	let venueAddress2: string = $state('');
	let venueCity: string = $state('');
	let venuePostcode: string = $state('');
	let venueWebsite: string = $state('');

	let venuesLoaded: boolean = $state(false);

	let newVenueHidden = $derived(venueId != 'new-venue' && venuesLoaded);
	let buttonDisabled = $derived(
		!sessionName && !sessionTimeStart && !sessionTimeFinish && !sessionDescription
	);

	// small wrapper around getVenues for better state management
	// (venuesLoaded controls whether the menu to add a venue is shown)
	const getVenuesAsList = async (): Promise<VenueFeature[]> => {
		const venues = await getVenues();
		if (venues && venues.features) {
			venuesLoaded = true;
			// sort alphabetically
			let alphabeticalVenues = venues.features.toSorted((a, b) =>
				a.properties.venue_name.localeCompare(b.properties.venue_name, undefined, {
					sensitivity: 'base'
				})
			);
			return alphabeticalVenues;
		}
		venuesLoaded = false;
		return [];
	};

	const onSubmit = async (ev: MouseEvent) => {
		ev.preventDefault();

		let d = new Date(sessionDate);
		let sessionParams: SessionProperties | SessionPropertiesWithVenue = {
			session_name: sessionName,
			description: sessionDescription,
			interval: sessionInterval!,
			start_time_utc: new Date(
				d.getFullYear(),
				d.getMonth(),
				d.getDate(),
				...sessionTimeStart.split(':').map((i) => parseInt(i))
			).toISOString(),
			duration_minutes: minutesBetweenTimestamps(sessionTimeStart, sessionTimeFinish),
			genres: Array.from(document.querySelectorAll('.genre-checkbox:checked')).map(
				(i) => i.id.replace('session-genre-', '') as Genre
			),
			session_website: sessionWebsite
		};
		// add venue if necessary
		if (venueId == 'new-venue') {
			if (venueName === '') {
				alert('Please add the name of the venue');
				return;
			}
			if (venueAddress1 === '') {
				alert('Please add the address of the venue');
			}
			if (venueCity === '') {
				alert('Please add the city of the venue');
			}
			if (venuePostcode === '') {
				alert('Please add the city of the venue');
			}
			if (venueWebsite === '') {
				alert('Please add the website of the venue');
			}
			let venueParams: VenueProperties = {
				venue_name: venueName,
				address_first_line: venueAddress1,
				address_second_line: venueAddress2 ? venueAddress2 : undefined,
				city: venueCity,
				postcode: venuePostcode,
				venue_website: venueWebsite,
				backline: Array.from(document.querySelectorAll('.backline-checkbox:checked')).map(
					(i) => i.id.replace('venue-backline-', '') as Backline
				)
			};
			Object.assign(sessionParams, venueParams);
		} else {
			sessionParams['venue'] = parseInt(venueId);
		}
		try {
			await postSession(sessionParams);
			alert(
				"Thank you for submitting a new session! We'll review your suggestions and apply the changes. If there is anything else, you can email felix.schott@proton.me"
			);
			$addSessionPopupVisible = false;
		} catch (e) {
			alert(
				`An error occured trying to communicate with the server (${e}). Please try again or email felix.schott@proton.me`
			);
		}
	};
</script>

{#await getVenuesAsList() then venues}
	{#if venues}
		<Modal
			isVisible={() => $addSessionPopupVisible}
			hide={() => {
				$addSessionPopupVisible = false;
			}}
		>
			<form bind:this={form}>
				<h2>Add new session to the database</h2>
				<div class="card">
					<h3>Venue</h3>
					{#if venuesLoaded}
						<select title="Select venue" bind:value={venueId}>
							{#each venues as venue, idx}
								{#if idx === 0}
									<option value={venue.properties.venue_id} selected
										>{venue.properties.venue_name}</option
									>
								{:else}
									<option value={venue.properties.venue_id}>{venue.properties.venue_name}</option>
								{/if}
							{/each}
							<option value="new-venue">None of the above</option>
						</select>
					{/if}
					<div class:hidden={newVenueHidden} style="margin-top: 1em;">
						<div class="vertical">
							<div style="justify-content: center; display: flex; margin-top: 1em;">
								<b>Add new venue to the database</b>
							</div>
							<label for="venue-name"
								>Name of the venue <input
									id="venue-name"
									bind:value={venueName}
									type="text"
									required
								/></label
							>
							<label for="venue-address-first-line"
								>Address 1st line <input
									type="text"
									bind:value={venueAddress1}
									id="venue-address-first-line"
									required
								/></label
							>
							<label for="venue-address-second-line"
								>Address 2nd line <input
									type="text"
									bind:value={venueAddress2}
									id="venue-address-second-line"
								/></label
							>
							<label for="venue-address-city"
								>City <input
									type="text"
									id="venue-address-city"
									bind:value={venueCity}
									required
								/></label
							>
							<label for="venue-address-postcode"
								>Postcode <input
									type="text"
									id="venue-address-postcode"
									bind:value={venuePostcode}
									required
								/></label
							>
							<label for="venue-website"
								>Website <input
									type="url"
									id="venue-website"
									bind:value={venueWebsite}
									required
								/></label
							>
						</div>
						Backline provided
						<div class="checkboxes">
							{#each Object.values(Backline) as backline}
								<label for="venue-backline-{backline}"
									><input
										type="checkbox"
										class="backline-checkbox"
										id="venue-backline-{backline}"
										name={backline}
									/>{backline.replace('_', ' ')}</label
								>
							{/each}
						</div>
					</div>
				</div>
				<div class="card">
					<h3>Session details</h3>
					<div class="vertical">
						<label for="session-name"
							>Name of the session <input
								id="session-name"
								bind:value={sessionName}
								type="text"
								required
							/></label
						>
						<div>
							<label for="session-date"
								>Next date of the session <input
									type="date"
									bind:value={sessionDate}
									required
								/></label
							>
							<label for="session-time-start"
								>From <input
									type="time"
									id="session-time-start"
									bind:value={sessionTimeStart}
									required
								/></label
							>
							<label for="session-time-finish"
								>To <input
									type="time"
									id="session-time-finish"
									bind:value={sessionTimeFinish}
									min={sessionTimeStart}
									required
								/></label
							>
							<div>How often does the session happen?</div>
							<select title="Select interval" bind:value={sessionInterval}>
								{#if sessionDate !== ''}
									{#each Object.values(Interval) as interval, idx}
										{#if idx === 0}
											<option value={interval} selected
												>{@html constructIntervalString(interval, new Date(sessionDate))}</option
											>
										{:else}
											<option value={interval}
												>{@html constructIntervalString(interval, new Date(sessionDate))}</option
											>
										{/if}
									{/each}
								{/if}
							</select>
						</div>
						<label for="description" class="inline"
							>Description <textarea
								id="description"
								required
								bind:value={sessionDescription}
								style="height: 6em;"
								placeholder="Who is hosting? Do you have to pay to attend? Anything else worth mentioning?"
							></textarea></label
						>
					</div>
					Main genre(s)
					<div class="checkboxes">
						{#each Object.values(Genre) as genre}
							{#if genre != 'Any'}
								<label for="session-genre-{genre}"
									><input
										type="checkbox"
										class="genre-checkbox"
										id="session-genre-{genre}"
										name={genre}
									/>{genre.replace('_', ' ')}</label
								>
							{/if}
						{/each}
					</div>
					<div class="vertical">
						<label for="session-website"
							><div style="display: flex; align-items: center;">
								Website <InfoIcon
									style="margin-left: 0.5em;"
									title="May be the same as the venue website."
								/>
							</div>
							<input type="url" id="session-website" bind:value={sessionWebsite} required /></label
						>
					</div>
				</div>
				<div style="display: flex; justify-content: center; margin-top: 1em;">
					<button disabled={buttonDisabled} type="submit" onclick={onSubmit}
						><span>Submit</span></button
					>
				</div>
			</form>
		</Modal>
	{/if}
{/await}

<style>
	.vertical {
		display: flex;
		flex-direction: column;
	}

	.vertical label {
		display: flex;
		justify-content: space-between;
	}

	.vertical b {
		margin-bottom: 1em;
	}

	.vertical label {
		margin-bottom: 1em;
	}

	input {
		max-width: 12em;
		margin-left: 1em;
	}

	@media (max-width: 480px) {
		input {
			margin-left: 0;
			max-width: 9em;
		}

		textarea {
			width: 100%;
		}

		label.inline {
			display: inline-block;
		}
	}

	select {
		margin-bottom: 1em;
		font-size: unset;
		background: white;
	}

	.card {
		margin-top: 1em;
	}

	h3 {
		margin-block-start: 0;
	}
</style>
