import { Player } from "../../player";
import {
	BattleElixir,
	Class,
	Conjured,
	Consumes,
	Explosive,
	Flask,
	Food,
	GuardianElixir,
	PetFood,
	Potions,
	Spec,
	Stat,
} from "../../proto/common";
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

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void
	showWhen?: (player: Player<any>) => boolean
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
					actionId: option.config.actionId,
					showWhen: (player: Player<any>) =>
						(!option.config.showWhen || option.config.showWhen(player)) &&
						(option.config.faction || player.getFaction()) == player.getFaction()
				} as IconEnumValueConfig<Player<any>, T>;
				if (option.config.value) rtn.value = option.config.value

				return rtn
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,	
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.gearChangeEmitter]),
			showWhen: (player: Player<any>) => !args.showWhen || args.showWhen(player),
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
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

export const ConjuredDarkRune = { actionId: ActionId.fromItemId(12662), value: Conjured.ConjuredDarkRune };
export const ConjuredFlameCap = { actionId: ActionId.fromItemId(22788), value: Conjured.ConjuredFlameCap };
export const ConjuredHealthstone = { actionId: ActionId.fromItemId(22105), value: Conjured.ConjuredHealthstone };
export const ConjuredRogueThistleTea = {
  actionId: ActionId.fromItemId(7676),
  value: Conjured.ConjuredRogueThistleTea,
  showWhen: (player: Player<Spec>) => player.getClass() == Class.ClassRogue
};

export const CONJURED_CONFIG = [
  { config: ConjuredRogueThistleTea, stats: [] },
	{ config: ConjuredHealthstone, stats: [Stat.StatStamina] },
  { config: ConjuredDarkRune, stats: [Stat.StatIntellect] },
  { config: ConjuredFlameCap, stats: [] },
] as ConsumableStatOption<Conjured>[]

export const makeConjuredInput = makeConsumeInputFactory({consumesFieldName: 'defaultConjured'});

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const ExplosiveSaroniteBomb    = { actionId: ActionId.fromItemId(41119), value: Explosive.ExplosiveSaroniteBomb };
export const ExplosiveCobaltFragBomb  = { actionId: ActionId.fromItemId(40771), value: Explosive.ExplosiveCobaltFragBomb };

export const EXPLOSIVES_CONFIG = [
	{ config: ExplosiveSaroniteBomb, stats: [] },
	{ config: ExplosiveCobaltFragBomb, stats: [] },
] as ConsumableStatOption<Explosive>[];

export const makeExplosivesInput = makeConsumeInputFactory({consumesFieldName: 'fillerExplosive'});

export const ThermalSapper = makeBooleanConsumeInput({actionId: ActionId.fromItemId(42641), fieldName: 'thermalSapper'});
export const ExplosiveDecoy = makeBooleanConsumeInput({actionId: ActionId.fromItemId(40536), fieldName: 'explosiveDecoy'});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS + ELIXIRS
///////////////////////////////////////////////////////////////////////////

// Flasks
export const FlaskOfTheFrostWyrm      = { actionId: ActionId.fromItemId(46376), value: Flask.FlaskOfTheFrostWyrm };
export const FlaskOfEndlessRage       = { actionId: ActionId.fromItemId(46377), value: Flask.FlaskOfEndlessRage };
export const FlaskOfPureMojo          = { actionId: ActionId.fromItemId(46378), value: Flask.FlaskOfPureMojo };
export const FlaskOfStoneblood        = { actionId: ActionId.fromItemId(46379), value: Flask.FlaskOfStoneblood };
export const LesserFlaskOfToughness   = { actionId: ActionId.fromItemId(40079), value: Flask.LesserFlaskOfToughness };
export const LesserFlaskOfResistance  = { actionId: ActionId.fromItemId(44939), value: Flask.LesserFlaskOfResistance };
export const FlaskOfBlindingLight     = { actionId: ActionId.fromItemId(22861), value: Flask.FlaskOfBlindingLight };
export const FlaskOfMightyRestoration = { actionId: ActionId.fromItemId(22853), value: Flask.FlaskOfMightyRestoration };
export const FlaskOfPureDeath         = { actionId: ActionId.fromItemId(22866), value: Flask.FlaskOfPureDeath };
export const FlaskOfRelentlessAssault = { actionId: ActionId.fromItemId(22854), value: Flask.FlaskOfRelentlessAssault };
export const FlaskOfSupremePower      = { actionId: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower };
export const FlaskOfFortification     = { actionId: ActionId.fromItemId(22851), value: Flask.FlaskOfFortification };
export const FlaskOfChromaticWonder   = { actionId: ActionId.fromItemId(33208), value: Flask.FlaskOfChromaticWonder };

