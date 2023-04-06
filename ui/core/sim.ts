import { ArmorType, SimDatabase } from './proto/common.js';
import { Class, Faction } from './proto/common.js';
import { Consumes } from './proto/common.js';
import { Encounter as EncounterProto } from './proto/common.js';
import { EquipmentSpec } from './proto/common.js';
import { GemColor } from './proto/common.js';
import { ItemQuality } from './proto/common.js';
import { ItemSlot } from './proto/common.js';
import { ItemSpec } from './proto/common.js';
import { ItemType } from './proto/common.js';
import { Profession } from './proto/common.js';
import { Race } from './proto/common.js';
import { RaidTarget } from './proto/common.js';
import { Spec } from './proto/common.js';
import { Stat, PseudoStat } from './proto/common.js';
import { RangedWeaponType, WeaponType } from './proto/common.js';
import { BulkSimRequest, BulkSimResult, BulkSettings, Raid as RaidProto } from './proto/api.js';
import { ComputeStatsRequest, ComputeStatsResult } from './proto/api.js';
import { RaidSimRequest, RaidSimResult } from './proto/api.js';
import { SimOptions } from './proto/api.js';
import { StatWeightsRequest, StatWeightsResult } from './proto/api.js';
import {
	DatabaseFilters,
	SimSettings as SimSettingsProto,
	SourceFilterOption,
	RaidFilterOption,
} from './proto/ui.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from './proto/ui.js';

import { Database } from './proto_utils/database.js';
import { EquippedItem } from './proto_utils/equipped_item.js';
import { Gear } from './proto_utils/gear.js';
import { SimResult } from './proto_utils/sim_result.js';
import { Stats } from './proto_utils/stats.js';
import { SpecRotation } from './proto_utils/utils.js';
import { SpecTalents } from './proto_utils/utils.js';
import { SpecTypeFunctions } from './proto_utils/utils.js';
import { specTypeFunctions } from './proto_utils/utils.js';
import { SpecOptions } from './proto_utils/utils.js';
import { specToClass } from './proto_utils/utils.js';
import { specToEligibleRaces } from './proto_utils/utils.js';
import { getEligibleItemSlots } from './proto_utils/utils.js';
import { playerToSpec } from './proto_utils/utils.js';

import { getBrowserLanguageCode, setLanguageCode } from './constants/lang.js';
import { Encounter } from './encounter.js';
import { Player } from './player.js';
import { Raid } from './raid.js';
import { Listener } from './typed_event.js';
import { EventID, TypedEvent } from './typed_event.js';
import { getEnumValues } from './utils.js';
import { sum } from './utils.js';
import { wait } from './utils.js';
import { WorkerPool } from './worker_pool.js';

