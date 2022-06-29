import { Stat } from '/tbc/core/proto/common.js';
import { getEnumValues } from '/tbc/core/utils.js';

const STATS_LEN = getEnumValues(Stat).length;

/**
 * Represents values for all character stats (stam, agi, spell power, hit raiting, etc).
 *
 * This is an immutable type.
 */
export class Stats {
	private readonly stats: Array<number>;

	constructor(stats?: Array<number>) {
		this.stats = stats?.slice(0, STATS_LEN) || [];

		if (this.stats.length < STATS_LEN) {
			this.stats = this.stats.concat(new Array(STATS_LEN - (stats?.length || 0)).fill(0));
		}

		for (let i = 0; i < STATS_LEN; i++) {
			if (this.stats[i] == null)
				this.stats[i] = 0;
		}
	}

	equals(other: Stats): boolean {
		return this.stats.every((newStat, statIdx) => newStat == other.getStat(statIdx));
	}

	getStat(stat: Stat): number {
		return this.stats[stat];
	}

	withStat(stat: Stat, value: number): Stats {
		const newStats = this.stats.slice();
		newStats[stat] = value;
		return new Stats(newStats);
	}

	addStat(stat: Stat, value: number): Stats {
		return this.withStat(stat, this.getStat(stat) + value);
	}

	add(other: Stats): Stats {
		return new Stats(this.stats.map((value, stat) => value + other.stats[stat]));
	}

	subtract(other: Stats): Stats {
		return new Stats(this.stats.map((value, stat) => value - other.stats[stat]));
	}

	computeEP(epWeights: Stats): number {
		let total = 0;
		this.stats.forEach((stat, idx) => {
			total += stat * epWeights.stats[idx];
		});
		return total;
	}

	asArray(): Array<number> {
		return this.stats.slice();
	}

	toJson(): Object {
		return this.asArray();
	}

	static fromJson(obj: any): Stats {
		return new Stats(obj as Array<number>);
	}

	static fromMap(statsMap: Partial<Record<Stat, number>>): Stats {
		const statsArr = new Array(STATS_LEN).fill(0);
		Object.entries(statsMap).forEach(entry => {
			const [statStr, value] = entry;
			statsArr[Number(statStr)] = value;
		});
		return new Stats(statsArr);
	}
}
