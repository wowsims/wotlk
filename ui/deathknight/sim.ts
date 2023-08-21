import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';

import { Deathknight, Deathknight_Rotation as DeathKnightRotation, DeathknightTalents as DeathKnightTalents, Deathknight_Options as DeathKnightOptions } from '../core/proto/deathknight.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as DeathKnightInputs from './inputs.js';
import * as Presets from './presets.js';

export class DeathknightSimUI extends IndividualSimUI<Spec.SpecDeathknight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDeathknight>) {
		super(parentElem, player, {
			cssClass: 'deathknight-sim-ui',
			cssScheme: 'death-knight',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatArmor,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
			],
			epPseudoStats: [
				PseudoStat.PseudoStatMainHandDps,
				PseudoStat.PseudoStatOffHandDps,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.P2_UNHOLY_DW_BIS_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStrength]: 3.22,
					[Stat.StatAgility]: 0.62,
					[Stat.StatArmor]: 0.01,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 1.13,
					[Stat.StatMeleeHaste]: 1.85,
					[Stat.StatMeleeHit]: 1.92,
					[Stat.StatMeleeCrit]: 0.76,
					[Stat.StatArmorPenetration]: 0.77,
					[Stat.StatSpellHit]: 0.80,
					[Stat.StatSpellCrit]: 0.34,
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 3.10,
					[PseudoStat.PseudoStatOffHandDps]: 1.79,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultUnholyRotation,
				// Default talents.
				talents: Presets.UnholyDualWieldTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultUnholyOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					swiftRetribution: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					abominationsMight: true,
					leaderOfThePack: TristateEffect.TristateEffectRegular,
					sanctifiedRetribution: true,
					bloodlust: true,
					devotionAura: TristateEffect.TristateEffectImproved,
					stoneskinTotem: TristateEffect.TristateEffectImproved,
					moonkinAura: TristateEffect.TristateEffectRegular,
					wrathOfAirTotem: true,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					heroicPresence: false,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					bloodFrenzy: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					sunderArmor: true,
					ebonPlaguebringer: true,
					mangle: true,
					heartOfTheCrusader: true,
					shadowMastery: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DeathKnightInputs.DeathKnightRotationConfig,
			petConsumeInputs: [
				IconInputs.SpicedMammothTreats,
			],
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.SpellDamageDebuff,
				IconInputs.StaminaBuff,
			],
			excludeBuffDebuffInputs: [
				IconInputs.AttackPowerDebuff,
				IconInputs.DamageReductionPercentBuff,
				IconInputs.MeleeAttackSpeedDebuff,
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					DeathKnightInputs.DiseaseDowntime,
					DeathKnightInputs.DrwPestiApply,
					DeathKnightInputs.SelfUnholyFrenzy,
					DeathKnightInputs.StartingRunicPower,
					DeathKnightInputs.PetUptime,

					DeathKnightInputs.PrecastGhoulFrenzy,
					DeathKnightInputs.PrecastHornOfWinter,

					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.BloodTalents,
					Presets.FrostTalents,
					Presets.FrostUnholyTalents,
					Presets.UnholyDualWieldTalents,
					Presets.UnholyDualWieldSSTalents,
					Presets.Unholy2HTalents,
					Presets.UnholyAoeTalents,
				],
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.BLOOD_ROTATION_PRESET_LEGACY_DEFAULT,
					Presets.FROST_ROTATION_PRESET_LEGACY_DEFAULT,
					Presets.UNHOLY_DW_ROTATION_PRESET_LEGACY_DEFAULT,
					Presets.BLOOD_PESTI_ROTATION_PRESET_DEFAULT,
					Presets.FROST_BL_PESTI_ROTATION_PRESET_DEFAULT,
					Presets.FROST_UH_PESTI_ROTATION_PRESET_DEFAULT,
					Presets.UNHOLY_DW_ROTATION_PRESET_DEFAULT,
					Presets.UNHOLY_2H_ROTATION_PRESET_DEFAULT,
					Presets.UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_BLOOD_BIS_PRESET,
					Presets.P2_BLOOD_BIS_PRESET,
					Presets.P3_BLOOD_BIS_PRESET,
					Presets.P1_FROST_BIS_PRESET,
					Presets.P2_FROST_BIS_PRESET,
					Presets.P3_FROST_BIS_PRESET,
					Presets.P1_UNHOLY_DW_BIS_PRESET,
					Presets.P2_UNHOLY_DW_BIS_PRESET,
					Presets.P3_UNHOLY_DW_BIS_PRESET,
					// Not needed anymore just filling ui Space
					// Disabled on purpose
					//Presets.P1_FROSTSUBUNH_BIS_PRESET,
					//Presets.P1_FROST_PRE_BIS_PRESET,
					//Presets.P1_UNHOLY_DW_PRERAID_PRESET,
					//Presets.P1_UNHOLY_2H_PRERAID_PRESET,
					//Presets.P1_UNHOLY_2H_BIS_PRESET,
				],
			},
		});
	}
}
