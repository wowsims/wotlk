import { MAX_CHARACTER_LEVEL } from "../../constants/mechanics";
import { Player } from "../../player";
import {
	AgilityElixir,
	Conjured,
	Consumes,
	Explosive,
	Faction,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	Potions,
	ShadowPowerBuff,
	SpellPowerBuff,
	Stat,
	StrengthBuff,
	WeaponImbue } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";
import { EventID, TypedEvent } from "../../typed_event";

import { IconEnumValueConfig } from "../icon_enum_picker";
import { makeBooleanConsumeInput } from "../icon_inputs";

import { ActionInputConfig, ItemStatOption } from "./stat_options";

import * as InputHelpers from '../input_helpers';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T,
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>
}

function makePotionInputFactory(consumesFieldName: keyof Consumes): (options: Array<Potions>, tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, Potions> {
	return makeConsumeInputFactory({
		consumesFieldName: consumesFieldName,
		allOptions: [
			{ id: ActionId.fromItemId(3385), value: Potions.LesserManaPotion },
			{ id: ActionId.fromItemId(3827), value: Potions.ManaPotion },
		] as Array<IconEnumValueConfig<Player<any>, Potions>>,
	});
}
export const makePotionsInput = makePotionInputFactory('defaultPotion');

// TODO: Classic? 
export const makeConjuredInput = makeConsumeInputFactory({
	consumesFieldName: 'defaultConjured',
	allOptions: [
		{ id: ActionId.fromItemId(4381), value: Conjured.ConjuredMinorRecombobulator, showWhen: (player: Player<any>) => player.getGear().hasTrinket(4381) },
		{ id: ActionId.fromItemId(12662), value: Conjured.ConjuredDemonicRune, showWhen: (p) => p.getLevel() >= 40 },
	] as Array<IconEnumValueConfig<Player<any>, Conjured>>
});

// TODO: Classic
export const PetScrollOfAgilityV = makeBooleanConsumeInput(ActionId.fromItemId(27498), 'petScrollOfAgility', 5);
export const PetScrollOfStrengthV = makeBooleanConsumeInput(ActionId.fromItemId(27503), 'petScrollOfStrength', 5);

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventID: EventID, p: Player<any>, newValue: T) => void
	minLevel?: number,
	maxLevel?: number,
	faction?: Faction,
}

function makeConsumeInputFactory<T extends number>(args: ConsumeInputFactoryArgs<T>): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [
				{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
			].concat(options.map(option => {
				const rtn = {
					id: option.config.id,
					showWhen: (player: Player<any>) =>
						(option.config.minLevel || 0) <= player.getLevel() &&
						(option.config.maxLevel || MAX_CHARACTER_LEVEL) >= player.getLevel() &&
						(option.config.faction || player.getFaction()) == player.getFaction()
				} as IconEnumValueConfig<Player<any>, T>;
				if (option.config.value) rtn.value = option.config.value

				return rtn
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,	
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter, player.gearChangeEmitter]),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();

				if (newConsumes[args.consumesFieldName] === newValue){
					return;
				}

				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}







///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const ExplosiveDenseDynamite = { id: ActionId.fromItemId(18641), value: Explosive.ExplosiveDenseDynamite, minLevel: 40 };
export const ExplosiveThoriumGrenade = { id: ActionId.fromItemId(15993), value: Explosive.ExplosiveThoriumGrenade, minLevel: 40 };

export const EXPLOSIVES_CONFIG = [
	{ config: ExplosiveDenseDynamite, stats: [] },
	{ config: ExplosiveThoriumGrenade, stats: [] },
] as ConsumableStatOption<Explosive>[];

export const makeExplosivesInput = makeConsumeInputFactory({consumesFieldName: 'fillerExplosive'});

export const Sapper = makeBooleanConsumeInput({id: ActionId.fromItemId(10646), fieldName: 'sapper', minLevel: 40});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

export const FlaskOfTheTitans = { id: ActionId.fromItemId(13510), value: Flask.FlaskOfTheTitans, minLevel: 50 };
export const FlaskOfDistilledWisdom = { id: ActionId.fromItemId(13511), value: Flask.FlaskOfDistilledWisdom, minLevel: 50 };
export const FlaskOfSupremePower = { id: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower, minLevel: 50 };
export const FlaskOfChromaticResistance = { id: ActionId.fromItemId(13513), value: Flask.FlaskOfChromaticResistance, minLevel: 50 };

export const FLASKS_CONFIG = [
	{ config: FlaskOfTheTitans, stats: [Stat.StatStamina] },
	{ config: FlaskOfDistilledWisdom, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
] as ConsumableStatOption<Flask>[];

export const makeFlasksInput = makeConsumeInputFactory({consumesFieldName: 'flask'});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const SmokedSagefish = { id: ActionId.fromItemId(21072), value: Food.FoodSmokedSagefish, minLevel: 10 };
export const HotWolfRibs = { id: ActionId.fromItemId(13851), value: Food.FoodHotWolfRibs, minLevel: 25 };
export const TenderWolfSteak = { id: ActionId.fromItemId(22480), value: Food.FoodTenderWolfSteak, minLevel: 40 };
export const NightfinSoup = { id: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup, minLevel: 35 };
export const GrilledSquid = { id: ActionId.fromItemId(13928), value: Food.FoodGrilledSquid, minLevel: 35 };
export const SmokedDesertDumpling = { id: ActionId.fromItemId(20452), value: Food.FoodSmokedDesertDumpling, minLevel: 45 };
export const RunnTumTuberSurprise = { id: ActionId.fromItemId(18254), value: Food.FoodRunnTumTuberSurprise, minLevel: 45 };
export const BlessedSunfruitJuice = { id: ActionId.fromItemId(13813), value: Food.FoodBlessedSunfruitJuice, minLevel: 45 };
export const BlessSunfruit = { id: ActionId.fromItemId(13810), value: Food.FoodBlessSunfruit, minLevel: 45 };
export const DirgesKickChimaerokChops = { id: ActionId.fromItemId(21023), value: Food.FoodDirgesKickChimaerokChops, minLevel: 55 };

export const FOOD_CONFIG = [
	{ config: HotWolfRibs, stats: [Stat.StatSpirit] },
	{ config: SmokedSagefish, stats: [Stat.StatMP5] },
	{ config: NightfinSoup, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: GrilledSquid, stats: [Stat.StatAgility] },
	{ config: SmokedDesertDumpling, stats: [Stat.StatStrength] },
	{ config: RunnTumTuberSurprise, stats: [Stat.StatIntellect] },
	{ config: DirgesKickChimaerokChops, stats: [Stat.StatStamina] },
	{ config: BlessSunfruit, stats: [Stat.StatStrength] },
	{ config: BlessedSunfruitJuice, stats: [Stat.StatSpirit] },
] as ConsumableStatOption<Food>[];

export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});

