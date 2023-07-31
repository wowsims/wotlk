import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { UnitReference } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';

import {
	HealingPriest,
	HealingPriest_Rotation as Rotation,
	HealingPriest_Options as Options,
} from '../core/proto/priest.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as HealingPriestInputs from './inputs.js';
import * as Presets from './presets.js';

export class HealingPriestSimUI extends IndividualSimUI<Spec.SpecHealingPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHealingPriest>) {
		super(parentElem, player, {
			cssClass: 'healing-priest-sim-ui',
			cssScheme: 'priest',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				'Talents that apply to, "friendly targets at or below 50% health" are not implemented.',
				'Prayer of Mending always bounces the maximum number of times.',
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
				gear: Presets.DISC_P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 2.73,
					[Stat.StatSpirit]: 1.63,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellCrit]: 0.75,
					[Stat.StatSpellHaste]: 0.28,
					[Stat.StatMP5]: 2.05,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DiscDefaultRotation,
				// Default talents.
				talents: Presets.DiscTalents.data,
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
				HealingPriestInputs.SelfPowerInfusion,
				HealingPriestInputs.InnerFire,
				HealingPriestInputs.Shadowfiend,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: HealingPriestInputs.HealingPriestRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					HealingPriestInputs.RapturesPerMinute,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.DiscTalents,
					Presets.HolyTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.DISC_PRERAID_PRESET,
					Presets.DISC_P1_PRESET,
					Presets.DISC_P2_PRESET,
					Presets.HOLY_PRERAID_PRESET,
					Presets.HOLY_P1_PRESET,
					Presets.HOLY_P2_PRESET,
				],
			},
		});
	}
}
