<script lang="ts">
	import Modal from './Modal.svelte';
	import type { SessionComment, SessionProperties, VenueProperties } from '../types';
	import LocationIcon from './icons/LocationIcon.svelte';
	import ShareIcon from './icons/ShareIcon.svelte';
	import TimeIcon from './icons/TimeIcon.svelte';
	import FileTrayIcon from './icons/FileTrayIcon.svelte';
	import MicrophoneIcon from './icons/MicrophoneIcon.svelte';
	import EditIcon from './icons/EditIcon.svelte';
	import Rating from './Rating.svelte';
	import EditSession from './EditSessionPopup.svelte';
	import PlusIcon from './icons/PlusIcon.svelte';
	import { postCommentForSessionById } from '../api';
	import { constructTimeString } from './timeUtils';
	import InfoIcon from './icons/InfoIcon.svelte';
	import SelectRating from './SelectRating.svelte';
	import { editingSession } from '../stores';

	// props
	export let sessionProperties: SessionProperties;
	export let venueProperties: VenueProperties;
	export let sessionComments: SessionComment[];

	// state management
	let newCommentHidden: boolean = true;

	let newCommentContent: string;
	let newCommentAuthor: string;
	let newRating: number = 0;

	// share functionality
	const onShare = async () => {
		const shareData = {
			title: sessionProperties.session_name,
			text: 'Check out this jam session',
			url: window.location.href
		};
		try {
			await navigator.share(shareData);
		} catch (err) {
			console.log('Could not share:', err);
		}
	};

	// event handler
	const onSubmitNewComment = async () => {
		try {
			await postCommentForSessionById(sessionProperties.session_id!, {
				author: newCommentAuthor,
				content: newCommentContent
			});
			alert(
				'Thanks for submitting a comment! All content is moderated, so please bare with us while we review your comment. If there is anything else, get in touch at felix.schott@proton.me'
			);
			$editingSession = false;
		} catch (e) {
			alert('An error occured when trying to post a new comment: ' + (e as Error).message);
			throw e;
		}
	};
</script>

