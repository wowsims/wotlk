import { ActionMetrics as ActionMetricsProto } from '../proto/api.js';
import { AuraMetrics as AuraMetricsProto } from '../proto/api.js';
import { DistributionMetrics as DistributionMetricsProto } from '../proto/api.js';
import { Encounter as EncounterProto } from '../proto/common.js';
import { EncounterMetrics as EncounterMetricsProto } from '../proto/api.js';
import { Party as PartyProto } from '../proto/api.js';
import { PartyMetrics as PartyMetricsProto } from '../proto/api.js';
import { Player as PlayerProto } from '../proto/api.js';
import { UnitMetrics as UnitMetricsProto } from '../proto/api.js';
import { Raid as RaidProto } from '../proto/api.js';
import { RaidMetrics as RaidMetricsProto } from '../proto/api.js';
import { ResourceMetrics as ResourceMetricsProto, ResourceType } from '../proto/api.js';
import { Target as TargetProto } from '../proto/common.js';
import { TargetedActionMetrics as TargetedActionMetricsProto } from '../proto/api.js';
import { RaidSimRequest, RaidSimResult } from '../proto/api.js';
import { Class } from '../proto/common.js';
import { Spec } from '../proto/common.js';
import { SimRun } from '../proto/ui.js';
import { ActionId, defaultTargetIcon } from '../proto_utils/action_id.js';
import { classColors } from '../proto_utils/utils.js';
import { getTalentTreeIcon } from '../proto_utils/utils.js';
import { playerToSpec } from '../proto_utils/utils.js';
import { specToClass } from '../proto_utils/utils.js';
import { bucket } from '../utils.js';
import { sum } from '../utils.js';

import {
	AuraUptimeLog,
	CastLog,
	DamageDealtLog,
	DpsLog,
	Entity,
	MajorCooldownUsedLog,
	ResourceChangedLogGroup,
	SimLog,
	ThreatLogGroup,
} from './logs_parser.js';
import { MAX_PARTY_SIZE } from '../party.js';

export interface SimResultFilter {
	// Raid index of the player to display, or null for all players.
	player?: number | null;

	// Target index of the target to display, or null for all targets.
	target?: number | null;
}

class SimResultData {
	readonly request: RaidSimRequest;
	readonly result: RaidSimResult;

	constructor(request: RaidSimRequest, result: RaidSimResult) {
		this.request = request;
		this.result = result;
	}

	get iterations() {
		return this.request.simOptions?.iterations || 1;
	}

	get duration() {
		return this.result.avgIterationDuration || 1;
	}

	get firstIterationDuration() {
		return this.result.firstIterationDuration || 1;
	}
}

// Holds all the data from a simulation call, and provides helper functions
// for parsing it.
export class SimResult {
	readonly request: RaidSimRequest;
	readonly result: RaidSimResult;

	readonly raidMetrics: RaidMetrics;
	readonly encounterMetrics: EncounterMetrics;
	readonly logs: Array<SimLog>;

	private players: Array<UnitMetrics>;
	private units: Array<UnitMetrics>;

	private constructor(request: RaidSimRequest, result: RaidSimResult, raidMetrics: RaidMetrics, encounterMetrics: EncounterMetrics, logs: Array<SimLog>) {
		this.request = request;
		this.result = result;
		this.raidMetrics = raidMetrics;
		this.encounterMetrics = encounterMetrics;
		this.logs = logs;

		this.players = raidMetrics.parties.map(party => party.players).flat();
		this.units = this.players.concat(encounterMetrics.targets);
	}

	get iterations() {
		return this.request.simOptions?.iterations || 1;
	}

	get duration() {
		return this.result.avgIterationDuration || 1;
	}

	get firstIterationDuration() {
		return this.result.firstIterationDuration || 1;
	}

	getPlayers(filter?: SimResultFilter): Array<UnitMetrics> {
		if (filter?.player || filter?.player === 0) {
			const player = this.getUnitWithIndex(filter.player);
			return player ? [player] : [];
		} else {
			return this.raidMetrics.parties.map(party => party.players).flat();
		}
	}

	// Returns the first player, regardless of which party / raid slot its in.
	getFirstPlayer(): UnitMetrics | null {
		return this.getPlayers()[0] || null;
	}

