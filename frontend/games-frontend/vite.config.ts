import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	server: {
		strictPort: true,
		port: 5174,
		host: '0.0.0.0',
		hmr: {
			host: '192.168.1.18',
			protocol: 'ws',
		},
		proxy: {
			'/images': {
				target: 'http://192.168.1.18:8085',
				changeOrigin: true,
			},
		},
	},
});
