import path from "path";
import glob from "glob";
import { defineConfig } from 'vite'

export default defineConfig(({ command, mode }) => ({
	base: "/wotlk/",
	root: path.join(__dirname, "ui"),
	build: {
		outDir: path.join(__dirname, "dist", "wotlk"),
		minify: mode === "development" ? false : "terser",
		sourcemap: command === "serve" ? "inline" : "false",
		target: ["es2020"],
		rollupOptions: {
			input: glob.sync(path.resolve(__dirname, "ui", "**/index.html").replace(/\\/g, "/")),
			output: {
				assetFileNames: () => "bundle/[name]-[hash].style.css",
				entryFileNames: () => "bundle/[name]-[hash].entry.js",
				chunkFileNames: () => "bundle/[name]-[hash].chunk.js",
			},
		},
		server: {
			origin: 'http://localhost:3000',
		},
	}
}));
