import { Class, Faction } from '/tbc/core/proto/common.js';
import { Consumes } from '/tbc/core/proto/common.js';
import { Enchant } from '/tbc/core/proto/common.js';
import { Encounter as EncounterProto } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';
import { ItemQuality } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { ItemType } from '/tbc/core/proto/common.js';
import { Item } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { Raid as RaidProto } from '/tbc/core/proto/api.js';
import { PresetEncounter, PresetTarget } from '/tbc/core/proto/api.js';
import { ComputeStatsRequest, ComputeStatsResult } from '/tbc/core/proto/api.js';
import { GearListRequest, GearListResult } from '/tbc/core/proto/api.js';
import { RaidSimRequest, RaidSimResult } from '/tbc/core/proto/api.js';
import { SimOptions } from '/tbc/core/proto/api.js';
import { StatWeightsRequest, StatWeightsResult } from '/tbc/core/proto/api.js';
import { SimSettings as SimSettingsProto } from '/tbc/core/proto/ui.js';

import { EquippedItem } from '/tbc/core/proto_utils/equipped_item.js';
import { Gear } from '/tbc/core/proto_utils/gear.js';
import { SimResult } from '/tbc/core/proto_utils/sim_result.js';
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { gemEligibleForSocket } from '/tbc/core/proto_utils/gems.js';
import { gemMatchesSocket } from '/tbc/core/proto_utils/gems.js';
import { SpecRotation } from '/tbc/core/proto_utils/utils.js';
import { SpecTalents } from '/tbc/core/proto_utils/utils.js';
import { SpecTypeFunctions } from '/tbc/core/proto_utils/utils.js';
import { specTypeFunctions } from '/tbc/core/proto_utils/utils.js';
import { SpecOptions } from '/tbc/core/proto_utils/utils.js';
import { specToClass } from '/tbc/core/proto_utils/utils.js';
import { specToEligibleRaces } from '/tbc/core/proto_utils/utils.js';
import { getEligibleItemSlots } from '/tbc/core/proto_utils/utils.js';
import { getEligibleEnchantSlots } from '/tbc/core/proto_utils/utils.js';
import { playerToSpec } from '/tbc/core/proto_utils/utils.js';

import { Encounter } from './encounter.js';
import { Player } from './player.js';
import { Raid } from './raid.js';
import { Listener } from './typed_event.js';
import { EventID, TypedEvent } from './typed_event.js';
import { sum } from './utils.js';
import { wait } from './utils.js';
import { WorkerPool } from './worker_pool.js';

import * as OtherConstants from '/tbc/core/constants/other.js';

declare var pako: any;

export type RaidSimData = {
	request: RaidSimRequest,
	result: RaidSimResult,
};

export type StatWeightsData = {
	request: StatWeightsRequest,
	result: StatWeightsResult,
};

// Core Sim module which deals only with api types, no UI-related stuff.
export class Sim {
	private readonly workerPool: WorkerPool;

	private iterations: number = 3000;
	private phase: number = OtherConstants.CURRENT_PHASE;
	private faction: Faction = Faction.Alliance;
	private fixedRngSeed: number = 0;
	private show1hWeapons: boolean = true;
	private show2hWeapons: boolean = true;
	private showMatchingGems: boolean = true;
	private showThreatMetrics: boolean = false;
	private showExperimental: boolean = false;

	readonly raid: Raid;
	readonly encounter: Encounter;

	// Database
	private items: Record<number, Item> = {};
	private enchants: Record<number, Enchant> = {};
	private gems: Record<number, Gem> = {};
	private presetEncounters: Record<string, PresetEncounter> = {};
	private presetTargets: Record<string, PresetTarget> = {};

	readonly iterationsChangeEmitter = new TypedEvent<void>();
	readonly phaseChangeEmitter = new TypedEvent<void>();
	readonly factionChangeEmitter = new TypedEvent<void>();
	readonly fixedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly lastUsedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly show1hWeaponsChangeEmitter = new TypedEvent<void>();
	readonly show2hWeaponsChangeEmitter = new TypedEvent<void>();
	readonly showMatchingGemsChangeEmitter = new TypedEvent<void>();
	readonly showThreatMetricsChangeEmitter = new TypedEvent<void>();
	readonly showExperimentalChangeEmitter = new TypedEvent<void>();

