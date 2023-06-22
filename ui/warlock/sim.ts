import {
	RaidBuffs,
	PartyBuffs,
	IndividualBuffs,
	Debuffs,
	Spec,
	Stat,
	TristateEffect,
	Race,
} from '../core/proto/common.js';

import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { TypedEvent } from '../core/typed_event.js';

import {
	Warlock,
	Warlock_Rotation as WarlockRotation,
	WarlockTalents as WarlockTalents,
	Warlock_Options as WarlockOptions,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options_WeaponImbue as WeaponImbue,
} from '../core/proto/warlock.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, {
			cssClass: 'warlock-sim-ui',
			cssScheme: 'warlock',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"Several secondary spells need to be implemented.",
				"Rotations will be optimized.",
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatStamina,
			],
			// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
				Stat.StatStamina,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P3_Preset_Affliction_Horde.gear,

				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.18,
					[Stat.StatSpirit]: 0.54,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellHit]: 0.93,
					[Stat.StatSpellCrit]: 0.53,
					[Stat.StatSpellHaste]: 0.81,
					[Stat.StatStamina]: 0.01,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,

				// Default rotation settings.
				rotation: Presets.AfflictionRotation,
				// Default talents.
				talents: Presets.AfflictionTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.AfflictionOptions,

				// Default buffs and debuffs settings.
				raidBuffs: Presets.DefaultRaidBuffs,

				partyBuffs: PartyBuffs.create({}),

				individualBuffs: Presets.DefaultIndividualBuffs,

				debuffs: Presets.DefaultDebuffs,

				other: Presets.OtherDefaults,
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				WarlockInputs.PetInput,
				WarlockInputs.ArmorInput,
				WarlockInputs.WeaponImbueInput,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationIconInputs: [
				WarlockInputs.PrimarySpellInput,
				WarlockInputs.CorruptionSpell,
				WarlockInputs.SecondaryDotInput,
				WarlockInputs.SpecSpellInput,
			],
			rotationInputs: WarlockInputs.WarlockRotationConfig,

			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.ReplenishmentBuff,
				IconInputs.MajorArmorDebuff,
				IconInputs.MinorArmorDebuff,
				IconInputs.PhysicalDamageDebuff,
				IconInputs.MeleeHasteBuff,
				IconInputs.MeleeCritBuff,
				IconInputs.MP5Buff,
				IconInputs.AttackPowerPercentBuff,
				IconInputs.AttackPowerBuff,
				IconInputs.StrengthAndAgilityBuff,
				IconInputs.StaminaBuff,
			],
			excludeBuffDebuffInputs: [
			],
			petConsumeInputs: [
				IconInputs.SpicedMammothTreats,
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.DistanceFromTarget,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.AfflictionTalents,
					Presets.DemonologyTalents,
					Presets.DestructionTalents,
				],
				//Preset gear configurations that the user can quickly select.
				gear: [
					// Presets.Naked,
					Presets.SWP_BIS,
					Presets.P1_PreBiS_11,
					Presets.P1_PreBiS_14,
					Presets.P1_Preset_Affliction,
					Presets.P1_Preset_Demo_Destro,
					Presets.P2_Preset_Affliction,
					Presets.P2_Preset_Demo_Destro,
					Presets.P3_Preset_Affliction_Horde,
					Presets.P3_Preset_Affliction_Alliance,
					Presets.P3_Preset_Demo_Horde,
					Presets.P3_Preset_Demo_Alliance,
					Presets.P3_Preset_Destro_Horde,
					Presets.P3_Preset_Destro_Alliance,
				],
			},
		});
	}
}
