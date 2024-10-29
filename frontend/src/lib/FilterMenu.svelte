<script lang="ts">
	import { loading, filterMenuVisible, selectedSessions, visibleLayer } from '../stores';
	import Modal from './Modal.svelte';
	import { Backline, Genre, type SessionFeatureCollection, MapLayer } from '../types';
	import { getSessions } from '../api';
	import MultiSelect from './MultiSelect.svelte';
	import MicrophoneIcon from './icons/MicrophoneIcon.svelte';
	import FileTrayIcon from './icons/FileTrayIcon.svelte';
	import ResetIcon from './icons/ResetIcon.svelte';

	let selectedGenre: string;
	let selectedBackline: string[] = [];

	let onChangeBackline = () => {
		selectedBackline = Array.from(document.querySelectorAll('#backline-select > option:checked'))
			.filter((i) => i)
			.map((i) => (i as HTMLOptionElement).value as Backline);
		window.sessionStorage.setItem('selectedBackline', selectedBackline.join(','));
	};

	let onSubmit = async () => {
		$filterMenuVisible = false;
		$loading = true;
		try {
			$selectedSessions = (await getSessions({
				date: new Date(window.sessionStorage.getItem('selectedDateStr')!),
				backline: window.sessionStorage.getItem('selectedBackline')?.split(',') as Backline[],
				genre: window.sessionStorage.getItem('selectedGenre') as Genre
			})) as SessionFeatureCollection;
		} catch (e) {
			alert('An error occured when waiting for data from the server: ' + (e as Error).message);
			throw e;
		}
		$visibleLayer = MapLayer.SESSIONS;
		$loading = false;
	};

	let onReset = () => {
		selectedGenre = Genre.ANY;
		selectedBackline = [];
		window.sessionStorage.setItem('selectedBackline', '');
		window.sessionStorage.setItem('selectedGenre', Genre.ANY);
	};
</script>

<Modal
	isVisible={() => $filterMenuVisible}
	hide={() => {
		$filterMenuVisible = false;
	}}
>
	<table>
		<tbody>
			<tr>
				<td><FileTrayIcon title="Select genre" class="icon-auto" /></td>
				<td
					><select
						title="Select genre"
						name="genre"
						bind:value={selectedGenre}
						onchange={() => {
							window.sessionStorage.setItem('selectedGenre', selectedGenre);
						}}
					>
						{#each Object.values(Genre) as opt}
							{#if window.sessionStorage.getItem('selectedGenre') === opt}
								<option value={opt} selected>{opt.replace('_', ' ')}</option>
							{:else}
								<option value={opt}>{opt.replace('_', ' ')}</option>
							{/if}
						{/each}
					</select></td
				>
			</tr>
			<tr>
				<td><MicrophoneIcon title="Select backline provided by venue" /></td>
				<td>
					<MultiSelect
						title="Select backline provided by genre"
						id="backline-select"
						name="backline"
						onchange={onChangeBackline}
					>
						{#each Object.values(Backline) as opt}
							{#if selectedBackline.includes(opt)}
								<option value={opt} selected>{opt.replace('_', ' ')}</option>
							{:else}
								<option value={opt}>{opt.replace('_', ' ')}</option>
							{/if}
						{/each}
					</MultiSelect>
				</td>
			</tr>
		</tbody>
	</table>
	<div
		style="display: flex; width: 100%; height: 100%; justify-content: space-between; align-items: center; margin-top: 1em;"
	>
		<ResetIcon title="Unset filters" onclick={onReset} style="cursor: pointer;" />
		<button title="Apply filters and load sessions" onclick={onSubmit}>Apply</button>
		<ResetIcon title="" style="visibility: hidden;" />
	</div>
</Modal>

<style>
	td {
		vertical-align: top;
		padding: 0.5em;
	}

	td:first-child {
		font-weight: bold;
	}
</style>
