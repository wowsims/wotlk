import {
	AuraStats as AuraStatsProto,
	SpellStats as SpellStatsProto,
	UnitMetadata as UnitMetadataProto,
} from './proto/api.js';
import {
	APLRotation,
	APLRotation_Type as APLRotationType,
	APLValue,
	SimpleRotation,
} from './proto/apl.js';
import {
	Class,
	Consumes,
	Cooldowns,
	Faction,
	Glyphs,
	HandType,
	HealingModel,
	IndividualBuffs,
	ItemSlot,
	Profession,
	PseudoStat,
	Race,
	SimDatabase,
	Spec,
	Stat,
	UnitReference,
	UnitReference_Type,
	UnitStats,
} from './proto/common.js';
import {
	DungeonDifficulty,
	UIEnchant as Enchant,
	Expansion,
	UIItem as Item,
	RaidFilterOption,
	UIRune as Rune,
	SourceFilterOption,
	UIItem_FactionRestriction,
} from './proto/ui.js';

import { LaunchStatus, aplLaunchStatuses } from './launched_sims';
import { Player as PlayerProto, PlayerStats, StatWeightsResult } from './proto/api.js';
import { ActionId } from './proto_utils/action_id.js';
import { EquippedItem, getWeaponDPS } from './proto_utils/equipped_item.js';

import { Gear, ItemSwapGear } from './proto_utils/gear.js';
import { Stats } from './proto_utils/stats.js';
import { playerTalentStringToProto } from './talents/factory.js';

import {
	AL_CATEGORY_HARD_MODE,
	ClassSpecs,
	ShamanSpecs,
	SpecOptions,
	SpecRotation,
	SpecTalents,
	SpecTypeFunctions,
	canEquipEnchant,
	canEquipItem,
	classColors,
	emptyUnitReference,
	enchantAppliesToItem,
	getTalentTree,
	getTalentTreeIcon,
	getTalentTreePoints,
	isTankSpec,
	newUnitReference,
	raceToFaction,
	specToClass,
	specToEligibleRaces,
	specTypeFunctions,
	withSpecProto,
} from './proto_utils/utils.js';

import { getLanguageCode } from './constants/lang.js';
import * as Mechanics from './constants/mechanics.js';
import { MAX_PARTY_SIZE, Party } from './party.js';
import { ElementalShaman_Options, ElementalShaman_Options_ThunderstormRange, ElementalShaman_Rotation, ElementalShaman_Rotation_BloodlustUse, EnhancementShaman_Rotation, EnhancementShaman_Rotation_BloodlustUse, RestorationShaman_Rotation, RestorationShaman_Rotation_BloodlustUse } from './proto/shaman.js';
import { Raid } from './raid.js';
import { Sim } from './sim.js';
import { EventID, TypedEvent } from './typed_event.js';
import { stringComparator } from './utils.js';

export interface AuraStats {
	data: AuraStatsProto,
	id: ActionId,
}
export interface SpellStats {
	data: SpellStatsProto,
	id: ActionId,
}

export class UnitMetadata {
	private name: string;
	private auras: Array<AuraStats>;
	private spells: Array<SpellStats>;

	constructor() {
		this.name = '';
		this.auras = [];
		this.spells = [];
	}

	getName(): string {
		return this.name;
	}

	getAuras(): Array<AuraStats> {
		return this.auras.slice();
	}

	getSpells(): Array<SpellStats> {
		return this.spells.slice();
	}

	// Returns whether any updates were made.
	async update(metadata: UnitMetadataProto): Promise<boolean> {
		let newSpells = metadata!.spells.map(spell => {
			return {
				data: spell,
				id: ActionId.fromProto(spell.id!),
			};
		});
		let newAuras = metadata!.auras.map(aura => {
			return {
				data: aura,
				id: ActionId.fromProto(aura.id!),
			};
		});

		await Promise.all([...newSpells, ...newAuras].map(newSpell => newSpell.id.fill().then(newId => newSpell.id = newId)));

		newSpells = newSpells.sort((a, b) => stringComparator(a.id.name, b.id.name))
		newAuras = newAuras.sort((a, b) => stringComparator(a.id.name, b.id.name))

		let anyUpdates = false;
		if (metadata.name != this.name) {
			this.name = metadata.name;
			anyUpdates = true;
		}
		if (newSpells.length != this.spells.length || newSpells.some((newSpell, i) => !newSpell.id.equals(this.spells[i].id))) {
			this.spells = newSpells;
			anyUpdates = true;
		}
		if (newAuras.length != this.auras.length || newAuras.some((newAura, i) => !newAura.id.equals(this.auras[i].id))) {
			this.auras = newAuras;
			anyUpdates = true;
		}

		return anyUpdates;
	}
}

export class UnitMetadataList {
	private metadatas: Array<UnitMetadata>;

	constructor() {
		this.metadatas = [];
	}

	async update(newMetadatas: Array<UnitMetadataProto>): Promise<boolean> {
		const oldLen = this.metadatas.length;

		if (newMetadatas.length > oldLen) {
			for (let i = oldLen; i < newMetadatas.length; i++) {
				this.metadatas.push(new UnitMetadata());
			}
		} else if (newMetadatas.length < oldLen) {
			this.metadatas = this.metadatas.slice(0, newMetadatas.length);
		}

		const anyUpdates = await Promise.all(newMetadatas.map((metadata, i) => this.metadatas[i].update(metadata)));

		return oldLen != this.metadatas.length || anyUpdates.some(v => v);
	}

	asList(): Array<UnitMetadata> {
		return this.metadatas.slice();
	}
}

export interface MeleeCritCapInfo {
	meleeCrit: number,
	meleeHit: number,
	expertise: number,
	suppression: number,
	glancing: number,
	debuffCrit: number,
	hasOffhandWeapon: boolean,
	meleeHitCap: number,
	expertiseCap: number,
	remainingMeleeHitCap: number,
	remainingExpertiseCap: number,
	baseCritCap: number,
	specSpecificOffset: number,
	playerCritCapDelta: number
}

export type AutoRotationGenerator<SpecType extends Spec> = (player: Player<SpecType>) => APLRotation;
export type SimpleRotationGenerator<SpecType extends Spec> = (player: Player<SpecType>, simpleRotation: SpecRotation<SpecType>, cooldowns: Cooldowns) => APLRotation;

// Manages all the gear / consumes / other settings for a single Player.
export class Player<SpecType extends Spec> {
	readonly sim: Sim;
	private party: Party | null;
	private raid: Raid | null;