	// Emits when any of the settings change (but not the raid / encounter).
	readonly settingsChangeEmitter: TypedEvent<void>;

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	// Fires when a raid sim API call completes.
	readonly simResultEmitter = new TypedEvent<SimResult>();

	private readonly _initPromise: Promise<void>;
	private lastUsedRngSeed: number = 0;

	// These callbacks are needed so we can apply BuffBot modifications automatically before sending requests.
	private modifyRaidProto: ((raidProto: RaidProto) => void) = () => { };

	constructor() {
		this.workerPool = new WorkerPool(3);

		this._initPromise = this.workerPool.getGearList(GearListRequest.create()).then(result => {
			result.items.forEach(item => this.items[item.id] = item);
			result.enchants.forEach(enchant => this.enchants[enchant.id] = enchant);
			result.gems.forEach(gem => this.gems[gem.id] = gem);
			result.encounters.forEach(encounter => this.presetEncounters[encounter.path] = encounter);
			result.encounters.map(e => e.targets).flat().forEach(target => this.presetTargets[target.path] = target);
		});

		this.raid = new Raid(this);
		this.encounter = new Encounter(this);

		this.settingsChangeEmitter = TypedEvent.onAny([
			this.iterationsChangeEmitter,
			this.phaseChangeEmitter,
			this.fixedRngSeedChangeEmitter,
			this.show1hWeaponsChangeEmitter,
			this.show2hWeaponsChangeEmitter,
			this.showMatchingGemsChangeEmitter,
			this.showThreatMetricsChangeEmitter,
			this.showExperimentalChangeEmitter,
		]);

		this.changeEmitter = TypedEvent.onAny([
			this.settingsChangeEmitter,
			this.raid.changeEmitter,
			this.encounter.changeEmitter,
		]);

		this.raid.changeEmitter.on(eventID => this.updateCharacterStats(eventID));
	}

	waitForInit(): Promise<void> {
		return this._initPromise;
	}

	setModifyRaidProto(newModFn: (raidProto: RaidProto) => void) {
		this.modifyRaidProto = newModFn;
	}
	getModifiedRaidProto(): RaidProto {
		const raidProto = this.raid.toProto();
		this.modifyRaidProto(raidProto);

		// Remove any inactive meta gems, since the backend doesn't have its own validation.
		raidProto.parties.forEach(party => {
			party.players.forEach(player => {
				if (!player.equipment) {
					return;
				}

				const gear = this.lookupEquipmentSpec(player.equipment);
				if (gear.hasInactiveMetaGem()) {
					player.equipment = gear.withoutMetaGem().asSpec();
				}
			});
		});

		return raidProto;
	}

	private makeRaidSimRequest(debug: boolean): RaidSimRequest {
		const raid = this.getModifiedRaidProto();
		const encounter = this.encounter.toProto();
		const hunters = raid.parties.map(party => party.players).flat().filter(player => player.name && playerToSpec(player) == Spec.SpecHunter);
		if (hunters.some(hunter => (specTypeFunctions[Spec.SpecHunter]!.talentsFromPlayer(hunter) as SpecTalents<Spec.SpecHunter>).exposeWeakness > 0)) {
			if (raid.debuffs) {
				raid.debuffs.exposeWeaknessUptime = 0;
			}
		}

		return RaidSimRequest.create({
			raid: raid,
			encounter: encounter,
			simOptions: SimOptions.create({
				iterations: debug ? 1 : this.getIterations(),
				randomSeed: BigInt(this.nextRngSeed()),
				debugFirstIteration: true,
			}),
		});
	}

