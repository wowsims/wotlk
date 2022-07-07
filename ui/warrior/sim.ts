import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Encounter } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';

import { Alchohol } from '/wotlk/core/proto/common.js';
import { BattleElixir } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { GuardianElixir } from '/wotlk/core/proto/common.js';
import { Conjured } from '/wotlk/core/proto/common.js';
import { Drums } from '/wotlk/core/proto/common.js';
import { PetFood } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';

import { Warrior, Warrior_Rotation as WarriorRotation, WarriorTalents as WarriorTalents, Warrior_Options as WarriorOptions } from '/wotlk/core/proto/warrior.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as WarriorInputs from './inputs.js';
import * as Presets from './presets.js';

export class WarriorSimUI extends IndividualSimUI<Spec.SpecWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarrior>) {
		super(parentElem, player, {
			cssClass: 'warrior-sim-ui',
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
				Stat.StatStamina,
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
				gear: Presets.P1_FURY_PRESET.gear,
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
				talents: Presets.FuryTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					bloodlust: 1,
					drums: Drums.DrumsOfBattle,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					windfuryTotem: TristateEffect.TristateEffectImproved,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					blessingOfSalvation: true,
					unleashedRage: true,
				}),
				debuffs: Debuffs.create({
					mangle: true,
					sunderArmor: true,
					curseOfWeakness: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					improvedSealOfTheCrusader: true,
					huntersMark: TristateEffect.TristateEffectImproved,
					exposeWeaknessUptime: 0.95,
					exposeWeaknessHunterAgility: 1200,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				WarriorInputs.ShoutPicker,
				WarriorInputs.Weakness,
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.GiftOfTheWild,
			],
			partyBuffInputs: [
				IconInputs.DrumsOfBattleBuff,
				IconInputs.Bloodlust,
				IconInputs.StrengthOfEarthTotem,

				IconInputs.WindfuryTotem,
				IconInputs.BattleShout,
				IconInputs.LeaderOfThePack,
				IconInputs.FerociousInspiration,
				IconInputs.TrueshotAura,
				IconInputs.SanctityAura,
				IconInputs.HeroicPresence,
				IconInputs.BraidedEterniumChain,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfMight,
				IconInputs.BlessingOfSalvation,
				IconInputs.UnleashedRage,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.BloodFrenzy,
				IconInputs.Mangle,
				IconInputs.ImprovedSealOfTheCrusader,
				IconInputs.HuntersMark,
				IconInputs.FaerieFire,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.CurseOfWeakness,
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
					WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
					WeaponImbue.WeaponImbueAdamantiteWeightstone,
					WeaponImbue.WeaponImbueElementalSharpeningStone,
					WeaponImbue.WeaponImbueRighteousWeaponCoating,
				],
				other: [
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfStrengthV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: WarriorInputs.WarriorRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					WarriorInputs.StartingRage,
					WarriorInputs.PrecastShout,
					WarriorInputs.PrecastShoutWithSapphire,
					WarriorInputs.PrecastShoutWithT2,
					OtherInputs.ExposeWeaknessUptime,
					OtherInputs.ExposeWeaknessHunterAgility,
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
				showExecuteProportion: true,
			},

			// If true, the talents on the talents tab will not be individually modifiable by the user.
			// Note that the use can still pick between preset talents, if there is more than 1.
			freezeTalents: false,

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.ArmsSlamTalents,
					Presets.ArmsDWTalents,
					Presets.FuryTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_FURY_PRESET,
					Presets.P2_FURY_PRESET,
					Presets.P3_FURY_PRESET,
					Presets.P4_FURY_PRESET,
					Presets.P5_FURY_PRESET,
					Presets.P1_ARMS_PRESET,
					Presets.P2_ARMS_PRESET,
					Presets.P3_ARMS_PRESET,
					Presets.P4_ARMS_PRESET,
					Presets.P5_ARMS_PRESET,
				],
			},
		});
	}
}