import * as OtherConstants from './constants/other.js';

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
	private filters: DatabaseFilters = Sim.defaultFilters();
	private showDamageMetrics: boolean = true;
	private showThreatMetrics: boolean = false;
	private showHealingMetrics: boolean = false;
	private showExperimental: boolean = false;
	private language: string = '';

	readonly raid: Raid;
	readonly encounter: Encounter;

	private db_: Database|null = null;

	readonly iterationsChangeEmitter = new TypedEvent<void>();
	readonly phaseChangeEmitter = new TypedEvent<void>();
	readonly factionChangeEmitter = new TypedEvent<void>();
	readonly fixedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly lastUsedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly filtersChangeEmitter = new TypedEvent<void>();
	readonly showDamageMetricsChangeEmitter = new TypedEvent<void>();
	readonly showThreatMetricsChangeEmitter = new TypedEvent<void>();
	readonly showHealingMetricsChangeEmitter = new TypedEvent<void>();
	readonly showExperimentalChangeEmitter = new TypedEvent<void>();
	readonly languageChangeEmitter = new TypedEvent<void>();
	readonly crashEmitter = new TypedEvent<SimError>();

	// Emits when any of the settings change (but not the raid / encounter).
	readonly settingsChangeEmitter: TypedEvent<void>;

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	// Fires when a raid sim API call completes.
	readonly simResultEmitter = new TypedEvent<SimResult>();

	// Fires when a bulk sim API call starts.
	readonly bulkSimStartEmitter = new TypedEvent<BulkSimRequest>();
	// Fires when a bulk sim API call completes..
	readonly bulkSimResultEmitter = new TypedEvent<BulkSimResult>();

	private readonly _initPromise: Promise<any>;
	private lastUsedRngSeed: number = 0;

	// These callbacks are needed so we can apply BuffBot modifications automatically before sending requests.
	private modifyRaidProto: ((raidProto: RaidProto) => void) = () => { };

	constructor() {
		this.workerPool = new WorkerPool(1);
		this._initPromise = Database.get().then(db => {
			this.db_ = db;
		});

		this.raid = new Raid(this);
		this.encounter = new Encounter(this);

		this.settingsChangeEmitter = TypedEvent.onAny([
			this.iterationsChangeEmitter,
			this.phaseChangeEmitter,
			this.fixedRngSeedChangeEmitter,
			this.filtersChangeEmitter,
			this.showDamageMetricsChangeEmitter,
			this.showThreatMetricsChangeEmitter,
			this.showHealingMetricsChangeEmitter,
			this.showExperimentalChangeEmitter,
			this.languageChangeEmitter,
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

	get db(): Database {
		return this.db_!;
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

				let gear = this.db.lookupEquipmentSpec(player.equipment);
				let gearChanged = false;

				const isBlacksmith = [player.profession1, player.profession2].includes(Profession.Blacksmithing);

				// Disable meta gem if inactive.
				if (gear.hasInactiveMetaGem(isBlacksmith)) {
					gear = gear.withoutMetaGem();
					gearChanged = true;
				}

				// Remove bonus sockets if not blacksmith.
				if (!isBlacksmith) {
					gear = gear.withoutBlacksmithSockets();
					gearChanged = true;
				}

				if (gearChanged) {
					player.equipment = gear.asSpec();
				}
			});
		});

		return raidProto;
	}

	private makeRaidSimRequest(debug: boolean): RaidSimRequest {
		const raid = this.getModifiedRaidProto();
		const encounter = this.encounter.toProto();

		// TODO: remove any replenishment from sim request here? probably makes more sense to do it inside the sim to protect against accidents

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

	async runBulkSim(bulkSettings: BulkSettings, bulkItemsDb: SimDatabase, onProgress: Function): Promise<BulkSimResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.getNumTargets() < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = BulkSimRequest.create({
			baseSettings: this.makeRaidSimRequest(false),
			bulkSettings: bulkSettings,
		});

		if (!request.baseSettings?.raid || request.baseSettings?.raid?.parties.length == 0 || request.baseSettings?.raid?.parties[0].players.length == 0) {
			throw new Error('Raid must contain exactly 1 player for bulk sim.');
		}

		// Attach the extra database to the player.
		const playerDatabase = request.baseSettings.raid.parties[0].players[0].database;
		playerDatabase?.items.push(...bulkItemsDb.items);
		playerDatabase?.enchants.push(...bulkItemsDb.enchants);
		playerDatabase?.gems.push(...bulkItemsDb.gems);

		this.bulkSimStartEmitter.emit(TypedEvent.nextEventID(), request);
		
		var result = await this.workerPool.bulkSimAsync(request, onProgress);
		if (result.errorResult != "") {
			throw new SimError(result.errorResult);
		}

		this.bulkSimResultEmitter.emit(TypedEvent.nextEventID(), result);
		return result;
	}

	async runRaidSim(eventID: EventID, onProgress: Function): Promise<SimResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.getNumTargets() < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = this.makeRaidSimRequest(false);

		var result = await this.workerPool.raidSimAsync(request, onProgress);
		if (result.errorResult != "") {
			throw new SimError(result.errorResult);
		}
		const simResult = await SimResult.makeNew(request, result);
		this.simResultEmitter.emit(eventID, simResult);
		return simResult;
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
		if (result.errorResult != "") {
			throw new SimError(result.errorResult);
		}
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
			this.crashEmitter.emit(eventID, new SimError(result.errorResult));
			return;
		}

		TypedEvent.freezeAllAndDo(() => {
			result.raidStats!.parties
				.forEach((partyStats, partyIndex) =>
					partyStats.players.forEach((playerStats, playerIndex) =>
						players[partyIndex * 5 + playerIndex]?.setCurrentStats(eventID, playerStats)));
		});
	}

	async statWeights(player: Player<any>, epStats: Array<Stat>, epPseudoStats: Array<PseudoStat>, epReferenceStat: Stat, onProgress: Function): Promise<StatWeightsResult> {
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
				pseudoStatsToWeigh: epPseudoStats,
				epReferenceStat: epReferenceStat,
			});
			var result = await this.workerPool.statWeightsAsync(request, onProgress);
			return result;
		}
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

	getFilters(): DatabaseFilters {
		// Make a defensive copy
		return DatabaseFilters.clone(this.filters);
	}
	setFilters(eventID: EventID, newFilters: DatabaseFilters) {
		if (DatabaseFilters.equals(newFilters, this.filters)) {
			return;
		}

		// Make a defensive copy
		this.filters = DatabaseFilters.clone(newFilters);
		this.filtersChangeEmitter.emit(eventID);
	}

	getShowDamageMetrics(): boolean {
		return this.showDamageMetrics;
	}
	setShowDamageMetrics(eventID: EventID, newShowDamageMetrics: boolean) {
		if (newShowDamageMetrics != this.showDamageMetrics) {
			this.showDamageMetrics = newShowDamageMetrics;
			this.showDamageMetricsChangeEmitter.emit(eventID);
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

	getShowHealingMetrics(): boolean {
		return this.showHealingMetrics;
	}
	setShowHealingMetrics(eventID: EventID, newShowHealingMetrics: boolean) {
		if (newShowHealingMetrics != this.showHealingMetrics) {
			this.showHealingMetrics = newShowHealingMetrics;
			this.showHealingMetricsChangeEmitter.emit(eventID);
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

	getLanguage(): string {
		return this.language;
	}
	setLanguage(eventID: EventID, newLanguage: string) {
		newLanguage = newLanguage || getBrowserLanguageCode();
		if (newLanguage != this.language) {
			this.language = newLanguage;
			setLanguageCode(this.language);
			this.languageChangeEmitter.emit(eventID);
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

	static readonly ALL_ARMOR_TYPES = (getEnumValues(ArmorType) as Array<ArmorType>).filter(v => v != 0);
	static readonly ALL_WEAPON_TYPES = (getEnumValues(WeaponType) as Array<WeaponType>).filter(v => v != 0);
	static readonly ALL_RANGED_WEAPON_TYPES = (getEnumValues(RangedWeaponType) as Array<RangedWeaponType>).filter(v => v != 0);
	static readonly ALL_SOURCES = (getEnumValues(SourceFilterOption) as Array<SourceFilterOption>).filter(v => v != 0);
	static readonly ALL_RAIDS = (getEnumValues(RaidFilterOption) as Array<RaidFilterOption>).filter(v => v != 0);

	toProto(): SimSettingsProto {
		const filters = this.getFilters();
		if (filters.armorTypes.length == Sim.ALL_ARMOR_TYPES.length) {
			filters.armorTypes = [];
		}
		if (filters.weaponTypes.length == Sim.ALL_WEAPON_TYPES.length) {
			filters.weaponTypes = [];
		}
		if (filters.rangedWeaponTypes.length == Sim.ALL_RANGED_WEAPON_TYPES.length) {
			filters.rangedWeaponTypes = [];
		}
		if (filters.sources.length == Sim.ALL_SOURCES.length) {
			filters.sources = [];
		}
		if (filters.raids.length == Sim.ALL_RAIDS.length) {
			filters.raids = [];
		}

		return SimSettingsProto.create({
			iterations: this.getIterations(),
			phase: this.getPhase(),
			fixedRngSeed: BigInt(this.getFixedRngSeed()),
			showDamageMetrics: this.getShowDamageMetrics(),
			showThreatMetrics: this.getShowThreatMetrics(),
			showHealingMetrics: this.getShowHealingMetrics(),
			showExperimental: this.getShowExperimental(),
			language: this.getLanguage(),
			faction: this.getFaction(),
			filters: filters,
		});
	}

	fromProto(eventID: EventID, proto: SimSettingsProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setIterations(eventID, proto.iterations || 3000);
			this.setPhase(eventID, proto.phase || OtherConstants.CURRENT_PHASE);
			this.setFixedRngSeed(eventID, Number(proto.fixedRngSeed));
			this.setShowDamageMetrics(eventID, proto.showDamageMetrics);
			this.setShowThreatMetrics(eventID, proto.showThreatMetrics);
			this.setShowHealingMetrics(eventID, proto.showHealingMetrics);
			this.setShowExperimental(eventID, proto.showExperimental);
			this.setLanguage(eventID, proto.language);
			this.setFaction(eventID, proto.faction || Faction.Alliance)

			const filters = proto.filters || Sim.defaultFilters();
			if (filters.armorTypes.length == 0) {
				filters.armorTypes = Sim.ALL_ARMOR_TYPES.slice();
			}
			if (filters.weaponTypes.length == 0) {
				filters.weaponTypes = Sim.ALL_WEAPON_TYPES.slice();
			}
			if (filters.rangedWeaponTypes.length == 0) {
				filters.rangedWeaponTypes = Sim.ALL_RANGED_WEAPON_TYPES.slice();
			}
			if (filters.sources.length == 0) {
				filters.sources = Sim.ALL_SOURCES.slice();
			}
			if (filters.raids.length == 0) {
				filters.raids = Sim.ALL_RAIDS.slice();
			}
			this.setFilters(eventID, filters);
		});
	}

	applyDefaults(eventID: EventID, isTankSim: boolean, isHealingSim: boolean) {
		this.fromProto(eventID, SimSettingsProto.create({
			iterations: 3000,
			phase: OtherConstants.CURRENT_PHASE,
			faction: Faction.Alliance,
			showDamageMetrics: !isHealingSim,
			showThreatMetrics: isTankSim,
			showHealingMetrics: isHealingSim,
			language: this.getLanguage(), // Don't change language.
			filters: Sim.defaultFilters(),
		}));
	}

	static defaultFilters(): DatabaseFilters {
		return DatabaseFilters.create({
			oneHandedWeapons: true,
			twoHandedWeapons: true,
		});
	}
}

export class SimError extends Error {
	readonly errorStr: string;

	constructor(errorStr: string) {
		super(errorStr);
		this.errorStr = errorStr;
	}
}
