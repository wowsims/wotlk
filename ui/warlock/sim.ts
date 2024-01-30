import {
	Class,
	Faction,
	ItemSlot,
	PartyBuffs,
	Race,
	Spec,
	Stat,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';

import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarlock, {
	cssClass: 'warlock-sim-ui',
	cssScheme: 'warlock',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
		"Drain Soul is currently disabled for APL rotations"
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
		gear: Presets.P3_AFFLICTION_HORDE_PRESET.gear,

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

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.MinorArmorDebuff,
		BuffDebuffInputs.PhysicalDamageDebuff,
		BuffDebuffInputs.MeleeHasteBuff,
		BuffDebuffInputs.MeleeCritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.AttackPowerBuff,
		BuffDebuffInputs.StrengthAndAgilityBuff,
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [
	],
	petConsumeInputs: [
		ConsumablesInputs.SpicedMammothTreats,
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarlockInputs.DetonateSeed,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
			OtherInputs.nibelungAverageCasts,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
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
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.APL_Affliction_Default,
			Presets.APL_Demo_Default,
			Presets.APL_Destro_Default,
		],

		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.SWP_BIS,
			Presets.PRERAID_AFFLICTION_PRESET,
			Presets.P1_AFFLICTION_PRESET,
			Presets.P2_AFFLICTION_PRESET,
			Presets.P3_AFFLICTION_ALLIANCE_PRESET,
			Presets.P3_AFFLICTION_HORDE_PRESET,
			Presets.P4_AFFLICTION_PRESET,
			Presets.PRERAID_DEMODESTRO_PRESET,
			Presets.P1_DEMODESTRO_PRESET,
			Presets.P2_DEMODESTRO_PRESET,
			Presets.P3_DEMO_ALLIANCE_PRESET,
			Presets.P3_DEMO_HORDE_PRESET,
			Presets.P4_DEMO_PRESET,
			Presets.P3_DESTRO_ALLIANCE_PRESET,
			Presets.P3_DESTRO_HORDE_PRESET,
			Presets.P4_DESTRO_PRESET,
		],
	},

	autoRotation: (player: Player<Spec.SpecWarlock>): APLRotation => {
		const talentTree = player.getTalentTree();
		if (talentTree == 0) {
			return Presets.APL_Affliction_Default.rotation.rotation!;
		} else if (talentTree == 1) {
			return Presets.APL_Demo_Default.rotation.rotation!;
		} else {
			return Presets.APL_Destro_Default.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Affliction Warlock',
			defaultName: 'Affliction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 0),

			talents: Presets.AfflictionTalents.data,
			specOptions: Presets.AfflictionOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_AFFLICTION_PRESET.gear,
					2: Presets.P2_AFFLICTION_PRESET.gear,
					3: Presets.P3_AFFLICTION_ALLIANCE_PRESET.gear,
					4: Presets.P4_AFFLICTION_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_AFFLICTION_PRESET.gear,
					2: Presets.P2_AFFLICTION_PRESET.gear,
					3: Presets.P3_AFFLICTION_HORDE_PRESET.gear,
					4: Presets.P4_AFFLICTION_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Demonology Warlock',
			defaultName: 'Demonology',
			iconUrl: getSpecIcon(Class.ClassWarlock, 1),

			talents: Presets.DemonologyTalents.data,
			specOptions: Presets.DemonologyOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_DEMODESTRO_PRESET.gear,
					2: Presets.P2_DEMODESTRO_PRESET.gear,
					3: Presets.P3_DEMO_ALLIANCE_PRESET.gear,
					4: Presets.P4_DEMO_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_DEMODESTRO_PRESET.gear,
					2: Presets.P2_DEMODESTRO_PRESET.gear,
					3: Presets.P3_DEMO_HORDE_PRESET.gear,
					4: Presets.P4_DEMO_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Destruction Warlock',
			defaultName: 'Destruction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 2),

			talents: Presets.DestructionTalents.data,
			specOptions: Presets.DestructionOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_DEMODESTRO_PRESET.gear,
					2: Presets.P2_DEMODESTRO_PRESET.gear,
					3: Presets.P3_DESTRO_ALLIANCE_PRESET.gear,
					4: Presets.P4_DESTRO_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_DEMODESTRO_PRESET.gear,
					2: Presets.P2_DEMODESTRO_PRESET.gear,
					3: Presets.P3_DESTRO_HORDE_PRESET.gear,
					4: Presets.P4_DESTRO_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
