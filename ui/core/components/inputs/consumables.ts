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
	Potions,
	Scroll,
	SpellPowerBuff,
	Stat,
	StrengthBuff,
	WeaponImbue } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";
import { EventID, TypedEvent } from "../../typed_event";

import { IconEnumValueConfig } from "../icon_enum_picker";
import { makeBooleanConsumeInput, makeBooleanDebuffInput, withLabel } from "../icon_inputs";

import { ItemInputConfig, ItemStatOption } from "./stat_options";

import * as InputHelpers from '../input_helpers';

export interface ConsumableInputConfig<T> extends ItemInputConfig {
	value: T,
}

export interface ConsumableStatOption<T> extends ItemStatOption {
	item: ConsumableInputConfig<T>
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

export const makeAgilityConsumeInput = makeConsumeInputFactory({consumesFieldName: 'agilityElixir'})
export const makeFlasksInput = makeConsumeInputFactory({consumesFieldName: 'flask'});
export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});
export const makeOffHandImbuesInput = makeConsumeInputFactory({consumesFieldName: 'offHandImbue'});
export const makeMainHandImbuesInput = makeConsumeInputFactory({consumesFieldName: 'mainHandImbue'});
export const makeStrengthConsumeInput = makeConsumeInputFactory({consumesFieldName: 'strengthBuff'})

export const StrengthBuffInput = makeConsumeInput('strengthBuff', [
	{ id: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower, minLevel: 46 },
	{ id: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants, minLevel: 46 },
  { id: ActionId.fromItemId(3391), value: StrengthBuff.ElixirOfOgresStrength, minLevel: 20},
	{ id: ActionId.fromItemId(10310), value: StrengthBuff.ScrollOfStrength },
] as Array<IconEnumValueConfig<Player<any>, StrengthBuff>>, (p) => p.getLevel() >= 20);

export const SpellDamageBuff = makeConsumeInput('spellPowerBuff', [
	{ id: ActionId.fromItemId(9155), value: SpellPowerBuff.ArcaneElixir, minLevel: 37 },
	{ id: ActionId.fromItemId(13454), value: SpellPowerBuff.GreaterArcaneElixir, minLevel: 46 },
] as Array<IconEnumValueConfig<Player<any>, SpellPowerBuff>>, (p) => p.getLevel() >= 37);

export const FireDamageBuff = makeConsumeInput('firePowerBuff', [
	{ id: ActionId.fromItemId(6373), value: FirePowerBuff.ElixirOfFirepower, minLevel: 18 },
	{ id: ActionId.fromItemId(21546), value: FirePowerBuff.ElixirOfGreaterFirepower, minLevel: 40 },
] as Array<IconEnumValueConfig<Player<any>, FirePowerBuff>>, (p) => p.getLevel() >= 18);

export const ShadowDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(9264), fieldName: 'shadowPowerBuff', minLevel: 40});
export const FrostDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(17708), fieldName: 'frostPowerBuff', minLevel: 40});

export const FillerExplosiveInput = makeConsumeInput('fillerExplosive', [
	{ id: ActionId.fromItemId(18641), value: Explosive.ExplosiveDenseDynamite, showWhen: (p) => p.getLevel() >= 40 },
	{ id: ActionId.fromItemId(15993), value: Explosive.ExplosiveThoriumGrenade, showWhen: (p) => p.getLevel() >= 40 },
] as Array<IconEnumValueConfig<Player<any>, Explosive>>);

// TODO: These should be moved to consumes
export const GiftOfArthas = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(11374), fieldName: 'giftOfArthas'}),
	'Gift of Arthas',
);
export const CrystalYield = withLabel(
	makeBooleanDebuffInput({id: ActionId.fromSpellId(15235), fieldName: 'crystalYield'}),
	'Crystal Yield',
);

Consumes
export const Sapper = makeBooleanConsumeInput({id: ActionId.fromItemId(10646), fieldName: 'sapper', minLevel: 40});

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
					id: option.item.id,
					showWhen: (player: Player<any>) =>
						(option.item.minLevel || 0) <= player.getLevel() &&
						(option.item.maxLevel || MAX_CHARACTER_LEVEL) >= player.getLevel() &&
						(option.item.faction || player.getFaction()) == player.getFaction()
				} as IconEnumValueConfig<Player<any>, T>;
				if (option.item.value) rtn.value = option.item.value

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