	getPlayerWithIndex(unitIndex: number): UnitMetrics | null {
		return this.players.find(player => player.unitIndex == unitIndex) || null;
	}
	getPlayerWithRaidIndex(raidIndex: number): UnitMetrics | null {
		return this.players.find(player => player.index == raidIndex) || null;
	}

	getTargets(filter?: SimResultFilter): Array<UnitMetrics> {
		if (filter?.target || filter?.target === 0) {
			const target = this.getUnitWithIndex(filter.target);
			return target ? [target] : [];
		} else {
			return this.encounterMetrics.targets.slice();
		}
	}

	getTargetWithIndex(unitIndex: number): UnitMetrics | null {
		return this.getTargets().find(target => target.unitIndex == unitIndex) || null;
	}
	getUnitWithIndex(unitIndex: number): UnitMetrics | null {
		return this.units.find(unit => unit.unitIndex == unitIndex) || null;
	}

	getDamageMetrics(filter: SimResultFilter): DistributionMetricsProto {
		if (filter.player || filter.player === 0) {
			return this.getPlayerWithIndex(filter.player)?.dps || DistributionMetricsProto.create();
		}

		return this.raidMetrics.dps;
	}

	getActionMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return ActionMetrics.joinById(this.getPlayers(filter).map(player => player.getPlayerAndPetActions().map(action => action.forTarget(filter))).flat());
	}

	getSpellMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && !e.isMeleeAction);
	}

	getMeleeMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && e.isMeleeAction);
	}

	getResourceMetrics(resourceType: ResourceType, filter?: SimResultFilter): Array<ResourceMetrics> {
		return ResourceMetrics.joinById(this.getPlayers(filter).map(player => player.resources.filter(resource => resource.type == resourceType)).flat());
	}

	getBuffMetrics(filter?: SimResultFilter): Array<AuraMetrics> {
		return AuraMetrics.joinById(this.getPlayers(filter).map(player => player.auras).flat());
	}

	getDebuffMetrics(filter?: SimResultFilter): Array<AuraMetrics> {
		return AuraMetrics.joinById(this.getTargets(filter).map(target => target.auras).flat()).filter(aura => aura.uptimePercent != 0);
	}

	toProto(): SimRun {
		return SimRun.create({
			request: this.request,
			result: this.result,
		});
	}

	static async fromProto(proto: SimRun): Promise<SimResult> {
		return SimResult.makeNew(proto.request || RaidSimRequest.create(), proto.result || RaidSimResult.create());
	}

	static async makeNew(request: RaidSimRequest, result: RaidSimResult): Promise<SimResult> {
		const resultData = new SimResultData(request, result);
		const logs = await SimLog.parseAll(result);

		const raidPromise = RaidMetrics.makeNew(resultData, request.raid!, result.raidMetrics!, logs);
		const encounterPromise = EncounterMetrics.makeNew(resultData, request.encounter!, result.encounterMetrics!, logs);

		const raidMetrics = await raidPromise;
		const encounterMetrics = await encounterPromise;

		return new SimResult(request, result, raidMetrics, encounterMetrics, logs);
	}
}

export class RaidMetrics {
	private readonly raid: RaidProto;
	private readonly metrics: RaidMetricsProto;

	readonly dps: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly parties: Array<PartyMetrics>;

	private constructor(raid: RaidProto, metrics: RaidMetricsProto, parties: Array<PartyMetrics>) {
		this.raid = raid;
		this.metrics = metrics;
		this.dps = this.metrics.dps!;
		this.hps = this.metrics.hps!;
		this.parties = parties;
	}

	static async makeNew(resultData: SimResultData, raid: RaidProto, metrics: RaidMetricsProto, logs: Array<SimLog>): Promise<RaidMetrics> {
		const numParties = Math.min(raid.parties.length, metrics.parties.length);

		const parties = await Promise.all(
			[...new Array(numParties).keys()]
				.map(i => PartyMetrics.makeNew(
					resultData,
					raid.parties[i],
					metrics.parties[i],
					i,
					logs)));

		return new RaidMetrics(raid, metrics, parties);
	}
}

