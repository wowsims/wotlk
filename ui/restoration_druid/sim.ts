import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

export class RestorationDruidSimUI extends IndividualSimUI<Spec.SpecRestorationDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRestorationDruid>) {
		super(parentElem, player, {
			cssClass: 'restoration-druid-sim-ui',
			cssScheme: 'druid',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatMana,
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.38,
					[Stat.StatSpirit]: 0.34,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellCrit]: 0.69,
					[Stat.StatSpellHaste]: 0.77,
					[Stat.StatMP5]: 0.00,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.CelestialFocusTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: Presets.DefaultRaidBuffs,

				partyBuffs: Presets.DefaultPartyBuffs,

				individualBuffs: Presets.DefaultIndividualBuffs,

				debuffs: Presets.DefaultDebuffs,

				other: Presets.OtherDefaults,
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				DruidInputs.SelfInnervate,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DruidInputs.RestorationDruidRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
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
					Presets.CelestialFocusTalents,
					Presets.ThiccRestoTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PRERAID_PRESET,
					Presets.P1_PRESET,
					Presets.P2_PRESET,
				],
			},

			autoRotation: (_player: Player<Spec.SpecRestorationDruid>): APLRotation => {
				return APLRotation.create();
			},
		});
	}
}