	async runRaidSim(eventID: EventID, onProgress: Function) {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.getNumTargets() < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = this.makeRaidSimRequest(false);

		var result = await this.workerPool.raidSimAsync(request, onProgress);
		if (result.errorResult != "") {
			this.handleError(result.errorResult, this.encodeSimReq(request));
			return;
		}
		const simResult = await SimResult.makeNew(request, result);
		this.simResultEmitter.emit(eventID, simResult);
	}

	encodeSimReq(req: RaidSimRequest): string {
		const protoBytes = RaidSimRequest.toBinary(req);
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		return btoa(String.fromCharCode(...deflated));
	}

	handleError(errorStr: string, extra: string) {
		if (window.confirm("Simulation Failure:\n" + errorStr + "\nPress Ok to file crash report")) {
			// Splice out just the line numbers
			var filteredError = errorStr.substring(0, errorStr.indexOf("Stack Trace:"));
			const rExp: RegExp = /(.*\.go:\d+)/g;
			filteredError += errorStr.match(rExp)?.join(" ");
			var hash = this.hashCode(filteredError);
			fetch('https://api.github.com/search/issues?q=is:issue+is:open+repo:wowsims/tbc+' + hash).then(resp => {
				resp.json().then((issues) => {
					if (issues.total_count > 0) {
						window.open(issues.items[0].html_url, "_blank");
					} else {
						window.open("https://github.com/wowsims/tbc/issues/new?assignees=&labels=&title=Crash%20Report%20" + hash + "&body=" + encodeURIComponent(errorStr + "\n\nRequest:\n" + extra), '_blank');
					}
				});
			}).catch(fetchErr => {
				alert("Failed to file report... try again another time:" + fetchErr);
			});
		}
		return;
	}

	hashCode(str: string): number {
		let hash = 0;
		for (let i = 0, len = str.length; i < len; i++) {
			let chr = str.charCodeAt(i);
			hash = (hash << 5) - hash + chr;
			hash |= 0; // Convert to 32bit integer
		}
		return hash;
	}

	async runRaidSimWithLogs(eventID: EventID): Promise<SimResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.getNumTargets() < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = this.makeRaidSimRequest(true);
		const result = await this.workerPool.raidSimAsync(request, () => { });