{#if !$editingSession}
	<div>
		<h2>
			<!-- https://stackoverflow.com/a/24357132 -->
			<span class="line">{sessionProperties?.session_name} </span>
			<span class="line">
				<ShareIcon
					style="cursor: pointer; margin-left: 0.3em; vertical-align: middle;"
					title="Share link to session"
					onclick={onShare}
				/>
				<EditIcon
					style="cursor: pointer; margin-left: 0.3em; vertical-align: middle;"
					title="Suggest changes to this page"
					onclick={() => {
						$editingSession = true;
					}}
				/></span
			>
		</h2>
	</div>
	<table>
		<tbody>
			<tr>
				<td><LocationIcon title="Address of venue" class="icon-auto" /></td>
				<td>
					<a href={venueProperties?.venue_website} target="_blank">{venueProperties?.venue_name}</a
					><br />
					{venueProperties?.address_first_line}<br />
					{#if venueProperties.address_second_line}
						{venueProperties?.address_second_line}<br />
					{/if}
					{venueProperties?.city}<br />
					{venueProperties?.postcode}<br />
					<a
						target="_blank"
						href="https://www.google.com/maps/place/{venueProperties.address_first_line.replaceAll(
							' ',
							'+'
						)},+{venueProperties.city.replaceAll(' ', '+')}+{venueProperties.postcode.replaceAll(
							' ',
							'+'
						)}/">View on Google Maps</a
					>
				</td>
			</tr>
			<tr>
				<td><TimeIcon title="Time of event" class="icon-auto" /></td>
				<td>{constructTimeString(sessionProperties)}</td>
			</tr>
			{#if sessionProperties.genres && sessionProperties.genres.length !== 0}
				<tr>
					<td><FileTrayIcon title="Primary genre of event" class="icon-auto" /></td>
					<td>{sessionProperties?.genres.map((i) => i.replace('_', ' ')).join(', ')}</td>
				</tr>
			{/if}
			{#if venueProperties.backline && venueProperties.backline.length !== 0}
				<tr>
					<td><MicrophoneIcon title="Backline provided by venue" class="icon-auto" /></td>
					<td
						>{venueProperties?.backline
							.slice(0, -1)
							.map((i) => i.replace('_', ' '))
							.join(', ') +
							' and ' +
							venueProperties?.backline.slice(-1)[0].replace('_', ' ')}</td
					>
				</tr>
			{/if}
		</tbody>
	</table>
	<p>
		{sessionProperties?.description}
	</p>
	<p
		style="border-radius: 6px; padding: 0.5em; background-color: var(--accent-color); margin-bottom: 2em;"
	>
		The data may be inaccurate or outdated, and sessions can be cancelled at short notice. Please
		always check the <a target="_blank" href={sessionProperties?.session_website}
			>website of the organiser</a
		>.
	</p>
	<hr />
	<div style="display: flex; align-items: center; margin-top: 0.3em;">
		<b style="font-size: larger;">Community voices</b>
		<Rating
			size="1.2em"
			style="margin-left: 0.5em;"
			n={sessionProperties.rating ? sessionProperties.rating : 0}
		/>
	</div>
	<div class="comments">
		<ul>
			{#each sessionComments as comment}
				<li>
					<span
						><Rating style="margin-right: 0.3em;" n={comment.rating ? comment.rating : 0} /></span
					>
					<span>{comment.content}</span>
					<span
						><i>
							&mdash; {comment.author}, {new Date(comment.dt_posted).toLocaleDateString()}</i
						></span
					>
				</li>
			{/each}
			<div
				class:horizontal-center={newCommentHidden}
				style="margin-top: 1em;"
				class:hidden={!newCommentHidden}
			>
				<button
					title="Add comment"
					style="display: flex; align-items: center; padding: 0.3em 0.6em; font-size: smaller;"
					onclick={() => {
						newCommentHidden = false;
					}}
					><PlusIcon
						title="Add comment"
						style="cursor: pointer; margin-right: 0.3em;"
						class="icon-auto"
					/> Add comment</button
				>
			</div>
			<div
				style="background-color: lightgrey; border-radius: 10px; padding: 0.5em 1em 1em; margin-top: 1em;"
				class:hidden={newCommentHidden}
			>
				<span style="display: inline-flex; margin-top: 0.7em;"><b>Add comment</b></span>
				<span
					title="Close new comment section"
					role="button"
					style="cursor: pointer; float: right; color: red;"
					onclick={() => {
						newCommentHidden = true;
					}}>×</span
				>
				<p style="font-size: smaller;">
					If you want to report inaccurate data instead, please <span
						onclick={() => {
							$editingSession = true;
						}}
						style="color: #646cff; cursor: pointer;
	">click here</span
					> instead.
				</p>
				<textarea
					placeholder="Describe your experience at the jam session and provide useful information for others."
					bind:value={newCommentContent}
					id="new-comment"
					style="width: 100%; height: 3em; margin-top: 1em;"
				></textarea>
				<div style="margin-top: 0.5em; display: flex; align-items: center;">
					Rate your experience: <SelectRating
						style="margin-left: 0.3em;"
						onchange={(rating) => {
							newRating = rating;
						}}
					/>
				</div>
				<label style="display: flex; align-items: center; margin-top: 0.5em;"
					>Your name: <InfoIcon
						style="margin-right: 0.3em;"
						title="This is the name that appears next to your comment - it can be your first name, your nickname or however you wish to present."
					/><input bind:value={newCommentAuthor} id="new-comment-author" type="text" /></label
				>
				<div class="horizontal-center">
					<button
						style="font-size: smaller; margin-top: 0.5em; background-color: white; font-color: black;"
						onclick={onSubmitNewComment}>Submit</button
					>
				</div>
			</div>
		</ul>
	</div>
{:else}
	<EditSession properties={sessionProperties} />
{/if}

<style>
	.horizontal-center {
		display: flex;
		justify-content: center;
	}

	table {
		margin-left: 1em;
	}

	td {
		vertical-align: top;
		padding: 0.5em;
	}

	td:first-child {
		padding-top: 0.5em;
	}

	.comments {
		background-color: whitesmoke;
		padding: 0.5em;
		border-radius: 10px;
		margin-top: 1em;
	}

	ul {
		padding-inline-start: 0;
		margin-block-start: 0;
	}

	li {
		list-style-type: none;
		margin: 0.5em;
		margin-left: 1em;
		display: flex;
		align-items: center;
	}

	@media (max-width: 480px) {
		li {
			flex-direction: column;
			align-items: unset;
		}

		li span:not(:first-child) {
			margin-top: 0.2em;
		}
	}

	.line {
		display: inline-block;
	}
</style>
