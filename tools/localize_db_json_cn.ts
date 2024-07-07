/**
 * 编译后的中文db.json再通过wowhead进行本地化处理以避免过于复杂的修改源代码逻辑
 */

// ts-node ./tools/localize_db_json_cn.ts

const fs = require("fs");
const path = require("path");
const axios = require("axios");

type Tooltip = {
	name: string;
	quality: number;
	icon: string;
	tooltip: string;
	spell: any[];
};

type Database = {
	items: any[];
	enchants: any[];
	gems: any[];
	zones: any[];
	npcs: any[];
	itemIcons: any[];
	spellIcons: any[];
	encounters: any[];
	glyphIds: any[];
};

const url = "https://nether.wowhead.com/wotlk/cn/tooltip";

const processDatabase = async () => {

	const dbPath = "./assets/database/db_en.json";
	const content = fs.readFileSync(path.resolve(dbPath), "utf-8");
	const dbData = JSON.parse(content) as Database;

	if (!dbData) {
		console.log("Failed to read database file.");
		return;
	}

	const newDbData: Database = {
		items: [],
		enchants: [],
		gems: [],
		zones: [],
		npcs: [],
		itemIcons: [],
		spellIcons: [],
		encounters: [],
		glyphIds: [],
	};

	for (const key of Object.keys(dbData) as Array<keyof Database>) {
		let prefix = "";
		if (key === "items" || key === "gems" || key === "itemIcons") prefix = "item";
		if (key === "enchants" || key === "spellIcons") prefix = "spell";
		if (key === "zones" || key === "npcs") prefix = key.slice(0, -1);
		console.log(`Processing ${key}`);
		if (!prefix.length) {
			newDbData[key] = [...dbData[key]];
			continue;
		}
		const entries = dbData[key];
		let count = 0;
		for (const entry of entries) {
			count++;
			if (count % 500 === 0) console.log(`--- ${count}/${entries.length}`);
			const id = entry.id || entry.spellId;
			const response = await axios.get(`${url}/${prefix}/${id}`, {
				headers: {
					"Content-Type": "application/json",
				},
				responseType: "json",
			});
			const tooltip = response.data as Tooltip;
			if (tooltip && tooltip.name && tooltip.name.length !== 0) {
				newDbData[key].push({...entry, name: tooltip.name});
			} else {
				console.log(`failed at ${prefix}/${id}`);
				newDbData[key].push(entry);
			}

		}
	}

	// Write newDbData to a new file or update existing file
	const newDbPath = "./assets/database/db.json";
	fs.writeFileSync(newDbPath, JSON.stringify(newDbData, null, 2));
	console.log("Database localization completed.");
};

// Call the async function to start processing
processDatabase().catch(error => {
	console.error("Error processing database:", error);
});
