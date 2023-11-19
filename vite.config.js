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
			input: glob.sync(path.resolve(__dirname, "ui", "shadow_priest/index.html").replace(/\\/g, "/")),
			output: {
				assetFileNames: () => "bundle/[name]-[hash].style.css",
				entryFileNames: () => "bundle/[name]-[hash].entry.js",
				chunkFileNames: () => "bundle/[name]-[hash].chunk.js",
			},
			external: (id) => {
				return [
					"ui/deathknight",
					"ui/elemental_shaman",
					"ui/enhancement_shaman",
					"ui/feral_druid",
					"ui/feral_tank_druid",
					"ui/healing_priest",
					"ui/holy_paladin",
					"ui/hunter",
					"ui/mage",
					"ui/protection_paladin",
					"ui/protection_warrior",
					"ui/restoration_druid",
					"ui/restoration_shaman",
					"ui/retribution_paladin",
					"ui/rogue",
					"ui/tank_deathknight",
					"ui/warlock",
					"ui/warrior"
				].some(path => id.includes(path));
			},
		},
		server: {
			origin: 'http://localhost:3000',
		},
	}
}));
