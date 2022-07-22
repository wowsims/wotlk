module.exports = {
	collectCoverage: true,
	roots: ["<rootDir>"],
	transform: { "^.+\\.tsx?$": "ts-jest" },
	testRegex: "ui/tests/.*(feature|test|spec).tsx?$",
	moduleFileExtensions: ["ts", "tsx", "js", "mjs", "json"],
	moduleDirectories: ["node_modules"],
	moduleNameMapper: {
		"^/protobuf-ts$": ["<rootDir>/node_modules/@protobuf-ts/runtime"],
		"^/wotlk/(.*)$": ["<rootDir>/ui/$1"],
	},
	testEnvironment: "node",
	globals: {
		"ts-jest": {
			tsConfig: "ui/tsconfig-base.json",
		},
	},
};