export class PartyMetrics {
	private readonly party: PartyProto;
	private readonly metrics: PartyMetricsProto;

	readonly partyIndex: number;
	readonly dps: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly players: Array<UnitMetrics>;

	private constructor(party: PartyProto, metrics: PartyMetricsProto, partyIndex: number, players: Array<UnitMetrics>) {
		this.party = party;
		this.metrics = metrics;
		this.partyIndex = partyIndex;
		this.dps = this.metrics.dps!;
		this.hps = this.metrics.hps!;
		this.players = players;
	}

	static async makeNew(resultData: SimResultData, party: PartyProto, metrics: PartyMetricsProto, partyIndex: number, logs: Array<SimLog>): Promise<PartyMetrics> {
		const numPlayers = Math.min(party.players.length, metrics.players.length);
		const players = await Promise.all(
			[...new Array(numPlayers).keys()]
				.filter(i => party.players[i].class != Class.ClassUnknown)
				.map(i => UnitMetrics.makeNewPlayer(
					resultData,
					party.players[i],
					metrics.players[i],
					partyIndex * 5 + i,
					false,
					logs)));

		return new PartyMetrics(party, metrics, partyIndex, players);
	}
}

export class UnitMetrics {
	// If this Unit is a pet, player is the owner. If it's a target, player is null.
	private readonly player: PlayerProto | null;
	private readonly target: TargetProto | null;
	private readonly metrics: UnitMetricsProto;

	readonly index: number;
	readonly unitIndex: number;
	readonly name: string;
	readonly spec: Spec;
	readonly petActionId: ActionId | null;
	readonly iconUrl: string;
	readonly classColor: string;
	readonly dps: DistributionMetricsProto;
	readonly dpasp: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly tps: DistributionMetricsProto;
	readonly dtps: DistributionMetricsProto;
	readonly tmi: DistributionMetricsProto;
	readonly tto: DistributionMetricsProto;
	readonly actions: Array<ActionMetrics>;
	readonly auras: Array<AuraMetrics>;
	readonly resources: Array<ResourceMetrics>;
	readonly pets: Array<UnitMetrics>;
	private readonly iterations: number;
	private readonly duration: number;

	readonly logs: Array<SimLog>;
	readonly damageDealtLogs: Array<DamageDealtLog>;
	readonly groupedResourceLogs: Record<ResourceType, Array<ResourceChangedLogGroup>>;
	readonly dpsLogs: Array<DpsLog>;
	readonly auraUptimeLogs: Array<AuraUptimeLog>;
	readonly majorCooldownLogs: Array<MajorCooldownUsedLog>;
	readonly castLogs: Array<CastLog>;
	readonly threatLogs: Array<ThreatLogGroup>;

	// Aura uptime logs, filtered to include only auras that correspond to a
	// major cooldown.
	readonly majorCooldownAuraUptimeLogs: Array<AuraUptimeLog>;

