import {
	IconData,
	UIDatabase,
} from '../proto/ui.js';

//const dbUrl = '/wotlk/assets/database/db.json';
const dbUrl = '/wotlk/assets/database/db.bin';

export class Database {
	private static loadPromise: Promise<Database>|null = null;
	static get(): Promise<Database> {
		if (Database.loadPromise == null) {
			Database.loadPromise = fetch(dbUrl)
				// For reading JSON db.
				//.then(response => response.json())
				//.then(json => new Database(UIDatabase.fromJson(json)));
				.then(response => response.arrayBuffer())
				.then(buffer => new Database(UIDatabase.fromBinary(new Uint8Array(buffer))));
		}
		return Database.loadPromise;
	}

	private readonly db: UIDatabase;
	private readonly itemIcons: Record<number, IconData>;
	private readonly spellIcons: Record<number, IconData>;

	private constructor(db: UIDatabase) {
		this.db = db;

		this.itemIcons = {};
		this.spellIcons = {};
		db.itemIcons.forEach(data => this.itemIcons[data.id] = data);
		db.spellIcons.forEach(data => this.spellIcons[data.id] = data);
	}

	static async getItemIconData(itemId: number): Promise<IconData> {
		const db = await Database.get();
		return db.itemIcons[itemId] || IconData.create();
	}

	static async getSpellIconData(spellId: number): Promise<IconData> {
		const db = await Database.get();
		return db.spellIcons[spellId] || IconData.create();
	}

	//private static async getWowheadTooltipDataHelper(id: number, tooltipPostfix: string, cache: Map<number, Promise<any>>): Promise<any> {
	//	if (!cache.has(id)) {
	//		const url = `https://wowhead.com/wotlk/tooltip/${tooltipPostfix}/${id}`;
	//		try {
	//			const response = await fetch(url);
	//			cache.set(id, response.json());
	//		} catch (e) {
	//			console.error('Error while fetching url: ' + url + '\n\n' + e);
	//			cache.set(id, Promise.resolve({
	//				name: '',
	//				icon: '',
	//				tooltip: '',
	//			}));
	//		}
	//	}

	//	return cache.get(id) as Promise<any>;
	//}
}
