<script lang="ts">
	import { type SessionProperties, Genre, Backline } from '../types';
	import ButtonGroup from '$lib/ButtonGroup.svelte';
	import {
		deleteSessionById,
		deleteVenueById,
		patchSessionById,
		patchVenueById,
		postSuggestionForSessionById
	} from '../api';

	export let properties: SessionProperties;

	// state management for form inputs
	let addressChanged = false;
	let venueClosed = false;
	let timeChanged = false;
	let sessionClosed = false;
	let otherChanges = false;
	let genreChanged = false;
	let backlineChanged = false;

	// enable submit button when at least one of the inputs have data
	$: buttonEnabled =
		addressChanged ||
		venueClosed ||
		timeChanged ||
		sessionClosed ||
		genreChanged ||
		backlineChanged ||
		otherChanges;

	const onSubmit = async (ev: MouseEvent) => {
		ev.preventDefault();

		let success = true;

		if (addressChanged) {
			try {
				await patchVenueById(properties.venue!, {
					address_first_line: (
						document.getElementById('new-address-first-line') as HTMLInputElement
					).value,
					address_second_line: (
						document.getElementById('new-address-second-line') as HTMLInputElement
					).value,
					city: (document.getElementById('new-address-city') as HTMLInputElement).value,
					postcode: (document.getElementById('new-address-postcode') as HTMLInputElement).value
				});
			} catch (e) {
				success = false;
				alert(
					'An error occurred when trying to suggest a new address: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (venueClosed) {
			try {
				await deleteVenueById(properties.venue!);
			} catch (e) {
				success = false;
				alert(
					'An error occurred when trying to suggest that the venue is closed: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (timeChanged) {
			try {
				await patchSessionById(properties.session_id!, {
					start_time_utc: (document.getElementById('new-time-start') as HTMLInputElement).value,
					duration_minutes: Math.floor(
						(new Date(
							(document.getElementById('new-time-finish') as HTMLInputElement).value
						).getTime() -
							new Date(
								(document.getElementById('new-time-start') as HTMLInputElement).value
							).getTime() /
								1000) /
							60
					)
				});
			} catch (e) {
				success = false;
				alert(
					'An error occurred trying to suggest there is a different time for the session: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (sessionClosed) {
			try {
				await deleteSessionById(properties.session_id!);
			} catch (e) {
				success = false;
				alert(
					'An error occurred trying to propose the deletion of the session: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (genreChanged) {
			try {
				await patchSessionById(properties.session_id!, {
					genres: Array.from(document.querySelectorAll('.new-genre-checkbox:checked')).map(
						(i) => i.id.replace('new-genre-', '') as Genre
					)
				});
			} catch (e) {
				success = false;
				alert(
					'An error occurred trying to suggest there is a different genre list: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (backlineChanged) {
			try {
				await patchVenueById(properties.venue!, {
					backline: Array.from(document.querySelectorAll('.new-backline-checkbox:checked')).map(
						(i) => i.id.replace('new-backline-', '') as Genre
					)
				});
			} catch (e) {
				success = false;
				alert(
					'An error occurred trying to suggest there is a different backline: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (otherChanges) {
			try {
				await postSuggestionForSessionById(properties.session_id!, {
					content: (document.getElementById('other-suggestions') as HTMLTextAreaElement).value
				});
			} catch (e) {
				success = false;
				alert(
					'An error occurred trying to submit a change suggestion: ' +
						(e as Error).message +
						'\nPlease get in touch at felix.schott@proton.me instead.'
				);
			}
		}
		if (success === true) {
			alert(
				'Thank you! A member of the admin team will review your suggestions and apply the changes. For other feedback you can also get in touch at felix.schott@proton.me.'
			);
		}
	};
</script>

<h2>{properties.session_name} &mdash; what has changed?</h2>
<form>
	<div class="card">
		<div class="horizontal">
			<b>The address of the venue has changed.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption === 'yes') {
						addressChanged = true;
					} else {
						addressChanged = false;
					}
				}}
			/>
		</div>
		<div class:vertical={addressChanged} class:hidden={!addressChanged}>
			<label for="new-address-first-line"
				>Address first line: <input type="text" id="new-address-first-line" required={addressChanged} /></label
			>
			<label for="new-address-second-line"
				>Address second line: <input type="text" id="new-address-second-line" /></label
			>
			<label for="new-address-city"
				>City: <input type="text" id="new-address-city" required={addressChanged} /></label
			>
			<label for="new-address-postcode"
				>Postcode: <input type="text" id="new-address-postcode" required={addressChanged} /></label
			>
		</div>
	</div>

	<div class="card">
		<div class="horizontal">
			<b>The venue has closed down.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption == 'yes') {
						venueClosed = true;
					} else {
						venueClosed = false;
					}
				}}
			/>
		</div>
	</div>
	<div class="card">
		<div class="horizontal">
			<b>The time of the session has changed.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption === 'yes') {
						timeChanged = true;
					} else {
						timeChanged = false;
					}
				}}
			/>
		</div>
		<div class:vertical={timeChanged} class:hidden={!timeChanged}>
			<label for="new-time-start">From: <input type="time" id="new-time-start" required={timeChanged} /></label>
			<label for="new-time-finish">To: <input type="time" id="new-time-finish" required={timeChanged} /></label>
		</div>
	</div>
	<div class="card">
		<div class="horizontal">
			<b>The session doesn't exist anymore.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption == 'yes') {
						sessionClosed = true;
					} else {
						sessionClosed = false;
					}
				}}
			/>
		</div>
	</div>
	<div class="card">
		<div class="horizontal">
			<b>The list of genres is incomplete or wrong.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption === 'yes') {
						genreChanged = true;
					} else {
						genreChanged = false;
					}
				}}
			/>
		</div>
		<div class:vertical={genreChanged} class:hidden={!genreChanged}>
			Please select all genres that apply:
			<div class="checkboxes">
				{#each Object.values(Genre) as genre}
					{#if genre != 'ANY'}
						<label for="new-genre-{genre}"
							><input
								type="checkbox"
								class="new-genre-checkbox"
								id="new-genre-{genre}"
							/>{genre.replace('_', ' ')}</label
						>
					{/if}
				{/each}
			</div>
		</div>
	</div>
	<div class="card">
		<div class="horizontal">
			<b>The backline information is incomplete/wrong.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption === 'yes') {
						backlineChanged = true;
					} else {
						backlineChanged = false;
					}
				}}
			/>
		</div>
		<div class:vertical={backlineChanged} class:hidden={!backlineChanged}>
			Please select the full backline provided by the venue:
			<div class="checkboxes">
				{#each Object.values(Backline) as backline}
					<label for="new-backline-{backline}"
						><input
							type="checkbox"
							class="new-backline-checkbox"
							id="new-backline-{backline}"
						/>{backline.replace('_', ' ')}</label
					>
				{/each}
			</div>
		</div>
	</div>
	<div class="card">
		<div class="horizontal">
			<b>There have been other changes.</b>
			<ButtonGroup
				activeIndex={1}
				options={['yes', 'no']}
				on:change={(ev) => {
					if (ev.detail.activeOption === 'yes') {
						otherChanges = true;
					} else {
						otherChanges = false;
					}
				}}
			/>
		</div>
		<div class:vertical={otherChanges} class:hidden={!otherChanges}>
			<textarea id="other-suggestions" cols="40" rows="5" required={otherChanges}></textarea>
		</div>
	</div>

	<div style="display: flex; justify-content: center; margin-top: 2em;">
		<button disabled={!buttonEnabled} on:click={onSubmit}
			><span>Suggest changes</span><br /><span>to our team</span></button
		>
	</div>
</form>

<style>
	button:disabled {
		background-color: whitesmoke;
	}

	button span:first-child {
		font-size: larger;
		font-weight: bolder;
	}

	.horizontal {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	@media (max-width: 480px) {
		.horizontal {
			flex-direction: column;
		}

		:global(.horizontal > div) {
			margin-top: 0.5em;
		}
	}

	.vertical {
		display: flex;
		flex-direction: column;
		margin-top: 1.5em;
		margin-bottom: 1em;
	}

	.vertical > label {
		display: flex;
		justify-content: space-between;
	}

	.vertical > label:not(:last-child) {
		margin-bottom: 1em;
	}
</style>