	private constructor(
		player: PlayerProto | null,
		target: TargetProto | null,
		petActionId: ActionId | null,
		metrics: UnitMetricsProto,
		index: number,
		actions: Array<ActionMetrics>,
		auras: Array<AuraMetrics>,
		resources: Array<ResourceMetrics>,
		pets: Array<UnitMetrics>,
		logs: Array<SimLog>,
		resultData: SimResultData) {
		this.player = player;
		this.target = target;
		this.metrics = metrics;

		this.index = index;
		this.unitIndex = metrics.unitIndex;
		this.name = metrics.name;
		this.spec = player ? playerToSpec(player) : 0;
		this.petActionId = petActionId;
		this.iconUrl = this.isPlayer ? getTalentTreeIcon(this.spec, player!.talentsString) :
			(this.isTarget ? defaultTargetIcon : '');
		this.classColor = this.isTarget ? 'black' : classColors[specToClass[this.spec]];
		this.dps = this.metrics.dps!;
		this.dpasp = this.metrics.dpasp!;
		this.hps = this.metrics.hps!;
		this.tps = this.metrics.threat!;
		this.dtps = this.metrics.dtps!;
		this.tmi = this.metrics.tmi!;
		this.tto = this.metrics.tto!;
		this.actions = actions;
		this.auras = auras;
		this.resources = resources;
		this.pets = pets;
		this.logs = logs;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;

		this.damageDealtLogs = this.logs.filter((log): log is DamageDealtLog => log.isDamageDealt());
		this.dpsLogs = DpsLog.fromLogs(this.damageDealtLogs);
		this.castLogs = CastLog.fromLogs(this.logs);
		this.threatLogs = ThreatLogGroup.fromLogs(this.logs);

		this.auraUptimeLogs = AuraUptimeLog.fromLogs(this.logs, new Entity(this.name, '', this.index, this.target != null, this.isPet), resultData.firstIterationDuration);
		this.majorCooldownLogs = this.logs.filter((log): log is MajorCooldownUsedLog => log.isMajorCooldownUsed());

		this.groupedResourceLogs = ResourceChangedLogGroup.fromLogs(this.logs);
		AuraUptimeLog.populateActiveAuras(this.dpsLogs, this.auraUptimeLogs);
		AuraUptimeLog.populateActiveAuras(this.groupedResourceLogs[ResourceType.ResourceTypeMana], this.auraUptimeLogs);

		this.majorCooldownAuraUptimeLogs = this.auraUptimeLogs.filter(auraLog => this.majorCooldownLogs.find(mcdLog => mcdLog.actionId!.equals(auraLog.actionId!)));
	}

	get label() {
		if (this.target == null) {
			return `${this.name} (#${this.index + 1})`;
		} else {
			return this.name;
		}
	}

	get isPlayer() {
		return this.player != null;
	}

	get isTarget() {
		return this.target != null;
	}

	get isPet() {
		return this.petActionId != null;
	}

	// Returns the unit index of the target of this unit, as selected by the filter.
	getTargetIndex(filter?: SimResultFilter): number | null {
		if (!filter) {
			return null;
		}

		const index = this.isPlayer ? filter.target : filter.player;
		if (index == null || index == -1) {
			return null;
		}

		return index;
	}

	get inFrontOfTarget(): boolean {
		if (this.isTarget) {
			return true;
		} else if (this.isPlayer) {
			return this.player!.inFrontOfTarget;
		} else {
			return false; // TODO pets
		}
	}

	get chanceOfDeath(): number {
		return this.metrics.chanceOfDeath * 100;
	}

	get maxThreat() {
		return this.threatLogs[this.threatLogs.length - 1]?.threatAfter || 0;
	}

	get secondsOomAvg() {
		return this.metrics.secondsOomAvg
	}

	get totalDamage() {
		return this.dps.avg * this.duration;
	}

	getPlayerAndPetActions(): Array<ActionMetrics> {
		return this.actions.concat(this.pets.map(pet => pet.getPlayerAndPetActions()).flat());
	}

	private getActionsForDisplay(): Array<ActionMetrics> {
		return this.actions.filter(e => e.hitAttempts != 0 || e.tps != 0 || e.dps != 0);
	}

	getMeleeActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => e.isMeleeAction);
	}

	getSpellActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => !e.isMeleeAction);
	}

	getHealingActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay();
	}

	getResourceMetrics(resourceType: ResourceType): Array<ResourceMetrics> {
		return this.resources.filter(resource => resource.type == resourceType);
	}

	static async makeNewPlayer(resultData: SimResultData, player: PlayerProto, metrics: UnitMetricsProto, raidIndex: number, isPet: boolean, logs: Array<SimLog>): Promise<UnitMetrics> {
		const playerLogs = logs.filter(log => log.source && (!log.source.isTarget && (isPet == log.source.isPet) && log.source.index == raidIndex));

		const actionsPromise = Promise.all(metrics.actions.map(actionMetrics => ActionMetrics.makeNew(null, resultData, actionMetrics, raidIndex)));
		const aurasPromise = Promise.all(metrics.auras.map(auraMetrics => AuraMetrics.makeNew(null, resultData, auraMetrics, raidIndex)));
		const resourcesPromise = Promise.all(metrics.resources.map(resourceMetrics => ResourceMetrics.makeNew(null, resultData, resourceMetrics, raidIndex)));
		const petsPromise = Promise.all(metrics.pets.map(petMetrics => UnitMetrics.makeNewPlayer(resultData, player, petMetrics, raidIndex, true, playerLogs)));

		let petIdPromise: Promise<ActionId | null> = Promise.resolve(null);
		if (isPet) {
			petIdPromise = ActionId.fromPetName(metrics.name).fill(raidIndex);
		}

		const actions = await actionsPromise;
		const auras = await aurasPromise;
		const resources = await resourcesPromise;
		const pets = await petsPromise;
		const petActionId = await petIdPromise;

		const playerMetrics = new UnitMetrics(player, null, petActionId, metrics, raidIndex, actions, auras, resources, pets, playerLogs, resultData);
		actions.forEach(action => {
			action.unit = playerMetrics;
			action.resources = resources.filter(resourceMetrics => resourceMetrics.actionId.equals(action.actionId));
		});
		auras.forEach(aura => aura.unit = playerMetrics);
		resources.forEach(resource => resource.unit = playerMetrics);
		return playerMetrics;
	}

	static async makeNewTarget(resultData: SimResultData, target: TargetProto, metrics: UnitMetricsProto, index: number, logs: Array<SimLog>): Promise<UnitMetrics> {
		const targetLogs = logs.filter(log => log.source && (log.source.isTarget && log.source.index == index));

		const actionsPromise = Promise.all(metrics.actions.map(actionMetrics => ActionMetrics.makeNew(null, resultData, actionMetrics, index)));
		const aurasPromise = Promise.all(metrics.auras.map(auraMetrics => AuraMetrics.makeNew(null, resultData, auraMetrics)));

		const actions = await actionsPromise;
		const auras = await aurasPromise;

		const targetMetrics = new UnitMetrics(null, target, null, metrics, index, actions, auras, [], [], targetLogs, resultData);
		actions.forEach(action => action.unit = targetMetrics);
		auras.forEach(aura => aura.unit = targetMetrics);
		return targetMetrics;
	}
}

export class EncounterMetrics {
	private readonly encounter: EncounterProto;
	private readonly metrics: EncounterMetricsProto;

	readonly targets: Array<UnitMetrics>;

	private constructor(encounter: EncounterProto, metrics: EncounterMetricsProto, targets: Array<UnitMetrics>) {
		this.encounter = encounter;
		this.metrics = metrics;
		this.targets = targets;
	}

	static async makeNew(resultData: SimResultData, encounter: EncounterProto, metrics: EncounterMetricsProto, logs: Array<SimLog>): Promise<EncounterMetrics> {
		const numTargets = Math.min(encounter.targets.length, metrics.targets.length);
		const targets = await Promise.all(
			[...new Array(numTargets).keys()]
				.map(i => UnitMetrics.makeNewTarget(
					resultData,
					encounter.targets[i],
					metrics.targets[i],
					i,
					logs)));

		return new EncounterMetrics(encounter, metrics, targets);
	}

	get durationSeconds() {
		return this.encounter.duration;
	}
}

export class AuraMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: AuraMetricsProto;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: AuraMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.data = data;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
	}

	get uptimePercent() {
		return this.data.uptimeSecondsAvg / this.duration * 100;
	}

	get averageProcs() {
		return this.data.procsAvg
	}

	get ppm() {
		return this.data.procsAvg / (this.duration / 60);
	}

	static async makeNew(unit: UnitMetrics | null, resultData: SimResultData, auraMetrics: AuraMetricsProto, playerIndex?: number): Promise<AuraMetrics> {
		const actionId = await ActionId.fromProto(auraMetrics.id!).fill(playerIndex);
		return new AuraMetrics(unit, actionId, auraMetrics, resultData);
	}

	// Merges an array of metrics into a single metrics.
	static merge(auras: Array<AuraMetrics>, removeTag?: boolean, actionIdOverride?: ActionId): AuraMetrics {
		const firstAura = auras[0];
		const unit = auras.every(aura => aura.unit == firstAura.unit) ? firstAura.unit : null;
		let actionId = actionIdOverride || firstAura.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}
		return new AuraMetrics(
			unit,
			actionId,
			AuraMetricsProto.create({
				uptimeSecondsAvg: Math.max(...auras.map(a => a.data.uptimeSecondsAvg)),
			}),
			firstAura.resultData);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(auras: Array<AuraMetrics>, useTag?: boolean): Array<Array<AuraMetrics>> {
		if (useTag) {
			return Object.values(bucket(auras, aura => aura.actionId.toString()));
		} else {
			return Object.values(bucket(auras, aura => aura.actionId.toStringIgnoringTag()));
		}
	}

	// Merges aura metrics that have the same name/ID, adding their stats together.
	static joinById(auras: Array<AuraMetrics>, useTag?: boolean): Array<AuraMetrics> {
		return AuraMetrics.groupById(auras, useTag).map(aurasToJoin => AuraMetrics.merge(aurasToJoin));
	}
};

