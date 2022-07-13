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
import { TotemsSection } from '/wotlk/core/components/totem_inputs.js';

import { Alchohol } from '/wotlk/core/proto/common.js';
import { BattleElixir } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { GuardianElixir } from '/wotlk/core/proto/common.js';
import { Conjured } from '/wotlk/core/proto/common.js';
import { PetFood } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { DeathKnight, DeathKnight_Rotation as DeathKnightRotation, DeathKnightTalents as DeathKnightTalents, DeathKnight_Options as DeathKnightOptions } from '/wotlk/core/proto/deathknight.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as DeathKnightInputs from './inputs.js';
import * as Presets from './presets.js';

export class DeathKnightSimUI extends IndividualSimUI<Spec.SpecDeathKnight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDeathKnight>) {
		super(parentElem, player, {
			cssClass: 'deathknight-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.P1_FROST_BIS_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStrength]: 2.17,
					[Stat.StatAgility]: 1.4,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 3.29,
					[Stat.StatMeleeHit]: 0.41,
					[Stat.StatMeleeCrit]: 1.83,
					[Stat.StatMeleeHaste]: 2.07,
					[Stat.StatArmorPenetration]: 0.5,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.FrostTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectMissing,
					bloodlust: false,
					strengthOfEarthTotem: TristateEffect.TristateEffectMissing,
					windfuryTotem: TristateEffect.TristateEffectMissing,
					leaderOfThePack: TristateEffect.TristateEffectMissing,
					unleashedRage: false,
					abominationsMight: true,
					icyTalons: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: false,
					blessingOfWisdom: TristateEffect.TristateEffectMissing,
					blessingOfMight: TristateEffect.TristateEffectMissing,
				}),
				debuffs: Debuffs.create({
					bloodFrenzy: false,
					sunderArmor: false,
					curseOfWeakness: TristateEffect.TristateEffectMissing,
					curseOfElements: false,
					faerieFire: TristateEffect.TristateEffectMissing,
					judgementOfWisdom: false,
					misery: false,
					huntersMark: TristateEffect.TristateEffectMissing,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.GiftOfTheWild,
				IconInputs.Bloodlust,
				IconInputs.WrathOfAirTotem,
				IconInputs.TotemOfWrath,
				IconInputs.BattleShout,
				IconInputs.LeaderOfThePack,
				IconInputs.MoonkinAura,
				IconInputs.TrueshotAura,
			],
			partyBuffInputs: [
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
				IconInputs.HuntersMark,
				IconInputs.FaerieFire,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.CurseOfWeakness,
				IconInputs.CurseOfElements,
				IconInputs.Misery,
				IconInputs.ImprovedScorch,
				IconInputs.WintersChill,
				IconInputs.GiftOfArthas,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.HastePotion,
					Potions.InsaneStrengthPotion,
					Potions.MightyRagePotion,
				],
				conjured: [
					Conjured.ConjuredFlameCap,
				],
				flasks: [
					Flask.FlaskOfRelentlessAssault,
				],
				battleElixirs: [
					BattleElixir.ElixirOfDemonslaying,
					BattleElixir.ElixirOfMajorStrength,
					BattleElixir.ElixirOfMajorAgility,
					BattleElixir.ElixirOfTheMongoose,
					BattleElixir.FelStrengthElixir,
				],
				guardianElixirs: [
				],
				food: [
					Food.FoodRoastedClefthoof,
					Food.FoodGrilledMudfish,
					Food.FoodSpicyHotTalbuk,
					Food.FoodRavagerDog,
				],
				alcohol: [
				],
				weaponImbues: [
				],
				other: [
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfStrengthV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DeathKnightInputs.DeathKnightRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					DeathKnightInputs.StartingRunicPower,

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

			// If true, the talents on the talents tab will not be individually modifiable by the user.
			// Note that the use can still pick between preset talents, if there is more than 1.
			freezeTalents: false,

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.FrostTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_FROST_BIS_PRESET,
				],
			},
		});
	}
}
