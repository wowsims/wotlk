module.exports = {
	root: true,
	parser: '@typescript-eslint/parser',
	plugins: ['simple-import-sort'],
	extends: [
		'plugin:json/recommended',
		'plugin:import/errors',
		'plugin:import/warnings',
		'plugin:import/typescript',
		'plugin:@typescript-eslint/eslint-recommended',
		'plugin:@typescript-eslint/recommended',
		'plugin:prettier/recommended',
	],
	env: {
		es6: true,
		browser: true,
	},
	parserOptions: {
		ecmaVersion: 2021,
		sourceType: 'module',
		ecmaFeatures: {
			jsx: true,
		},
	},
	rules: {
		'@typescript-eslint/member-delimiter-style': 'off',
		'@typescript-eslint/explicit-function-return-type': 'off',
		'@typescript-eslint/explicit-module-boundary-types': 'off',
		'@typescript-eslint/no-non-null-assertion': 'off',
		'@typescript-eslint/no-explicit-any': 'off',
		'@typescript-eslint/no-use-before-define': 'off',
		'@typescript-eslint/indent': 'off',
		'@typescript-eslint/no-unused-vars': [
			'warn',
			{
				argsIgnorePattern: '^_',
				varsIgnorePattern: '^_',
			},
		],
		'@typescript-eslint/no-object-literal-type-assertion': 'off',
		'@typescript-eslint/explicit-member-accessibility': 'off',
		'@typescript-eslint/camelcase': 'off',
		'@typescript-eslint/no-empty-interface': 'off',
		'@typescript-eslint/ban-ts-comment': 'off',
		'prettier/prettier': 'off',
		'import/no-unresolved': 'off',
		'simple-import-sort/imports': 'warn',
		'import/named': 'off',
		'import/namespace': 'off',
		'arrow-parens': ['error', 'as-needed'],
	}
};
