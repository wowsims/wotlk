import { RaidBuffs, StrengthOfEarthType } from '/tbc/core/proto/common.js';
import { PartyBuffs } from '/tbc/core/proto/common.js';
import { IndividualBuffs } from '/tbc/core/proto/common.js';
import { Debuffs } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js'
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { Player } from '/tbc/core/player.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

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

import * as IconInputs from '/tbc/core/components/icon_inputs.js';
import * as OtherInputs from '/tbc/core/components/other_inputs.js';

import * as ProtectionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, {
			cssClass: 'protection-paladin-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatIntellect,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatExpertise,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHit,
				Stat.StatArmor,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatIntellect,
				Stat.StatMP5,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatExpertise,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatHolySpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.P4_PRESET.gear,
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
				talents: Presets.SanctityTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					powerWordFortitude: TristateEffect.TristateEffectRegular,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					bloodlust: 1,
					drums: Drums.DrumsOfBattle,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					wrathOfAirTotem: TristateEffect.TristateEffectRegular,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfSanctuary: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
				}),
				debuffs: Debuffs.create({
					improvedSealOfTheCrusader: true,
					misery: true,
					bloodFrenzy: true,
					exposeArmor: TristateEffect.TristateEffectImproved,
					sunderArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					curseOfRecklessness: true,
					huntersMark: TristateEffect.TristateEffectImproved,
					exposeWeaknessUptime: 0.95,
					exposeWeaknessHunterAgility: 1200,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.ArcaneBrilliance,
				IconInputs.PowerWordFortitude,
				IconInputs.ShadowProtection,
				IconInputs.DivineSpirit,
				IconInputs.GiftOfTheWild,
				IconInputs.Thorns,
			],
			partyBuffInputs: [
				IconInputs.DrumsOfBattleBuff,
				IconInputs.Bloodlust,
				IconInputs.ManaSpringTotem,
				IconInputs.WrathOfAirTotem,
				IconInputs.TotemOfWrath,
				IconInputs.WindfuryTotem,
				IconInputs.StrengthOfEarthTotem,
				IconInputs.GraceOfAirTotem,
				IconInputs.BattleShout,
				IconInputs.CommandingShout,
				IconInputs.DraeneiRacialCaster,
				IconInputs.DraeneiRacialMelee,
				IconInputs.LeaderOfThePack,
				IconInputs.FerociousInspiration,
				IconInputs.TrueshotAura,
				IconInputs.DevotionAura,
				IconInputs.RetributionAura,
				IconInputs.SanctityAura,
				IconInputs.BloodPact,
				IconInputs.BraidedEterniumChain,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfWisdom,
				IconInputs.BlessingOfMight,
				IconInputs.BlessingOfSanctuary,
				IconInputs.BlessingOfSalvation,
				IconInputs.UnleashedRage,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.JudgementOfWisdom,
				IconInputs.JudgementOfLight,
				IconInputs.ImprovedSealOfTheCrusader,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.BloodFrenzy,
				IconInputs.HuntersMark,
				IconInputs.FaerieFire,
				IconInputs.CurseOfRecklessness,
				IconInputs.Misery,
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
					Potions.SuperManaPotion,
					Potions.DestructionPotion,
					Potions.HastePotion,
				],
				conjured: [
					Conjured.ConjuredDarkRune,
					Conjured.ConjuredFlameCap,
					Conjured.ConjuredHealthstone,
				],
				flasks: [
					Flask.FlaskOfFortification,
					Flask.FlaskOfBlindingLight,
					Flask.FlaskOfRelentlessAssault,
					Flask.FlaskOfChromaticWonder,
				],
				battleElixirs: [
					BattleElixir.GreaterArcaneElixir,
					BattleElixir.ElixirOfMastery,
					BattleElixir.ElixirOfDemonslaying,
					BattleElixir.ElixirOfMajorAgility,
					BattleElixir.ElixirOfTheMongoose,
				],
				guardianElixirs: [
					GuardianElixir.ElixirOfMajorFortitude,
					GuardianElixir.ElixirOfMajorDefense,
					GuardianElixir.ElixirOfIronskin,
					GuardianElixir.GiftOfArthas,
					GuardianElixir.ElixirOfDraenicWisdom,
					GuardianElixir.ElixirOfMajorMageblood,
				],
				food: [
					Food.FoodRoastedClefthoof,
					Food.FoodGrilledMudfish,
					Food.FoodSpicyHotTalbuk,
					Food.FoodBlackenedBasilisk,
					Food.FoodFishermansFeast,
				],
				alcohol: [
					Alchohol.AlchoholKreegsStoutBeatdown,
				],
				weaponImbues: [
					WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
					WeaponImbue.WeaponImbueAdamantiteWeightstone,
					WeaponImbue.WeaponImbueBrilliantWizardOil,
					WeaponImbue.WeaponImbueSuperiorWizardOil,
					WeaponImbue.WeaponImbueRighteousWeaponCoating,
				],
				other: [
					IconInputs.ScrollOfStrengthV,
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfProtectionV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionPaladinInputs.ProtectionPaladinRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HpPercentForDefensives,
					ProtectionPaladinInputs.AuraSelection,
					ProtectionPaladinInputs.UseAvengingWrath,
					OtherInputs.ExposeWeaknessUptime,
					OtherInputs.ExposeWeaknessHunterAgility,
					OtherInputs.InspirationUptime,
					OtherInputs.SnapshotImprovedStrengthOfEarthTotem,
					OtherInputs.SnapshotImprovedWrathOfAirTotem,
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
					Presets.ArdentDefenderTalents,
					Presets.AvengersShieldTalents,
					Presets.SanctityTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_PRESET,
					Presets.P2_PRESET,
					Presets.P3_PRESET,
					Presets.P4_PRESET,
					Presets.P5_PRESET,
				],
			},
		});
	}
}
