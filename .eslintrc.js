module.exports = {
	"parser": "@typescript-eslint/parser",
	"plugins": [
		"@typescript-eslint",
		"unused-imports"
	],
	"env": {
		"browser": true,
		"es6": true,
		"node": true
	},
	"globals": {
		"Atomics": "readonly",
		"SharedArrayBuffer": "readonly"
	},
	"parserOptions": {
		"ecmaVersion": 2018,
		"sourceType": "module"
	},
	"rules": {
		"@typescript-eslint/no-unused-vars": "off",
		"unused-imports/no-unused-imports": "error",
		"unused-imports/no-unused-vars": [
			"error",
			{ "vars": "all", "varsIgnorePattern": "^_", "args": "after-used", "argsIgnorePattern": "^_" }
		],
	}
};
