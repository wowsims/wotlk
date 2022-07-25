import { ActionMetrics as ActionMetricsProto } from '/wotlk/core/proto/api.js';
import { AuraMetrics as AuraMetricsProto } from '/wotlk/core/proto/api.js';
import { DistributionMetrics as DistributionMetricsProto } from '/wotlk/core/proto/api.js';
import { ResourceMetrics as ResourceMetricsProto, ResourceType } from '/wotlk/core/proto/api.js';
import { TargetedActionMetrics as TargetedActionMetricsProto } from '/wotlk/core/proto/api.js';
import { RaidSimRequest, RaidSimResult } from '/wotlk/core/proto/api.js';
import { Class } from '/wotlk/core/proto/common.js';
import { SimRun } from '/wotlk/core/proto/ui.js';
import { ActionId, defaultTargetIcon } from '/wotlk/core/proto_utils/action_id.js';
import { classColors } from '/wotlk/core/proto_utils/utils.js';
import { getTalentTreeIcon } from '/wotlk/core/proto_utils/utils.js';
import { playerToSpec } from '/wotlk/core/proto_utils/utils.js';
import { specToClass } from '/wotlk/core/proto_utils/utils.js';
import { bucket } from '/wotlk/core/utils.js';
import { sum } from '/wotlk/core/utils.js';
import { AuraUptimeLog, CastLog, DpsLog, Entity, ResourceChangedLogGroup, SimLog, ThreatLogGroup, } from './logs_parser.js';
class SimResultData {
    constructor(request, result) {
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
    constructor(request, result, raidMetrics, encounterMetrics, logs) {
        this.request = request;
        this.result = result;
        this.raidMetrics = raidMetrics;
        this.encounterMetrics = encounterMetrics;
        this.logs = logs;
    }
    getPlayers(filter) {
        if (filter?.player || filter?.player === 0) {
            const player = this.getPlayerWithRaidIndex(filter.player);
            return player ? [player] : [];
        }
        else {
            return this.raidMetrics.parties.map(party => party.players).flat();
        }
    }
    // Returns the first player, regardless of which party / raid slot its in.
    getFirstPlayer() {
        return this.getPlayers()[0] || null;
    }
    getPlayerWithRaidIndex(raidIndex) {
        return this.getPlayers().find(player => player.index == raidIndex) || null;
    }
    getTargets(filter) {
        if (filter?.target || filter?.target === 0) {
            const target = this.getTargetWithIndex(filter.target);
            return target ? [target] : [];
        }
        else {
            return this.encounterMetrics.targets.slice();
        }
    }
    getTargetWithIndex(index) {
        return this.getTargets().find(target => target.index == index) || null;
    }
    getDamageMetrics(filter) {
        if (filter.player || filter.player === 0) {
            return this.getPlayerWithRaidIndex(filter.player)?.dps || DistributionMetricsProto.create();
        }
        return this.raidMetrics.dps;
    }
    getActionMetrics(filter) {
        return ActionMetrics.joinById(this.getPlayers(filter).map(player => player.getPlayerAndPetActions().map(action => action.forTarget(filter))).flat());
    }
    getSpellMetrics(filter) {
        return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && !e.isMeleeAction);
    }
    getMeleeMetrics(filter) {
        return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && e.isMeleeAction);
    }
    getResourceMetrics(resourceType, filter) {
        return ResourceMetrics.joinById(this.getPlayers(filter).map(player => player.resources.filter(resource => resource.type == resourceType)).flat());
    }
    getBuffMetrics(filter) {
        return AuraMetrics.joinById(this.getPlayers(filter).map(player => player.auras).flat());
    }
    getDebuffMetrics(filter) {
        return AuraMetrics.joinById(this.getTargets(filter).map(target => target.auras).flat()).filter(aura => aura.uptimePercent != 0);
    }
    toProto() {
        return SimRun.create({
            request: this.request,
            result: this.result,
        });
    }
    static async fromProto(proto) {
        return SimResult.makeNew(proto.request || RaidSimRequest.create(), proto.result || RaidSimResult.create());
    }
    static async makeNew(request, result) {
        const resultData = new SimResultData(request, result);
        const logs = await SimLog.parseAll(result);
        const raidPromise = RaidMetrics.makeNew(resultData, request.raid, result.raidMetrics, logs);
        const encounterPromise = EncounterMetrics.makeNew(resultData, request.encounter, result.encounterMetrics, logs);
        const raidMetrics = await raidPromise;
        const encounterMetrics = await encounterPromise;
        return new SimResult(request, result, raidMetrics, encounterMetrics, logs);
    }
}
export class RaidMetrics {
    constructor(raid, metrics, parties) {
        this.raid = raid;
        this.metrics = metrics;
        this.dps = this.metrics.dps;
        this.parties = parties;
    }
    static async makeNew(resultData, raid, metrics, logs) {
        const numParties = Math.min(raid.parties.length, metrics.parties.length);
        const parties = await Promise.all([...new Array(numParties).keys()]
            .map(i => PartyMetrics.makeNew(resultData, raid.parties[i], metrics.parties[i], i, logs)));
        return new RaidMetrics(raid, metrics, parties);
    }
}
export class PartyMetrics {
    constructor(party, metrics, partyIndex, players) {
        this.party = party;
        this.metrics = metrics;
        this.partyIndex = partyIndex;
        this.dps = this.metrics.dps;
        this.players = players;
    }
    static async makeNew(resultData, party, metrics, partyIndex, logs) {
        const numPlayers = Math.min(party.players.length, metrics.players.length);
        const players = await Promise.all([...new Array(numPlayers).keys()]
            .filter(i => party.players[i].class != Class.ClassUnknown)
            .map(i => UnitMetrics.makeNewPlayer(resultData, party.players[i], metrics.players[i], partyIndex * 5 + i, false, logs)));
        return new PartyMetrics(party, metrics, partyIndex, players);
    }
}
export class UnitMetrics {
    constructor(player, target, petActionId, metrics, index, actions, auras, resources, pets, logs, resultData) {
        this.player = player;
        this.target = target;
        this.metrics = metrics;
        this.index = index;
        this.name = metrics.name;
        this.spec = player ? playerToSpec(player) : 0;
        this.petActionId = petActionId;
        this.iconUrl = this.isPlayer ? getTalentTreeIcon(this.spec, player.talentsString) :
            (this.isTarget ? defaultTargetIcon : '');
        this.classColor = this.isTarget ? 'black' : classColors[specToClass[this.spec]];
        this.dps = this.metrics.dps;
        this.tps = this.metrics.threat;
        this.dtps = this.metrics.dtps;
        this.actions = actions;
        this.auras = auras;
        this.resources = resources;
        this.pets = pets;
        this.logs = logs;
        this.iterations = resultData.iterations;
        this.duration = resultData.duration;
        this.damageDealtLogs = this.logs.filter((log) => log.isDamageDealt());
        this.dpsLogs = DpsLog.fromLogs(this.damageDealtLogs);
        this.castLogs = CastLog.fromLogs(this.logs);
        this.threatLogs = ThreatLogGroup.fromLogs(this.logs);
        this.auraUptimeLogs = AuraUptimeLog.fromLogs(this.logs, new Entity(this.name, '', this.index, this.target != null, this.isPet), resultData.firstIterationDuration);
        this.majorCooldownLogs = this.logs.filter((log) => log.isMajorCooldownUsed());
        this.groupedResourceLogs = ResourceChangedLogGroup.fromLogs(this.logs);
        AuraUptimeLog.populateActiveAuras(this.dpsLogs, this.auraUptimeLogs);
        AuraUptimeLog.populateActiveAuras(this.groupedResourceLogs[ResourceType.ResourceTypeMana], this.auraUptimeLogs);
        this.majorCooldownAuraUptimeLogs = this.auraUptimeLogs.filter(auraLog => this.majorCooldownLogs.find(mcdLog => mcdLog.actionId.equals(auraLog.actionId)));
    }
    get label() {
        if (this.target == null) {
            return `${this.name} (#${this.index + 1})`;
        }
        else {
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
    // Returns the index of the target of this unit, as selected by the filter.
    getTargetIndex(filter) {
        if (!filter) {
            return null;
        }
        const index = this.isPlayer ? filter.target : filter.player;
        if (index == null || index == -1) {
            return null;
        }
        return index;
    }
    get inFrontOfTarget() {
        if (this.isTarget) {
            return true;
        }
        else if (this.isPlayer) {
            return this.player.inFrontOfTarget;
        }
        else {
            return false; // TODO pets
        }
    }
    get chanceOfDeath() {
        return this.metrics.chanceOfDeath * 100;
    }
    get maxThreat() {
        return this.threatLogs[this.threatLogs.length - 1]?.threatAfter || 0;
    }
    get secondsOomAvg() {
        return this.metrics.secondsOomAvg;
    }
    get totalDamage() {
        return this.dps.avg * this.duration;
    }
    getPlayerAndPetActions() {
        return this.actions.concat(this.pets.map(pet => pet.getPlayerAndPetActions()).flat());
    }
    getActionsForDisplay() {
        return this.actions.filter(e => e.hitAttempts != 0 || e.tps != 0 || e.dps != 0);
    }
    getMeleeActions() {
        return this.getActionsForDisplay().filter(e => e.isMeleeAction);
    }
    getSpellActions() {
        return this.getActionsForDisplay().filter(e => !e.isMeleeAction);
    }
    getResourceMetrics(resourceType) {
        return this.resources.filter(resource => resource.type == resourceType);
    }
    static async makeNewPlayer(resultData, player, metrics, raidIndex, isPet, logs) {
        const playerLogs = logs.filter(log => log.source && (!log.source.isTarget && (isPet == log.source.isPet) && log.source.index == raidIndex));
        const actionsPromise = Promise.all(metrics.actions.map(actionMetrics => ActionMetrics.makeNew(null, resultData, actionMetrics, raidIndex)));
        const aurasPromise = Promise.all(metrics.auras.map(auraMetrics => AuraMetrics.makeNew(null, resultData, auraMetrics, raidIndex)));
        const resourcesPromise = Promise.all(metrics.resources.map(resourceMetrics => ResourceMetrics.makeNew(null, resultData, resourceMetrics, raidIndex)));
        const petsPromise = Promise.all(metrics.pets.map(petMetrics => UnitMetrics.makeNewPlayer(resultData, player, petMetrics, raidIndex, true, playerLogs)));
        let petIdPromise = Promise.resolve(null);
        if (isPet) {
            petIdPromise = ActionId.fromPetName(metrics.name).fill(raidIndex);
        }
        const actions = await actionsPromise;
        const auras = await aurasPromise;
        const resources = await resourcesPromise;
        const pets = await petsPromise;
        const petActionId = await petIdPromise;
        const playerMetrics = new UnitMetrics(player, null, petActionId, metrics, raidIndex, actions, auras, resources, pets, playerLogs, resultData);
        actions.forEach(action => action.unit = playerMetrics);
        auras.forEach(aura => aura.unit = playerMetrics);
        resources.forEach(resource => resource.unit = playerMetrics);
        return playerMetrics;
    }
    static async makeNewTarget(resultData, target, metrics, index, logs) {
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
    constructor(encounter, metrics, targets) {
        this.encounter = encounter;
        this.metrics = metrics;
        this.targets = targets;
    }
    static async makeNew(resultData, encounter, metrics, logs) {
        const numTargets = Math.min(encounter.targets.length, metrics.targets.length);
        const targets = await Promise.all([...new Array(numTargets).keys()]
            .map(i => UnitMetrics.makeNewTarget(resultData, encounter.targets[i], metrics.targets[i], i, logs)));
        return new EncounterMetrics(encounter, metrics, targets);
    }
    get durationSeconds() {
        return this.encounter.duration;
    }
}
export class AuraMetrics {
    constructor(unit, actionId, data, resultData) {
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
    static async makeNew(unit, resultData, auraMetrics, playerIndex) {
        const actionId = await ActionId.fromProto(auraMetrics.id).fill(playerIndex);
        return new AuraMetrics(unit, actionId, auraMetrics, resultData);
    }
    // Merges an array of metrics into a single metrics.
    static merge(auras, removeTag, actionIdOverride) {
        const firstAura = auras[0];
        const unit = auras.every(aura => aura.unit == firstAura.unit) ? firstAura.unit : null;
        let actionId = actionIdOverride || firstAura.actionId;
        if (removeTag) {
            actionId = actionId.withoutTag();
        }
        return new AuraMetrics(unit, actionId, AuraMetricsProto.create({
            uptimeSecondsAvg: Math.max(...auras.map(a => a.data.uptimeSecondsAvg)),
        }), firstAura.resultData);
    }
    // Groups similar metrics, i.e. metrics with the same item/spell/other ID but
    // different tags, and returns them as separate arrays.
    static groupById(auras, useTag) {
        if (useTag) {
            return Object.values(bucket(auras, aura => aura.actionId.toString()));
        }
        else {
            return Object.values(bucket(auras, aura => aura.actionId.toStringIgnoringTag()));
        }
    }
    // Merges aura metrics that have the same name/ID, adding their stats together.
    static joinById(auras, useTag) {
        return AuraMetrics.groupById(auras, useTag).map(aurasToJoin => AuraMetrics.merge(aurasToJoin));
    }
}
;
export class ResourceMetrics {
    constructor(unit, actionId, data, resultData) {
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
    static async makeNew(unit, resultData, resourceMetrics, playerIndex) {
        const actionId = await ActionId.fromProto(resourceMetrics.id).fill(playerIndex);
        return new ResourceMetrics(unit, actionId, resourceMetrics, resultData);
    }
    // Merges an array of metrics into a single metrics.
    static merge(resources, removeTag, actionIdOverride) {
        const firstResource = resources[0];
        const unit = resources.every(resource => resource.unit == firstResource.unit) ? firstResource.unit : null;
        let actionId = actionIdOverride || firstResource.actionId;
        if (removeTag) {
            actionId = actionId.withoutTag();
        }
        return new ResourceMetrics(unit, actionId, ResourceMetricsProto.create({
            events: sum(resources.map(a => a.data.events)),
            gain: sum(resources.map(a => a.data.gain)),
            actualGain: sum(resources.map(a => a.data.actualGain)),
        }), firstResource.resultData);
    }
    // Groups similar metrics, i.e. metrics with the same item/spell/other ID but
    // different tags, and returns them as separate arrays.
    static groupById(resources, useTag) {
        if (useTag) {
            return Object.values(bucket(resources, resource => resource.actionId.toString()));
        }
        else {
            return Object.values(bucket(resources, resource => resource.actionId.toStringIgnoringTag()));
        }
    }
    // Merges resource metrics that have the same name/ID, adding their stats together.
    static joinById(resources, useTag) {
        return ResourceMetrics.groupById(resources, useTag).map(resourcesToJoin => ResourceMetrics.merge(resourcesToJoin));
    }
}
;
// Manages the metrics for a single unit action (e.g. Lightning Bolt).
export class ActionMetrics {
    constructor(unit, actionId, data, resultData) {
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
    get tps() {
        return this.combinedMetrics.tps;
    }
    get casts() {
        return this.combinedMetrics.casts;
    }
    get castsPerMinute() {
        return this.combinedMetrics.castsPerMinute;
    }
    get avgCast() {
        return this.combinedMetrics.avgCast;
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
    forTarget(filter) {
        const index = this.unit.getTargetIndex(filter);
        if (index == null) {
            return this;
        }
        else {
            const target = this.targets.find(target => target.data.unitIndex == index);
            if (target) {
                const targetData = ActionMetricsProto.clone(this.data);
                targetData.targets = [target.data];
                return new ActionMetrics(this.unit, this.actionId, targetData, this.resultData);
            }
            else {
                throw new Error('Could not find target with index ' + index);
            }
        }
    }
    static async makeNew(unit, resultData, actionMetrics, playerIndex) {
        const actionId = await ActionId.fromProto(actionMetrics.id).fill(playerIndex);
        return new ActionMetrics(unit, actionId, actionMetrics, resultData);
    }
    // Merges an array of metrics into a single metric.
    static merge(actions, removeTag, actionIdOverride) {
        const firstAction = actions[0];
        const unit = firstAction.unit;
        let actionId = actionIdOverride || firstAction.actionId;
        if (removeTag) {
            actionId = actionId.withoutTag();
        }
        const maxTargets = Math.max(...actions.map(action => action.targets.length));
        const mergedTargets = [...Array(maxTargets).keys()].map(i => TargetedActionMetrics.merge(actions.map(action => action.targets[i])));
        return new ActionMetrics(unit, actionId, ActionMetricsProto.create({
            isMelee: firstAction.isMeleeAction,
            targets: mergedTargets.map(t => t.data),
        }), firstAction.resultData);
    }
    // Groups similar metrics, i.e. metrics with the same item/spell/other ID but
    // different tags, and returns them as separate arrays.
    static groupById(actions, useTag) {
        if (useTag) {
            return Object.values(bucket(actions, action => action.actionId.toString()));
        }
        else {
            return Object.values(bucket(actions, action => action.actionId.toStringIgnoringTag()));
        }
    }
    // Merges action metrics that have the same name/ID, adding their stats together.
    static joinById(actions, useTag) {
        return ActionMetrics.groupById(actions, useTag).map(actionsToJoin => ActionMetrics.merge(actionsToJoin));
    }
}
// Manages the metrics for a single action applied to a specific target.
export class TargetedActionMetrics {
    constructor(iterations, duration, data) {
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
    get tps() {
        return this.data.threat / this.iterations / this.duration;
    }
    get casts() {
        return Math.max(this.data.casts, this.hitAttempts) / this.iterations;
    }
    get castsPerMinute() {
        return this.casts / (this.duration / 60);
    }
    get avgCast() {
        return (this.data.damage / this.iterations) / (this.casts || 1);
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
        return (this.data.misses / (this.hitAttempts || 1)) * 100;
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
    static merge(actions) {
        return new TargetedActionMetrics(actions[0]?.iterations || 1, actions[0]?.duration || 1, TargetedActionMetricsProto.create({
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
        }));
    }
}
