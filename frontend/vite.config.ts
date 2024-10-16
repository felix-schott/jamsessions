import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	define: {
		"process.env.API_ADDRESS": JSON.stringify(process.env.API_ADDRESS),
		"process.env.TILES_ADDRESS": JSON.stringify(process.env.TILES_ADDRESS),
	},
	plugins: [sveltekit()]
});