///////////////////////////////////////////////////////////////////////////
//                                 PHYSICAL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Agility
export const ElixirOfTheMongoose = { id: ActionId.fromItemId(13452), value: AgilityElixir.ElixirOfTheMongoose, minLevel: 46 };
export const ElixirOfGreaterAgility = { id: ActionId.fromItemId(9187), value: AgilityElixir.ElixirOfGreaterAgility, minLevel: 38 };
export const ElixirOfLesserAgility = { id: ActionId.fromItemId(3390), value: AgilityElixir.ElixirOfLesserAgility, minLevel: 18 };
export const ScrollOfAgility = { id: ActionId.fromItemId(10309), value: AgilityElixir.ScrollOfAgility };

export const AGILITY_CONSUMES_CONFIG = [
	{ config: ElixirOfTheMongoose, stats: [Stat.StatAgility] },
	{ config: ElixirOfGreaterAgility, stats: [Stat.StatAgility] },
	{ config: ElixirOfLesserAgility, stats: [Stat.StatAgility] },
	{ config: ScrollOfAgility, stats: [Stat.StatAgility] },
] as ConsumableStatOption<AgilityElixir>[];

export const makeAgilityConsumeInput = makeConsumeInputFactory({consumesFieldName: 'agilityElixir'});

// Strength
export const JujuPower = { id: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower, minLevel: 46 };
export const ElixirOfGiants = { id: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants, minLevel: 46 };
export const ElixirOfOgresStrength = { id: ActionId.fromItemId(3391), value: StrengthBuff.ElixirOfOgresStrength, minLevel: 20 };
export const ScrollOfStrength = { id: ActionId.fromItemId(10310), value: StrengthBuff.ScrollOfStrength };

export const STRENGTH_CONSUMES_CONFIG = [
	{ config: JujuPower, stats: [Stat.StatStrength] },
	{ config: ElixirOfGiants, stats: [Stat.StatStrength] },
	{ config: ElixirOfOgresStrength, stats: [Stat.StatStrength] },
	{ config: ScrollOfStrength, stats: [Stat.StatStrength] },
] as ConsumableStatOption<StrengthBuff>[];

export const makeStrengthConsumeInput = makeConsumeInputFactory({consumesFieldName: 'strengthBuff'});

// Other
export const BoglingRootDebuff = makeBooleanConsumeInput({id: ActionId.fromItemId(5206), fieldName: 'boglingRoot'});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const LesserManaPotion = { id: ActionId.fromItemId(3385), value: Potions.LesserManaPotion };
export const ManaPotion = { id: ActionId.fromItemId(3385), value: Potions.ManaPotion };