export class ResourceMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	readonly type: ResourceType;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: ResourceMetricsProto;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: ResourceMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.type = data.type;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
		this.data = data;
	}

	get events() {
		return this.data.events / this.iterations;
	}

	get gain() {
		return this.data.gain / this.iterations;
	}

	get gainPerSecond() {
		return this.data.gain / this.iterations / this.duration;
	}

	get avgGain() {
		return this.data.gain / this.data.events;
	}

	get wastedGain() {
		return (this.data.gain - this.data.actualGain) / this.iterations;
	}

	static async makeNew(unit: UnitMetrics | null, resultData: SimResultData, resourceMetrics: ResourceMetricsProto, playerIndex?: number): Promise<ResourceMetrics> {
		const actionId = await ActionId.fromProto(resourceMetrics.id!).fill(playerIndex);
		return new ResourceMetrics(unit, actionId, resourceMetrics, resultData);
	}

	// Merges an array of metrics into a single metrics.
	static merge(resources: Array<ResourceMetrics>, removeTag?: boolean, actionIdOverride?: ActionId): ResourceMetrics {
		const firstResource = resources[0];
		const unit = resources.every(resource => resource.unit == firstResource.unit) ? firstResource.unit : null;
		let actionId = actionIdOverride || firstResource.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}
		return new ResourceMetrics(
			unit,
			actionId,
			ResourceMetricsProto.create({
				events: sum(resources.map(a => a.data.events)),
				gain: sum(resources.map(a => a.data.gain)),
				actualGain: sum(resources.map(a => a.data.actualGain)),
			}),
			firstResource.resultData);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(resources: Array<ResourceMetrics>, useTag?: boolean): Array<Array<ResourceMetrics>> {
		if (useTag) {
			return Object.values(bucket(resources, resource => resource.actionId.toString()));
		} else {
			return Object.values(bucket(resources, resource => resource.actionId.toStringIgnoringTag()));
		}
	}

	// Merges resource metrics that have the same name/ID, adding their stats together.
	static joinById(resources: Array<ResourceMetrics>, useTag?: boolean): Array<ResourceMetrics> {
		return ResourceMetrics.groupById(resources, useTag).map(resourcesToJoin => ResourceMetrics.merge(resourcesToJoin));
	}
};

