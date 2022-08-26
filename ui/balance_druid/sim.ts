import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';

import { BalanceDruid, BalanceDruid_Rotation as BalanceDruidRotation, DruidTalents as DruidTalents, BalanceDruid_Options as BalanceDruidOptions } from '../core/proto/druid.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

export class BalanceDruidSimUI extends IndividualSimUI<Spec.SpecBalanceDruid> {
    constructor(parentElem: HTMLElement, player: Player<Spec.SpecBalanceDruid>) {
        super(parentElem, player, {
            cssClass: 'balance-druid-sim-ui',
            // List any known bugs / issues here and they'll be shown on the site.
            knownIssues: [
            ],

            // All stats for which EP should be calculated.
            epStats: [
                Stat.StatIntellect,
                Stat.StatSpirit,
                Stat.StatSpellPower,
                Stat.StatArcaneSpellPower,
                Stat.StatNatureSpellPower,
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
                Stat.StatArcaneSpellPower,
                Stat.StatNatureSpellPower,
                Stat.StatSpellHit,
                Stat.StatSpellCrit,
                Stat.StatSpellHaste,
                Stat.StatMP5,
            ],

			defaults: {
				// Default equipped gear.
                gear: Presets.PRE_RAID_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.52,
					[Stat.StatSpirit]: 0.3,
					[Stat.StatSpellPower]: 1,
					[Stat.StatArcaneSpellPower]: 0.45,
					[Stat.StatNatureSpellPower]: 0.50,
					[Stat.StatSpellCrit]: 0.61,
					[Stat.StatSpellHaste]: 0.67,
					[Stat.StatMP5]: 0.00,
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

				partyBuffs: Presets.DefaultPartyBuffs,

				individualBuffs: Presets.DefaultIndividualBuffs,

				debuffs: Presets.DefaultDebuffs,
			},

            // IconInputs to include in the 'Player' section on the settings tab.
            playerIconInputs: [
                DruidInputs.SelfInnervate,
            ],
            // Inputs to include in the 'Rotation' section on the settings tab.
            rotationInputs: DruidInputs.BalanceDruidRotationConfig,
            // Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
            includeBuffDebuffInputs: [
                IconInputs.MeleeHasteBuff,
                IconInputs.MeleeCritBuff,
                IconInputs.AttackPowerPercentBuff,
                IconInputs.AttackPowerBuff,
                IconInputs.MajorArmorDebuff,
                IconInputs.MinorArmorDebuff,
                IconInputs.PhysicalDamageDebuff,
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
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_PRESET,
                    Presets.PRE_RAID_PRESET,
				],
			},
		});
	}
}
