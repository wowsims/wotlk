import {
	Class,
	Consumes,
	Cooldowns,
	Faction,
	GemColor,
	Glyphs,
	HandType,
	HealingModel,
	IndividualBuffs,
	ItemSlot,
	Profession,
	PseudoStat,
	Race,
	UnitReference,
	SimDatabase,
	Spec,
	Stat,
	UnitStats,
} from './proto/common.js';
import {
	AuraStats as AuraStatsProto,
	SpellStats as SpellStatsProto,
	UnitMetadata as UnitMetadataProto,
} from './proto/api.js';
import {
	APLRotation,
	APLRotation_Type as APLRotationType,
	SimpleRotation,
} from './proto/apl.js';
import {
	DungeonDifficulty,
	Expansion,
	RaidFilterOption,
	SourceFilterOption,
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
	UIItem_FactionRestriction,
} from './proto/ui.js';

import { PlayerStats } from './proto/api.js';
import { Player as PlayerProto } from './proto/api.js';
import { StatWeightsResult } from './proto/api.js';
import { ActionId } from './proto_utils/action_id.js';
import { EquippedItem, getWeaponDPS } from './proto_utils/equipped_item.js';

import { playerTalentStringToProto } from './talents/factory.js';
import { Gear, ItemSwapGear } from './proto_utils/gear.js';
import {
	isUnrestrictedGem,
	gemMatchesSocket,
} from './proto_utils/gems.js';
import { Stats } from './proto_utils/stats.js';

import {
	AL_CATEGORY_HARD_MODE,
	ClassSpecs,
	SpecRotation,
	SpecTalents,
	SpecTypeFunctions,
	SpecOptions,
	canEquipEnchant,
	canEquipItem,
	classColors,
	emptyUnitReference,
	enchantAppliesToItem,
	getTalentTree,
	getTalentTreeIcon,
	getTalentTreePoints,
	getMetaGemEffectEP,
	isTankSpec,
	newUnitReference,
	raceToFaction,
	specToClass,
	specToEligibleRaces,
	specTypeFunctions,
	withSpecProto,
} from './proto_utils/utils.js';

import * as Mechanics from './constants/mechanics.js';
import { getLanguageCode } from './constants/lang.js';
import { EventID, TypedEvent } from './typed_event.js';
import { Party, MAX_PARTY_SIZE } from './party.js';
import { Raid } from './raid.js';
import { Sim, SimSettingCategories } from './sim.js';
import { stringComparator, sum } from './utils.js';
import { Database } from './proto_utils/database.js';

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

export interface PlayerConfig<SpecType extends Spec> {
	autoRotation: AutoRotationGenerator<SpecType>,
	simpleRotation?: SimpleRotationGenerator<SpecType>,
}

const SPEC_CONFIGS: Partial<Record<Spec, PlayerConfig<any>>> = {};

export function registerSpecConfig(spec: Spec, config: PlayerConfig<any>) {
	SPEC_CONFIGS[spec] = config;
}