export const FLASKS_CONFIG = [
  { config: FlaskOfTheFrostWyrm,      stats: [Stat.StatSpellPower] },
  { config: FlaskOfEndlessRage,       stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
  { config: FlaskOfPureMojo,          stats: [Stat.StatMP5] },
  { config: FlaskOfStoneblood,        stats: [Stat.StatStamina] },
  { config: LesserFlaskOfToughness,   stats: [Stat.StatResilience] },
  { config: LesserFlaskOfResistance,  stats: [Stat.StatArcaneResistance, Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance, Stat.StatShadowResistance] },
] as ConsumableStatOption<Flask>[];

export const makeFlasksInput = makeConsumeInputFactory({
  consumesFieldName: 'flask',
  onSet: (eventID: EventID, player: Player<any>, newValue: Flask) => {
    if (newValue) {
      const newConsumes = player.getConsumes();
      newConsumes.battleElixir = BattleElixir.BattleElixirUnknown;
      newConsumes.guardianElixir = GuardianElixir.GuardianElixirUnknown;
      player.setConsumes(eventID, newConsumes);
    }
  }
});

// Battle Elixirs
export const ElixirOfAccuracy         = { actionId: ActionId.fromItemId(44325), value: BattleElixir.ElixirOfAccuracy };
export const ElixirOfArmorPiercing    = { actionId: ActionId.fromItemId(44330), value: BattleElixir.ElixirOfArmorPiercing };
export const ElixirOfDeadlyStrikes    = { actionId: ActionId.fromItemId(44327), value: BattleElixir.ElixirOfDeadlyStrikes };
export const ElixirOfExpertise        = { actionId: ActionId.fromItemId(44329), value: BattleElixir.ElixirOfExpertise };
export const ElixirOfLightningSpeed   = { actionId: ActionId.fromItemId(44331), value: BattleElixir.ElixirOfLightningSpeed };
export const ElixirOfMightyAgility    = { actionId: ActionId.fromItemId(39666), value: BattleElixir.ElixirOfMightyAgility };
export const ElixirOfMightyStrength   = { actionId: ActionId.fromItemId(40073), value: BattleElixir.ElixirOfMightyStrength };
export const GurusElixir              = { actionId: ActionId.fromItemId(40076), value: BattleElixir.GurusElixir };
export const SpellpowerElixir         = { actionId: ActionId.fromItemId(40070), value: BattleElixir.SpellpowerElixir };
export const WrathElixir              = { actionId: ActionId.fromItemId(40068), value: BattleElixir.WrathElixir };
export const AdeptsElixir             = { actionId: ActionId.fromItemId(28103), value: BattleElixir.AdeptsElixir };
export const ElixirOfDemonslaying     = { actionId: ActionId.fromItemId(9224),  value: BattleElixir.ElixirOfDemonslaying };
export const ElixirOfMajorAgility     = { actionId: ActionId.fromItemId(22831), value: BattleElixir.ElixirOfMajorAgility };
export const ElixirOfMajorFirePower   = { actionId: ActionId.fromItemId(22833), value: BattleElixir.ElixirOfMajorFirePower };
export const ElixirOfMajorFrostPower  = { actionId: ActionId.fromItemId(22827), value: BattleElixir.ElixirOfMajorFrostPower };
export const ElixirOfMajorShadowPower = { actionId: ActionId.fromItemId(22835), value: BattleElixir.ElixirOfMajorShadowPower };
export const ElixirOfMajorStrength    = { actionId: ActionId.fromItemId(22824), value: BattleElixir.ElixirOfMajorStrength };
export const ElixirOfMastery          = { actionId: ActionId.fromItemId(28104), value: BattleElixir.ElixirOfMastery };
export const ElixirOfTheMongoose      = { actionId: ActionId.fromItemId(13452), value: BattleElixir.ElixirOfTheMongoose };
export const FelStrengthElixir        = { actionId: ActionId.fromItemId(31679), value: BattleElixir.FelStrengthElixir };
export const GreaterArcaneElixir      = { actionId: ActionId.fromItemId(13454), value: BattleElixir.GreaterArcaneElixir };

export const BATTLE_ELIXIRS_CONFIG = [
  { config: ElixirOfAccuracy,       stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
  { config: ElixirOfArmorPiercing,  stats: [Stat.StatArmorPenetration] },
  { config: ElixirOfDeadlyStrikes,  stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
  { config: ElixirOfExpertise,      stats: [Stat.StatExpertise] },
  { config: ElixirOfLightningSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: ElixirOfMightyAgility,  stats: [Stat.StatAgility] },
  { config: ElixirOfMightyStrength, stats: [Stat.StatStrength] },
  { config: GurusElixir,            stats: [Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatSpirit, Stat.StatIntellect] },
  { config: SpellpowerElixir,       stats: [Stat.StatSpellPower] },
  { config: WrathElixir,            stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
] as ConsumableStatOption<BattleElixir>[];

export const makeBattleElixirsInput = makeConsumeInputFactory({
  consumesFieldName: 'battleElixir',
  onSet: (eventID: EventID, player: Player<any>, newValue: BattleElixir) => {
    if (newValue) {
      const newConsumes = player.getConsumes();
      newConsumes.flask = Flask.FlaskUnknown;
      player.setConsumes(eventID, newConsumes);
    }
  }
});

// Guardian Elixirs
export const ElixirOfMightyDefense    = { actionId: ActionId.fromItemId(44328), value: GuardianElixir.ElixirOfMightyDefense };
export const ElixirOfMightyFortitude  = { actionId: ActionId.fromItemId(40078), value: GuardianElixir.ElixirOfMightyFortitude };
export const ElixirOfMightyMageblood  = { actionId: ActionId.fromItemId(40109), value: GuardianElixir.ElixirOfMightyMageblood };
export const ElixirOfMightyThoughts   = { actionId: ActionId.fromItemId(44332), value: GuardianElixir.ElixirOfMightyThoughts };
export const ElixirOfProtection       = { actionId: ActionId.fromItemId(40097), value: GuardianElixir.ElixirOfProtection };
export const ElixirOfSpirit           = { actionId: ActionId.fromItemId(40072), value: GuardianElixir.ElixirOfSpirit };
export const GiftOfArthas             = { actionId: ActionId.fromItemId(9088),  value: GuardianElixir.GiftOfArthas };
export const ElixirOfDraenicWisdom    = { actionId: ActionId.fromItemId(32067), value: GuardianElixir.ElixirOfDraenicWisdom };
export const ElixirOfIronskin         = { actionId: ActionId.fromItemId(32068), value: GuardianElixir.ElixirOfIronskin };
export const ElixirOfMajorDefense     = { actionId: ActionId.fromItemId(22834), value: GuardianElixir.ElixirOfMajorDefense };
export const ElixirOfMajorFortitude   = { actionId: ActionId.fromItemId(32062), value: GuardianElixir.ElixirOfMajorFortitude };
export const ElixirOfMajorMageblood   = { actionId: ActionId.fromItemId(22840), value: GuardianElixir.ElixirOfMajorMageblood };

export const GUARDIAN_ELIXIRS_CONFIG = [
  { config: ElixirOfMightyDefense,    stats: [Stat.StatDefense] },
  { config: ElixirOfMightyFortitude,  stats: [Stat.StatStamina] },
  { config: ElixirOfMightyMageblood,  stats: [Stat.StatMP5] },
  { config: ElixirOfMightyThoughts,   stats: [Stat.StatIntellect] },
  { config: ElixirOfProtection,       stats: [Stat.StatArmor] },
  { config: ElixirOfSpirit,           stats: [Stat.StatSpirit] },
  { config: GiftOfArthas,             stats: [Stat.StatStamina] },
] as ConsumableStatOption<GuardianElixir>[];

export const makeGuardianElixirsInput = makeConsumeInputFactory({
	consumesFieldName: 'guardianElixir',
	onSet: (eventID: EventID, player: Player<any>, newValue: GuardianElixir) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			newConsumes.flask = Flask.FlaskUnknown;
			player.setConsumes(eventID, newConsumes);
		}
	}
});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const FoodFishFeast              = { actionId: ActionId.fromItemId(43015), value: Food.FoodFishFeast };
export const FoodGreatFeast             = { actionId: ActionId.fromItemId(34753), value: Food.FoodGreatFeast };
export const FoodBlackenedDragonfin     = { actionId: ActionId.fromItemId(42999), value: Food.FoodBlackenedDragonfin };
export const FoodHeartyRhino            = { actionId: ActionId.fromItemId(42995), value: Food.FoodHeartyRhino };
export const FoodMegaMammothMeal        = { actionId: ActionId.fromItemId(34754), value: Food.FoodMegaMammothMeal };
export const FoodSpicedWormBurger       = { actionId: ActionId.fromItemId(34756), value: Food.FoodSpicedWormBurger };
export const FoodRhinoliciousWormsteak  = { actionId: ActionId.fromItemId(42994), value: Food.FoodRhinoliciousWormsteak };
export const FoodImperialMantaSteak     = { actionId: ActionId.fromItemId(34769), value: Food.FoodImperialMantaSteak };
export const FoodSnapperExtreme         = { actionId: ActionId.fromItemId(42996), value: Food.FoodSnapperExtreme };
export const FoodMightyRhinoDogs        = { actionId: ActionId.fromItemId(34758), value: Food.FoodMightyRhinoDogs };
export const FoodFirecrackerSalmon      = { actionId: ActionId.fromItemId(34767), value: Food.FoodFirecrackerSalmon };
export const FoodCuttlesteak            = { actionId: ActionId.fromItemId(42998), value: Food.FoodCuttlesteak };
export const FoodDragonfinFilet         = { actionId: ActionId.fromItemId(43000), value: Food.FoodDragonfinFilet };

export const FoodBlackenedBasilisk  = { actionId: ActionId.fromItemId(27657), value: Food.FoodBlackenedBasilisk };
export const FoodGrilledMudfish     = { actionId: ActionId.fromItemId(27664), value: Food.FoodGrilledMudfish };
export const FoodRavagerDog         = { actionId: ActionId.fromItemId(27655), value: Food.FoodRavagerDog };
export const FoodRoastedClefthoof   = { actionId: ActionId.fromItemId(27658), value: Food.FoodRoastedClefthoof };
export const FoodSpicyHotTalbuk     = { actionId: ActionId.fromItemId(33872), value: Food.FoodSpicyHotTalbuk };
export const FoodSkullfishSoup      = { actionId: ActionId.fromItemId(33825), value: Food.FoodSkullfishSoup };
export const FoodFishermansFeast    = { actionId: ActionId.fromItemId(33052), value: Food.FoodFishermansFeast };

export const FOOD_CONFIG = [
  { config: FoodFishFeast,              stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower] },
  { config: FoodGreatFeast,             stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower] },
  { config: FoodBlackenedDragonfin,     stats: [Stat.StatAgility] },
  { config: FoodDragonfinFilet,         stats: [Stat.StatStrength] },
  { config: FoodCuttlesteak,            stats: [Stat.StatSpirit] },
  { config: FoodMegaMammothMeal,        stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
  { config: FoodHeartyRhino,            stats: [Stat.StatArmorPenetration] },
  { config: FoodRhinoliciousWormsteak,  stats: [Stat.StatExpertise] },
  { config: FoodFirecrackerSalmon,      stats: [Stat.StatSpellPower] },
  { config: FoodSnapperExtreme,         stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
  { config: FoodSpicedWormBurger,       stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
  { config: FoodImperialMantaSteak,     stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: FoodMightyRhinoDogs,        stats: [Stat.StatMP5] },
] as ConsumableStatOption<Food>[];

export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

export const SpicedMammothTreats = makeBooleanConsumeInput({actionId: ActionId.fromItemId(43005), fieldName: 'petFood', value: PetFood.PetFoodSpicedMammothTreats});
export const PetScrollOfAgilityV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27498), fieldName: 'petScrollOfAgility', value: 5});
export const PetScrollOfStrengthV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27503), fieldName: 'petScrollOfStrength', value: 5});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const RunicHealingPotion   = { actionId: ActionId.fromItemId(33447), value: Potions.RunicHealingPotion };
export const RunicHealingInjector = { actionId: ActionId.fromItemId(41166), value: Potions.RunicHealingInjector };
export const RunicManaPotion      = { actionId: ActionId.fromItemId(33448), value: Potions.RunicManaPotion };
export const RunicManaInjector    = { actionId: ActionId.fromItemId(42545), value: Potions.RunicManaInjector };
export const IndestructiblePotion = { actionId: ActionId.fromItemId(40093), value: Potions.IndestructiblePotion };
export const PotionOfSpeed        = { actionId: ActionId.fromItemId(40211), value: Potions.PotionOfSpeed };
export const PotionOfWildMagic    = { actionId: ActionId.fromItemId(40212), value: Potions.PotionOfWildMagic };

