import {
	Class,
	Debuffs,
	Faction,
	HandType,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../core/proto/common.js';
import { APLRotation } from '../core/proto/apl.js';
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as DeathKnightInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecDeathknight, {
	cssClass: 'deathknight-sim-ui',
	cssScheme: 'death-knight',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatArmor,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
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
		Stat.StatArmor,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
	],
	defaults: {
		// Default equipped gear.
		gear: Presets.P2_UNHOLY_DW_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 3.22,
			[Stat.StatAgility]: 0.62,
			[Stat.StatArmor]: 0.01,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 1.13,
			[Stat.StatMeleeHaste]: 1.85,
			[Stat.StatMeleeHit]: 1.92,
			[Stat.StatMeleeCrit]: 0.76,
			[Stat.StatArmorPenetration]: 0.77,
			[Stat.StatSpellHit]: 0.80,
			[Stat.StatSpellCrit]: 0.34,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 3.10,
			[PseudoStat.PseudoStatOffHandDps]: 1.79,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.UnholyDualWieldTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultUnholyOptions,
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
			moonkinAura: TristateEffect.TristateEffectRegular,
			wrathOfAirTotem: true,
			powerWordFortitude: TristateEffect.TristateEffectImproved,
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
			faerieFire: TristateEffect.TristateEffectImproved,
			sunderArmor: true,
			ebonPlaguebringer: true,
			mangle: true,
			heartOfTheCrusader: true,
			shadowMastery: true,
		}),
	},

	autoRotation: (player: Player<Spec.SpecDeathknight>): APLRotation => {
		const talentTree = player.getTalentTree();
		const numTargets = player.sim.encounter.targets.length;
		switch (talentTree) {
			case 0: 
				if (player.getSpecOptions().drwPestiApply || numTargets > 1) {
					if (numTargets > 5) {
						return Presets.BLOOD_PESTI_AOE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
					} else {
						return Presets.BLOOD_DPS_ROTATION_PRESET_DEFAULT.rotation.rotation!;
					}
				} else {
					return Presets.BLOOD_DPS_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				}
			case 1: 
				const talentPoints = player.getTalentTreePoints()
				// TODO: Add Frost AOE rotation
				if (talentPoints[0] > talentPoints[2]) {
					return Presets.FROST_BL_PESTI_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				} else {
					return Presets.FROST_UH_PESTI_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				}
			default: 
				if (numTargets > 1) {
					return Presets.UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				} else {
					if (player.getEquippedItem(ItemSlot.ItemSlotMainHand)!.item.handType == HandType.HandTypeTwoHand) {
						return Presets.UNHOLY_2H_ROTATION_PRESET_DEFAULT.rotation.rotation!;
					} else {
						return Presets.UNHOLY_DW_ROTATION_PRESET_DEFAULT.rotation.rotation!;
					}
				}
		}
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
	],
	petConsumeInputs: [
		ConsumablesInputs.SpicedMammothTreats,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellDamageDebuff,
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.AttackPowerDebuff,
		BuffDebuffInputs.DamageReductionPercentBuff,
		BuffDebuffInputs.MeleeAttackSpeedDebuff,
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DeathKnightInputs.SelfUnholyFrenzy,
			DeathKnightInputs.StartingRunicPower,
			DeathKnightInputs.PetUptime,
			DeathKnightInputs.DrwPestiApply,
			DeathKnightInputs.UseAMSInput,
			DeathKnightInputs.AvgAMSSuccessRateInput,
			DeathKnightInputs.AvgAMSHitInput,

			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.BloodTalents,
			Presets.FrostTalents,
			Presets.FrostUnholyTalents,
			Presets.UnholyDualWieldTalents,
			Presets.UnholyDualWieldSSTalents,
			Presets.Unholy2HTalents,
			Presets.UnholyAoeTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.BLOOD_DPS_ROTATION_PRESET_DEFAULT,
			Presets.BLOOD_PESTI_AOE_ROTATION_PRESET_DEFAULT,
			Presets.FROST_BL_PESTI_ROTATION_PRESET_DEFAULT,
			Presets.FROST_UH_PESTI_ROTATION_PRESET_DEFAULT,
			Presets.UNHOLY_DW_ROTATION_PRESET_DEFAULT,
			Presets.UNHOLY_2H_ROTATION_PRESET_DEFAULT,
			Presets.UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_BLOOD_PRESET,
			Presets.P2_BLOOD_PRESET,
			Presets.P3_BLOOD_PRESET,
			Presets.P4_BLOOD_PRESET,
			Presets.P1_FROST_PRESET,
			Presets.P2_FROST_PRESET,
			Presets.P3_FROST_PRESET,
			Presets.P4_FROST_PRESET,
			Presets.P1_UNHOLY_DW_PRESET,
			Presets.P2_UNHOLY_DW_PRESET,
			Presets.P3_UNHOLY_DW_PRESET,
			Presets.P4_UNHOLY_DW_PRESET,
			Presets.P4_UNHOLY_2H_PRESET,
			// Not needed anymore just filling ui Space
			// Disabled on purpose
			//Presets.P1_FROSTSUBUNH_PRESET,
			//Presets.P1_FROST_PRE_BIS_PRESET,
			//Presets.PRERAID_UNHOLY_DW_PRESET,
			//Presets.PRERAID_UNHOLY_2H_PRESET,
			//Presets.P1_UNHOLY_2H_PRESET,
		],
	},

	raidSimPresets: [
		{
			spec: Spec.SpecDeathknight,
			tooltip: 'Frost Death Knight',
			defaultName: 'Frost',
			iconUrl: getSpecIcon(Class.ClassDeathknight, 1),

			talents: Presets.FrostTalents.data,
			specOptions: Presets.DefaultFrostOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_FROST_PRESET.gear,
					2: Presets.P2_FROST_PRESET.gear,
					3: Presets.P3_FROST_PRESET.gear,
					4: Presets.P4_FROST_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_FROST_PRESET.gear,
					2: Presets.P2_FROST_PRESET.gear,
					3: Presets.P3_FROST_PRESET.gear,
					4: Presets.P4_FROST_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecDeathknight,
			tooltip: 'Dual-Wield Unholy DK',
			defaultName: 'Unholy',
			iconUrl: getSpecIcon(Class.ClassDeathknight, 2),

			talents: Presets.UnholyDualWieldTalents.data,
			specOptions: Presets.DefaultUnholyOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_UNHOLY_DW_PRESET.gear,
					2: Presets.P2_UNHOLY_DW_PRESET.gear,
					3: Presets.P3_UNHOLY_DW_PRESET.gear,
					4: Presets.P4_UNHOLY_DW_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_UNHOLY_DW_PRESET.gear,
					2: Presets.P2_UNHOLY_DW_PRESET.gear,
					3: Presets.P3_UNHOLY_DW_PRESET.gear,
					4: Presets.P4_UNHOLY_DW_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class DeathknightSimUI extends IndividualSimUI<Spec.SpecDeathknight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDeathknight>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
