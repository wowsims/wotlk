import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { TotemsSection } from '../core/components/totem_inputs.js';
import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as ShamanInputs from './inputs.js';
import * as Presets from './presets.js';
import { FireElementalSection } from '../core/components/fire_elemental_inputs.js';

export class EnhancementShamanSimUI extends IndividualSimUI<Spec.SpecEnhancementShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecEnhancementShaman>) {
		super(parentElem, player, {
			cssClass: 'enhancement-shaman-sim-ui',
			cssScheme: 'shaman',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"Fire Elemental is in a alpha state",
				"Some things regarding weapon imbues need further testing and changes",
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHit,
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
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatIntellect,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatExpertise,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 1.48,
					[Stat.StatAgility]: 1.59,
					[Stat.StatStrength]: 1.1,
					[Stat.StatSpellPower]: 1.13,
					[Stat.StatSpellHit]: 0, //default EP assumes cap
					[Stat.StatSpellCrit]: 0.91,
					[Stat.StatSpellHaste]: 0.37,
					[Stat.StatAttackPower]: 1.0,
					[Stat.StatMeleeHit]: 1.38,
					[Stat.StatMeleeCrit]: 0.81,
					[Stat.StatMeleeHaste]: 1.61, //haste is complicated
					[Stat.StatArmorPenetration]: 0.48,
					[Stat.StatExpertise]: 0, //default EP assumes cap
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 5.21,
					[PseudoStat.PseudoStatOffHandDps]: 2.21,
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
				raidBuffs: Presets.DefaultRaidBuffs,
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					judgementsOfTheWise: true,
				}),
				debuffs: Presets.DefaultDebuffs,
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				ShamanInputs.ShamanShieldInput,
				ShamanInputs.Bloodlust,
				ShamanInputs.ShamanImbueMH,
				ShamanInputs.ShamanImbueOH,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ShamanInputs.EnhancementShamanRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.ReplenishmentBuff,
				IconInputs.MP5Buff,
				IconInputs.SpellHasteBuff,
				IconInputs.SpiritBuff,
			],
			excludeBuffDebuffInputs: [
				IconInputs.BleedDebuff,
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					ShamanInputs.SyncTypeInput,
					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			customSections: [
				TotemsSection,
				FireElementalSection
			],
			encounterPicker: {
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
					Presets.PreRaid_PRESET,
					Presets.P1_PRESET,
					Presets.P2_PRESET_FT,
					Presets.P2_PRESET_WF,
					Presets.P3_PRESET_ALLIANCE,
					Presets.P3_PRESET_HORDE,
				],
			},
		});
	}
}