		const simResult = await SimResult.makeNew(request, result);
		this.simResultEmitter.emit(eventID, simResult);
		return simResult;
	}

	// This should be invoked internally whenever stats might have changed.
	private async updateCharacterStats(eventID: EventID) {
		if (eventID == 0) {
			// Skip the first event ID because it interferes with the loaded stats.
			return;
		}
		eventID = TypedEvent.nextEventID();

		await this.waitForInit();

		// Capture the current players so we avoid issues if something changes while
		// request is in-flight.
		const players = this.raid.getPlayers();

		const req = ComputeStatsRequest.create({ raid: this.getModifiedRaidProto() });
		const result = await this.workerPool.computeStats(req);

		if (result.errorResult != "") {
			this.handleError(result.errorResult, this.encodeComputeStatsReq(req));
			return
		}

		TypedEvent.freezeAllAndDo(() => {
			result.raidStats!.parties
				.forEach((partyStats, partyIndex) =>
					partyStats.players.forEach((playerStats, playerIndex) =>
						players[partyIndex * 5 + playerIndex]?.setCurrentStats(eventID, playerStats)));
		});
	}


	encodeComputeStatsReq(req: ComputeStatsRequest): string {
		const protoBytes = ComputeStatsRequest.toBinary(req);
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		return btoa(String.fromCharCode(...deflated));
	}

	async statWeights(player: Player<any>, epStats: Array<Stat>, epReferenceStat: Stat, onProgress: Function): Promise<StatWeightsResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.getNumTargets() < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		if (player.getParty() == null) {
			console.warn('Trying to get stat weights without a party!');
			return StatWeightsResult.create();
		} else {
			const tanks = this.raid.getTanks().map(tank => tank.targetIndex).includes(player.getRaidIndex())
				? [RaidTarget.create({ targetIndex: 0 })]
				: [];
			const request = StatWeightsRequest.create({
				player: player.toProto(),
				raidBuffs: this.raid.getBuffs(),
				partyBuffs: player.getParty()!.getBuffs(),
				debuffs: this.raid.getDebuffs(),
				encounter: this.encounter.toProto(),
				simOptions: SimOptions.create({
					iterations: this.getIterations(),
					randomSeed: BigInt(this.nextRngSeed()),
					debug: false,
				}),
				tanks: tanks,

				statsToWeigh: epStats,
				epReferenceStat: epReferenceStat,
			});
			var result = await this.workerPool.statWeightsAsync(request, onProgress);
			return result;
		}
	}

	getItems(slot: ItemSlot | undefined): Array<Item> {
		let items = Object.values(this.items);
		if (slot != undefined) {
			items = items.filter(item => getEligibleItemSlots(item).includes(slot));
		}
		return items;
	}

	getEnchants(slot: ItemSlot | undefined): Array<Enchant> {
		let enchants = Object.values(this.enchants);
		if (slot != undefined) {
			enchants = enchants.filter(enchant => getEligibleEnchantSlots(enchant).includes(slot));
		}
		return enchants;
	}

	// ID can be the formula ID OR the effect ID.
	getEnchantFlexible(id: number): Enchant | null {
		return Object.values(this.enchants).find(enchant => enchant.id == id || enchant.effectId == id) || null;
	}

	getGems(socketColor: GemColor | undefined): Array<Gem> {
		let gems = Object.values(this.gems);
		if (socketColor) {
			gems = gems.filter(gem => gemEligibleForSocket(gem, socketColor));
		}
		return gems;
	}

	getMatchingGems(socketColor: GemColor): Array<Gem> {
		return Object.values(this.gems).filter(gem => gemMatchesSocket(gem, socketColor));
	}

	getPresetEncounter(path: string): PresetEncounter | null {
		return this.presetEncounters[path] || null;
	}
	getPresetTarget(path: string): PresetTarget | null {
		return this.presetTargets[path] || null;
	}
	getAllPresetEncounters(): Array<PresetEncounter> {
		return Object.values(this.presetEncounters);
	}
	getAllPresetTargets(): Array<PresetTarget> {
		return Object.values(this.presetTargets);
	}

	getPhase(): number {
		return this.phase;
	}
	setPhase(eventID: EventID, newPhase: number) {
		if (newPhase != this.phase && newPhase > 0) {
			this.phase = newPhase;
			this.phaseChangeEmitter.emit(eventID);
		}
	}

	getFaction(): Faction {
		return this.faction;
	}
	setFaction(eventID: EventID, newFaction: Faction) {
		if (newFaction != this.faction && !!newFaction) {
			this.faction = newFaction;
			this.factionChangeEmitter.emit(eventID);
		}
	}

	getFixedRngSeed(): number {
		return this.fixedRngSeed;
	}
	setFixedRngSeed(eventID: EventID, newFixedRngSeed: number) {
		if (newFixedRngSeed != this.fixedRngSeed) {
			this.fixedRngSeed = newFixedRngSeed;
			this.fixedRngSeedChangeEmitter.emit(eventID);
		}
	}

	static MAX_RNG_SEED = Math.pow(2, 32) - 1;
	private nextRngSeed(): number {
		let rngSeed = 0;
		if (this.fixedRngSeed) {
			rngSeed = this.fixedRngSeed;
		} else {
			rngSeed = Math.floor(Math.random() * Sim.MAX_RNG_SEED);
		}

		this.lastUsedRngSeed = rngSeed;
		this.lastUsedRngSeedChangeEmitter.emit(TypedEvent.nextEventID());
		return rngSeed;
	}
	getLastUsedRngSeed(): number {
		return this.lastUsedRngSeed;
	}


	getShow1hWeapons(): boolean {
		return this.show1hWeapons;
	}
	setShow1hWeapons(eventID: EventID, newShow1hWeapons: boolean) {
		if (newShow1hWeapons != this.show1hWeapons) {
			this.show1hWeapons = newShow1hWeapons;
			this.show1hWeaponsChangeEmitter.emit(eventID);
		}
	}

	getShow2hWeapons(): boolean {
		return this.show2hWeapons;
	}
	setShow2hWeapons(eventID: EventID, newShow2hWeapons: boolean) {
		if (newShow2hWeapons != this.show2hWeapons) {
			this.show2hWeapons = newShow2hWeapons;
			this.show2hWeaponsChangeEmitter.emit(eventID);
		}
	}

	getShowMatchingGems(): boolean {
		return this.showMatchingGems;
	}
	setShowMatchingGems(eventID: EventID, newShowMatchingGems: boolean) {
		if (newShowMatchingGems != this.showMatchingGems) {
			this.showMatchingGems = newShowMatchingGems;
			this.showMatchingGemsChangeEmitter.emit(eventID);
		}
	}

	getShowThreatMetrics(): boolean {
		return this.showThreatMetrics;
	}
	setShowThreatMetrics(eventID: EventID, newShowThreatMetrics: boolean) {
		if (newShowThreatMetrics != this.showThreatMetrics) {
			this.showThreatMetrics = newShowThreatMetrics;
			this.showThreatMetricsChangeEmitter.emit(eventID);
		}
	}

	getShowExperimental(): boolean {
		return this.showExperimental;
	}
	setShowExperimental(eventID: EventID, newShowExperimental: boolean) {
		if (newShowExperimental != this.showExperimental) {
			this.showExperimental = newShowExperimental;
			this.showExperimentalChangeEmitter.emit(eventID);
		}
	}

	getIterations(): number {
		return this.iterations;
	}
	setIterations(eventID: EventID, newIterations: number) {
		if (newIterations != this.iterations) {
			this.iterations = newIterations;
			this.iterationsChangeEmitter.emit(eventID);
		}
	}

	lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null {
		const item = this.items[itemSpec.id];
		if (!item)
			return null;

		const enchant = this.enchants[itemSpec.enchant] || null;
		const gems = itemSpec.gems.map(gemId => this.gems[gemId] || null);

		return new EquippedItem(item, enchant, gems);
	}

	lookupEquipmentSpec(equipSpec: EquipmentSpec): Gear {
		// EquipmentSpec is supposed to be indexed by slot, but here we assume
		// it isn't just in case.
		const gearMap: Partial<Record<ItemSlot, EquippedItem | null>> = {};

		equipSpec.items.forEach(itemSpec => {
			const item = this.lookupItemSpec(itemSpec);
			if (!item)
				return;

			const itemSlots = getEligibleItemSlots(item.item);

			const assignedSlot = itemSlots.find(slot => !gearMap[slot]);
			if (assignedSlot == null)
				throw new Error('No slots left to equip ' + Item.toJsonString(item.item));

			gearMap[assignedSlot] = item;
		});

		return new Gear(gearMap);
	}

	toProto(): SimSettingsProto {
		return SimSettingsProto.create({
			iterations: this.getIterations(),
			phase: this.getPhase(),
			fixedRngSeed: BigInt(this.getFixedRngSeed()),
			showThreatMetrics: this.getShowThreatMetrics(),
			showExperimental: this.getShowExperimental(),
			faction: this.getFaction(),
		});
	}

	fromProto(eventID: EventID, proto: SimSettingsProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setIterations(eventID, proto.iterations || 3000);
			this.setPhase(eventID, proto.phase || OtherConstants.CURRENT_PHASE);
			this.setFixedRngSeed(eventID, Number(proto.fixedRngSeed));
			this.setShowThreatMetrics(eventID, proto.showThreatMetrics);
			this.setShowExperimental(eventID, proto.showExperimental);
			this.setFaction(eventID, proto.faction || Faction.Alliance)
		});
	}

	applyDefaults(eventID: EventID, isTankSim: boolean) {
		this.fromProto(eventID, SimSettingsProto.create({
			iterations: 3000,
			phase: OtherConstants.CURRENT_PHASE,
			faction: Faction.Alliance,
			showThreatMetrics: isTankSim,
		}));
	}
}
