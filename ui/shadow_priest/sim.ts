import {
	PartyBuffs,
	Spec,
	Stat,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';

import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

import * as ShadowPriestInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecShadowPriest, {
	cssClass: 'shadow-priest-sim-ui',
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
	modifyDisplayStats: (player: Player<Spec.SpecShadowPriest>) => {
		let stats = new Stats();
		stats = stats.addStat(Stat.StatSpellHit, player.getTalents().shadowFocus * 1 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.11,
			[Stat.StatSpirit]: 0.47,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellHit]: 0.87,
			[Stat.StatSpellCrit]: 0.74,
			[Stat.StatSpellHaste]: 1.65,
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

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		ShadowPriestInputs.ArmorInput,
	],
	rotationIconInputs: [
		ShadowPriestInputs.MindBlastInput,
		ShadowPriestInputs.ShadowWordDeathInput,
		ShadowPriestInputs.ShadowfiendInput,
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: ShadowPriestInputs.ShadowPriestRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		IconInputs.ReplenishmentBuff,
		IconInputs.MeleeHasteBuff,
		IconInputs.MeleeCritBuff,
		IconInputs.MP5Buff,
		IconInputs.AttackPowerPercentBuff,
		IconInputs.AttackPowerBuff,
		IconInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
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
			Presets.StandardTalents,
			Presets.EnlightenmentTalents,
		],
		rotations: [
			Presets.ROTATION_PRESET_DEFAULT,
			Presets.ROTATION_PRESET_AOE24,
			Presets.ROTATION_PRESET_AOE4PLUS,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_PRESET,
			Presets.P2_PRESET,
			Presets.P3_PRESET,
			Presets.P4_PRESET,
		],
	},

	autoRotation: (player: Player<Spec.SpecShadowPriest>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 4) {
			return Presets.ROTATION_PRESET_AOE4PLUS.rotation.rotation!;
		} else if (numTargets > 1) {
			return Presets.ROTATION_PRESET_AOE24.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},
});

export class ShadowPriestSimUI extends IndividualSimUI<Spec.SpecShadowPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecShadowPriest>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