// Manages the metrics for a single unit action (e.g. Lightning Bolt).
export class ActionMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	readonly targets: Array<TargetedActionMetrics>;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: ActionMetricsProto;
	private readonly combinedMetrics: TargetedActionMetrics;
	resources: Array<ResourceMetrics>;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: ActionMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
		this.data = data;
		this.targets = data.targets.map(tam => new TargetedActionMetrics(this.iterations, this.duration, tam));
		this.combinedMetrics = TargetedActionMetrics.merge(this.targets);
		this.resources = [];
	}

	get isMeleeAction() {
		return this.data.isMelee;
	}

	get damage() {
		return this.combinedMetrics.damage;
	}

	get dps() {
		return this.combinedMetrics.dps;
	}

	get hps() {
		return this.combinedMetrics.hps;
	}

	get tps() {
		return this.combinedMetrics.tps;
	}

	get casts() {
		return this.combinedMetrics.casts;
	}

	get castsPerMinute() {
		return this.combinedMetrics.castsPerMinute;
	}

	get avgCastTimeMs() {
		return this.combinedMetrics.avgCastTimeMs;
	}

	get hpm() {
		const totalHealing = this.combinedMetrics.hps * this.duration;
		const manaMetrics = this.resources.find(r => r.type == ResourceType.ResourceTypeMana);
		if (manaMetrics) {
			return totalHealing / -manaMetrics.gain;
		}

		return 0;
	}

	get healingThroughput() {
		return this.combinedMetrics.healingThroughput;
	}

	get avgCast() {
		return this.combinedMetrics.avgCast;
	}

	get avgCastHealing() {
		return this.combinedMetrics.avgCastHealing;
	}

	get avgCastThreat() {
		return this.combinedMetrics.avgCastThreat;
	}

	get landedHits() {
		return this.combinedMetrics.landedHits;
	}

	get hitAttempts() {
		return this.combinedMetrics.hitAttempts;
	}

	get avgHit() {
		return this.combinedMetrics.avgHit;
	}

	get avgHitThreat() {
		return this.combinedMetrics.avgHitThreat;
	}

	get critPercent() {
		return this.combinedMetrics.critPercent;
	}

	get misses() {
		return this.combinedMetrics.misses;
	}

	get missPercent() {
		return this.combinedMetrics.missPercent;
	}

	get dodges() {
		return this.combinedMetrics.dodges;
	}

	get dodgePercent() {
		return this.combinedMetrics.dodgePercent;
	}

	get parries() {
		return this.combinedMetrics.parries;
	}

	get parryPercent() {
		return this.combinedMetrics.parryPercent;
	}

	get blocks() {
		return this.combinedMetrics.blocks;
	}

	get blockPercent() {
		return this.combinedMetrics.blockPercent;
	}

	get glances() {
		return this.combinedMetrics.glances;
	}

	get glancePercent() {
		return this.combinedMetrics.glancePercent;
	}

	forTarget(filter?: SimResultFilter): ActionMetrics {
		const unitIndex = this.unit!.getTargetIndex(filter);
		if (unitIndex == null) {
			return this;
		} else {
			const target = this.targets.find(target => target.data.unitIndex == unitIndex);
			if (target) {
				const targetData = ActionMetricsProto.clone(this.data);
				targetData.targets = [target.data];
				return new ActionMetrics(this.unit, this.actionId, targetData, this.resultData);
			} else {
				throw new Error('Could not find target with unitIndex ' + unitIndex);
			}
		}
	}

	static async makeNew(unit: UnitMetrics | null, resultData: SimResultData, actionMetrics: ActionMetricsProto, playerIndex?: number): Promise<ActionMetrics> {
		const actionId = await ActionId.fromProto(actionMetrics.id!).fill(playerIndex);
		return new ActionMetrics(unit, actionId, actionMetrics, resultData);
	}

	// Merges an array of metrics into a single metric.
	static merge(actions: Array<ActionMetrics>, removeTag?: boolean, actionIdOverride?: ActionId): ActionMetrics {
		const firstAction = actions[0];
		const unit = firstAction.unit
		let actionId = actionIdOverride || firstAction.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}

		const maxTargets = Math.max(...actions.map(action => action.targets.length));
		const mergedTargets = [...Array(maxTargets).keys()].map(i => TargetedActionMetrics.merge(actions.map(action => action.targets[i])));

		return new ActionMetrics(
			unit,
			actionId,
			ActionMetricsProto.create({
				isMelee: firstAction.isMeleeAction,
				targets: mergedTargets.map(t => t.data),
			}),
			firstAction.resultData);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(actions: Array<ActionMetrics>, useTag?: boolean): Array<Array<ActionMetrics>> {
		if (useTag) {
			return Object.values(bucket(actions, action => action.actionId.toString()));
		} else {
			return Object.values(bucket(actions, action => action.actionId.toStringIgnoringTag()));
		}
	}

	// Merges action metrics that have the same name/ID, adding their stats together.
	static joinById(actions: Array<ActionMetrics>, useTag?: boolean): Array<ActionMetrics> {
		return ActionMetrics.groupById(actions, useTag).map(actionsToJoin => ActionMetrics.merge(actionsToJoin));
	}
}

// Manages the metrics for a single action applied to a specific target.
export class TargetedActionMetrics {
	private readonly iterations: number;
	private readonly duration: number;
	readonly data: TargetedActionMetricsProto;

