<script lang="ts">
	import Modal from './Modal.svelte';
	import { infoVisible } from '../stores';
	import GithubIcon from './icons/GithubIcon.svelte';
</script>

<Modal
	isVisible={() => $infoVisible}
	hide={() => {
		$infoVisible = false;
	}}
>
	<h2>About</h2>
	<p>
		<i>Jam Sessions in London</i> is a free website that shows public jam sessions happening in London.
	</p>
	<h2>How it works</h2>
	<p>
		By default, sessions for the current date are shown on the map&mdash;you can easily change that
		by clicking on the date in the title bar of the website.
	</p>
	<p>
		Next to the date button, you'll find a button that opens a menu for you to adjust filters. This
		allows you to show sessions for multiple days at once, only display jam sessions on the selected
		date (or in the selected date range, respectively) that primarily feature a certain genre or
		where the venue provides a microphone and a PA. Not that this would also exclude sessions where
		that information is not known.
	</p>
	<p>To switch between the map view and a list view, click on the map/list button on the left.</p>
	<p>
		Sessions are symbolised by a black musical note icon on the map. Click on a note to see a popup
		with little more info about it, the name of the session and venue as well as the time of the
		event. If you click on that popup, a window will open with even more information about the
		event. Here you'll find what genre is primarily played at the session, what backline is provided
		by the event and how often the event happens.
	</p>
	<div style="display: flex; align-items: center; justify-content: center;">
		&copy;&nbsp;<a href="http://felixschott.net">Felix Schott</a>
		<a href="https://github.com/felix-schott/jamsessions" style="color: black;" target="_blank"
			><GithubIcon
				height="1.2em"
				width="1.2em"
				style="margin-left: 0.5em; vertical-align: text-bottom;"
				title="View source code on Github"
			/></a
		>
	</div></Modal
>
