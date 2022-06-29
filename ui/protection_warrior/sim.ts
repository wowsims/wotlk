import { RaidBuffs, StrengthOfEarthType } from '/tbc/core/proto/common.js';
import { PartyBuffs } from '/tbc/core/proto/common.js';
import { IndividualBuffs } from '/tbc/core/proto/common.js';
import { Class } from '/tbc/core/proto/common.js';
import { Consumes } from '/tbc/core/proto/common.js';
import { Debuffs } from '/tbc/core/proto/common.js';
import { Encounter } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { MobType } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js'
import { EquipmentSpec } from '/tbc/core/proto/common.js'
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { TypedEvent } from '/tbc/core/typed_event.js';

import { Alchohol } from '/tbc/core/proto/common.js';
import { BattleElixir } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { GuardianElixir } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { PetFood } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';

import { ProtectionWarrior, ProtectionWarrior_Rotation as ProtectionWarriorRotation, WarriorTalents as WarriorTalents, ProtectionWarrior_Options as ProtectionWarriorOptions } from '/tbc/core/proto/warrior.js';

import * as IconInputs from '/tbc/core/components/icon_inputs.js';
import * as OtherInputs from '/tbc/core/components/other_inputs.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

import * as ProtectionWarriorInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionWarriorSimUI extends IndividualSimUI<Spec.SpecProtectionWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionWarrior>) {
		super(parentElem, player, {
			cssClass: 'protection-warrior-sim-ui',
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
				Stat.StatArmor,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_BALANCED_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.05,
					[Stat.StatStamina]: 1,
					[Stat.StatStrength]: 0.33,
					[Stat.StatAgility]: 0.6,
					[Stat.StatAttackPower]: 0.06,
					[Stat.StatExpertise]: 0.67,
					[Stat.StatMeleeHit]: 0.67,
					[Stat.StatMeleeCrit]: 0.28,
					[Stat.StatMeleeHaste]: 0.21,
					[Stat.StatArmorPenetration]: 0.19,
					[Stat.StatBlock]: 0.35,
					[Stat.StatBlockValue]: 0.59,
					[Stat.StatDodge]: 0.7,
					[Stat.StatParry]: 0.58,
					[Stat.StatDefense]: 0.8,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.ImpaleProtTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					powerWordFortitude: TristateEffect.TristateEffectRegular,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					bloodlust: 1,
					drums: Drums.DrumsOfBattle,
					graceOfAirTotem: TristateEffect.TristateEffectImproved,
					strengthOfEarthTotem: StrengthOfEarthType.EnhancingTotems,
					windfuryTotemRank: 5,
					windfuryTotemIwt: 2,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
				}),
				debuffs: Debuffs.create({
					mangle: true,
					curseOfRecklessness: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					improvedSealOfTheCrusader: true,
					huntersMark: TristateEffect.TristateEffectImproved,
					exposeWeaknessUptime: 0.95,
					exposeWeaknessHunterAgility: 1200,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				ProtectionWarriorInputs.ShoutPicker,
				ProtectionWarriorInputs.ShieldWall,
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.PowerWordFortitude,
				IconInputs.ShadowProtection,
				IconInputs.GiftOfTheWild,
				IconInputs.Thorns,
			],
			partyBuffInputs: [
				IconInputs.DrumsOfBattleBuff,
				IconInputs.Bloodlust,
				IconInputs.StrengthOfEarthTotem,
				IconInputs.GraceOfAirTotem,
				IconInputs.WindfuryTotem,
				IconInputs.BattleShout,
				IconInputs.CommandingShout,
				IconInputs.LeaderOfThePack,
				IconInputs.FerociousInspiration,
				IconInputs.TrueshotAura,
				IconInputs.DevotionAura,
				IconInputs.RetributionAura,
				IconInputs.SanctityAura,
				IconInputs.DraeneiRacialMelee,
				IconInputs.BraidedEterniumChain,
				IconInputs.BloodPact,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfMight,
				IconInputs.BlessingOfSanctuary,
				IconInputs.BlessingOfSalvation,
				IconInputs.UnleashedRage,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.BloodFrenzy,
				IconInputs.Mangle,
				IconInputs.ImprovedSealOfTheCrusader,
				IconInputs.JudgementOfLight,
				IconInputs.JudgementOfWisdom,
				IconInputs.HuntersMark,
				IconInputs.FaerieFire,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.CurseOfRecklessness,
				IconInputs.GiftOfArthas,
				IconInputs.DemoralizingRoar,
				IconInputs.DemoralizingShout,
				IconInputs.Screech,
				IconInputs.ThunderClap,
				IconInputs.ShadowEmbrace,
				IconInputs.InsectSwarm,
				IconInputs.ScorpidSting,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.IronshieldPotion,
					Potions.HastePotion,
					Potions.MightyRagePotion,
					Potions.InsaneStrengthPotion,
				],
				conjured: [
					Conjured.ConjuredFlameCap,
					Conjured.ConjuredHealthstone,
				],
				flasks: [
					Flask.FlaskOfRelentlessAssault,
					Flask.FlaskOfFortification,
					Flask.FlaskOfChromaticWonder,
				],
				battleElixirs: [
					BattleElixir.ElixirOfDemonslaying,
					BattleElixir.ElixirOfMajorStrength,
					BattleElixir.ElixirOfMajorAgility,
					BattleElixir.ElixirOfTheMongoose,
					BattleElixir.ElixirOfMastery,
				],
				guardianElixirs: [
					GuardianElixir.ElixirOfMajorFortitude,
					GuardianElixir.ElixirOfMajorDefense,
					GuardianElixir.ElixirOfIronskin,
					GuardianElixir.GiftOfArthas,
				],
				food: [
					Food.FoodRoastedClefthoof,
					Food.FoodGrilledMudfish,
					Food.FoodSpicyHotTalbuk,
					Food.FoodRavagerDog,
					Food.FoodFishermansFeast,
				],
				alcohol: [
				],
				weaponImbues: [
					WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
					WeaponImbue.WeaponImbueAdamantiteWeightstone,
					WeaponImbue.WeaponImbueRighteousWeaponCoating,
				],
				other: [
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfStrengthV,
					IconInputs.ScrollOfProtectionV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionWarriorInputs.ProtectionWarriorRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HpPercentForDefensives,
					ProtectionWarriorInputs.StartingRage,
					ProtectionWarriorInputs.PrecastShout,
					ProtectionWarriorInputs.PrecastShoutWithSapphire,
					ProtectionWarriorInputs.PrecastShoutWithT2,
					OtherInputs.ExposeWeaknessUptime,
					OtherInputs.ExposeWeaknessHunterAgility,
					OtherInputs.InspirationUptime,
					OtherInputs.SnapshotImprovedStrengthOfEarthTotem,
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
					Presets.StandardTalents,
					Presets.ImpDemoTalents,
					Presets.ImpaleProtTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_BALANCED_PRESET,
					Presets.P2_BALANCED_PRESET,
					Presets.P3_BALANCED_PRESET,
					Presets.P4_BALANCED_PRESET,
					Presets.P5_BALANCED_PRESET,
				],
			},
		});
	}
}