	readonly spec: Spec;
	private name: string = '';
	private buffs: IndividualBuffs = IndividualBuffs.create();
	private consumes: Consumes = Consumes.create();
	private bonusStats: Stats = new Stats();
	private gear: Gear = new Gear({});
	//private bulkEquipmentSpec: BulkEquipmentSpec = BulkEquipmentSpec.create();
	private itemSwapGear: ItemSwapGear = new ItemSwapGear();
	private race: Race;
	private level: number;
	private profession1: Profession = 0;
	private profession2: Profession = 0;
	private rotation: SpecRotation<SpecType>;
	aplRotation: APLRotation = APLRotation.create();
	private talentsString: string = '';
	private glyphs: Glyphs = Glyphs.create();
	private specOptions: SpecOptions<SpecType>;
	private cooldowns: Cooldowns = Cooldowns.create();
	private reactionTime: number = 0;
	private channelClipDelay: number = 0;
	private inFrontOfTarget: boolean = false;
	private distanceFromTarget: number = 0;
	private healingModel: HealingModel = HealingModel.create();
	private healingEnabled: boolean = false;

	private autoRotationGenerator: AutoRotationGenerator<SpecType> | null = null;
	private simpleRotationGenerator: SimpleRotationGenerator<SpecType> | null = null;

	private itemEPCache = new Array<Map<number, number>>();
	private gemEPCache = new Map<number, number>();
	private enchantEPCache = new Map<number, number>();
	private talents: SpecTalents<SpecType> | null = null;

	readonly specTypeFunctions: SpecTypeFunctions<SpecType>;

	private static readonly numEpRatios = 6;
	private epRatios: Array<number> = new Array<number>(Player.numEpRatios).fill(0);
	private epWeights: Stats = new Stats();
	private currentStats: PlayerStats = PlayerStats.create();
	private metadata: UnitMetadata = new UnitMetadata();
	private petMetadatas: UnitMetadataList = new UnitMetadataList();

	readonly nameChangeEmitter = new TypedEvent<void>('PlayerName');
	readonly buffsChangeEmitter = new TypedEvent<void>('PlayerBuffs');
	readonly consumesChangeEmitter = new TypedEvent<void>('PlayerConsumes');
	readonly bonusStatsChangeEmitter = new TypedEvent<void>('PlayerBonusStats');
	readonly gearChangeEmitter = new TypedEvent<void>('PlayerGear');
	readonly professionChangeEmitter = new TypedEvent<void>('PlayerProfession');
	readonly raceChangeEmitter = new TypedEvent<void>('PlayerRace');
	readonly levelChangeEmitter = new TypedEvent<void>('PlayerLevel');
	readonly rotationChangeEmitter = new TypedEvent<void>('PlayerRotation');
	readonly talentsChangeEmitter = new TypedEvent<void>('PlayerTalents');
	readonly glyphsChangeEmitter = new TypedEvent<void>('PlayerGlyphs');
	readonly specOptionsChangeEmitter = new TypedEvent<void>('PlayerSpecOptions');
	readonly cooldownsChangeEmitter = new TypedEvent<void>('PlayerCooldowns');
	readonly inFrontOfTargetChangeEmitter = new TypedEvent<void>('PlayerInFrontOfTarget');
	readonly distanceFromTargetChangeEmitter = new TypedEvent<void>('PlayerDistanceFromTarget');
	readonly healingModelChangeEmitter = new TypedEvent<void>('PlayerHealingModel');
	readonly epWeightsChangeEmitter = new TypedEvent<void>('PlayerEpWeights');
	readonly miscOptionsChangeEmitter = new TypedEvent<void>('PlayerMiscOptions');

	readonly currentStatsEmitter = new TypedEvent<void>('PlayerCurrentStats');
	readonly epRatiosChangeEmitter = new TypedEvent<void>('PlayerEpRatios');
	readonly epRefStatChangeEmitter = new TypedEvent<void>('PlayerEpRefStat');

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	constructor(spec: Spec, sim: Sim) {
		this.sim = sim;
		this.party = null;
		this.raid = null;

		this.spec = spec;
		this.race = specToEligibleRaces[this.spec][0];
		this.level = Mechanics.MAX_CHARACTER_LEVEL;
		this.specTypeFunctions = specTypeFunctions[this.spec] as SpecTypeFunctions<SpecType>;
		this.rotation = this.specTypeFunctions.rotationCreate();
		this.specOptions = this.specTypeFunctions.optionsCreate();

		for(let i = 0; i < ItemSlot.ItemSlotRanged+1; ++i) {
			this.itemEPCache[i] = new Map();
		}

		this.changeEmitter = TypedEvent.onAny([
			this.nameChangeEmitter,
			this.buffsChangeEmitter,
			this.consumesChangeEmitter,
			this.bonusStatsChangeEmitter,
			this.gearChangeEmitter,
			this.professionChangeEmitter,
			this.raceChangeEmitter,
			this.levelChangeEmitter,
			this.rotationChangeEmitter,
			this.talentsChangeEmitter,
			this.glyphsChangeEmitter,
			this.specOptionsChangeEmitter,
			this.cooldownsChangeEmitter,
			this.miscOptionsChangeEmitter,
			this.inFrontOfTargetChangeEmitter,
			this.distanceFromTargetChangeEmitter,
			this.healingModelChangeEmitter,
			this.epWeightsChangeEmitter,
			this.epRatiosChangeEmitter,
			this.epRefStatChangeEmitter,
		], 'PlayerChange');
	}

	getSpecIcon(): string {
		return getTalentTreeIcon(this.spec, this.getTalentsString());
	}

	getClass(): Class {
		return specToClass[this.spec];
	}

	getClassColor(): string {
		return classColors[this.getClass()];
	}

	isSpec<T extends Spec>(spec: T): this is Player<T> {
		return this.spec == spec;
	}
	isClass<T extends Class>(clazz: T): this is Player<ClassSpecs<T>> {
		return this.getClass() == clazz;
	}

	getParty(): Party | null {
		return this.party;
	}

	getRaid(): Raid | null {
		return this.raid;
	}

	// Returns this player's index within its party [0-4].
	getPartyIndex(): number {
		if (this.party == null) {
			throw new Error('Can\'t get party index for player without a party!');
		}

		return this.party.getPlayers().indexOf(this);
	}

	// Returns this player's index within its raid [0-24].
	getRaidIndex(): number {
		if (this.party == null) {
			throw new Error('Can\'t get raid index for player without a party!');
		}

		return this.party.getIndex() * MAX_PARTY_SIZE + this.getPartyIndex();
	}

	// This should only ever be called from party.
	setParty(newParty: Party | null) {
		if (newParty == null) {
			this.party = null;
			this.raid = null;
		} else {
			this.party = newParty;
			this.raid = newParty.raid;
		}
	}