export const DestructionPotion    = { actionId: ActionId.fromItemId(22839), value: Potions.DestructionPotion };
export const HastePotion          = { actionId: ActionId.fromItemId(22838), value: Potions.HastePotion };
export const MightyRagePotion     = { actionId: ActionId.fromItemId(13442), value: Potions.MightyRagePotion };
export const SuperManaPotion      = { actionId: ActionId.fromItemId(22832), value: Potions.SuperManaPotion };
export const FelManaPotion        = { actionId: ActionId.fromItemId(31677), value: Potions.FelManaPotion };
export const InsaneStrengthPotion = { actionId: ActionId.fromItemId(22828), value: Potions.InsaneStrengthPotion };
export const IronshieldPotion     = { actionId: ActionId.fromItemId(22849), value: Potions.IronshieldPotion };
export const HeroicPotion         = { actionId: ActionId.fromItemId(22837), value: Potions.HeroicPotion };

export const POTIONS_CONFIG = [
  { config: RunicHealingPotion,   stats: [Stat.StatStamina] },
  { config: RunicHealingInjector, stats: [Stat.StatStamina] },
  { config: RunicManaPotion,      stats: [Stat.StatIntellect] },
  { config: RunicManaInjector,    stats: [Stat.StatIntellect] },
  { config: IndestructiblePotion, stats: [Stat.StatArmor] },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength] },
  { config: HeroicPotion,         stats: [Stat.StatStamina] },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
] as ConsumableStatOption<Potions>[];

export const PRE_POTIONS_CONFIG = [
  { config: IndestructiblePotion, stats: [Stat.StatArmor] },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength] },
  { config: HeroicPotion,         stats: [Stat.StatStamina] },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
] as ConsumableStatOption<Potions>[];

export const makePotionsInput = makeConsumeInputFactory({consumesFieldName: 'defaultPotion'});
export const makePrepopPotionsInput = makeConsumeInputFactory({consumesFieldName: 'prepopPotion'});

