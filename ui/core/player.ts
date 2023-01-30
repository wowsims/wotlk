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
	ItemSwap,
	Profession,
	PseudoStat,
	Race,
	RaidTarget,
	RangedWeaponType,
	SimDatabase,
	Spec,
	Stat,
	UnitStats,
	WeaponType,
} from './proto/common.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from './proto/ui.js';

import { PlayerStats } from './proto/api.js';
import { Player as PlayerProto } from './proto/api.js';
import { StatWeightsResult } from './proto/api.js';
import { EquippedItem, getWeaponDPS } from './proto_utils/equipped_item.js';

import { playerTalentStringToProto } from './talents/factory.js';
import { Gear, ItemSwapGear } from './proto_utils/gear.js';
import {
	isUnrestrictedGem,
	gemEligibleForSocket,
	gemMatchesSocket,
} from './proto_utils/gems.js';
import { Stats } from './proto_utils/stats.js';

import {
	ClassSpecs,
	SpecRotation,
	SpecTalents,
	SpecTypeFunctions,
	SpecOptions,
	canEquipEnchant,
	canEquipItem,
	classColors,
	emptyRaidTarget,
	enchantAppliesToItem,
	getEligibleEnchantSlots,
	getEligibleItemSlots,
	getTalentTree,
	getTalentTreeIcon,
	getMetaGemEffectEP,
	isTankSpec,
	newRaidTarget,
	playerToSpec,
	raceToFaction,
	specToClass,
	specToEligibleRaces,
	specTypeFunctions,
	withSpecProto,
} from './proto_utils/utils.js';

