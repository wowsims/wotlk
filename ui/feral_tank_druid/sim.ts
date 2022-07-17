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
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { TypedEvent } from '/wotlk/core/typed_event.js';

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
	DruidTalents as DruidTalents,
	FeralTankDruid,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Options as DruidOptions
} from '/wotlk/core/proto/druid.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

export class FeralTankDruidSimUI extends IndividualSimUI<Spec.SpecFeralTankDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralTankDruid>) {
		super(parentElem, player, {
			cssClass: 'feral-tank-druid-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
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
				Stat.StatDodge,
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
				Stat.StatDodge,
				Stat.StatResilience,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P4_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.59,
					[Stat.StatStamina]: 3.05,
					[Stat.StatStrength]: 2.266,
					[Stat.StatAgility]: 4.6,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 7.3,
					[Stat.StatMeleeHit]: 3.5,
					[Stat.StatMeleeCrit]: 1.0,
					[Stat.StatMeleeHaste]: 1.6,
					[Stat.StatArmorPenetration]: 0.34,
					[Stat.StatDefense]: 2.2,
					[Stat.StatDodge]: 1.7,
					[Stat.StatResilience]: 1.7,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.StandardTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					powerWordFortitude: TristateEffect.TristateEffectRegular,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					battleShout: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
				}),
				partyBuffs: PartyBuffs.create({
					braidedEterniumChain: true,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					bloodFrenzy: true,
					exposeArmor: TristateEffect.TristateEffectImproved,
					faerieFire: TristateEffect.TristateEffectImproved,
					sunderArmor: true,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
					thunderClap: TristateEffect.TristateEffectImproved,
					demoralizingShout: TristateEffect.TristateEffectImproved,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.PowerWordFortitude,
				IconInputs.ShadowProtection,
				IconInputs.GiftOfTheWild,
				IconInputs.Thorns,
				IconInputs.Bloodlust,
				IconInputs.StrengthOfEarthTotem,
				IconInputs.BattleShout,
				IconInputs.CommandingShout,
				IconInputs.DevotionAura,
				IconInputs.RetributionAura,
				IconInputs.TrueshotAura,
				IconInputs.BloodPact,
				IconInputs.UnleashedRage,
			],
			partyBuffInputs: [
				IconInputs.BraidedEterniumChain,
				IconInputs.HeroicPresence,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfMight,
				IconInputs.BlessingOfSanctuary,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.BloodFrenzy,
				IconInputs.Mangle,
				IconInputs.JudgementOfLight,
				IconInputs.JudgementOfWisdom,
				IconInputs.FaerieFire,
				IconInputs.SunderArmor,
				IconInputs.ExposeArmor,
				IconInputs.CurseOfWeakness,
				IconInputs.GiftOfArthas,
				IconInputs.DemoralizingRoar,
				IconInputs.DemoralizingShout,
				IconInputs.Screech,
				IconInputs.ThunderClap,
				IconInputs.InsectSwarm,
				IconInputs.ScorpidSting,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.IronshieldPotion,
					Potions.HastePotion,
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
				],
				other: [
					IconInputs.ScrollOfAgilityV,
					IconInputs.ScrollOfStrengthV,
					IconInputs.ScrollOfProtectionV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DruidInputs.FeralTankDruidRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HpPercentForDefensives,
					DruidInputs.StartingRage,
					OtherInputs.PrepopPotion,
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
					Presets.StandardTalents,
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
