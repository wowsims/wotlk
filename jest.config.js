module.exports = {
	collectCoverage: true,
	roots: [
		'<rootDir>',
	],
	testRegex: 'ui/tests/.*(feature|test|spec).tsx?$',
	moduleFileExtensions: [
		'js',
		'json',
		'mjs',
		'ts',
		'tsx',
	],
	moduleDirectories: [
		'node_modules',
	],
	testEnvironment: 'jsdom',
	preset: 'ts-jest',
}
