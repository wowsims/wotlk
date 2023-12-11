import {
	Class,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	Race,
	RaidBuffs,
	Spec,
	Stat, PseudoStat,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as OtherInputs from '../core/components/other_inputs.js';
import * as WarriorInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarrior, {
	cssClass: 'warrior-sim-ui',
	cssScheme: 'warrior',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatArmor,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatArmor,
	],
	modifyDisplayStats: (_: Player<Spec.SpecWarrior>) => {
		let stats = new Stats();

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.GearArmsDefault.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 2.72,
			[Stat.StatAgility]: 1.82,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 2.55,
			[Stat.StatMeleeHit]: 0.79,
			[Stat.StatMeleeCrit]: 2.12,
			[Stat.StatMeleeHaste]: 1.72,
			[Stat.StatArmorPenetration]: 2.17,
			[Stat.StatArmor]: 0.03,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 6.29,
			[PseudoStat.PseudoStatOffHandDps]: 3.58,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.Talent25.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
			leaderOfThePack: true,
			devotionAura: TristateEffect.TristateEffectImproved,
			stoneskinTotem: TristateEffect.TristateEffectImproved,
		}),
		partyBuffs: PartyBuffs.create({
			heroicPresence: false,
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			sunderArmor: true,
			curseOfWeakness: TristateEffect.TristateEffectRegular,
			faerieFire: TristateEffect.TristateEffectImproved,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		WarriorInputs.ShoutPicker,
		WarriorInputs.Recklessness,
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: WarriorInputs.WarriorRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarriorInputs.StartingRage,
			WarriorInputs.StanceSnapshot,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.Talent25
		],
		// Preset rotations that the user can quickly select.
		rotations: [
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
		],
	},

	autoRotation: (player: Player<Spec.SpecWarrior>): APLRotation => {
		const talentTree = player.getTalentTree();
		if (talentTree == 0) {
			return Presets.RotationArmsDefault.rotation.rotation!;
		} else {
			return Presets.RotationFuryDefault.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarrior,
			tooltip: 'Arms Warrior',
			defaultName: 'Arms',
			iconUrl: getSpecIcon(Class.ClassWarrior, 0),

			talents: Presets.Talent25.data,
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
					1: Presets.GearArmsDefault.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearArmsDefault.gear,
				},
			},
		},
		{
			spec: Spec.SpecWarrior,
			tooltip: 'Fury Warrior',
			defaultName: 'Fury',
			iconUrl: getSpecIcon(Class.ClassWarrior, 1),

			talents: Presets.Talent25.data,
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
					1: Presets.GearFuryDefault.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearFuryDefault.gear,
				},
			},
		},
	],
});

export class WarriorSimUI extends IndividualSimUI<Spec.SpecWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
