<script lang="ts">
	import Modal from './Modal.svelte';
	import { infoVisible } from '../stores';
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
		allows you to only display jam sessions on the selected date that primarily feature a certain
		genre or where the venue provides a microphone and a PA.
	</p>
	<p>
		Sessions are symbolised by a black musical note icon on the map. Click on a note to see a little
		more info about it, the name of the session and venue as well as the time of the event. If you
		click on little popup next to the note, a bigger popup will show with even more information
		about the event. Here you'll find explained what genre is primarily played at the session, what
		backline is provided by the event and how often the event happens.
	</p>
	<p>&copy; <a href="http://felixschott.net">Felix Schott</a></p></Modal
>