import { getLanguageCode } from './constants/lang.js';
import { Listener } from './typed_event.js';
import { EventID, TypedEvent } from './typed_event.js';
import { Party, MAX_PARTY_SIZE } from './party.js';
import { Raid } from './raid.js';
import { Sim } from './sim.js';
import { sum } from './utils.js';
import { wait } from './utils.js';
import { WorkerPool } from './worker_pool.js';
import { EnhancementShaman_Options } from './proto/shaman.js';

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
	private itemSwapGear: ItemSwapGear = new ItemSwapGear();
	private race: Race;
	private profession1: Profession = 0;
	private profession2: Profession = 0;
	private rotation: SpecRotation<SpecType>;
	private talentsString: string = '';
	private glyphs: Glyphs = Glyphs.create();
	private specOptions: SpecOptions<SpecType>;
	private cooldowns: Cooldowns = Cooldowns.create();
	private inFrontOfTarget: boolean = false;
	private distanceFromTarget: number = 0;
	private healingModel: HealingModel = HealingModel.create();
	private healingEnabled: boolean = false;

	private itemEPCache: Map<number, number> = new Map<number, number>();
	private gemEPCache: Map<number, number> = new Map<number, number>();
	private enchantEPCache: Map<number, number> = new Map<number, number>();
	private talents: SpecTalents<SpecType> | null = null;

	readonly specTypeFunctions: SpecTypeFunctions<SpecType>;

	private epWeights: Stats = new Stats();
	private currentStats: PlayerStats = PlayerStats.create();

	readonly nameChangeEmitter = new TypedEvent<void>('PlayerName');
	readonly buffsChangeEmitter = new TypedEvent<void>('PlayerBuffs');
	readonly consumesChangeEmitter = new TypedEvent<void>('PlayerConsumes');
	readonly bonusStatsChangeEmitter = new TypedEvent<void>('PlayerBonusStats');
	readonly gearChangeEmitter = new TypedEvent<void>('PlayerGear');
	readonly professionChangeEmitter = new TypedEvent<void>('PlayerProfession');
	readonly raceChangeEmitter = new TypedEvent<void>('PlayerRace');
	readonly rotationChangeEmitter = new TypedEvent<void>('PlayerRotation');
	readonly talentsChangeEmitter = new TypedEvent<void>('PlayerTalents');
	readonly glyphsChangeEmitter = new TypedEvent<void>('PlayerGlyphs');
	readonly specOptionsChangeEmitter = new TypedEvent<void>('PlayerSpecOptions');
	readonly cooldownsChangeEmitter = new TypedEvent<void>('PlayerCooldowns');
	readonly inFrontOfTargetChangeEmitter = new TypedEvent<void>('PlayerInFrontOfTarget');
	readonly distanceFromTargetChangeEmitter = new TypedEvent<void>('PlayerDistanceFromTarget');
	readonly healingModelChangeEmitter = new TypedEvent<void>('PlayerHealingModel');
	readonly epWeightsChangeEmitter = new TypedEvent<void>('PlayerEpWeights');

	readonly currentStatsEmitter = new TypedEvent<void>('PlayerCurrentStats');

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	constructor(spec: Spec, sim: Sim) {
		this.sim = sim;
		this.party = null;
		this.raid = null;

		this.spec = spec;
		this.race = specToEligibleRaces[this.spec][0];
		this.specTypeFunctions = specTypeFunctions[this.spec] as SpecTypeFunctions<SpecType>;
		this.rotation = this.specTypeFunctions.rotationCreate();
		this.specOptions = this.specTypeFunctions.optionsCreate();

		this.changeEmitter = TypedEvent.onAny([
			this.nameChangeEmitter,
			this.buffsChangeEmitter,
			this.consumesChangeEmitter,
			this.bonusStatsChangeEmitter,
			this.gearChangeEmitter,
			this.professionChangeEmitter,
			this.raceChangeEmitter,
			this.rotationChangeEmitter,
			this.talentsChangeEmitter,
			this.glyphsChangeEmitter,
			this.specOptionsChangeEmitter,
			this.cooldownsChangeEmitter,
			this.inFrontOfTargetChangeEmitter,
			this.distanceFromTargetChangeEmitter,
			this.healingModelChangeEmitter,
			this.epWeightsChangeEmitter,
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
		this.itemEPCache = new Map();
		this.enchantEPCache = new Map();
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

		//// Remove item cooldowns if there is no cooldown available for the item.
		//const availableCooldowns = this.currentStats.cooldowns;
		//const newCooldowns = this.getCooldowns();
		//newCooldowns.cooldowns = newCooldowns.cooldowns.filter(cd => {
		//	if (cd.id && 'itemId' in cd.id.rawId) {
		//		return availableCooldowns.find(acd => ActionIdProto.equals(acd, cd.id)) != null;
		//	} else {
		//		return true;
		//	}
		//});
		//// TODO: Reference the parent event ID
		//this.setCooldowns(TypedEvent.nextEventID(), newCooldowns);
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

	getCooldowns(): Cooldowns {
		// Make a defensive copy
		return Cooldowns.clone(this.cooldowns);
	}

	setCooldowns(eventID: EventID, newCooldowns: Cooldowns) {
		if (Cooldowns.equals(this.cooldowns, newCooldowns))
			return;

		// Make a defensive copy
		this.cooldowns = Cooldowns.clone(newCooldowns);
		this.cooldownsChangeEmitter.emit(eventID);
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

	getBonusStats(): Stats {
		return this.bonusStats;
	}

	setBonusStats(eventID: EventID, newBonusStats: Stats) {
		if (newBonusStats.equals(this.bonusStats))
			return;

		this.bonusStats = newBonusStats;
		this.bonusStatsChangeEmitter.emit(eventID);
	}

	getRotation(): SpecRotation<SpecType> {
		return this.specTypeFunctions.rotationCopy(this.rotation);
	}

	setRotation(eventID: EventID, newRotation: SpecRotation<SpecType>) {
		if (this.specTypeFunctions.rotationEquals(newRotation, this.rotation))
			return;

		this.rotation = this.specTypeFunctions.rotationCopy(newRotation);
		this.rotationChangeEmitter.emit(eventID);
	}

	getTalents(): SpecTalents<SpecType> {
		if (this.talents == null) {
			this.talents = playerTalentStringToProto(this.spec, this.talentsString) as SpecTalents<SpecType>;
		}
		return this.talents;
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

	enableHealing() {
		this.healingEnabled = true;
		var hm = this.getHealingModel();
		if (hm.cadenceSeconds == 0) {
			hm.cadenceSeconds = 2;
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
		// If we have enabled healing model and try to set 0s cadence, default to 2s.
		if (this.healingModel.cadenceSeconds == 0 && this.healingEnabled) {
			this.healingModel.cadenceSeconds = 2;
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

		if (this.itemEPCache.has(item.id)) {
			return this.itemEPCache.get(item.id)!;
		}

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

		this.itemEPCache.set(item.id, ep);
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

		elem.setAttribute('data-wowhead', parts.join('&'));
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

	filterItemData<T>(itemData: Array<T>, getItemFunc: (val: T) => Item, slot: ItemSlot): Array<T> {
		const filters = this.sim.getFilters();

		if (Player.ARMOR_SLOTS.includes(slot)) {
			return itemData.filter(itemElem => {
				const item = getItemFunc(itemElem);

				if (!filters.armorTypes.includes(item.armorType)) {
					return false;
				}

				return true;
			});
		} else if (Player.WEAPON_SLOTS.includes(slot)) {
			return itemData.filter(itemElem => {
				const item = getItemFunc(itemElem);

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
			return itemData.filter(itemElem => {
				const item = getItemFunc(itemElem);

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
		} else {
			return itemData;
		}
	}

	filterEnchantData<T>(enchantData: Array<T>, getEnchantFunc: (val: T) => Enchant, slot: ItemSlot, currentEquippedItem: EquippedItem|null): Array<T> {
		if (!currentEquippedItem) {
			return enchantData;
		}

		const filters = this.sim.getFilters();

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

		return gemData.filter(gemElem => {
			const gem = getGemFunc(gemElem);

			if (filters.matchingGemsOnly && !gemMatchesSocket(gem, socketColor)) {
				return false;
			}

			return true;
		});
	}

	makeRaidTarget(): RaidTarget {
		if (this.party == null) {
			return emptyRaidTarget();
		} else {
			return newRaidTarget(this.getRaidIndex());
		}
	}

	private toDatabase(): SimDatabase {
		const dbGear =  this.getGear().toDatabase()
		const dbItemSwapGear = this.getItemSwapGear().toDatabase();
		return SimDatabase.create({
			items: dbGear.items.concat(dbItemSwapGear.items),
			enchants: dbGear.enchants.concat(dbItemSwapGear.enchants),
			gems: dbGear.gems.concat(dbItemSwapGear.gems),
		})
	}

	toProto(forExport?: boolean): PlayerProto {
		const gear = this.getGear();
		return withSpecProto(
			this.spec,
			PlayerProto.create({
				name: this.getName(),
				race: this.getRace(),
				class: this.getClass(),
				equipment: gear.asSpec(),
				consumes: this.getConsumes(),
				bonusStats: this.getBonusStats().toProto(),
				buffs: this.getBuffs(),
				cooldowns: this.getCooldowns(),
				talentsString: this.getTalentsString(),
				glyphs: this.getGlyphs(),
				profession1: this.getProfession1(),
				profession2: this.getProfession2(),
				inFrontOfTarget: this.getInFrontOfTarget(),
				distanceFromTarget: this.getDistanceFromTarget(),
				healingModel: this.getHealingModel(),
				database: forExport ? SimDatabase.create() : this.toDatabase(),
			}),
			this.getRotation(),
			this.getSpecOptions());
	}

	fromProto(eventID: EventID, proto: PlayerProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setName(eventID, proto.name);
			this.setRace(eventID, proto.race);
			this.setGear(eventID, proto.equipment ? this.sim.db.lookupEquipmentSpec(proto.equipment) : new Gear({}));
			this.setConsumes(eventID, proto.consumes || Consumes.create());
			this.setBonusStats(eventID, Stats.fromProto(proto.bonusStats || UnitStats.create()));
			this.setBuffs(eventID, proto.buffs || IndividualBuffs.create());
			this.setCooldowns(eventID, proto.cooldowns || Cooldowns.create());
			this.setTalentsString(eventID, proto.talentsString);
			this.setGlyphs(eventID, proto.glyphs || Glyphs.create());
			this.setProfession1(eventID, proto.profession1);
			this.setProfession2(eventID, proto.profession2);
			this.setInFrontOfTarget(eventID, proto.inFrontOfTarget);
			this.setDistanceFromTarget(eventID, proto.distanceFromTarget);
			this.setHealingModel(eventID, proto.healingModel || HealingModel.create());
			this.setRotation(eventID, this.specTypeFunctions.rotationFromPlayer(proto));
			this.setSpecOptions(eventID, this.specTypeFunctions.optionsFromPlayer(proto));
		});
	}

	clone(eventID: EventID): Player<SpecType> {
		const newPlayer = new Player<SpecType>(this.spec, this.sim);
		newPlayer.fromProto(eventID, this.toProto());
		return newPlayer;
	}

	applySharedDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			this.setInFrontOfTarget(eventID, isTankSpec(this.spec));
			this.setHealingModel(eventID, HealingModel.create({
				burstWindow: isTankSpec(this.spec) ? 6 : 0,
			}));
			this.setCooldowns(eventID, Cooldowns.create({
				hpPercentForDefensives: isTankSpec(this.spec) ? 0.35 : 0,
			}));
			this.setBonusStats(eventID, new Stats());
		});
	}
}
