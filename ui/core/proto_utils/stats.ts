import { Class, Stat, PseudoStat, UnitStats } from '../proto/common.js';
import { getEnumValues } from '../utils.js';
import { getClassStatName, pseudoStatNames } from './names.js';

const STATS_LEN = getEnumValues(Stat).length;
const PSEUDOSTATS_LEN = getEnumValues(PseudoStat).length;

export class UnitStat {
	private readonly stat: Stat | null;
	private readonly pseudoStat: PseudoStat | null;

	private constructor(stat: Stat | null, pseudoStat: PseudoStat | null) {
		this.stat = stat;
		this.pseudoStat = pseudoStat;
	}

	isStat(): boolean {
		return this.stat != null;
	}
	isPseudoStat(): boolean {
		return this.pseudoStat != null;
	}

	getStat(): Stat {
		if (!this.isStat()) {
			throw new Error('Not a stat!');
		}
		return this.stat!;
	}
	getPseudoStat(): PseudoStat {
		if (!this.isPseudoStat()) {
			throw new Error('Not a pseudo stat!');
		}
		return this.pseudoStat!;
	}

	equals(other: UnitStat): boolean {
		return this.stat == other.stat && this.pseudoStat == other.pseudoStat;
	}

	getName(clazz: Class): string {
		if (this.isStat()) {
			return getClassStatName(this.stat!, clazz);
		} else {
			return pseudoStatNames.get(this.pseudoStat!)!;
		}
	}

	getProtoValue(proto: UnitStats): number {
		if (this.isStat()) {
			return proto.stats[this.stat!];
		} else {
			return proto.pseudoStats[this.pseudoStat!];
		}
	}

	setProtoValue(proto: UnitStats, val: number) {
		if (this.isStat()) {
			proto.stats[this.stat!] = val;
		} else {
			proto.pseudoStats[this.pseudoStat!] = val;
		}
	}

	static fromStat(stat: Stat): UnitStat {
		return new UnitStat(stat, null);
	}
	static fromPseudoStat(pseudoStat: PseudoStat): UnitStat {
		return new UnitStat(null, pseudoStat);
	}

	static getAll(): Array<UnitStat> {
		const allStats = (getEnumValues(Stat) as Array<Stat>).filter(stat => ![Stat.StatEnergy, Stat.StatRage].includes(stat));
		const allPseudoStats = getEnumValues(PseudoStat) as Array<PseudoStat>;
		return [
			allStats.map(stat => UnitStat.fromStat(stat)),
			allPseudoStats.map(stat => UnitStat.fromPseudoStat(stat)),
		].flat();
	}
}

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
	getUnitStat(stat: UnitStat): number {
		if (stat.isStat()) {
			return this.stats[stat.getStat()];
		} else {
			return this.pseudoStats[stat.getPseudoStat()];
		}
	}

	withStat(stat: Stat, value: number): Stats {
		const newStats = this.stats.slice();
		newStats[stat] = value;
		return new Stats(newStats, this.pseudoStats);
	}
	withPseudoStat(stat: PseudoStat, value: number): Stats {
		const newStats = this.pseudoStats.slice();
		newStats[stat] = value;
		return new Stats(this.stats, newStats);
	}
	withUnitStat(stat: UnitStat, value: number): Stats {
		if (stat.isStat()) {
			return this.withStat(stat.getStat(), value);
		} else {
			return this.withPseudoStat(stat.getPseudoStat(), value);
		}
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

	scale(scalar: number): Stats {
		return new Stats(
			this.stats.map((value, stat) => value * scalar),
			this.pseudoStats.map((value, stat) => value * scalar));
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

	belowCaps(statCaps: Stats): boolean {
		for (const [idx, stat] of this.stats.entries()) {
			if ((statCaps.stats[idx] > 0) && (stat > statCaps.stats[idx])) {
				return false;
			}
		}

		return true;
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

	static fromMap(statsMap: Partial<Record<Stat, number>>, pseudoStatsMap?: Partial<Record<PseudoStat, number>>): Stats {
		const statsArr = new Array(STATS_LEN).fill(0);
		Object.entries(statsMap).forEach(entry => {
			const [statStr, value] = entry;
			statsArr[Number(statStr)] = value;
		});

		const pseudoStatsArr = new Array(PSEUDOSTATS_LEN).fill(0);
		if (pseudoStatsMap) {
			Object.entries(pseudoStatsMap).forEach(entry => {
				const [pseudoStatstr, value] = entry;
				pseudoStatsArr[Number(pseudoStatstr)] = value;
			});
		}

		return new Stats(statsArr, pseudoStatsArr);
	}

	static fromProto(unitStats?: UnitStats): Stats {
		if (unitStats) {
			return new Stats(unitStats.stats, unitStats.pseudoStats);
		} else {
			return new Stats();
		}
	}
}