	readonly landedHitsRaw: number;
	readonly hitAttempts: number;

	constructor(iterations: number, duration: number, data: TargetedActionMetricsProto) {
		this.iterations = iterations;
		this.duration = duration;
		this.data = data;

		this.landedHitsRaw = this.data.hits + this.data.crits + this.data.blocks + this.data.glances;

		this.hitAttempts = this.data.misses
			+ this.data.dodges
			+ this.data.parries
			+ this.data.blocks
			+ this.data.glances
			+ this.data.crits
			+ this.data.hits;
	}

	get damage() {
		return this.data.damage;
	}

	get dps() {
		return this.data.damage / this.iterations / this.duration;
	}

	get hps() {
		return (this.data.healing + this.data.shielding) / this.iterations / this.duration;
	}

	get tps() {
		return this.data.threat / this.iterations / this.duration;
	}

	get casts() {
		return (this.data.casts || this.hitAttempts) / this.iterations;
	}

	get castsPerMinute() {
		return this.casts / (this.duration / 60);
	}

	get avgCastTimeMs() {
		return this.data.castTimeMs / this.iterations / this.casts;
	}

	get healingThroughput() {
		if (this.avgCastTimeMs) {
			return this.hps / (this.avgCastTimeMs / 1000);
		} else {
			return 0;
		}
	}

	get timeSpentCastingMs() {
		return this.data.castTimeMs / this.iterations;
	}

	get avgCast() {
		return (this.data.damage / this.iterations) / (this.casts || 1);
	}

	get avgCastHealing() {
		return ((this.data.healing + this.data.shielding) / this.iterations) / (this.casts || 1);
	}

	get avgCastThreat() {
		return (this.data.threat / this.iterations) / (this.casts || 1);
	}

	get landedHits() {
		return this.landedHitsRaw / this.iterations;
	}

	get avgHit() {
		const lhr = this.landedHitsRaw;
		return lhr == 0 ? 0 : this.data.damage / lhr;
	}

	get avgHitThreat() {
		const lhr = this.landedHitsRaw;
		return lhr == 0 ? 0 : this.data.threat / lhr;
	}

	get critPercent() {
		return (this.data.crits / (this.hitAttempts || 1)) * 100;
	}

	get misses() {
		return this.data.misses / this.iterations;
	}

	get missPercent() {
		return (this.data.misses / (this.data.casts || 1)) * 100;
	}

	get dodges() {
		return this.data.dodges / this.iterations;
	}

	get dodgePercent() {
		return (this.data.dodges / (this.hitAttempts || 1)) * 100;
	}

	get parries() {
		return this.data.parries / this.iterations;
	}

	get parryPercent() {
		return (this.data.parries / (this.hitAttempts || 1)) * 100;
	}

	get blocks() {
		return this.data.blocks / this.iterations;
	}

	get blockPercent() {
		return (this.data.blocks / (this.hitAttempts || 1)) * 100;
	}

	get glances() {
		return this.data.glances / this.iterations;
	}

	get glancePercent() {
		return (this.data.glances / (this.hitAttempts || 1)) * 100;
	}

	// Merges an array of metrics into a single metric.
	static merge(actions: Array<TargetedActionMetrics>): TargetedActionMetrics {
		return new TargetedActionMetrics(
			actions[0]?.iterations || 1,
			actions[0]?.duration || 1,
			TargetedActionMetricsProto.create({
				casts: sum(actions.map(a => a.data.casts)),
				hits: sum(actions.map(a => a.data.hits)),
				crits: sum(actions.map(a => a.data.crits)),
				misses: sum(actions.map(a => a.data.misses)),
				dodges: sum(actions.map(a => a.data.dodges)),
				parries: sum(actions.map(a => a.data.parries)),
				blocks: sum(actions.map(a => a.data.blocks)),
				glances: sum(actions.map(a => a.data.glances)),
				damage: sum(actions.map(a => a.data.damage)),
				threat: sum(actions.map(a => a.data.threat)),
				healing: sum(actions.map(a => a.data.healing)),
				shielding: sum(actions.map(a => a.data.shielding)),
				castTimeMs: sum(actions.map(a => a.data.castTimeMs)),
			}));
	}
}
