import { RaidSimRequest, RaidSimResult } from '/tbc/core/proto/api.js';
import { ResourceType } from '/tbc/core/proto/api.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { resourceNames, stringToResourceType } from '/tbc/core/proto_utils/names.js';
import { bucket, getEnumValues, stringComparator, sum } from '/tbc/core/utils.js';

export class Entity {
	readonly name: string;
	readonly ownerName: string; // Blank if not a pet.

	// Either target index, player index, or owner index depending on what kind
	// of entity this is.
	readonly index: number;

	readonly isTarget: boolean;
	readonly isPet: boolean;

	constructor(name: string, ownerName: string, index: number, isTarget: boolean, isPet: boolean) {
		this.name = name;
		this.ownerName = ownerName;
		this.index = index;
		this.isTarget = isTarget;
		this.isPet = isPet;
	}

	equals(other: Entity) {
		return this.isTarget == other.isTarget && this.isPet == other.isPet && this.index == other.index && this.name == other.name;
	}

	toString(): string {
		if (this.isTarget) {
			return 'Target ' + (this.index + 1);
		} else if (this.isPet) {
			return `${this.ownerName} (#${this.index + 1}) - ${this.name}`;
		} else {
			return `${this.name} (#${this.index + 1})`;
		}
	}

	// Parses one or more Entities from a string.
	// Each entity label should be one of:
	//   'Target 1' if a target,
	//   'PlayerName (#1)' if a player, or
	//   'PlayerName (#1) - PetName' if a pet.
	static parseRegex = /\[(Target (\d+))|(([a-zA-Z0-9]+) \(#(\d+)\) - ([a-zA-Z0-9\s]+))|(([a-zA-Z0-9]+) \(#(\d+)\))\]/g;
	static parseAll(str: string): Array<Entity> {
		return Array.from(str.matchAll(Entity.parseRegex)).map(match => {
			if (match[1]) {
				return new Entity(match[1], '', parseInt(match[2]) - 1, true, false);
			} else if (match[3]) {
				return new Entity(match[6], match[4], parseInt(match[5]) - 1, false, true);
			} else if (match[7]) {
				return new Entity(match[8], '', parseInt(match[9]) - 1, false, false);
			} else {
				throw new Error('Invalid Entity match');
			}
		});
	}
}

interface SimLogParams {
	raw: string,
	logIndex: number,
	timestamp: number,
	source: Entity | null,
	target: Entity | null,
	actionId: ActionId | null,
	threat: number,
}

export class SimLog {
	readonly raw: string;

	// Index of this log within the full log output.
	// When comparing timestamps this should be used instead of timestamp, because
	// timestamp is scraped from log text and doesn't have enough precision.
	readonly logIndex: number;

	// Time in seconds from the encounter start.
	readonly timestamp: number;

	readonly source: Entity | null;
	readonly target: Entity | null;
	readonly actionId: ActionId | null;

	// Amount of threat generated from this event. Note that not all events generate threat, so this will be 0.
	readonly threat: number;

	// Logs for auras that were active at this timestamp.
	// This is only filled if populateActiveAuras() is called.
	activeAuras: Array<AuraUptimeLog>;

	constructor(params: SimLogParams) {
		this.raw = params.raw;
		this.logIndex = params.logIndex;
		this.timestamp = params.timestamp;
		this.source = params.source;
		this.target = params.target;
		this.actionId = params.actionId;
		this.threat = params.threat;
		this.activeAuras = [];
	}

	toString(): string {
		return this.raw;
	}

	toStringPrefix(): string {
		const timestampStr = `[${this.timestamp.toFixed(2)}]`;
		if (this.source) {
			return `${timestampStr} [${this.source}]`;
		} else {
			return timestampStr;
		}
	}

	static async parseAll(result: RaidSimResult): Promise<Array<SimLog>> {
		const lines = result.logs.split('\n');

		return Promise.all(lines.map((line, lineIndex) => {
			const params: SimLogParams = {
				raw: line,
				logIndex: lineIndex,
				timestamp: 0,
				source: null,
				target: null,
				actionId: null,
				threat: 0,
			};

			const threatMatch = line.match(/ \(Threat: (-?[0-9]+\.[0-9]+)\)/);
			if (threatMatch) {
				params.threat = parseFloat(threatMatch[1]);
				line = line.substring(0, threatMatch.index);
			}

			let match = line.match(/\[([0-9]+\.[0-9]+)\]\w*(.*)/);
			if (!match || !match[1]) {
				return new SimLog(params);
			}

			params.timestamp = parseFloat(match[1]);
			let remainder = match[2];

			const entities = Entity.parseAll(remainder);
			params.source = entities[0] || null;
			params.target = entities[1] || null;

			// Order from most to least common to reduce number of checks.
			return DamageDealtLog.parse(params)
				|| ResourceChangedLog.parse(params)
				|| AuraEventLog.parse(params)
				|| AuraStacksChangeLog.parse(params)
				|| MajorCooldownUsedLog.parse(params)
				|| CastBeganLog.parse(params)
				|| CastCompletedLog.parse(params)
				|| StatChangeLog.parse(params)
				|| Promise.resolve(new SimLog(params));
		}));
	}

	isDamageDealt(): this is DamageDealtLog {
		return this instanceof DamageDealtLog;
	}

	isResourceChanged(): this is ResourceChangedLog {
		return this instanceof ResourceChangedLog;
	}

	isAuraEvent(): this is AuraEventLog {
		return this instanceof AuraEventLog;
	}

	isAuraStacksChange(): this is AuraStacksChangeLog {
		return this instanceof AuraStacksChangeLog;
	}

	isMajorCooldownUsed(): this is MajorCooldownUsedLog {
		return this instanceof MajorCooldownUsedLog;
	}

	isCastBegan(): this is CastBeganLog {
		return this instanceof CastBeganLog;
	}

	isCastCompleted(): this is CastCompletedLog {
		return this instanceof CastCompletedLog;
	}

	isStatChange(): this is StatChangeLog {
		return this instanceof StatChangeLog;
	}

	// Group events that happen at the same time.
	static groupDuplicateTimestamps<LogType extends SimLog>(logs: Array<LogType>): Array<Array<LogType>> {
		const grouped: Array<Array<LogType>> = [];
		let curGroup: Array<LogType> = [];

		logs.forEach(log => {
			if (curGroup.length == 0 || log.timestamp == curGroup[0].timestamp) {
				curGroup.push(log);
			} else {
				grouped.push(curGroup);
				curGroup = [log];
			}
		});
		if (curGroup.length > 0) {
			grouped.push(curGroup);
		}

		return grouped;
	}
}

export class DamageDealtLog extends SimLog {
	readonly amount: number;
	readonly miss: boolean;
	readonly hit: boolean;
	readonly crit: boolean;
	readonly crush: boolean;
	readonly glance: boolean;
	readonly dodge: boolean;
	readonly parry: boolean;
	readonly block: boolean;
	readonly tick: boolean;
	readonly partialResist1_4: boolean;
	readonly partialResist2_4: boolean;
	readonly partialResist3_4: boolean;

	constructor(params: SimLogParams, amount: number, miss: boolean, crit: boolean, crush: boolean, glance: boolean, dodge: boolean, parry: boolean, block: boolean, tick: boolean, partialResist1_4: boolean, partialResist2_4: boolean, partialResist3_4: boolean) {
		super(params);
		this.amount = amount;
		this.miss = miss;
		this.glance = glance;
		this.dodge = dodge;
		this.parry = parry;
		this.block = block;
		this.hit = !miss && !crit;
		this.crit = crit;
		this.crush = crush;
		this.tick = tick;
		this.partialResist1_4 = partialResist1_4;
		this.partialResist2_4 = partialResist2_4;
		this.partialResist3_4 = partialResist3_4;
	}

	resultString(): string {
		let result = this.miss ? 'Miss'
			: this.dodge ? 'Dodge'
				: this.parry ? 'Parry'
					: this.glance ? 'Glance'
						: this.block ? (this.crit ? 'Critical Block' : 'Block')
							: this.crit ? 'Crit'
								: this.crush ? 'Crush'
									: this.tick ? 'Tick'
										: 'Hit';
		result += ' ' + this.target;
		if (!this.miss && !this.dodge && !this.parry) {
			result += ` for ${this.amount.toFixed(2)}`;
			if (this.partialResist1_4) {
				result += ' (25% Resist)';
			} else if (this.partialResist2_4) {
				result += ' (50% Resist)';
			} else if (this.partialResist3_4) {
				result += ' (75% Resist)';
			}
			result += '.'
		}
		return result;
	}

	toString(): string {
		const threatPostfix = this.source?.isTarget ? '' : ` (${this.threat.toFixed(2)} Threat)`;
		return `${this.toStringPrefix()} ${this.actionId!.name} ${this.resultString()}${threatPostfix}`;
	}

	static parse(params: SimLogParams): Promise<DamageDealtLog> | null {
		const match = params.raw.match(/] (.*?) (tick )?((Miss)|(Hit)|(CriticalBlock)|(Crit)|(Crush)|(Glance)|(Dodge)|(Parry)|(Block))( \((\d+)% Resist\))?( for (\d+\.\d+) damage)?/);
		if (match) {
			return ActionId.fromLogString(match[1]).fill(params.source?.index).then(cause => {
				params.actionId = cause;

				let amount = 0;
				if (match[16]) {
					amount = parseFloat(match[16]);
				}

				return new DamageDealtLog(
					params,
					amount,
					match[3] == 'Miss',
					match[3] == 'Crit' || match[3] == 'CriticalBlock',
					match[3] == 'Crush',
					match[3] == 'Glance',
					match[3] == 'Dodge',
					match[3] == 'Parry',
					match[3] == 'Block' || match[3] == 'CriticalBlock',
					Boolean(match[2]) && match[2].includes('tick'),
					match[14] == '25',
					match[14] == '50',
					match[14] == '75');
			});
		} else {
			return null;
		}
	}
}

export class DpsLog extends SimLog {
	readonly dps: number;

	// Damage events that occurred at the same time as this log.
	readonly damageLogs: Array<DamageDealtLog>;

	constructor(params: SimLogParams, dps: number, damageLogs: Array<DamageDealtLog>) {
		super(params);
		this.dps = dps;
		this.damageLogs = damageLogs;
	}

	static DPS_WINDOW = 15; // Window over which to calculate DPS.
	static fromLogs(damageDealtLogs: Array<DamageDealtLog>): Array<DpsLog> {
		const groupedDamageLogs = SimLog.groupDuplicateTimestamps(damageDealtLogs);

		let curDamageLogs: Array<DamageDealtLog> = [];
		let curDamageTotal = 0;

		return groupedDamageLogs.map(ddLogGroup => {
			ddLogGroup.forEach(ddLog => {
				curDamageLogs.push(ddLog);
				curDamageTotal += ddLog.amount;
			});

			const newStartIdx = curDamageLogs.findIndex(curLog => {
				const inWindow = curLog.timestamp > ddLogGroup[0].timestamp - DpsLog.DPS_WINDOW;
				if (!inWindow) {
					curDamageTotal -= curLog.amount;
				}
				return inWindow;
			});
			if (newStartIdx == -1) {
				curDamageLogs = [];
			} else {
				curDamageLogs = curDamageLogs.slice(newStartIdx);
			}

			const dps = curDamageTotal / DpsLog.DPS_WINDOW;
			if (isNaN(dps)) {
				console.warn('NaN dps!');
			}

			return new DpsLog({
				raw: '',
				logIndex: ddLogGroup[0].logIndex,
				timestamp: ddLogGroup[0].timestamp,
				source: ddLogGroup[0].source,
				target: null,
				actionId: null,
				threat: 0,
			}, dps, ddLogGroup);
		});
	}
}

export class ThreatLogGroup extends SimLog {
	readonly threatBefore: number;
	readonly threatAfter: number;
	readonly logs: Array<SimLog>;

	constructor(params: SimLogParams, threatBefore: number, threatAfter: number, logs: Array<SimLog>) {
		super(params);
		this.threatBefore = threatBefore;
		this.threatAfter = threatAfter;
		this.logs = logs;
	}

	static fromLogs(logs: Array<SimLog>): Array<ThreatLogGroup> {
		const groupedLogs = SimLog.groupDuplicateTimestamps(logs.filter(log => log.threat != 0));
		let curThreat = 0;
		return groupedLogs.map(logGroup => {
			const newThreat = sum(logGroup.map(log => log.threat));
			const threatLog = new ThreatLogGroup(
				{
					raw: '',
					logIndex: logGroup[0].logIndex,
					timestamp: logGroup[0].timestamp,
					source: logGroup[0].source,
					target: logGroup[0].target,
					actionId: null,
					threat: newThreat,
				},
				curThreat,
				curThreat + newThreat,
				logGroup);

			curThreat += newThreat;
			return threatLog;
		});
	}
}

export class AuraEventLog extends SimLog {
	readonly isGained: boolean;
	readonly isFaded: boolean;
	readonly isRefreshed: boolean;

	constructor(params: SimLogParams, isGained: boolean, isFaded: boolean, isRefreshed: boolean) {
		super(params);
		this.isGained = isGained;
		this.isFaded = isFaded;
		this.isRefreshed = isRefreshed;
	}

	toString(): string {
		return `${this.toStringPrefix()} Aura ${this.isGained ? 'gained' : this.isFaded ? 'faded' : 'refreshed'}: ${this.actionId!.name}.`;
	}

	static parse(params: SimLogParams): Promise<AuraEventLog> | null {
		const match = params.raw.match(/Aura ((gained)|(faded)|(refreshed)): (.*)/);
		if (match && match[5]) {
			return ActionId.fromLogString(match[5]).fill(params.source?.index).then(aura => {
				params.actionId = aura;
				const event = match[1];
				return new AuraEventLog(params, event == 'gained', event == 'faded', event == 'refreshed');
			});
		} else {
			return null;
		}
	}
}

export class AuraStacksChangeLog extends SimLog {
	readonly oldStacks: number;
	readonly newStacks: number;

	constructor(params: SimLogParams, oldStacks: number, newStacks: number) {
		super(params);
		this.oldStacks = oldStacks;
		this.newStacks = newStacks;
	}

	toString(): string {
		return `${this.toStringPrefix()} ${this.actionId!.name} stacks: ${this.oldStacks} --> ${this.newStacks}.`;
	}

	static parse(params: SimLogParams): Promise<AuraStacksChangeLog> | null {
		const match = params.raw.match(/(.*) stacks: ([0-9]+) --> ([0-9]+)/);
		if (match && match[1]) {
			return ActionId.fromLogString(match[1]).fill(params.source?.index).then(aura => {
				params.actionId = aura;
				return new AuraStacksChangeLog(params, parseInt(match[2]), parseInt(match[3]));
			});
		} else {
			return null;
		}
	}
}

export class AuraUptimeLog extends SimLog {
	readonly gainedAt: number;
	readonly fadedAt: number;

	constructor(params: SimLogParams, fadedAt: number) {
		super(params);
		this.gainedAt = params.timestamp;
		this.fadedAt = fadedAt;
	}

	static fromLogs(logs: Array<SimLog>, entity: Entity, encounterDuration: number): Array<AuraUptimeLog> {
		let unmatchedGainedLogs: Array<AuraEventLog> = [];
		const uptimeLogs: Array<AuraUptimeLog> = [];

		logs.forEach((log: SimLog) => {
			if (!log.source || !log.source.equals(entity) || !log.isAuraEvent()) {
				return;
			}

			if (log.isGained) {
				unmatchedGainedLogs.push(log);
				return;
			}

			const matchingGainedIdx = unmatchedGainedLogs.findIndex(gainedLog => gainedLog.actionId!.equals(log.actionId!));
			if (matchingGainedIdx == -1) {
				console.warn('Unmatched aura faded log: ' + log.actionId!.name);
				return;
			}
			const gainedLog = unmatchedGainedLogs.splice(matchingGainedIdx, 1)[0];

			uptimeLogs.push(new AuraUptimeLog({
				raw: log.raw,
				logIndex: gainedLog.logIndex,
				timestamp: gainedLog.timestamp,
				source: log.source,
				target: log.target,
				actionId: gainedLog.actionId,
				threat: gainedLog.threat,
			}, log.timestamp));

			if (log.isRefreshed) {
				unmatchedGainedLogs.push(log);
			}
		});

		// Auras active at the end won't have a faded log, so need to add them separately.
		unmatchedGainedLogs.forEach(gainedLog => {
			uptimeLogs.push(new AuraUptimeLog({
				raw: gainedLog.raw,
				logIndex: gainedLog.logIndex,
				timestamp: gainedLog.timestamp,
				source: gainedLog.source,
				target: gainedLog.target,
				actionId: gainedLog.actionId,
				threat: gainedLog.threat,
			}, encounterDuration));
		});

		uptimeLogs.sort((a, b) => a.gainedAt - b.gainedAt);
		return uptimeLogs;
	}

	// Populates the activeAuras field for all logs using the provided auras.
	static populateActiveAuras(logs: Array<SimLog>, auraLogs: Array<AuraUptimeLog>) {
		let curAuras: Array<AuraUptimeLog> = [];
		let auraLogsIndex = 0;

		logs.forEach(log => {
			while (auraLogsIndex < auraLogs.length && auraLogs[auraLogsIndex].gainedAt <= log.timestamp) {
				curAuras.push(auraLogs[auraLogsIndex]);
				auraLogsIndex++;
			}
			curAuras = curAuras.filter(curAura => curAura.fadedAt > log.timestamp);

			const activeAuras = curAuras.slice();
			activeAuras.sort((a, b) => stringComparator(a.actionId!.name, b.actionId!.name));
			log.activeAuras = activeAuras;
		});
	}
}

export class ResourceChangedLog extends SimLog {
	readonly resourceType: ResourceType;
	readonly valueBefore: number;
	readonly valueAfter: number;
	readonly isSpend: boolean;

	constructor(params: SimLogParams, resourceType: ResourceType, valueBefore: number, valueAfter: number, isSpend: boolean) {
		super(params);
		this.resourceType = resourceType;
		this.valueBefore = valueBefore;
		this.valueAfter = valueAfter;
		this.isSpend = isSpend;
	}

	toString(): string {
		const signedDiff = (this.valueAfter - this.valueBefore) * (this.isSpend ? -1 : 1);

		const isHealth = this.resourceType == ResourceType.ResourceTypeHealth;
		const verb = isHealth ? (this.isSpend ? 'Lost' : 'Recovered') : (this.isSpend ? 'Spent' : 'Gained');

		return `${this.toStringPrefix()} ${verb} ${signedDiff.toFixed(1)} ${resourceNames[this.resourceType]} from ${this.actionId!.name}. (${this.valueBefore.toFixed(1)} --> ${this.valueAfter.toFixed(1)})`;
	}

	resultString(): string {
		const delta = this.valueAfter - this.valueBefore;
		if (delta < 0) {
			return delta.toFixed(1);
		} else {
			return '+' + delta.toFixed(1);
		}
	}

	static parse(params: SimLogParams): Promise<ResourceChangedLog> | null {
		const match = params.raw.match(/((Gained)|(Spent)) \d+\.?\d* ((health)|(mana)|(energy)|(focus)|(rage)|(combo points)) from (.*) \((\d+\.?\d*) --> (\d+\.?\d*)\)/);
		if (match) {
			const resourceType = stringToResourceType(match[4]);
			return ActionId.fromLogString(match[11]).fill(params.source?.index).then(cause => {
				params.actionId = cause;
				return new ResourceChangedLog(params, resourceType, parseFloat(match[12]), parseFloat(match[13]), match[1] == 'Spent');
			});
		} else {
			return null;
		}
	}
}

export class ResourceChangedLogGroup extends SimLog {
	readonly resourceType: ResourceType;
	readonly valueBefore: number;
	readonly valueAfter: number;
	readonly logs: Array<ResourceChangedLog>;

	constructor(params: SimLogParams, resourceType: ResourceType, valueBefore: number, valueAfter: number, logs: Array<ResourceChangedLog>) {
		super(params);
		this.resourceType = resourceType;
		this.valueBefore = valueBefore;
		this.valueAfter = valueAfter;
		this.logs = logs;
	}

	toString(): string {
		return `${this.toStringPrefix()} ${resourceNames[this.resourceType]}: ${this.valueBefore.toFixed(1)} --> ${this.valueAfter.toFixed(1)}`;
	}

	static fromLogs(logs: Array<SimLog>): Record<ResourceType, Array<ResourceChangedLogGroup>> {
		const allResourceChangedLogs = logs.filter((log): log is ResourceChangedLog => log.isResourceChanged());

		const results: Partial<Record<ResourceType, Array<ResourceChangedLogGroup>>> = {};
		const resourceTypes = (getEnumValues(ResourceType) as Array<ResourceType>).filter(val => val != ResourceType.ResourceTypeNone);
		resourceTypes.forEach(resourceType => {
			const resourceChangedLogs = allResourceChangedLogs.filter(log => log.resourceType == resourceType);

			const groupedLogs = SimLog.groupDuplicateTimestamps(resourceChangedLogs);
			results[resourceType] = groupedLogs.map(logGroup => new ResourceChangedLogGroup(
				{
					raw: '',
					logIndex: logGroup[0].logIndex,
					timestamp: logGroup[0].timestamp,
					source: logGroup[0].source,
					target: logGroup[0].target,
					actionId: null,
					threat: 0,
				},
				resourceType,
				logGroup[0].valueBefore,
				logGroup[logGroup.length - 1].valueAfter,
				logGroup));
		});

		return results as Record<ResourceType, Array<ResourceChangedLogGroup>>;
	}
}

export class MajorCooldownUsedLog extends SimLog {
	constructor(params: SimLogParams) {
		super(params);
	}

	toString(): string {
		return `${this.toStringPrefix()} Major cooldown used: ${this.actionId!.name}.`;
	}

	static parse(params: SimLogParams): Promise<MajorCooldownUsedLog> | null {
		const match = params.raw.match(/Major cooldown used: (.*)/);
		if (match) {
			return ActionId.fromLogString(match[1]).fill(params.source?.index).then(cooldownId => {
				params.actionId = cooldownId;
				return new MajorCooldownUsedLog(params);
			});
		} else {
			return null;
		}
	}
}

export class CastBeganLog extends SimLog {
	readonly manaCost: number;
	readonly castTime: number;

	constructor(params: SimLogParams, manaCost: number, castTime: number) {
		super(params);
		this.manaCost = manaCost;
		this.castTime = castTime;
	}

	toString(): string {
		return `${this.toStringPrefix()} Casting ${this.actionId!.name} (Cast time = ${this.castTime.toFixed(2)}s, Cost = ${this.manaCost.toFixed(1)}).`;
	}

	static parse(params: SimLogParams): Promise<CastBeganLog> | null {
		const match = params.raw.match(/Casting (.*) \(Cost = (\d+\.?\d*), Cast Time = (\d+\.?\d*)(m?s)\)/);
		if (match) {
			let castTime = parseFloat(match[3]);
			if (match[4] == 'ms') {
				castTime /= 1000;
			}
			return ActionId.fromLogString(match[1]).fill(params.source?.index).then(castId => {
				params.actionId = castId;
				return new CastBeganLog(params, parseFloat(match[2]), castTime);
			});
		} else {
			return null;
		}
	}
}

export class CastCompletedLog extends SimLog {
	constructor(params: SimLogParams) {
		super(params);
	}

	toString(): string {
		return `${this.toStringPrefix()} Completed cast ${this.actionId!.name}.`;
	}

	static parse(params: SimLogParams): Promise<CastCompletedLog> | null {
		const match = params.raw.match(/Completed cast (.*)/);
		if (match) {
			return ActionId.fromLogString(match[1]).fill(params.source?.index).then(castId => {
				params.actionId = castId;
				return new CastCompletedLog(params);
			});
		} else {
			return null;
		}
	}
}

export class CastLog extends SimLog {
	readonly castTime: number;

	readonly castBeganLog: CastBeganLog;
	readonly castCompletedLog: CastCompletedLog | null;

	// All instances of damage dealt from the completion of this cast until the completion of the next cast.
	readonly damageDealtLogs: Array<DamageDealtLog>;

	constructor(castBeganLog: CastBeganLog, castCompletedLog: CastCompletedLog | null, damageDealtLogs: Array<DamageDealtLog>) {
		super({
			raw: castBeganLog.raw,
			logIndex: castBeganLog.logIndex,
			timestamp: castBeganLog.timestamp,
			source: castBeganLog.source,
			target: castBeganLog.target,
			actionId: castCompletedLog?.actionId || castBeganLog.actionId, // Use completed log because of arcane blast
			threat: castCompletedLog?.threat || castBeganLog.threat,
		});
		this.castTime = castBeganLog.castTime;
		this.castBeganLog = castBeganLog;
		this.castCompletedLog = castCompletedLog;
		this.damageDealtLogs = damageDealtLogs;
	}

	toString(): string {
		return `${this.toStringPrefix()} Casting ${this.actionId!.name} (Cast time = ${this.castTime.toFixed(2)}s).`;
	}

	static fromLogs(logs: Array<SimLog>): Array<CastLog> {
		const castBeganLogs = logs.filter((log): log is CastBeganLog => log.isCastBegan());
		const castCompletedLogs = logs.filter((log): log is CastCompletedLog => log.isCastCompleted());
		const damageDealtLogs = logs.filter((log): log is DamageDealtLog => log.isDamageDealt());

		const toBucketKey = (actionId: ActionId) => {
			if (actionId.spellId == 30451) {
				// Arcane Blast is unique because it can finish its cast as a different
				// spell than it started (if stacks drop).
				return actionId.toStringIgnoringTag();
			} else {
				return actionId.toString();
			}
		};
		const castBeganLogsByAbility = bucket(castBeganLogs, log => toBucketKey(log.actionId!));
		const castCompletedLogsByAbility = bucket(castCompletedLogs, log => toBucketKey(log.actionId!));
		const damageDealtLogsByAbility = bucket(damageDealtLogs, log => toBucketKey(log.actionId!));

		const castLogs: Array<CastLog> = [];
		Object.keys(castBeganLogsByAbility).forEach(bucketKey => {
			const abilityCastsBegan = castBeganLogsByAbility[bucketKey]!;
			const abilityCastsCompleted = castCompletedLogsByAbility[bucketKey];
			const abilityDamageDealt = damageDealtLogsByAbility[bucketKey];
			const actionId = abilityCastsBegan[0].actionId!;

			const getCastCompleted = (cbIndex: number) => {
				if (!abilityCastsCompleted) {
					return null;
				}
				if (cbIndex >= abilityCastsBegan.length) {
					return null;
				}
				const nextBeganIndex = abilityCastsBegan[cbIndex + 1]?.logIndex || null;
				return abilityCastsCompleted.find(ccl =>
					ccl.logIndex > abilityCastsBegan[cbIndex].logIndex
					&& (nextBeganIndex == null || ccl.logIndex < nextBeganIndex)
					&& toBucketKey(ccl.actionId!) == toBucketKey(actionId)) || null;
			};

			if (!abilityDamageDealt) {
				abilityCastsBegan.forEach((castBegan, i) => castLogs.push(new CastLog(castBegan, getCastCompleted(i), [])));
				return;
			}

			let curCCL: CastCompletedLog | null = null;
			let nextCCL = getCastCompleted(0);
			let curDamageDealtLogs: Array<DamageDealtLog> = [];
			let curCbIndex = -1;
			abilityDamageDealt.forEach(ddLog => {
				if (curCbIndex >= abilityCastsBegan.length) {
					return;
				}
				if (nextCCL == null || ddLog.logIndex < nextCCL.logIndex) {
					curDamageDealtLogs.push(ddLog);
				} else {
					if (curCbIndex != -1) {
						castLogs.push(new CastLog(abilityCastsBegan[curCbIndex], curCCL, curDamageDealtLogs));
					}
					curDamageDealtLogs = [ddLog];
					curCbIndex++;
					if (curCbIndex >= abilityCastsBegan.length) {
						return;
					}
					curCCL = nextCCL;
					nextCCL = getCastCompleted(curCbIndex + 1);
				}
			});
			castLogs.push(new CastLog(abilityCastsBegan[curCbIndex], curCCL, curDamageDealtLogs));
		});

		castLogs.sort((a, b) => a.timestamp - b.timestamp);
		return castLogs;
	}
}

export class StatChangeLog extends SimLog {
	readonly isGain: boolean;
	readonly stats: string;

	constructor(params: SimLogParams, isGain: boolean, stats: string) {
		super(params);
		this.isGain = isGain;
		this.stats = stats;
	}

	toString(): string {
		if (this.isGain) {
			return `${this.toStringPrefix()} Gained ${this.stats} from ${this.actionId!.name}.`;
		} else {
			return `${this.toStringPrefix()} Lost ${this.stats} from fading ${this.actionId!.name}.`;
		}
	}

	static parse(params: SimLogParams): Promise<StatChangeLog> | null {
		const match = params.raw.match(/((Gained)|(Lost)) ({.*}) from (fading )?(.*)/);
		if (match) {
			return ActionId.fromLogString(match[6]).fill(params.source?.index).then(effectId => {
				params.actionId = effectId;
				const sign = match[1] == 'Lost' ? -1 : 1;
				return new StatChangeLog(params, sign == 1, match[4]);
			});
		} else {
			return null;
		}
	}
}
