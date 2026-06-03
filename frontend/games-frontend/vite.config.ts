import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		strictPort: true,
		port: 5174,
		host: '0.0.0.0',
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:8080',
				changeOrigin: true,
			},
		},
	},
});