export const POTIONS_CONFIG = [
	{ config: LesserManaPotion, stats: [Stat.StatIntellect] },
	{ config: ManaPotion, stats: [Stat.StatIntellect] },
] as ConsumableStatOption<Potions>[];

///////////////////////////////////////////////////////////////////////////
//                                 SPELL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Arcane
export const ArcaneElixir = { id: ActionId.fromItemId(9155), value: SpellPowerBuff.ArcaneElixir, minLevel: 37 };
export const GreaterArcaneElixir = { id: ActionId.fromItemId(13454), value: SpellPowerBuff.GreaterArcaneElixir, minLevel: 46 };

export const SPELL_POWER_CONFIG = [
	{
		config: ArcaneElixir,
		stats: [Stat.StatSpellPower],
	},
	{
		config: GreaterArcaneElixir,
		stats: [Stat.StatSpellPower],
	}
] as ConsumableStatOption<SpellPowerBuff>[];

export const makeSpellPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'spellPowerBuff'})

// Fire
export const ElixirOfFirepower = { id: ActionId.fromItemId(6373), value: FirePowerBuff.ElixirOfFirepower, minLevel: 18 };
export const ElixirOfGreaterFirepower = { id: ActionId.fromItemId(21546), value: FirePowerBuff.ElixirOfGreaterFirepower, minLevel: 40 };

export const FIRE_POWER_CONFIG = [
	{
		config: ElixirOfFirepower,
		stats: [Stat.StatFirePower],
	},
	{
		config: ElixirOfGreaterFirepower,
		stats: [Stat.StatFirePower],
	}
] as ConsumableStatOption<FirePowerBuff>[];

export const makeFirePowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'firePowerBuff'})

// Frost
export const ElixirOfFrostPower = {id: ActionId.fromItemId(17708), value: FrostPowerBuff.ElixirOfFrostPower, minLevel: 40 };

export const FROST_POWER_CONFIG = [
	{
		config: ElixirOfFrostPower,
		stats: [Stat.StatFrostPower],
	},
] as ConsumableStatOption<FrostPowerBuff>[];

export const makeFrostPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'frostPowerBuff'})

// Shadow
export const ElixirOfShadowPower = {id: ActionId.fromItemId(9264), value: ShadowPowerBuff.ElixirOfShadowPower, minLevel: 40 };

export const SHADOW_POWER_CONFIG = [
	{
		config: ElixirOfShadowPower,
		stats: [Stat.StatShadowPower],
	},
] as ConsumableStatOption<ShadowPowerBuff>[];

export const makeshadowPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'shadowPowerBuff'})

///////////////////////////////////////////////////////////////////////////
//                                 Weapon Imbues
///////////////////////////////////////////////////////////////////////////

export const BrillianWizardOil = { id: ActionId.fromItemId(20749), value: WeaponImbue.BrillianWizardOil, minLevel: 45 };
export const BrilliantManaOil = { id: ActionId.fromItemId(20748), value: WeaponImbue.BrilliantManaOil, minLevel: 45 };
export const DenseSharpeningStone = { id: ActionId.fromItemId(12404), value: WeaponImbue.DenseSharpeningStone, minLevel: 35 };
export const ElementalSharpeningStone = { id: ActionId.fromItemId(18262), value: WeaponImbue.ElementalSharpeningStone, minLevel: 50 };
export const BlackfathomManaOil = { id: ActionId.fromItemId(211848), value: WeaponImbue.BlackfathomManaOil, minLevel: 25 };
export const BlackfathomSharpeningStone = { id: ActionId.fromItemId(211845), value: WeaponImbue.BlackfathomSharpeningStone };
export const WildStrikes = { id: ActionId.fromSpellId(407975), value: WeaponImbue.WildStrikes };

export const WEAPON_IMBUES_OH_CONFIG = [
	{ config: BrillianWizardOil, stats: [Stat.StatSpellPower] },
	{ config: BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
	{ config: DenseSharpeningStone, stats: [Stat.StatAttackPower] },
	{ config: ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
	{ config: BlackfathomManaOil, stats: [Stat.StatSpellPower, Stat.StatMP5] },
	{ config: BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
] as ConsumableStatOption<WeaponImbue>[];

export const WEAPON_IMBUES_MH_CONFIG = [
	...WEAPON_IMBUES_OH_CONFIG,
	{ config: WildStrikes, stats: [Stat.StatMeleeHit] },
] as ConsumableStatOption<WeaponImbue>[];

export const makeMainHandImbuesInput = makeConsumeInputFactory({consumesFieldName: 'mainHandImbue'});
export const makeOffHandImbuesInput = makeConsumeInputFactory({consumesFieldName: 'offHandImbue'});