	getOtherPartyMembers(): Array<Player<any>> {
		if (this.party == null) {
			return [];
		}

		return this.party.getPlayers().filter(player => player != null && player != this) as Array<Player<any>>;
	}

	// Returns all items that this player can wear in the given slot.
	getItems(slot: ItemSlot): Array<Item> {
		return this.sim.db.getItems(slot).filter(item => canEquipItem(this, item, slot));
	}

	// Returns all enchants that this player can wear in the given slot.
	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.sim.db.getEnchants(slot).filter(enchant => canEquipEnchant(enchant, this.spec));
	}

	getRunes(slot: ItemSlot): Array<Rune> {
		return this.sim.db.getRunes(slot, this.getClass());
	}

	getEpWeights(): Stats {
		return this.epWeights;
	}

	setEpWeights(eventID: EventID, newEpWeights: Stats) {
		this.epWeights = newEpWeights;
		this.epWeightsChangeEmitter.emit(eventID);

		this.gemEPCache = new Map();
		this.enchantEPCache = new Map();
		for(let i = 0; i < ItemSlot.ItemSlotRanged+1; ++i) {
			this.itemEPCache[i] = new Map();
		}
	}

	getDefaultEpRatios(isTankSpec: boolean, isHealingSpec: boolean): Array<number> {
		const defaultRatios = new Array(Player.numEpRatios).fill(0);
		if (isHealingSpec) {
			// By default only value HPS EP for healing spec
			defaultRatios[1] = 1;
		} else if (isTankSpec) {
			// By default value TPS and DTPS EP equally for tanking spec
			defaultRatios[2] = 1;
			defaultRatios[3] = 1;
		} else {
			// By default only value DPS EP
			defaultRatios[0] = 1;
		}
		return defaultRatios;
	}

	getEpRatios() {
		return this.epRatios.slice();
	}

	setEpRatios(eventID: EventID, newRatios: Array<number>) {
		this.epRatios = newRatios;
		this.epRatiosChangeEmitter.emit(eventID);
	}

	async computeStatWeights(eventID: EventID, epStats: Array<Stat>, epPseudoStats: Array<PseudoStat>, epReferenceStat: Stat, onProgress: Function): Promise<StatWeightsResult> {
		const result = await this.sim.statWeights(this, epStats, epPseudoStats, epReferenceStat, onProgress);
		return result;
	}

	getCurrentStats(): PlayerStats {
		return PlayerStats.clone(this.currentStats);
	}

	setCurrentStats(eventID: EventID, newStats: PlayerStats) {
		this.currentStats = newStats;
		this.currentStatsEmitter.emit(eventID);
	}

	getMetadata(): UnitMetadata {
		return this.metadata;
	}

	getPetMetadatas(): UnitMetadataList {
		return this.petMetadatas;
	}

	async updateMetadata(): Promise<boolean> {
		const playerPromise = this.metadata.update(this.currentStats.metadata!);
		const petsPromise = this.petMetadatas.update(this.currentStats.pets.map(p => p.metadata!));
		const playerUpdated = await playerPromise;
		const petsUpdated = await petsPromise;
		return playerUpdated || petsUpdated;
	}

	getName(): string {
		return this.name;
	}
	setName(eventID: EventID, newName: string) {
		if (newName != this.name) {
			this.name = newName;
			this.nameChangeEmitter.emit(eventID);
		}
	}

	getLabel(): string {
		if (this.party) {
			return `${this.name} (#${this.getRaidIndex() + 1})`;
		} else {
			return this.name;
		}
	}

	getRace(): Race {
		return this.race;
	}
	setRace(eventID: EventID, newRace: Race) {
		if (newRace != this.race) {
			this.race = newRace;
			this.raceChangeEmitter.emit(eventID);
		}
	}
	getLevel(): number {
		return this.level;
	}
	setLevel(eventID: EventID, newLevel: number){
		if (newLevel != this.level){
			this.level = newLevel;
			this.levelChangeEmitter.emit(eventID)
		}
	}

	getProfession1(): Profession {
		return this.profession1;
	}
	setProfession1(eventID: EventID, newProfession: Profession) {
		if (newProfession != this.profession1) {
			this.profession1 = newProfession;
			this.professionChangeEmitter.emit(eventID);
		}
	}
	getProfession2(): Profession {
		return this.profession2;
	}
	setProfession2(eventID: EventID, newProfession: Profession) {
		if (newProfession != this.profession2) {
			this.profession2 = newProfession;
			this.professionChangeEmitter.emit(eventID);
		}
	}
	getProfessions(): Array<Profession> {
		return [this.profession1, this.profession2].filter(p => p != Profession.ProfessionUnknown);
	}
	setProfessions(eventID: EventID, newProfessions: Array<Profession>) {
		TypedEvent.freezeAllAndDo(() => {
			this.setProfession1(eventID, newProfessions[0] || Profession.ProfessionUnknown);
			this.setProfession2(eventID, newProfessions[1] || Profession.ProfessionUnknown);
		});
	}
	hasProfession(prof: Profession): boolean {
		return this.getProfessions().includes(prof);
	}
	isBlacksmithing(): boolean {
		return this.hasProfession(Profession.Blacksmithing);
	}

	getFaction(): Faction {
		return raceToFaction[this.getRace()];
	}

	getBuffs(): IndividualBuffs {
		// Make a defensive copy
		return IndividualBuffs.clone(this.buffs);
	}

	setBuffs(eventID: EventID, newBuffs: IndividualBuffs) {
		if (IndividualBuffs.equals(this.buffs, newBuffs))
			return;

		// Make a defensive copy
		this.buffs = IndividualBuffs.clone(newBuffs);
		this.buffsChangeEmitter.emit(eventID);
	}

	getConsumes(): Consumes {
		// Make a defensive copy
		return Consumes.clone(this.consumes);
	}

	setConsumes(eventID: EventID, newConsumes: Consumes) {
		if (Consumes.equals(this.consumes, newConsumes))
			return;

		// Make a defensive copy
		this.consumes = Consumes.clone(newConsumes);
		this.consumesChangeEmitter.emit(eventID);
	}

	getCooldowns(): Cooldowns {
		// Make a defensive copy
		if (aplLaunchStatuses[this.spec] == LaunchStatus.Launched) {
			return Cooldowns.clone(this.aplRotation.simple?.cooldowns || Cooldowns.create());
		} else {
			return Cooldowns.clone(this.cooldowns);
		}
	}

	setCooldowns(eventID: EventID, newCooldowns: Cooldowns) {
		if (Cooldowns.equals(this.getCooldowns(), newCooldowns))
			return;

		if (aplLaunchStatuses[this.spec] == LaunchStatus.Launched) {
			if (!this.aplRotation.simple) {
				this.aplRotation.simple = SimpleRotation.create();
			}
			this.aplRotation.simple.cooldowns = newCooldowns;
			this.rotationChangeEmitter.emit(eventID);
		} else {
			// Make a defensive copy
			this.cooldowns = Cooldowns.clone(newCooldowns);
		}

		this.cooldownsChangeEmitter.emit(eventID);
	}

	canDualWield2H(): boolean {
		return false;
	}

	equipItem(eventID: EventID, slot: ItemSlot, newItem: EquippedItem | null) {
		this.setGear(eventID, this.gear.withEquippedItem(slot, newItem, this.canDualWield2H()));
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear.getEquippedItem(slot);
	}

	getGear(): Gear {
		return this.gear;
	}

	getItemSwapGear(): ItemSwapGear {
		return this.itemSwapGear;
	}

	setGear(eventID: EventID, newGear: Gear) {
		if (newGear.equals(this.gear))
			return;

		// Commented for now because the UI for this is weird.
		//// If trinkets have changed and there were cooldowns assigned for those trinkets,
		//// try to match them up and switch to the new trinkets.
		//const newCooldowns = this.getCooldowns();
		//const oldTrinketIds = this.gear.getTrinkets().map(trinket => trinket?.asActionIdProto() || ActionIdProto.create());
		//const newTrinketIds = newGear.getTrinkets().map(trinket => trinket?.asActionIdProto() || ActionIdProto.create());

		//for (let i = 0; i < 2; i++) {
		//	const oldTrinketId = oldTrinketIds[i];
		//	const newTrinketId = newTrinketIds[i];
		//	if (ActionIdProto.equals(oldTrinketId, ActionIdProto.create())) {
		//		continue;
		//	}
		//	if (ActionIdProto.equals(newTrinketId, ActionIdProto.create())) {
		//		continue;
		//	}
		//	if (ActionIdProto.equals(oldTrinketId, newTrinketId)) {
		//		continue;
		//	}
		//	newCooldowns.cooldowns.forEach(cd => {
		//		if (ActionIdProto.equals(cd.id, oldTrinketId)) {
		//			cd.id = newTrinketId;
		//		}
		//	});
		//}

		TypedEvent.freezeAllAndDo(() => {
			this.gear = newGear;
			this.gearChangeEmitter.emit(eventID);
			//this.setCooldowns(eventID, newCooldowns);
		});
	}

	/*
	setBulkEquipmentSpec(eventID: EventID, newBulkEquipmentSpec: BulkEquipmentSpec) {
		if (BulkEquipmentSpec.equals(this.bulkEquipmentSpec, newBulkEquipmentSpec))
			return;

		TypedEvent.freezeAllAndDo(() => {
			this.bulkEquipmentSpec = newBulkEquipmentSpec;
			this.bulkGearChangeEmitter.emit(eventID);
		});
	}

	getBulkEquipmentSpec(): BulkEquipmentSpec {
		return BulkEquipmentSpec.clone(this.bulkEquipmentSpec);
	}
	*/

	getBonusStats(): Stats {
		return this.bonusStats;
	}

	setBonusStats(eventID: EventID, newBonusStats: Stats) {
		if (newBonusStats.equals(this.bonusStats))
			return;

		this.bonusStats = newBonusStats;
		this.bonusStatsChangeEmitter.emit(eventID);
	}

	getMeleeCritCapInfo(): MeleeCritCapInfo {
		const meleeCrit = (this.currentStats.finalStats?.stats[Stat.StatMeleeCrit] || 0.0) / Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE;
		const meleeHit = (this.currentStats.finalStats?.stats[Stat.StatMeleeHit] || 0.0) / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE;
		const expertise = (this.currentStats.finalStats?.stats[Stat.StatExpertise] || 0.0) / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4;
		const suppression = 4.8;
		const glancing = 24.0;

		const hasOffhandWeapon = this.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed !== undefined;
		// Due to warrior HS bug, hit cap for crit cap calculation should be 8% instead of 27%
		const meleeHitCap = hasOffhandWeapon && this.spec != Spec.SpecWarrior ? 27.0 : 8.0;
		const dodgeCap = 6.5
		const parryCap = this.getInFrontOfTarget() ? 14.0 : 0
		const expertiseCap = dodgeCap + parryCap

		const remainingMeleeHitCap = Math.max(meleeHitCap - meleeHit, 0.0);
		const remainingDodgeCap = Math.max(dodgeCap - expertise, 0.0)
		const remainingParryCap = Math.max(parryCap - expertise, 0.0)
		const remainingExpertiseCap = remainingDodgeCap + remainingParryCap

		let specSpecificOffset = 0.0;

		if(this.spec === Spec.SpecEnhancementShaman) {
			// Elemental Devastation uptime is near 100%
			const ranks = (this as Player<Spec.SpecEnhancementShaman>).getTalents().elementalDevastation;
			specSpecificOffset = 3.0 * ranks;
		}

		let debuffCrit = 0.0;

		this.sim.raid.getDebuffs();

		const baseCritCap = 100.0 - glancing + suppression - remainingMeleeHitCap - remainingExpertiseCap - specSpecificOffset;
		const playerCritCapDelta = meleeCrit - baseCritCap + debuffCrit;

		return {
			meleeCrit,
			meleeHit,
			expertise,
			suppression,
			glancing,
			debuffCrit,
			hasOffhandWeapon,
			meleeHitCap,
			expertiseCap,
			remainingMeleeHitCap,
			remainingExpertiseCap,
			baseCritCap,
			specSpecificOffset,
			playerCritCapDelta
		};
	}

	getMeleeCritCap() {
		return this.getMeleeCritCapInfo().playerCritCapDelta
	}

	getRotation(): SpecRotation<SpecType> {
		if (aplLaunchStatuses[this.spec] == LaunchStatus.Launched) {
			const jsonStr = this.aplRotation.simple?.specRotationJson || '';
			if (!jsonStr) {
				return this.specTypeFunctions.rotationCreate();
			}

			try {
				const json = JSON.parse(jsonStr);
				return this.specTypeFunctions.rotationFromJson(json);
			} catch (e) {
				console.warn(`Error parsing rotation spec options: ${e}\n\nSpec options: '${jsonStr}'`);
				return this.specTypeFunctions.rotationCreate();
			}
		} else {
			return this.specTypeFunctions.rotationCopy(this.rotation);
		}
	}

	setRotation(eventID: EventID, newRotation: SpecRotation<SpecType>) {
		if (this.specTypeFunctions.rotationEquals(newRotation, this.getRotation()))
			return;

		if (!this.aplRotation.simple) {
			this.aplRotation.simple = SimpleRotation.create();
		}
		this.aplRotation.simple.specRotationJson = JSON.stringify(this.specTypeFunctions.rotationToJson(newRotation));

		this.rotationChangeEmitter.emit(eventID);
	}

	setAplRotation(eventID: EventID, newRotation: APLRotation) {
		if (APLRotation.equals(newRotation, this.aplRotation))
			return;

		this.aplRotation = APLRotation.clone(newRotation);
		this.rotationChangeEmitter.emit(eventID);
	}

	getRotationType(): APLRotationType {
		if (this.aplRotation.enabled) {
			return APLRotationType.TypeAPL;
		} else if (this.aplRotation.type == APLRotationType.TypeUnknown) {
			return APLRotationType.TypeLegacy;
		} else {
			return this.aplRotation.type;
		}
	}

	setAutoRotationGenerator(generator: AutoRotationGenerator<SpecType>) {
		this.autoRotationGenerator = generator;
	}

	setSimpleRotationGenerator(generator: SimpleRotationGenerator<SpecType>) {
		this.simpleRotationGenerator = generator;
	}

	hasSimpleRotationGenerator(): boolean {
		return this.simpleRotationGenerator != null;
	}

	getResolvedAplRotation(): APLRotation {
		const type = this.getRotationType();
		if (type == APLRotationType.TypeAuto && this.autoRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			return APLRotation.clone(this.autoRotationGenerator(this));
		} else if (type == APLRotationType.TypeSimple && this.simpleRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			const rot = APLRotation.clone(this.simpleRotationGenerator(this, this.getRotation(), this.getCooldowns()));
			rot.type = APLRotationType.TypeAPL; // Set this here for convenience, so the generator functions don't need to.
			return rot;
		} else {
			return this.aplRotation;
		}
	}

	getTalents(): SpecTalents<SpecType> {
		if (this.talents == null) {
			this.talents = playerTalentStringToProto(this.spec, this.talentsString) as SpecTalents<SpecType>;
		}
		return this.talents!;
	}

	getTalentsString(): string {
		return this.talentsString;
	}

	setTalentsString(eventID: EventID, newTalentsString: string) {
		if (newTalentsString == this.talentsString)
			return;

		this.talentsString = newTalentsString;
		this.talents = null;
		this.talentsChangeEmitter.emit(eventID);
	}

	getTalentTree(): number {
		return getTalentTree(this.getTalentsString());
	}

	getTalentTreePoints(): Array<number> {
		return getTalentTreePoints(this.getTalentsString())
	}

	getTalentTreeIcon(): string {
		return getTalentTreeIcon(this.spec, this.getTalentsString());
	}

	getGlyphs(): Glyphs {
		// Make a defensive copy
		return Glyphs.clone(this.glyphs);
	}

	setGlyphs(eventID: EventID, newGlyphs: Glyphs) {
		if (Glyphs.equals(this.glyphs, newGlyphs))
			return;

		// Make a defensive copy
		this.glyphs = Glyphs.clone(newGlyphs);
		this.glyphsChangeEmitter.emit(eventID);
	}

	getMajorGlyphs(): Array<number> {
		return [
			this.glyphs.major1,
			this.glyphs.major2,
			this.glyphs.major3,
		].filter(glyph => glyph != 0);
	}

	getMinorGlyphs(): Array<number> {
		return [
			this.glyphs.minor1,
			this.glyphs.minor2,
			this.glyphs.minor3,
		].filter(glyph => glyph != 0);
	}

	getAllGlyphs(): Array<number> {
		return this.getMajorGlyphs().concat(this.getMinorGlyphs());
	}

	getSpecOptions(): SpecOptions<SpecType> {
		return this.specTypeFunctions.optionsCopy(this.specOptions);
	}

	setSpecOptions(eventID: EventID, newSpecOptions: SpecOptions<SpecType>) {
		if (this.specTypeFunctions.optionsEquals(newSpecOptions, this.specOptions))
			return;

		this.specOptions = this.specTypeFunctions.optionsCopy(newSpecOptions);
		this.specOptionsChangeEmitter.emit(eventID);
	}

	getReactionTime(): number {
		return this.reactionTime;
	}

	setReactionTime(eventID: EventID, newReactionTime: number) {
		if (newReactionTime == this.reactionTime)
			return;

		this.reactionTime = newReactionTime;
		this.miscOptionsChangeEmitter.emit(eventID);
	}

	getChannelClipDelay(): number {
		return this.channelClipDelay;
	}

	setChannelClipDelay(eventID: EventID, newChannelClipDelay: number) {
		if (newChannelClipDelay == this.channelClipDelay)
			return;

		this.channelClipDelay = newChannelClipDelay;
		this.miscOptionsChangeEmitter.emit(eventID);
	}

	getInFrontOfTarget(): boolean {
		return this.inFrontOfTarget;
	}

	setInFrontOfTarget(eventID: EventID, newInFrontOfTarget: boolean) {
		if (newInFrontOfTarget == this.inFrontOfTarget)
			return;

		this.inFrontOfTarget = newInFrontOfTarget;
		this.inFrontOfTargetChangeEmitter.emit(eventID);
	}

	getDistanceFromTarget(): number {
		return this.distanceFromTarget;
	}

	setDistanceFromTarget(eventID: EventID, newDistanceFromTarget: number) {
		if (newDistanceFromTarget == this.distanceFromTarget)
			return;

		this.distanceFromTarget = newDistanceFromTarget;
		this.distanceFromTargetChangeEmitter.emit(eventID);
	}

	setDefaultHealingParams(hm: HealingModel) {
		var boss = this.sim.encounter.primaryTarget;
		var dualWield = boss.dualWield;
		if (hm.cadenceSeconds == 0) {
			hm.cadenceSeconds = 1.5 * boss.swingSpeed;
			if (dualWield) {
				hm.cadenceSeconds /= 2;
			}
		}
		if (hm.hps == 0) {
			hm.hps = 0.175 * boss.minBaseDamage / boss.swingSpeed;
			if (dualWield) {
				hm.hps *= 1.5;
			}
		}
	}

	enableHealing() {
		this.healingEnabled = true;
		var hm = this.getHealingModel();
		if (hm.cadenceSeconds == 0 || hm.hps == 0) {
			this.setDefaultHealingParams(hm)
			this.setHealingModel(0, hm)
		}
	}

	getHealingModel(): HealingModel {
		// Make a defensive copy
		return HealingModel.clone(this.healingModel);
	}

	setHealingModel(eventID: EventID, newHealingModel: HealingModel) {
		if (HealingModel.equals(this.healingModel, newHealingModel))
			return;

		// Make a defensive copy
		this.healingModel = HealingModel.clone(newHealingModel);
		// If we have enabled healing model and try to set 0s cadence or 0 incoming HPS, then set intelligent defaults instead based on boss parameters.
		if (this.healingEnabled) {
			this.setDefaultHealingParams(this.healingModel)
		}
		this.healingModelChangeEmitter.emit(eventID);
	}

	computeStatsEP(stats?: Stats): number {
		if (stats == undefined) {
			return 0;
		}
		return stats.computeEP(this.epWeights);
	}

	computeEnchantEP(enchant: Enchant): number {
		if (this.enchantEPCache.has(enchant.effectId)) {
			return this.enchantEPCache.get(enchant.effectId)!;
		}

		let ep = this.computeStatsEP(new Stats(enchant.stats));
		this.enchantEPCache.set(enchant.effectId, ep);
		return ep
	}

	computeItemEP(item: Item, slot: ItemSlot): number {
		if (item == null)
			return 0;

		let cached = this.itemEPCache[slot].get(item.id);
		if (cached !== undefined)
			return cached;

		let itemStats = new Stats(item.stats);
		if (item.weaponSpeed > 0) {
			const weaponDps = getWeaponDPS(item);
			if (slot == ItemSlot.ItemSlotMainHand) {
				itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatMainHandDps, weaponDps);
			} else if (slot == ItemSlot.ItemSlotOffHand) {
				itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatOffHandDps, weaponDps);
			} else if (slot == ItemSlot.ItemSlotRanged) {
				itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatRangedDps, weaponDps);
			}
		}

		let ep = itemStats.computeEP(this.epWeights);

		// unique items are slightly worse than non-unique because you can have only one.
		if (item.unique) {
			ep -= 0.01;
		}

		this.itemEPCache[slot].set(item.id, ep);
		return ep;
	}

	setWowheadData(equippedItem: EquippedItem, elem: HTMLElement) {
		const parts = [];

		const lang = getLanguageCode();
		const langPrefix = lang ? lang + '.' : '';
		parts.push(`domain=${langPrefix}classic`);

		if (equippedItem.enchant != null) {
			parts.push('ench=' + equippedItem.enchant.effectId);
		}
		parts.push('pcs=' + this.gear.asArray().filter(ei => ei != null).map(ei => ei!.item.id).join(':'));

		elem.dataset.wowhead = parts.join('&');
		elem.dataset.whtticon = 'false';
	}

	static ARMOR_SLOTS: Array<ItemSlot> = [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotWrist,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotWaist,
		ItemSlot.ItemSlotFeet,
	];

	static WEAPON_SLOTS: Array<ItemSlot> = [
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
	];

	static readonly DIFFICULTY_SRCS: Partial<Record<SourceFilterOption, DungeonDifficulty>> = {
		[SourceFilterOption.SourceDungeon]: DungeonDifficulty.DifficultyNormal,
		[SourceFilterOption.SourceDungeonH]: DungeonDifficulty.DifficultyHeroic,
		[SourceFilterOption.SourceRaid10]: DungeonDifficulty.DifficultyRaid10,
		[SourceFilterOption.SourceRaid10H]: DungeonDifficulty.DifficultyRaid10H,
		[SourceFilterOption.SourceRaid25]: DungeonDifficulty.DifficultyRaid25,
		[SourceFilterOption.SourceRaid25H]: DungeonDifficulty.DifficultyRaid25H,
	};

	static readonly HEROIC_TO_NORMAL: Partial<Record<DungeonDifficulty, DungeonDifficulty>> = {
		[DungeonDifficulty.DifficultyHeroic]: DungeonDifficulty.DifficultyNormal,
		[DungeonDifficulty.DifficultyRaid10H]: DungeonDifficulty.DifficultyRaid10,
		[DungeonDifficulty.DifficultyRaid25H]: DungeonDifficulty.DifficultyRaid25,
	};

	static readonly RAID_IDS: Partial<Record<RaidFilterOption, number>> = {
		[RaidFilterOption.RaidNaxxramas]: 3456,
		[RaidFilterOption.RaidEyeOfEternity]: 4500,
		[RaidFilterOption.RaidObsidianSanctum]: 4493,
		[RaidFilterOption.RaidVaultOfArchavon]: 4603,
		[RaidFilterOption.RaidUlduar]: 4273,
		[RaidFilterOption.RaidTrialOfTheCrusader]: 4722,
		[RaidFilterOption.RaidOnyxiasLair]: 2159,
		[RaidFilterOption.RaidIcecrownCitadel]: 4812,
		[RaidFilterOption.RaidRubySanctum]: 4987,
	};

	filterItemData<T>(itemData: Array<T>, getItemFunc: (val: T) => Item, slot: ItemSlot): Array<T> {
		const filters = this.sim.getFilters();

		const filterItems = (itemData: Array<T>, filterFunc: (item: Item) => boolean) => {
			return itemData.filter(itemElem => filterFunc(getItemFunc(itemElem)));
		};

		if (filters.factionRestriction != UIItem_FactionRestriction.UNSPECIFIED) {
			itemData = filterItems(itemData, item => item.factionRestriction == filters.factionRestriction || item.factionRestriction == UIItem_FactionRestriction.UNSPECIFIED);
		}

		if (!filters.sources.includes(SourceFilterOption.SourceCrafting)) {
			itemData = filterItems(itemData, item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'crafted'));
		}
		if (!filters.sources.includes(SourceFilterOption.SourceQuest)) {
			itemData = filterItems(itemData, item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'quest'));
		}

		for (const [srcOptionStr, difficulty] of Object.entries(Player.DIFFICULTY_SRCS)) {
			const srcOption = parseInt(srcOptionStr) as SourceFilterOption;
			if (!filters.sources.includes(srcOption)) {
				itemData = filterItems(itemData, item =>
					!item.sources.some(itemSrc =>
						itemSrc.source.oneofKind == 'drop' && itemSrc.source.drop.difficulty == difficulty));

				if (difficulty == DungeonDifficulty.DifficultyRaid10H || difficulty == DungeonDifficulty.DifficultyRaid25H) {
					const normalDifficulty = Player.HEROIC_TO_NORMAL[difficulty];
					itemData = filterItems(itemData, item =>
						!item.sources.some(itemSrc =>
							itemSrc.source.oneofKind == 'drop' && itemSrc.source.drop.difficulty == normalDifficulty && itemSrc.source.drop.category == AL_CATEGORY_HARD_MODE));
				}
			}
		}

		if (!filters.raids.includes(RaidFilterOption.RaidVanilla)) {
			itemData = filterItems(itemData, item => item.expansion != Expansion.ExpansionVanilla);
		}
		if (!filters.raids.includes(RaidFilterOption.RaidTbc)) {
			itemData = filterItems(itemData, item => item.expansion != Expansion.ExpansionTbc);
		}
		for (const [raidOptionStr, zoneId] of Object.entries(Player.RAID_IDS)) {
			const raidOption = parseInt(raidOptionStr) as RaidFilterOption;
			if (!filters.raids.includes(raidOption)) {
				itemData = filterItems(itemData, item =>
					!item.sources.some(itemSrc =>
						itemSrc.source.oneofKind == 'drop' && itemSrc.source.drop.zoneId == zoneId));
			}
		}

		if (Player.ARMOR_SLOTS.includes(slot)) {
			itemData = filterItems(itemData, item => {
				if (!filters.armorTypes.includes(item.armorType)) {
					return false;
				}

				return true;
			});
		} else if (Player.WEAPON_SLOTS.includes(slot)) {
			itemData = filterItems(itemData, item => {
				if (!filters.weaponTypes.includes(item.weaponType)) {
					return false;
				}
				if (!filters.oneHandedWeapons && item.handType != HandType.HandTypeTwoHand) {
					return false;
				}
				if (!filters.twoHandedWeapons && item.handType == HandType.HandTypeTwoHand) {
					return false;
				}

				const minSpeed = slot == ItemSlot.ItemSlotMainHand ? filters.minMhWeaponSpeed : filters.minOhWeaponSpeed;
				const maxSpeed = slot == ItemSlot.ItemSlotMainHand ? filters.maxMhWeaponSpeed : filters.maxOhWeaponSpeed;
				if (minSpeed > 0 && item.weaponSpeed < minSpeed) {
					return false;
				}
				if (maxSpeed > 0 && item.weaponSpeed > maxSpeed) {
					return false;
				}

				return true;
			});
		} else if (slot == ItemSlot.ItemSlotRanged) {
			itemData = filterItems(itemData, item => {
				if (!filters.rangedWeaponTypes.includes(item.rangedWeaponType)) {
					return false;
				}

				const minSpeed = filters.minRangedWeaponSpeed;
				const maxSpeed = filters.maxRangedWeaponSpeed;
				if (minSpeed > 0 && item.weaponSpeed < minSpeed) {
					return false;
				}
				if (maxSpeed > 0 && item.weaponSpeed > maxSpeed) {
					return false;
				}

				return true;
			});
		}
		return itemData;
	}

	filterEnchantData<T>(enchantData: Array<T>, getEnchantFunc: (val: T) => Enchant, slot: ItemSlot, currentEquippedItem: EquippedItem | null): Array<T> {
		if (!currentEquippedItem) {
			return enchantData;
		}

		//const filters = this.sim.getFilters();

		return enchantData.filter(enchantElem => {
			const enchant = getEnchantFunc(enchantElem);

			if (!enchantAppliesToItem(enchant, currentEquippedItem.item)) {
				return false;
			}

			return true;
		});
	}

	makeUnitReference(): UnitReference {
		if (this.party == null) {
			return emptyUnitReference();
		} else {
			return newUnitReference(this.getRaidIndex());
		}
	}

	private toDatabase(): SimDatabase {
		const dbGear = this.getGear().toDatabase()
		const dbItemSwapGear = this.getItemSwapGear().toDatabase();
		return SimDatabase.create({
			items: dbGear.items.concat(dbItemSwapGear.items).filter(function(elem, index, self) {
				return index === self.findIndex((t) => (
					t.id === elem.id
				));
			}),
			enchants: dbGear.enchants.concat(dbItemSwapGear.enchants).filter(function(elem, index, self) {
				return index === self.findIndex((t) => (
					t.effectId === elem.effectId
				));
			}),
		})
	}

	toProto(forExport?: boolean, forSimming?: boolean): PlayerProto {
		const aplIsLaunched = aplLaunchStatuses[this.spec] == LaunchStatus.Launched;
		const gear = this.getGear();
		const aplRotation = forSimming ? this.getResolvedAplRotation() : this.aplRotation;
		return withSpecProto(
			this.spec,
			PlayerProto.create({
				name: this.getName(),
				race: this.getRace(),
				level: this.getLevel(),
				class: this.getClass(),
				equipment: gear.asSpec(),
				consumes: this.getConsumes(),
				bonusStats: this.getBonusStats().toProto(),
				buffs: this.getBuffs(),
				cooldowns: (aplIsLaunched || (forSimming && aplRotation.type == APLRotationType.TypeAPL))
					? Cooldowns.create({ hpPercentForDefensives: this.getCooldowns().hpPercentForDefensives })
					: this.getCooldowns(),
				talentsString: this.getTalentsString(),
				glyphs: this.getGlyphs(),
				rotation: aplRotation,
				profession1: this.getProfession1(),
				profession2: this.getProfession2(),
				reactionTimeMs: this.getReactionTime(),
				channelClipDelayMs: this.getChannelClipDelay(),
				inFrontOfTarget: this.getInFrontOfTarget(),
				distanceFromTarget: this.getDistanceFromTarget(),
				healingModel: this.getHealingModel(),
				database: forExport ? SimDatabase.create() : this.toDatabase(),
			}),
			(aplIsLaunched || (forSimming && aplRotation.type == APLRotationType.TypeAPL))
				? this.specTypeFunctions.rotationCreate()
				: this.getRotation(),
			this.getSpecOptions());
	}

	fromProto(eventID: EventID, proto: PlayerProto) {
		if (proto.rotation) {
			proto.rotation.prepullActions.forEach(ppa => {
				if (ppa.doAt) {
					ppa.doAtValue = APLValue.create({
						value: {oneofKind: 'const', const: { val: ppa.doAt }}
					});
					ppa.doAt = '';
				}
			});
			if (proto.rotation.enabled) {
				proto.rotation.enabled = false;
				proto.rotation.type = APLRotationType.TypeAPL;
			}
		}

		const rot = this.specTypeFunctions.rotationFromPlayer(proto);
		if (rot && !this.specTypeFunctions.rotationEquals(rot, this.specTypeFunctions.rotationCreate())) {
			if (proto.rotation?.type == APLRotationType.TypeAPL) {
				// Do nothing
			} else if (this.simpleRotationGenerator) {
				proto.rotation = APLRotation.create({
					type: APLRotationType.TypeSimple,
					simple: {
						specRotationJson: JSON.stringify(this.specTypeFunctions.rotationToJson(rot)),
						cooldowns: proto.cooldowns,
					},
				});
			} else {
				proto.rotation = APLRotation.create({
					type: APLRotationType.TypeAuto,
				});
			}
		}

		TypedEvent.freezeAllAndDo(() => {
			this.setName(eventID, proto.name);
			this.setRace(eventID, proto.race);
			this.setLevel(eventID, proto.level);
			this.setGear(eventID, proto.equipment ? this.sim.db.lookupEquipmentSpec(proto.equipment) : new Gear({}));
			//this.setBulkEquipmentSpec(eventID, BulkEquipmentSpec.create()); // Do not persist the bulk equipment settings.
			this.setConsumes(eventID, proto.consumes || Consumes.create());
			this.setBonusStats(eventID, Stats.fromProto(proto.bonusStats || UnitStats.create()));
			this.setBuffs(eventID, proto.buffs || IndividualBuffs.create());
			this.setTalentsString(eventID, proto.talentsString);
			this.setGlyphs(eventID, proto.glyphs || Glyphs.create());
			this.setProfession1(eventID, proto.profession1);
			this.setProfession2(eventID, proto.profession2);
			this.setReactionTime(eventID, proto.reactionTimeMs);
			this.setChannelClipDelay(eventID, proto.channelClipDelayMs);
			this.setInFrontOfTarget(eventID, proto.inFrontOfTarget);
			this.setDistanceFromTarget(eventID, proto.distanceFromTarget);
			this.setHealingModel(eventID, proto.healingModel || HealingModel.create());
			this.setSpecOptions(eventID, this.specTypeFunctions.optionsFromPlayer(proto));

			if (aplLaunchStatuses[this.spec] == LaunchStatus.Launched) {
				if (proto.rotation?.type == APLRotationType.TypeUnknown || proto.rotation?.type == APLRotationType.TypeLegacy) {
					if (!proto.rotation) {
						proto.rotation = APLRotation.create();
					}
					proto.rotation.type = APLRotationType.TypeAuto;
				}
			} else {
				this.setCooldowns(eventID, proto.cooldowns || Cooldowns.create());
				this.setRotation(eventID, this.specTypeFunctions.rotationFromPlayer(proto));
			}
			this.setAplRotation(eventID, proto.rotation || APLRotation.create())

			const options = this.getSpecOptions();
			for (let key in options) {
				if ((options[key] as any)?.['targetIndex']) {
					const targetIndex = (options[key] as any)['targetIndex'] as number;
					if (targetIndex == -1) {
						(options[key] as any) = UnitReference.create();
					} else {
						(options[key] as any) = UnitReference.create({type: UnitReference_Type.Player, index: targetIndex});
					}
					this.setSpecOptions(eventID, options);
					break;
				}
			}

			if (this.spec == Spec.SpecBalanceDruid) {
				const rot = this.getRotation() as SpecRotation<Spec.SpecBalanceDruid>;
				if (rot.okfPpm) {
					rot.okfPpm = 0;
					this.setRotation(eventID, rot as SpecRotation<SpecType>);
				}
			}

			if (this.spec == Spec.SpecHunter) {
				const rot = this.getRotation() as SpecRotation<Spec.SpecHunter>;
				if (rot.timeToTrapWeaveMs) {
					const options = this.getSpecOptions() as SpecOptions<Spec.SpecHunter>;
					options.timeToTrapWeaveMs = rot.timeToTrapWeaveMs;
					this.setSpecOptions(eventID, options as SpecOptions<SpecType>);
					rot.timeToTrapWeaveMs = 0;
					this.setRotation(eventID, rot as SpecRotation<SpecType>);
				}
			}

			if (this.spec == Spec.SpecShadowPriest) {
				// no-op
			}

			if ([Spec.SpecEnhancementShaman, Spec.SpecRestorationShaman, Spec.SpecElementalShaman].includes(this.spec)) {
				const rot = this.getRotation() as SpecRotation<ShamanSpecs>;
				if (rot.totems) {
					const options = this.getSpecOptions() as SpecOptions<ShamanSpecs>;
					options.totems = rot.totems;
					this.setSpecOptions(eventID, options as SpecOptions<SpecType>);
					rot.totems = undefined;
					this.setRotation(eventID, rot as SpecRotation<SpecType>);
				}
				const opt = this.getSpecOptions() as SpecOptions<ShamanSpecs>;

				// Update Bloodlust to be part of rotation instead of options to support APL casting bloodlust.
				if (opt.bloodlust) {
					opt.bloodlust = false;

					var tRot = this.getRotation();
					if (this.spec == Spec.SpecElementalShaman) {
						(tRot as ElementalShaman_Rotation).bloodlust = ElementalShaman_Rotation_BloodlustUse.UseBloodlust;
					} else if (this.spec == Spec.SpecEnhancementShaman) {
						(tRot as EnhancementShaman_Rotation).bloodlust = EnhancementShaman_Rotation_BloodlustUse.UseBloodlust;
					} else if (this.spec == Spec.SpecRestorationShaman) {
						(tRot as RestorationShaman_Rotation).bloodlust = RestorationShaman_Rotation_BloodlustUse.UseBloodlust;
					}

					this.setRotation(eventID, tRot as SpecRotation<SpecType>);
					this.setSpecOptions(eventID, opt as SpecOptions<SpecType>);
				}

				// Update Ele TS range option.
				if (this.spec == Spec.SpecElementalShaman) {
					var eleOpt = this.getSpecOptions() as ElementalShaman_Options;
					var eleRot = this.getRotation() as ElementalShaman_Rotation;
					if (eleRot.inThunderstormRange) {
						eleOpt.thunderstormRange = ElementalShaman_Options_ThunderstormRange.TSInRange;
						eleRot.inThunderstormRange = false;
						this.setRotation(eventID, eleRot as SpecRotation<SpecType>);
						this.setSpecOptions(eventID, eleOpt as SpecOptions<SpecType>);
					}
				}
			}
		});
	}

	clone(eventID: EventID): Player<SpecType> {
		const newPlayer = new Player<SpecType>(this.spec, this.sim);
		newPlayer.fromProto(eventID, this.toProto());
		return newPlayer;
	}

	applySharedDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			this.setReactionTime(eventID, 200);
			this.setInFrontOfTarget(eventID, isTankSpec(this.spec));
			this.setHealingModel(eventID, HealingModel.create({
				burstWindow: isTankSpec(this.spec) ? 6 : 0,
			}));
			this.setCooldowns(eventID, Cooldowns.create({
				hpPercentForDefensives: isTankSpec(this.spec) ? 0.35 : 0,
			}));
			this.setBonusStats(eventID, new Stats());

			if (aplLaunchStatuses[this.spec] >= LaunchStatus.Beta) {
				this.setAplRotation(eventID, APLRotation.create({
					type: APLRotationType.TypeAuto,
				}))
			}
		});
	}
}
