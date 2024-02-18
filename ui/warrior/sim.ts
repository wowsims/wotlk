import {
	Class,
	Debuffs,
	Faction,
	GemColor,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	Race,
	RaidBuffs,
	Spec,
	Stat, PseudoStat,
	Profession,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Gear } from '../core/proto_utils/gear.js';
import { PhysicalDPSGemOptimizer } from '../core/components/suggest_gems_action.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

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
	modifyDisplayStats: (player: Player<Spec.SpecWarrior>) => {
		let stats = new Stats();
		if (!player.getInFrontOfTarget()) {
			// When behind target, dodge is the only outcome affected by Expertise.
			stats = stats.addStat(Stat.StatExpertise, player.getTalents().weaponMastery * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
		}
		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.P3_FURY_PRESET_ALLIANCE.gear,
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
		talents: Presets.FuryTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			swiftRetribution: true,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
			icyTalons: true,
			abominationsMight: true,
			leaderOfThePack: TristateEffect.TristateEffectRegular,
			sanctifiedRetribution: true,
			bloodlust: true,
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
			bloodFrenzy: true,
			heartOfTheCrusader: true,
			mangle: true,
			sunderArmor: true,
			curseOfWeakness: TristateEffect.TristateEffectRegular,
			faerieFire: TristateEffect.TristateEffectImproved,
			ebonPlaguebringer: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		WarriorInputs.ShoutPicker,
		WarriorInputs.Recklessness,
		WarriorInputs.ShatteringThrow,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		// just for Bryntroll
		BuffDebuffInputs.SpellDamageDebuff,
		BuffDebuffInputs.SpellHitDebuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarriorInputs.StartingRage,
			WarriorInputs.StanceSnapshot,
			WarriorInputs.DisableExpertiseGemming,
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
			Presets.ArmsTalents,
			Presets.FuryTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_FURY,
			Presets.ROTATION_FURY_SUNDER,
			Presets.ROTATION_ARMS,
			Presets.ROTATION_ARMS_SUNDER,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_FURY_PRESET,
			Presets.P1_FURY_PRESET,
			Presets.P2_FURY_PRESET,
			Presets.P3_FURY_PRESET_ALLIANCE,
			Presets.P3_FURY_PRESET_HORDE,
			Presets.P4_FURY_PRESET_ALLIANCE,
			Presets.P4_FURY_PRESET_HORDE,
			Presets.PRERAID_ARMS_PRESET,
			Presets.P1_ARMS_PRESET,
			Presets.P2_ARMS_PRESET,
			Presets.P3_ARMS_2P_PRESET_ALLIANCE,
			Presets.P3_ARMS_4P_PRESET_ALLIANCE,
			Presets.P3_ARMS_2P_PRESET_HORDE,
			Presets.P3_ARMS_4P_PRESET_HORDE,
			Presets.P4_ARMS_PRESET_ALLIANCE,
			Presets.P4_ARMS_PRESET_HORDE,
		],
	},

	autoRotation: (player: Player<Spec.SpecWarrior>): APLRotation => {
		const talentTree = player.getTalentTree();
		if (talentTree == 0) {
			return Presets.ROTATION_ARMS_SUNDER.rotation.rotation!;
		} else {
			return Presets.ROTATION_FURY_SUNDER.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarrior,
			tooltip: 'Arms Warrior',
			defaultName: 'Arms',
			iconUrl: getSpecIcon(Class.ClassWarrior, 0),

			talents: Presets.ArmsTalents.data,
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
					1: Presets.P1_ARMS_PRESET.gear,
					2: Presets.P2_ARMS_PRESET.gear,
					3: Presets.P3_ARMS_4P_PRESET_ALLIANCE.gear,
					4: Presets.P4_ARMS_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_ARMS_PRESET.gear,
					2: Presets.P2_ARMS_PRESET.gear,
					3: Presets.P3_ARMS_4P_PRESET_HORDE.gear,
					4: Presets.P4_ARMS_PRESET_HORDE.gear,
				},
			},
		},
		{
			spec: Spec.SpecWarrior,
			tooltip: 'Fury Warrior',
			defaultName: 'Fury',
			iconUrl: getSpecIcon(Class.ClassWarrior, 1),

			talents: Presets.FuryTalents.data,
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
					1: Presets.P1_FURY_PRESET.gear,
					2: Presets.P2_FURY_PRESET.gear,
					3: Presets.P3_FURY_PRESET_ALLIANCE.gear,
					4: Presets.P4_FURY_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_FURY_PRESET.gear,
					2: Presets.P2_FURY_PRESET.gear,
					3: Presets.P3_FURY_PRESET_HORDE.gear,
					4: Presets.P4_FURY_PRESET_HORDE.gear,
				},
			},
		},
	],
});

export class WarriorSimUI extends IndividualSimUI<Spec.SpecWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
		const gemOptimizer = new WarriorGemOptimizer(this);
	}
}

class WarriorGemOptimizer extends PhysicalDPSGemOptimizer {
	readonly player: Player<Spec.SpecWarrior>;

	constructor(simUI: IndividualSimUI<Spec.SpecWarrior>) {
		super(simUI, true, true, false, true);
		this.player = simUI.player;
	}

	updateGemPriority(ungemmedGear: Gear, passiveStats: Stats) {
		this.useExpGems = !this.player.getSpecOptions().disableExpertiseGemming;
		super.updateGemPriority(ungemmedGear, passiveStats);
	}

	calcExpTarget(): number {
		let expTarget = super.calcExpTarget();
		const weaponMastery = this.player.getTalents().weaponMastery;
		const hasWeaponMasteryTalent = !!weaponMastery;
		
		if (hasWeaponMasteryTalent) {
			expTarget -= weaponMastery * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION;
		}

		return expTarget;
	}
}