export function getSpecConfig<SpecType extends Spec>(spec: SpecType): PlayerConfig<SpecType> {
	const config = SPEC_CONFIGS[spec] as PlayerConfig<SpecType>;
	if (!config) {
		throw new Error('No config registered for Spec: ' + spec);
	}
	return config;
}

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
	private enableItemSwap: boolean = false;
	private itemSwapGear: ItemSwapGear = new ItemSwapGear({});
	private race: Race;
	private profession1: Profession = 0;
	private profession2: Profession = 0;
	aplRotation: APLRotation = APLRotation.create();
	private talentsString: string = '';
	private glyphs: Glyphs = Glyphs.create();
	private specOptions: SpecOptions<SpecType>;
	private reactionTime: number = 0;
	private channelClipDelay: number = 0;
	private inFrontOfTarget: boolean = false;
	private distanceFromTarget: number = 0;
	private nibelungAverageCasts: number = 11;
	private nibelungAverageCastsSet: boolean = false;
	private healingModel: HealingModel = HealingModel.create();
	private healingEnabled: boolean = false;

	private readonly autoRotationGenerator: AutoRotationGenerator<SpecType> | null = null;
	private readonly simpleRotationGenerator: SimpleRotationGenerator<SpecType> | null = null;

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
	readonly itemSwapChangeEmitter = new TypedEvent<void>('PlayerItemSwap');
	readonly professionChangeEmitter = new TypedEvent<void>('PlayerProfession');
	readonly raceChangeEmitter = new TypedEvent<void>('PlayerRace');
	readonly rotationChangeEmitter = new TypedEvent<void>('PlayerRotation');
	readonly talentsChangeEmitter = new TypedEvent<void>('PlayerTalents');
	readonly glyphsChangeEmitter = new TypedEvent<void>('PlayerGlyphs');
	readonly specOptionsChangeEmitter = new TypedEvent<void>('PlayerSpecOptions');
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
		this.specTypeFunctions = specTypeFunctions[this.spec] as SpecTypeFunctions<SpecType>;
		this.specOptions = this.specTypeFunctions.optionsCreate();

		const specConfig = SPEC_CONFIGS[this.spec] as PlayerConfig<SpecType>;
		if (!specConfig) {
			throw new Error('Could not find spec config for spec: ' + this.spec);
		}
		this.autoRotationGenerator = specConfig.autoRotation;
		if (specConfig.simpleRotation) {
			this.simpleRotationGenerator = specConfig.simpleRotation;
		} else {
			this.simpleRotationGenerator = null;
		}

		for(let i = 0; i < ItemSlot.ItemSlotRanged+1; ++i) {
			this.itemEPCache[i] = new Map();
		}

		this.changeEmitter = TypedEvent.onAny([
			this.nameChangeEmitter,
			this.buffsChangeEmitter,
			this.consumesChangeEmitter,
			this.bonusStatsChangeEmitter,
			this.gearChangeEmitter,
			this.itemSwapChangeEmitter,
			this.professionChangeEmitter,
			this.raceChangeEmitter,
			this.rotationChangeEmitter,
			this.talentsChangeEmitter,
			this.glyphsChangeEmitter,
			this.specOptionsChangeEmitter,
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
		return this.sim.db.getItems(slot).filter(item => canEquipItem(item, this.spec, slot));
	}

	// Returns all enchants that this player can wear in the given slot.
	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.sim.db.getEnchants(slot).filter(enchant => canEquipEnchant(enchant, this.spec));
	}

	// Returns all gems that this player can wear of the given color.
	getGems(socketColor?: GemColor): Array<Gem> {
		return this.sim.db.getGems(socketColor);
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

	canDualWield2H(): boolean {
		return this.getClass() == Class.ClassWarrior && (this.getTalents() as SpecTalents<Spec.SpecWarrior>).titansGrip;
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

	setGear(eventID: EventID, newGear: Gear) {
		if (newGear.equals(this.gear))
			return;

		this.gear = newGear;
		this.gearChangeEmitter.emit(eventID);
	}

	getEnableItemSwap(): boolean {
		return this.enableItemSwap;
	}

	setEnableItemSwap(eventID: EventID, newEnableItemSwap: boolean) {
		if (newEnableItemSwap == this.enableItemSwap)
			return;

		this.enableItemSwap = newEnableItemSwap;
		this.itemSwapChangeEmitter.emit(eventID);
	}

	equipItemSwapitem(eventID: EventID, slot: ItemSlot, newItem: EquippedItem | null) {
		this.setItemSwapGear(eventID, this.itemSwapGear.withEquippedItem(slot, newItem, this.canDualWield2H()));
	}

	getItemSwapItem(slot: ItemSlot): EquippedItem | null {
		return this.itemSwapGear.getEquippedItem(slot);
	}

	getItemSwapGear(): ItemSwapGear {
		return this.itemSwapGear;
	}

	setItemSwapGear(eventID: EventID, newItemSwapGear: ItemSwapGear) {
		if (newItemSwapGear.equals(this.itemSwapGear))
			return;

		this.itemSwapGear = newItemSwapGear;
		this.itemSwapChangeEmitter.emit(eventID);
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
		//const agility = (this.currentStats.finalStats?.stats[Stat.StatAgility] || 0.0) / this.getClass();
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

		const debuffs = this.sim.raid.getDebuffs();
		if (debuffs.totemOfWrath || debuffs.heartOfTheCrusader || debuffs.masterPoisoner) {
			debuffCrit = 3.0;
		}

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

	setAplRotation(eventID: EventID, newRotation: APLRotation) {
		if (APLRotation.equals(newRotation, this.aplRotation))
			return;

		this.aplRotation = APLRotation.clone(newRotation);
		this.rotationChangeEmitter.emit(eventID);
	}

	getSimpleRotation(): SpecRotation<SpecType> {
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
	}

	setSimpleRotation(eventID: EventID, newRotation: SpecRotation<SpecType>) {
		if (this.specTypeFunctions.rotationEquals(newRotation, this.getSimpleRotation()))
			return;

		if (!this.aplRotation.simple) {
			this.aplRotation.simple = SimpleRotation.create();
		}
		this.aplRotation.simple.specRotationJson = JSON.stringify(this.specTypeFunctions.rotationToJson(newRotation));

		this.rotationChangeEmitter.emit(eventID);
	}

	getSimpleCooldowns(): Cooldowns {
		// Make a defensive copy
		return Cooldowns.clone(this.aplRotation.simple?.cooldowns || Cooldowns.create());
	}

	setSimpleCooldowns(eventID: EventID, newCooldowns: Cooldowns) {
		if (Cooldowns.equals(this.getSimpleCooldowns(), newCooldowns))
			return;

		if (!this.aplRotation.simple) {
			this.aplRotation.simple = SimpleRotation.create();
		}
		this.aplRotation.simple.cooldowns = newCooldowns;
		this.rotationChangeEmitter.emit(eventID);
	}

	getRotationType(): APLRotationType {
		if (this.aplRotation.type == APLRotationType.TypeUnknown) {
			return APLRotationType.TypeAPL;
		} else {
			return this.aplRotation.type;
		}
	}

	hasSimpleRotationGenerator(): boolean {
		return this.simpleRotationGenerator != null;
	}

	getResolvedAplRotation(): APLRotation {
		const type = this.getRotationType();
		if (type == APLRotationType.TypeAuto && this.autoRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			const rot = APLRotation.clone(this.autoRotationGenerator(this));
			rot.type = APLRotationType.TypeAuto;
			return rot;
		} else if (type == APLRotationType.TypeSimple && this.simpleRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			const simpleRot = this.getSimpleRotation();
			const rot = APLRotation.clone(this.simpleRotationGenerator(this, simpleRot, this.getSimpleCooldowns()));
			rot.simple = this.aplRotation.simple;
			rot.type = APLRotationType.TypeSimple;
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

	getNibelungAverageCasts(): number {
		return this.nibelungAverageCasts;
	}

	setNibelungAverageCastsSet(eventID: EventID, newnibelungAverageCastsSet: boolean) {
		if (newnibelungAverageCastsSet == this.nibelungAverageCastsSet)
			return;

		this.nibelungAverageCastsSet = newnibelungAverageCastsSet;
	}

	setNibelungAverageCasts(eventID: EventID, newnibelungAverageCasts: number) {
		if (newnibelungAverageCasts == this.nibelungAverageCasts)
			return;

		this.nibelungAverageCasts = Math.min(newnibelungAverageCasts, 16);
		this.miscOptionsChangeEmitter.emit(eventID);
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

	computeGemEP(gem: Gem): number {
		if (this.gemEPCache.has(gem.id)) {
			return this.gemEPCache.get(gem.id)!;
		}

		const epFromStats = this.computeStatsEP(new Stats(gem.stats));
		const epFromEffect = getMetaGemEffectEP(this.spec, gem, Stats.fromProto(this.currentStats.finalStats));
		let bonusEP = 0;
		// unique items are slightly worse than non-unique because you can have only one.
		if (gem.unique) {
			bonusEP -= 0.01;
		}

		let ep = epFromStats + epFromEffect + bonusEP;
		this.gemEPCache.set(gem.id, ep);
		return ep;
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

		// Compare whether its better to match sockets + get socket bonus, or just use best gems.
		const bestGemEPNotMatchingSockets = sum(item.gemSockets.map(socketColor => {
			const gems = this.sim.db.getGems(socketColor).filter(gem => isUnrestrictedGem(gem, this.sim.getPhase()));
			if (gems.length > 0) {
				return Math.max(...gems.map(gem => this.computeGemEP(gem)));
			} else {
				return 0;
			}
		}));

		const bestGemEPMatchingSockets = sum(item.gemSockets.map(socketColor => {
			const gems = this.sim.db.getGems(socketColor).filter(gem => isUnrestrictedGem(gem, this.sim.getPhase()) && gemMatchesSocket(gem, socketColor));
			if (gems.length > 0) {
				return Math.max(...gems.map(gem => this.computeGemEP(gem)));
			} else {
				return 0;
			}
		})) + this.computeStatsEP(new Stats(item.socketBonus));

		ep += Math.max(bestGemEPMatchingSockets, bestGemEPNotMatchingSockets);

		this.itemEPCache[slot].set(item.id, ep);
		return ep;
	}

	setWowheadData(equippedItem: EquippedItem, elem: HTMLElement) {
		const parts = [];

		const lang = getLanguageCode();
		const langPrefix = lang ? lang + '.' : '';
		parts.push(`domain=${langPrefix}wotlk`);

		const isBlacksmithing = this.hasProfession(Profession.Blacksmithing);
		if (equippedItem.gems.length > 0) {
			parts.push('gems=' + equippedItem.curGems(isBlacksmithing).map(gem => gem ? gem.id : 0).join(':'));
		}
		if (equippedItem.enchant != null) {
			parts.push('ench=' + equippedItem.enchant.effectId);
		}
		parts.push('pcs=' + this.gear.asArray().filter(ei => ei != null).map(ei => ei!.item.id).join(':'));

		if (equippedItem.hasExtraSocket(isBlacksmithing)) {
			parts.push('sock');
		}

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

	filterGemData<T>(gemData: Array<T>, getGemFunc: (val: T) => Gem, slot: ItemSlot, socketColor: GemColor): Array<T> {
		const filters = this.sim.getFilters();

		const isJewelcrafting = this.hasProfession(Profession.Jewelcrafting);
		return gemData.filter(gemElem => {
			const gem = getGemFunc(gemElem);
			if (!isJewelcrafting && gem.requiredProfession == Profession.Jewelcrafting) {
				return false;
			}

			if (filters.matchingGemsOnly && !gemMatchesSocket(gem, socketColor)) {
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
		return Database.mergeSimDatabases(dbGear, dbItemSwapGear);
	}

	toProto(forExport?: boolean, forSimming?: boolean, exportCategories?: Array<SimSettingCategories>): PlayerProto {
		const exportCategory = (cat: SimSettingCategories) =>
				!exportCategories
				|| exportCategories.length == 0
				|| exportCategories.includes(cat);

		const gear = this.getGear();
		const aplRotation = forSimming ? this.getResolvedAplRotation() : this.aplRotation;

		let player = PlayerProto.create({
			class: this.getClass(),
			database: forExport ? undefined : this.toDatabase(),
		});
		if (exportCategory(SimSettingCategories.Gear)) {
			PlayerProto.mergePartial(player, {
				equipment: gear.asSpec(),
				bonusStats: this.getBonusStats().toProto(),
				enableItemSwap: this.getEnableItemSwap(),
				itemSwap: this.getItemSwapGear().toProto(),
			});
		}
		if (exportCategory(SimSettingCategories.Talents)) {
			PlayerProto.mergePartial(player, {
				talentsString: this.getTalentsString(),
				glyphs: this.getGlyphs(),
			});
		}
		if (exportCategory(SimSettingCategories.Rotation)) {
			PlayerProto.mergePartial(player, {
				cooldowns: Cooldowns.create({ hpPercentForDefensives: this.getSimpleCooldowns().hpPercentForDefensives }),
				rotation: aplRotation,
			});
		}
		if (exportCategory(SimSettingCategories.Consumes)) {
			PlayerProto.mergePartial(player, {
				consumes: this.getConsumes(),
			});
		}
		if (exportCategory(SimSettingCategories.Miscellaneous)) {
			PlayerProto.mergePartial(player, {
				name: this.getName(),
				race: this.getRace(),
				profession1: this.getProfession1(),
				profession2: this.getProfession2(),
				reactionTimeMs: this.getReactionTime(),
				channelClipDelayMs: this.getChannelClipDelay(),
				inFrontOfTarget: this.getInFrontOfTarget(),
				distanceFromTarget: this.getDistanceFromTarget(),
				healingModel: this.getHealingModel(),
				nibelungAverageCasts: this.getNibelungAverageCasts(),
				nibelungAverageCastsSet: this.nibelungAverageCastsSet,
			});
			player = withSpecProto(this.spec, player, this.getSpecOptions());
		}
		if (exportCategory(SimSettingCategories.External)) {
			PlayerProto.mergePartial(player, {
				buffs: this.getBuffs(),
			});
		}
		return player;
	}

	fromProto(eventID: EventID, proto: PlayerProto, includeCategories?: Array<SimSettingCategories>) {
		const loadCategory = (cat: SimSettingCategories) =>
				!includeCategories
				|| includeCategories.length == 0
				|| includeCategories.includes(cat);

		// For backwards compatibility with legacy rotations (removed on 2024/01/15).
		if (proto.rotation?.type == APLRotationType.TypeLegacy) {
			proto.rotation.type = APLRotationType.TypeAuto;
		}

		TypedEvent.freezeAllAndDo(() => {
			if (loadCategory(SimSettingCategories.Gear)) {
				this.setGear(eventID, proto.equipment ? this.sim.db.lookupEquipmentSpec(proto.equipment) : new Gear({}));
				this.setEnableItemSwap(eventID, proto.enableItemSwap);
				this.setItemSwapGear(eventID, proto.itemSwap ? this.sim.db.lookupItemSwap(proto.itemSwap) : new ItemSwapGear({}));
				this.setBonusStats(eventID, Stats.fromProto(proto.bonusStats || UnitStats.create()));
				//this.setBulkEquipmentSpec(eventID, BulkEquipmentSpec.create()); // Do not persist the bulk equipment settings.
			}
			if (loadCategory(SimSettingCategories.Talents)) {
				this.setTalentsString(eventID, proto.talentsString);
				this.setGlyphs(eventID, proto.glyphs || Glyphs.create());
			}
			if (loadCategory(SimSettingCategories.Rotation)) {
				if (proto.rotation?.type == APLRotationType.TypeUnknown || proto.rotation?.type == APLRotationType.TypeLegacy) {
					if (!proto.rotation) {
						proto.rotation = APLRotation.create();
					}
					proto.rotation.type = APLRotationType.TypeAuto;
				}
				this.setAplRotation(eventID, proto.rotation || APLRotation.create())
			}
			if (loadCategory(SimSettingCategories.Consumes)) {
				this.setConsumes(eventID, proto.consumes || Consumes.create());
			}
			if (loadCategory(SimSettingCategories.Miscellaneous)) {
				this.setSpecOptions(eventID, this.specTypeFunctions.optionsFromPlayer(proto));
				this.setName(eventID, proto.name);
				this.setRace(eventID, proto.race);
				this.setProfession1(eventID, proto.profession1);
				this.setProfession2(eventID, proto.profession2);
				this.setReactionTime(eventID, proto.reactionTimeMs);
				this.setChannelClipDelay(eventID, proto.channelClipDelayMs);
				this.setInFrontOfTarget(eventID, proto.inFrontOfTarget);
				this.setDistanceFromTarget(eventID, proto.distanceFromTarget);
				this.setNibelungAverageCastsSet(eventID, proto.nibelungAverageCastsSet);
				if (this.nibelungAverageCastsSet) {
					this.setNibelungAverageCasts(eventID, proto.nibelungAverageCasts);
				}
				this.setHealingModel(eventID, proto.healingModel || HealingModel.create());
			}
			if (loadCategory(SimSettingCategories.External)) {
				this.setBuffs(eventID, proto.buffs || IndividualBuffs.create());
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
			this.setEnableItemSwap(eventID, false);
			this.setItemSwapGear(eventID, new ItemSwapGear({}));
			this.setReactionTime(eventID, 200);
			this.setInFrontOfTarget(eventID, isTankSpec(this.spec));
			this.setHealingModel(eventID, HealingModel.create({
				burstWindow: isTankSpec(this.spec) ? 6 : 0,
			}));
			this.setSimpleCooldowns(eventID, Cooldowns.create({
				hpPercentForDefensives: isTankSpec(this.spec) ? 0.35 : 0,
			}));
			this.setBonusStats(eventID, new Stats());

			this.setAplRotation(eventID, APLRotation.create({
				type: APLRotationType.TypeAuto,
			}))
		});
	}
}