function makeConsumeInput<T extends number>(consumesFieldName: keyof Consumes, allOptions: Array<IconEnumValueConfig<Player<any>, T>>, showWhen?: (obj: Player<any>) => boolean, onSet?: (eventID: EventID, p: Player<any>, newValue: T) => void): InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	// const factory = makeConsumeInputFactory({
	// 	consumesFieldName: consumesFieldName,
	// 	allOptions: allOptions,
	// 	onSet: onSet,
	// 	showWhen: showWhen,
	// });
	// return factory(allOptions.map(option => option.value));
}



///////////////////////////////////////////////////////////////////////////
//                                 AGILITY CONSUMES
///////////////////////////////////////////////////////////////////////////

export const ElixirOfTheMongoose = { id: ActionId.fromItemId(13452), value: AgilityElixir.ElixirOfTheMongoose, minLevel: 46 }
export const ElixirOfGreaterAgility = { id: ActionId.fromItemId(9187), value: AgilityElixir.ElixirOfGreaterAgility, minLevel: 38 }
export const ElixirOfLesserAgility = { id: ActionId.fromItemId(3390), value: AgilityElixir.ElixirOfLesserAgility, minLevel: 18 }
export const ScrollOfAgility = { id: ActionId.fromItemId(10309), value: AgilityElixir.ScrollOfAgility}

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

export const FlaskOfTheTitans = { id: ActionId.fromItemId(13510), value: Flask.FlaskOfTheTitans, minLevel: 50 }
export const FlaskOfDistilledWisdom = { id: ActionId.fromItemId(13511), value: Flask.FlaskOfDistilledWisdom, minLevel: 50 }
export const FlaskOfSupremePower = { id: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower, minLevel: 50 }
export const FlaskOfChromaticResistance = { id: ActionId.fromItemId(13513), value: Flask.FlaskOfChromaticResistance, minLevel: 50 }

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const SmokedSagefish = { id: ActionId.fromItemId(21072), value: Food.FoodSmokedSagefish, minLevel: 10 }
export const HotWolfRibs = { id: ActionId.fromItemId(13851), value: Food.FoodHotWolfRibs, minLevel: 25 }
export const TenderWolfSteak = { id: ActionId.fromItemId(22480), value: Food.FoodTenderWolfSteak, minLevel: 40 }
export const NightfinSoup = { id: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup, minLevel: 35 }
export const GrilledSquid = { id: ActionId.fromItemId(13928), value: Food.FoodGrilledSquid, minLevel: 35 }
export const SmokedDesertDumpling = { id: ActionId.fromItemId(20452), value: Food.FoodSmokedDesertDumpling, minLevel: 45 }
export const RunnTumTuberSurprise = { id: ActionId.fromItemId(18254), value: Food.FoodRunnTumTuberSurprise, minLevel: 45 }
export const BlessedSunfruitJuice = { id: ActionId.fromItemId(13813), value: Food.FoodBlessedSunfruitJuice, minLevel: 45 }
export const BlessSunfruit = { id: ActionId.fromItemId(13810), value: Food.FoodBlessSunfruit, minLevel: 45 }
export const DirgesKickChimaerokChops = { id: ActionId.fromItemId(21023), value: Food.FoodDirgesKickChimaerokChops, minLevel: 55 }

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const LesserManaPotion = { id: ActionId.fromItemId(3385), value: Potions.LesserManaPotion }
export const ManaPotion = { id: ActionId.fromItemId(3385), value: Potions.ManaPotion }

///////////////////////////////////////////////////////////////////////////
//                                 STRENGTH CONSUMES
///////////////////////////////////////////////////////////////////////////

export const JujuPower = { id: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower, minLevel: 46 }
export const ElixirOfGiants = { id: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants, minLevel: 46 }
export const ElixirOfOgresStrength = { id: ActionId.fromItemId(3391), value: StrengthBuff.ElixirOfOgresStrength, minLevel: 20}
export const ScrollOfStrength = { id: ActionId.fromItemId(10310), value: StrengthBuff.ScrollOfStrength }

