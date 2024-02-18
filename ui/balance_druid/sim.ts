import {
	Class,
	Faction,
	Race,
	Spec,
	Stat,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { Player } from '../core/player.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBalanceDruid, {
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
		gear: Presets.P3_PRESET_HORDE.gear,
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
		// Default talents.
		talents: Presets.Phase3Talents.data,
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
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.MeleeHasteBuff,
		BuffDebuffInputs.MeleeCritBuff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.AttackPowerBuff,
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.MinorArmorDebuff,
		BuffDebuffInputs.PhysicalDamageDebuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DruidInputs.OkfUptime,
			OtherInputs.TankAssignment,
			OtherInputs.ReactionTime,
			OtherInputs.DistanceFromTarget,
			OtherInputs.nibelungAverageCasts,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.Phase1Talents,
			Presets.Phase2Talents,
			Presets.Phase3Talents,
			Presets.Phase4Talents,
		],
		rotations: [
			Presets.ROTATION_PRESET_P3_APL,
			Presets.ROTATION_PRESET_P4_FOCUS_APL,
			Presets.ROTATION_PRESET_P4_STARFIRE_APL,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_PRESET,
			Presets.P2_PRESET,
			Presets.P3_PRESET_HORDE,
			Presets.P3_PRESET_ALLI,
			Presets.P4_PRESET_HORDE,
			Presets.P4_PRESET_ALLI,
		],
	},

	autoRotation: (_player: Player<Spec.SpecBalanceDruid>): APLRotation => {
		return Presets.ROTATION_PRESET_P3_APL.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBalanceDruid,
			tooltip: specNames[Spec.SpecBalanceDruid],
			defaultName: 'Balance',
			iconUrl: getSpecIcon(Class.ClassDruid, 0),

			talents: Presets.Phase2Talents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET_ALLI.gear,
					4: Presets.P4_PRESET_ALLI.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET_HORDE.gear,
					4: Presets.P4_PRESET_HORDE.gear,
				},
			},
		},
	],
});

export class BalanceDruidSimUI extends IndividualSimUI<Spec.SpecBalanceDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBalanceDruid>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
