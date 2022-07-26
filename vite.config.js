import path from "path";
import glob from "glob";
import { defineConfig } from 'vite'

console.log(glob.sync(path.resolve(__dirname, "ui", "*/index.html")))
export default defineConfig(({ command, mode }) => ({
	base: "/wotlk/",
	root: path.join(__dirname, "ui"),
	build: {
		outDir: path.join(__dirname, "dist", "wotlk"),
		emptyOutDir: true,
		minify: mode === "development" ? false : "terser",
		sourcemap: command === "serve" ? "inline" : "false",
		target: ["es2020"],
		rollupOptions: {
			input: glob.sync(path.resolve(__dirname, "ui", "**/index.html")),
		},
		server: {
			origin: 'http://localhost:3000',
		},
	}
}));
