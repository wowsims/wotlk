import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Encounter } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Player } from '/wotlk/core/player.js';
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { Alchohol } from '/wotlk/core/proto/common.js';
import { BattleElixir } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { GuardianElixir } from '/wotlk/core/proto/common.js';
import { Conjured } from '/wotlk/core/proto/common.js';

import { PetFood } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';

import {
	Hunter,
	Hunter_Rotation as HunterRotation,
	Hunter_Options as HunterOptions,
	Hunter_Options_PetType as PetType,
} from '/wotlk/core/proto/hunter.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, {
			cssClass: 'hunter-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],
			warnings: [
				// Warning when using exotic pet without BM talented.
				(simUI: IndividualSimUI<Spec.SpecHunter>) => {
					return {
						updateOn: TypedEvent.onAny([simUI.player.talentsChangeEmitter, simUI.player.specOptionsChangeEmitter]),
						getContent: () => {
							const petIsExotic = [
								PetType.Chimaera,
								PetType.CoreHound,
								PetType.Devilsaur,
								PetType.Silithid,
								PetType.SpiritBeast,
								PetType.Worm,
							].includes(simUI.player.getSpecOptions().petType);

							const isBM = simUI.player.getTalents().beastMastery;

							if (petIsExotic && !isBM) {
								return 'Cannot use exotic pets without the Beast Mastery talent.';
							} else {
								return '';
							}
						},
					};
				},
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatRangedAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],
			// Reference stat against which to calculate EP.
			epReferenceStat: Stat.StatRangedAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatIntellect,
				Stat.StatAttackPower,
				Stat.StatRangedAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.01,
					[Stat.StatAgility]: 2.5,
					[Stat.StatStrength]: 0.15,
					[Stat.StatAttackPower]: 0.15,
					[Stat.StatRangedAttackPower]: 1.0,
					[Stat.StatMeleeHit]: 0.3,
					[Stat.StatMeleeCrit]: 2.3,
					[Stat.StatMeleeHaste]: 1.97,
					[Stat.StatArmorPenetration]: 0.4,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.SurvivalTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					windfuryTotem: TristateEffect.TristateEffectImproved,
					battleShout: TristateEffect.TristateEffectImproved,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: 2,
					blessingOfMight: 2,
				}),
				debuffs: Debuffs.create({
					sunderArmor: true,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
					faerieFire: TristateEffect.TristateEffectImproved,
					judgementOfWisdom: true,
					curseOfElements: true,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				HunterInputs.WeaponAmmo,
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.ArcaneBrilliance,
				IconInputs.DivineSpirit,
				IconInputs.GiftOfTheWild,
				IconInputs.Bloodlust,
				IconInputs.WindfuryTotem,
				IconInputs.StrengthOfEarthTotem,
				IconInputs.ManaSpringTotem,
				IconInputs.BattleShout,
				IconInputs.LeaderOfThePack,
				IconInputs.TrueshotAura,
				IconInputs.UnleashedRage,
			],
			partyBuffInputs: [
				IconInputs.HeroicPresence,
				IconInputs.BraidedEterniumChain,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfWisdom,
				IconInputs.BlessingOfMight,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.BloodFrenzy,
				IconInputs.JudgementOfWisdom,
				IconInputs.FaerieFire,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.CurseOfWeakness,
				IconInputs.CurseOfElements,
				IconInputs.Misery,
				IconInputs.GiftOfArthas,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.HastePotion,
					Potions.SuperManaPotion,
					Potions.FelManaPotion,
				],
				conjured: [
					Conjured.ConjuredDarkRune,
					Conjured.ConjuredFlameCap,
				],
				flasks: [
					Flask.FlaskOfRelentlessAssault,
				],
				battleElixirs: [
					BattleElixir.ElixirOfDemonslaying,
					BattleElixir.ElixirOfMajorAgility,
					BattleElixir.ElixirOfTheMongoose,
				],
				guardianElixirs: [
					GuardianElixir.ElixirOfDraenicWisdom,
					GuardianElixir.ElixirOfMajorMageblood,
				],
				food: [
					Food.FoodGrilledMudfish,
					Food.FoodRavagerDog,
					Food.FoodSpicyHotTalbuk,
					Food.FoodRoastedClefthoof,
				],
				alcohol: [
				],
				weaponImbues: [
					WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
					WeaponImbue.WeaponImbueAdamantiteWeightstone,
					WeaponImbue.WeaponImbueRighteousWeaponCoating,
				],
				pet: [
					IconInputs.KiblersBits,
					IconInputs.PetScrollOfAgilityV,
					IconInputs.PetScrollOfStrengthV,
				],
				other: [
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfStrengthV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: HunterInputs.HunterRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					HunterInputs.PetTypeInput,
					HunterInputs.PetUptime,
					//HunterInputs.PetSingleAbility,
					HunterInputs.SniperTrainingUptime,
					//HunterInputs.LatencyMs,
					OtherInputs.PrepopPotion,
					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			encounterPicker: {
				// Target stats to show for 'Simple' encounters.
				simpleTargetStats: [
					Stat.StatArmor,
				],
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.BeastMasteryTalents,
					Presets.MarksmanTalents,
					Presets.SurvivalTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_PRESET,
				],
			},
		});
	}
}
