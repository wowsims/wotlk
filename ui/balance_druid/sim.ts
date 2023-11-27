import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';

import * as OtherInputs from '../core/components/other_inputs.js';

import { AgilityBuffInput, FireDamageBuff, FrostDamageBuff, ShadowDamageBuff, StrengthBuffInput } from '../core/components/icon_inputs.js';
import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

// noinspection TypeScriptValidateTypes
export class BalanceDruidSimUI extends IndividualSimUI<Spec.SpecBalanceDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBalanceDruid>) {
		super(parentElem, player, {
			cssClass: 'balance-druid-sim-ui',
			cssScheme: 'druid',
			// List any known bugs / issues here, and they'll be shown on the site.
			knownIssues: [
			],
			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
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
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.DEFAULT_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.43,
					[Stat.StatSpirit]: 0.34,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellCrit]: 0.82,
					[Stat.StatSpellHaste]: 0.80,
					[Stat.StatMP5]: 0.00,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.BalanceTalents.data,
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
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DruidInputs.BalanceDruidRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
				AgilityBuffInput,
				StrengthBuffInput,
				FireDamageBuff,
				FrostDamageBuff,
				ShadowDamageBuff,
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.ReactionTime,
					OtherInputs.DistanceFromTarget,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.BalanceTalents,
				],
				rotations: [
					Presets.DEFAULT_APL,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.DEFAULT_PRESET,
				],
			},

			autoRotation: (): APLRotation => {
				return Presets.DEFAULT_APL.rotation.rotation!;
			},
		});
	}
}
