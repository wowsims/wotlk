import { PartyBuffs } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';

import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

import * as SmitePriestInputs from './inputs.js';
import * as Presets from './presets.js';

export class SmitePriestSimUI extends IndividualSimUI<Spec.SpecSmitePriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecSmitePriest>) {
		super(parentElem, player, {
			cssClass: 'smite-priest-sim-ui',
			cssScheme: 'priest',
			// List any known bugs / issues here and they'll be shown on the site.
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
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			modifyDisplayStats: (player: Player<Spec.SpecSmitePriest>) => {
				let stats = new Stats();
				stats = stats.addStat(Stat.StatSpellHit, player.getTalents().shadowFocus * 1 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.38,
					[Stat.StatSpirit]: 0.38,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellHit]: 1.65,
					[Stat.StatSpellCrit]: 0.32,
					[Stat.StatSpellHaste]: 0.78,
					[Stat.StatMP5]: 0.35,
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
				partyBuffs: PartyBuffs.create({}),
				individualBuffs: Presets.DefaultIndividualBuffs,
				debuffs: Presets.DefaultDebuffs,
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				SmitePriestInputs.SelfPowerInfusion,
				SmitePriestInputs.InnerFire,
				SmitePriestInputs.Shadowfiend,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: SmitePriestInputs.SmitePriestRotationConfig,
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
					Presets.StandardTalents,
				],
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.ROTATION_PRESET_LEGACY_DEFAULT,
					Presets.ROTATION_PRESET_APL,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PRERAID_PRESET,
					Presets.P1_PRESET,
				],
			},
		});
	}
}
