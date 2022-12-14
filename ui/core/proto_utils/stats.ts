import { Stat, PseudoStat, UnitStats } from '../proto/common.js';
import { getEnumValues } from '../utils.js';

const STATS_LEN = getEnumValues(Stat).length;
const PSEUDOSTATS_LEN = getEnumValues(PseudoStat).length;

/**
 * Represents values for all character stats (stam, agi, spell power, hit raiting, etc).
 *
 * This is an immutable type.
 */
export class Stats {
	private readonly stats: Array<number>;
	private readonly pseudoStats: Array<number>;

	constructor(stats?: Array<number>, pseudoStats?: Array<number>) {
		this.stats = Stats.initStatsArray(STATS_LEN, stats);
		this.pseudoStats = Stats.initStatsArray(PSEUDOSTATS_LEN, pseudoStats);
	}

	private static initStatsArray(expectedLen: number, newStats?: Array<number>): Array<number> {
		let stats = newStats?.slice(0, expectedLen) || [];

		if (stats.length < expectedLen) {
			stats = stats.concat(new Array(expectedLen - (newStats?.length || 0)).fill(0));
		}

		for (let i = 0; i < expectedLen; i++) {
			if (stats[i] == null)
				stats[i] = 0;
		}
		return stats;
	}

	equals(other: Stats): boolean {
		return this.stats.every((newStat, statIdx) => newStat == other.getStat(statIdx)) &&
				this.pseudoStats.every((newStat, statIdx) => newStat == other.getPseudoStat(statIdx))
	}

	getStat(stat: Stat): number {
		return this.stats[stat];
	}
	getPseudoStat(stat: PseudoStat): number {
		return this.pseudoStats[stat];
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
		return new Stats(
			this.stats.map((value, stat) => value + other.stats[stat]),
			this.pseudoStats.map((value, stat) => value + other.pseudoStats[stat]));
	}

	subtract(other: Stats): Stats {
		return new Stats(
			this.stats.map((value, stat) => value - other.stats[stat]),
			this.pseudoStats.map((value, stat) => value - other.pseudoStats[stat]));
	}

	computeEP(epWeights: Stats): number {
		let total = 0;
		this.stats.forEach((stat, idx) => {
			total += stat * epWeights.stats[idx];
		});
		this.pseudoStats.forEach((stat, idx) => {
			total += stat * epWeights.pseudoStats[idx];
		});
		return total;
	}

	asArray(): Array<number> {
		return this.stats.slice();
	}

	toJson(): Object {
		return UnitStats.toJson(this.toProto()) as Object;
	}

	toProto(): UnitStats {
		return UnitStats.create({
			stats: this.stats.slice(),
			pseudoStats: this.pseudoStats.slice(),
		});
	}

	static fromJson(obj: any): Stats {
		return Stats.fromProto(UnitStats.fromJson(obj));
	}

	static fromMap(statsMap: Partial<Record<Stat, number>>): Stats {
		const statsArr = new Array(STATS_LEN).fill(0);
		Object.entries(statsMap).forEach(entry => {
			const [statStr, value] = entry;
			statsArr[Number(statStr)] = value;
		});
		return new Stats(statsArr);
	}

	static fromProto(unitStats?: UnitStats): Stats {
		if (unitStats) {
			return new Stats(unitStats.stats, unitStats.pseudoStats);
		} else {
			return new Stats();
		}
	}
}
