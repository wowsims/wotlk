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

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarlock, {
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
		Stat.StatMP5,
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatSpellPenetration,

		// Pet Related Stats
		Stat.StatStrength,
		Stat.StatStamina,
		Stat.StatAttackPower,
		Stat.StatArmorPenetration,
		Stat.StatAgility,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHit,
		Stat.StatMeleeHaste,
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
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatSpellPenetration,
		Stat.StatMP5,
		Stat.StatStamina,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.GearDestructionDefault.gear,

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
		talents: Presets.DestroTalents.data,
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

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.BlessingOfWisdom,
		BuffDebuffInputs.ManaSpringTotem,
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.JudgementOfWisdom,
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.ArmorBuff,
		BuffDebuffInputs.AttackPowerDebuff,
		BuffDebuffInputs.BleedDebuff,
		BuffDebuffInputs.DefensiveCooldownBuff,
		BuffDebuffInputs.MeleeAttackSpeedDebuff,
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.SpellWintersChillDebuff,
		...ConsumablesInputs.FROST_POWER_CONFIG,
	],
	petConsumeInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.DistanceFromTarget,
			OtherInputs.ChannelClipDelay,
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
			Presets.DestroTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.RotationDestructionDefault,
		],

		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.GearDestructionDefault,
		],
	},

	autoRotation: (player: Player<Spec.SpecWarlock>): APLRotation => {
		const talentTree = player.getTalentTree();
		if (talentTree == 0) {
			return Presets.RotationAfflictionDefault.rotation.rotation!;
		} else if (talentTree == 1) {
			return Presets.RotationDemonologyDefault.rotation.rotation!;
		} else {
			return Presets.RotationDestructionDefault.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Affliction DPS',
			defaultName: 'Affliction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 0),

			talents: Presets.DestroTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearAfflictionDefault.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearAfflictionDefault.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Demonology DPS',
			defaultName: 'Demonology',
			iconUrl: getSpecIcon(Class.ClassWarlock, 1),

			talents: Presets.DestroTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearDemonologyDefault.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearDemonologyDefault.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Destruction DPS',
			defaultName: 'Destruction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 2),

			talents: Presets.DestroTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearDestructionDefault.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearDestructionDefault.gear,
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
