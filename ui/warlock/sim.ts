import {
	PartyBuffs,
	Spec,
	Stat,
} from '../core/proto/common.js';

import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, {
			cssClass: 'warlock-sim-ui',
			cssScheme: 'warlock',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"Most abilities and pets are work in progress"
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
				Stat.StatMana,
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
				gear: Presets.DEFAULT_GEAR.gear,

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
				// Default talents.
				talents: Presets.DefaultTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,

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
			rotationInputs: WarlockInputs.WarlockRotationConfig,

			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.MP5Buff,
				IconInputs.StaminaBuff,
			],
			excludeBuffDebuffInputs: [
			],
			petConsumeInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.DistanceFromTarget,
					OtherInputs.TankAssignment,
					OtherInputs.ChannelClipDelay,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.DefaultTalents,
				],
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.ROTATION_PRESET_DEFAULT,
				],

				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.DEFAULT_GEAR,
				],
			},
		});
	}
}
