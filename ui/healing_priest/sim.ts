import {
	Class,
	Faction,
	PartyBuffs,
	Race,
	Spec,
	Stat,
} from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import {
	APLRotation,
} from '../core/proto/apl.js';

import * as HealingPriestInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHealingPriest, {
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
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_DISC,
			Presets.ROTATION_PRESET_HOLY,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.DISC_PRERAID_PRESET,
			Presets.DISC_P1_PRESET,
			Presets.DISC_P2_PRESET,
			Presets.DISC_P3_PRESET,
			Presets.DISC_P4_PRESET,
			Presets.HOLY_PRERAID_PRESET,
			Presets.HOLY_P1_PRESET,
			Presets.HOLY_P2_PRESET,
			Presets.HOLY_P3_PRESET,
			Presets.HOLY_P4_PRESET,
		],
	},

	autoRotation: (player: Player<Spec.SpecHealingPriest>): APLRotation => {
		const talentTree = player.getTalentTree();
		if (talentTree == 0) {
			return Presets.ROTATION_PRESET_DISC.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_HOLY.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecHealingPriest,
			tooltip: 'Discipline Priest',
			defaultName: 'Discipline',
			iconUrl: getSpecIcon(Class.ClassPriest, 0),

			talents: Presets.DiscTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDwarf,
				[Faction.Horde]: Race.RaceUndead,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.DISC_P1_PRESET.gear,
					2: Presets.DISC_P2_PRESET.gear,
					3: Presets.DISC_P3_PRESET.gear,
					4: Presets.DISC_P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.DISC_P1_PRESET.gear,
					2: Presets.DISC_P2_PRESET.gear,
					3: Presets.DISC_P3_PRESET.gear,
					4: Presets.DISC_P4_PRESET.gear,
				},
			},
		},
		{
			spec: Spec.SpecHealingPriest,
			tooltip: 'Holy Priest',
			defaultName: 'Holy',
			iconUrl: getSpecIcon(Class.ClassPriest, 1),

			talents: Presets.HolyTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDwarf,
				[Faction.Horde]: Race.RaceUndead,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.HOLY_P1_PRESET.gear,
					2: Presets.HOLY_P2_PRESET.gear,
					3: Presets.HOLY_P3_PRESET.gear,
					4: Presets.HOLY_P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.HOLY_P1_PRESET.gear,
					2: Presets.HOLY_P2_PRESET.gear,
					3: Presets.HOLY_P3_PRESET.gear,
					4: Presets.HOLY_P4_PRESET.gear,
				},
			},
		},
	],
});

export class HealingPriestSimUI extends IndividualSimUI<Spec.SpecHealingPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHealingPriest>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
