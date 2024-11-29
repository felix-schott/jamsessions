<script lang="ts">
	import {
		Interval,
		type CommentBody,
		type SessionComment,
		type SessionProperties,
		type SessionPropertiesWithVenue
	} from '../types';
	import LocationIcon from './icons/LocationIcon.svelte';
	import HomeIcon from './icons/HomeIcon.svelte';
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
	import { extractDomain, sanitisePathElement, processVenueAndSessionUrl } from './uriUtils';
	import InfoIcon from './icons/InfoIcon.svelte';
	import SelectRating from './SelectRating.svelte';
	import { editingSession } from '../stores';
	import { onMount } from 'svelte';

	// props
	interface Props {
		sessionProperties: SessionPropertiesWithVenue;
		sessionComments: SessionComment[];
	}
	let { sessionProperties, sessionComments }: Props = $props();

	// state management
	let newCommentHidden: boolean = $state(true);

	let newCommentContent: string = $state('');
	let newCommentAuthor: string = $state('');
	let newRating: number = $state(0);

	let sessionUrl: string = $state('');
	let venueUrl: string = $state('');

	let isMobile = $state(false);

	onMount(() => {
		isMobile = window.matchMedia('(max-width: 480px)').matches;

		let { venueWebsite, sessionWebsite } = processVenueAndSessionUrl(
			sessionProperties.venue_website,
			sessionProperties.session_website
		);
		venueUrl = venueWebsite;
		sessionUrl = sessionWebsite ? sessionWebsite : '';
	});

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
			let commentBody: CommentBody = {
				author: newCommentAuthor,
				content: newCommentContent
			};
			if (newRating !== 0) {
				commentBody['rating'] = newRating;
			}
			await postCommentForSessionById(sessionProperties.session_id!, commentBody);
			alert(
				'Thanks for submitting a comment! All content is moderated, so please bare with us while we review your comment. If there is anything else, get in touch at felix.schott@proton.me'
			);
			newCommentHidden = true;
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
				{#if isMobile}
					<!-- only works on mobile devices -->
					<ShareIcon
						style="cursor: pointer; margin-left: 0.3em; vertical-align: text-bottom;"
						height="1.1em"
						width="1.1em"
						title="Share link to session"
						onclick={onShare}
					/>
				{/if}
				<EditIcon
					style="cursor: pointer; margin-left: 0.3em; vertical-align: text-bottom;"
					title="Suggest changes to this page"
					height="1.1em"
					width="1.1em"
					onclick={() => {
						$editingSession = true;
					}}
				/></span
			>
		</h2>
	</div>
	<table class="session-info">
		<tbody>
			<tr>
				<td><LocationIcon title="Address of venue" class="icon-auto" /></td>
				<td>
					<a
						href={`/${sanitisePathElement(sessionProperties.venue_name)}-${sessionProperties.venue_id}`}
						>{sessionProperties?.venue_name}</a
					><br />
					{sessionProperties?.address_first_line}<br />
					{#if sessionProperties.address_second_line}
						{sessionProperties?.address_second_line}<br />
					{/if}
					{sessionProperties?.city}<br />
					{sessionProperties?.postcode}<br />
					<a
						target="_blank"
						href="https://www.google.com/maps/place/{sessionProperties.address_first_line.replaceAll(
							' ',
							'+'
						)},+{sessionProperties.city.replaceAll(
							' ',
							'+'
						)}+{sessionProperties.postcode.replaceAll(' ', '+')}/">View on Google Maps</a
					>
				</td>
			</tr>
			<tr>
				<td><HomeIcon title="Homepage" class="icon-auto" /></td>
				<td
					><a style="margin-top: 0.5em;" href={sessionProperties?.venue_website} target="_blank"
						>{venueUrl}</a
					>
					{#if sessionUrl !== ''}
						<br /><a
							style="margin-top: 0.5em;"
							href={sessionProperties?.session_website}
							target="_blank">{sessionUrl}</a
						>
					{/if}
				</td>
			</tr>
			<tr>
				<td><TimeIcon title="Time of event" class="icon-auto" /></td>
				<td>{@html constructTimeString(sessionProperties)}</td>
			</tr>
			{#if sessionProperties.genres && sessionProperties.genres.length !== 0}
				<tr>
					<td><FileTrayIcon title="Primary genre of event" class="icon-auto" /></td>
					<td>{sessionProperties?.genres.map((i) => i.replace('_', ' ')).join(', ')}</td>
				</tr>
			{/if}
			{#if sessionProperties.backline && sessionProperties.backline.length !== 0}
				<tr>
					<td><MicrophoneIcon title="Backline provided by venue" class="icon-auto" /></td>
					<td
						>{sessionProperties?.backline
							.slice(0, -1)
							.map((i) => i.replace('_', ' '))
							.join(', ') +
							' and ' +
							sessionProperties?.backline.slice(-1)[0].replace('_', ' ')}</td
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
		{#if sessionProperties?.interval === Interval.IRREGULARWEEKLY}
			IMPORTANT: This session doesn't operate on a regular schedule! If it takes place, it normally
			happens on this weekday but there is no guarantee for that.
		{:else}
			The data may be inaccurate or outdated, and sessions can be cancelled at short notice.
		{/if}
		Please always check the
		<a target="_blank" href={sessionProperties?.session_website}>website of the organiser</a>.
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
		{#each sessionComments as comment}
			<div class="comment">
				<Rating
					style="margin-right: 0.3em; padding-top: 0.3em; padding-right: 0.5em;"
					n={comment.rating ? comment.rating : 0}
				/>
				<div>
					{comment.content}
					<i> &mdash; {comment.author}, {new Date(comment.dt_posted).toLocaleDateString()}</i>
				</div>
			</div>
		{/each}
		<div
			class:horizontal-center={newCommentHidden}
			style="margin-top: 1em;"
			class:hidden={!newCommentHidden}
		>
			<button
				title="Add comment"
				style="display: flex; align-items: center; margin-bottom: 1em; padding: 0.3em 0.6em; font-size: smaller;"
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
				If you want to report inaccurate data, please <span
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
					style="font-size: smaller; margin-top: 1em; background-color: white; font-color: black;"
					onclick={onSubmitNewComment}>Submit</button
				>
			</div>
		</div>
	</div>
{:else}
	<EditSession properties={sessionProperties} />
{/if}

<style>
	.horizontal-center {
		display: flex;
		justify-content: center;
	}

	table.session-info {
		margin-left: 1em;
	}

	td {
		vertical-align: top;
	}

	table.session-info td {
		padding: 0.5em;
	}

	table.session-info td:first-child {
		padding-top: 0.5em;
	}

	.comments {
		background-color: whitesmoke;
		padding: 0.5em;
		border-radius: 10px;
		margin-top: 1em;
	}

	div.comment {
		list-style-type: none;
		margin: 0.5em;
		margin-left: 1em;
		margin-bottom: 1em;
		display: flex;
		flex-direction: row;
	}

	@media (max-width: 480px) {
		div.comment {
			flex-direction: column;
		}
	}

	.line {
		display: inline-block;
	}
</style>