///////////////////////////////////////////////////////////////////////////
//                                 Weapon Imbues
///////////////////////////////////////////////////////////////////////////

export const BrillianWizardOil = { id: ActionId.fromItemId(20749), value: WeaponImbue.BrillianWizardOil, minLevel: 45 }
export const BrilliantManaOil = { id: ActionId.fromItemId(20748), value: WeaponImbue.BrilliantManaOil, minLevel: 45 }
export const DenseSharpeningStone = { id: ActionId.fromItemId(12404), value: WeaponImbue.DenseSharpeningStone, minLevel: 35 }
export const ElementalSharpeningStone = { id: ActionId.fromItemId(18262), value: WeaponImbue.ElementalSharpeningStone, minLevel: 50 }
export const BlackfathomManaOil = { id: ActionId.fromItemId(211848), value: WeaponImbue.BlackfathomManaOil, minLevel: 25 }
export const BlackfathomSharpeningStone = { id: ActionId.fromItemId(211845), value: WeaponImbue.BlackfathomSharpeningStone }
export const WildStrikes = { id: ActionId.fromSpellId(407975), value: WeaponImbue.WildStrikes }

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const AGILITY_CONSUMES_CONFIG = [
	{ item: ElixirOfTheMongoose, stats: [Stat.StatAgility] },
	{ item: ElixirOfGreaterAgility, stats: [Stat.StatAgility] },
	{ item: ElixirOfLesserAgility, stats: [Stat.StatAgility] },
	{ item: ScrollOfAgility, stats: [Stat.StatAgility] },
] as ConsumableStatOption<AgilityElixir>[];

export const FLASKS_CONFIG = [
	{ item: FlaskOfTheTitans, stats: [Stat.StatStamina] },
	{ item: FlaskOfDistilledWisdom, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ item: FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ item: FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
] as ConsumableStatOption<Flask>[];

export const FOOD_CONFIG = [
	{ item: HotWolfRibs, stats: [Stat.StatSpirit] },
	{ item: SmokedSagefish, stats: [Stat.StatMP5] },
	{ item: NightfinSoup, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ item: GrilledSquid, stats: [Stat.StatAgility] },
	{ item: SmokedDesertDumpling, stats: [Stat.StatStrength] },
	{ item: RunnTumTuberSurprise, stats: [Stat.StatIntellect] },
	{ item: DirgesKickChimaerokChops, stats: [Stat.StatStamina] },
	{ item: BlessSunfruit, stats: [Stat.StatStrength] },
	{ item: BlessedSunfruitJuice, stats: [Stat.StatSpirit] },
] as ConsumableStatOption<Food>[];

export const POTIONS_CONFIG = [
	{ item: LesserManaPotion, stats: [Stat.StatIntellect] },
	{ item: ManaPotion, stats: [Stat.StatIntellect] },
] as ConsumableStatOption<Potions>[];

export const STRENGTH_CONSUMES_CONFIG = [
	{ item: JujuPower, stats: [Stat.StatStrength] },
	{ item: ElixirOfGiants, stats: [Stat.StatStrength] },
	{ item: ElixirOfOgresStrength, stats: [Stat.StatStrength] },
	{ item: ScrollOfStrength, stats: [Stat.StatStrength] },
]

export const WEAPON_IMBUES_MH_CONFIG = [
	{ item: BrillianWizardOil, stats: [Stat.StatSpellPower] },
	{ item: BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
	{ item: DenseSharpeningStone, stats: [Stat.StatAttackPower] },
	{ item: ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
	{ item: BlackfathomManaOil, stats: [Stat.StatSpellPower, Stat.StatMP5] },
	{ item: BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
	{ item: WildStrikes, stats: [Stat.StatMeleeHit] },
] as ConsumableStatOption<WeaponImbue>[];

export const WEAPON_IMBUES_OF_CONFIG = [
	{ item: BrillianWizardOil, stats: [Stat.StatSpellPower] },
	{ item: BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
	{ item: DenseSharpeningStone, stats: [Stat.StatAttackPower] },
	{ item: ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
	{ item: BlackfathomManaOil, stats: [Stat.StatSpellPower, Stat.StatMP5] },
	{ item: BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
] as ConsumableStatOption<WeaponImbue>[];
