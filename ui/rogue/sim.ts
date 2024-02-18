import {
	Class, 
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
	WeaponType
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import {
	Rogue_Options_PoisonImbue,
} from '../core/proto/rogue.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as RogueInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRogue, {
	cssClass: 'rogue-sim-ui',
	cssScheme: 'rogue',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
		'Rotations are not fully optimized, especially for non-standard setups.',
	],
	warnings: [
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.sim.encounter.changeEmitter,
				getContent: () => {
					let hasNoArmor = false
					for (const target of simUI.sim.encounter.targets) {
						if (new Stats(target.stats).getStat(Stat.StatArmor) <= 0) {
							hasNoArmor = true
							break
						}
					}
					if (hasNoArmor) {
						return 'One or more targets have no armor. Check advanced encounter settings.';
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (
						simUI.player.getTalents().mutilate &&
						(simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType != WeaponType.WeaponTypeDagger ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType != WeaponType.WeaponTypeDagger)
					) {
						return '"Mutilate" talent selected, but daggers not equipped in both hands.';
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getTalents().hackAndSlash) {
						if (simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeSword ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeAxe ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeSword ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeAxe) {
							return '';
						} else {
							return '"Hack and Slash" talent selected, but swords or axes not equipped.';
						}
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getTalents().closeQuartersCombat) {
						if (simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeFist ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeDagger ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeFist ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeDagger) {
							return '';
						} else {
							return '"Close Quarters Combat" talent selected, but fists or daggers not equipped.';
						}
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getTalents().maceSpecialization) {
						if (simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeMace ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeMace) {
							return '';
						} else {
							return '"Mace Specialization" talent selected, but maces not equipped.';
						}
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					const mhWeaponSpeed = simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
					const ohWeaponSpeed = simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
					const mhImbue = simUI.player.getSpecOptions().mhImbue;
					const ohImbue = simUI.player.getSpecOptions().ohImbue;
					if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined' || !simUI.player.getSpecOptions().applyPoisonsManually) {
						return '';
					}
					if (mhWeaponSpeed < ohWeaponSpeed && ohImbue == Rogue_Options_PoisonImbue.DeadlyPoison) {
						return 'Deadly poison applied to slower (off hand) weapon.';
					}
					if (ohWeaponSpeed < mhWeaponSpeed && mhImbue == Rogue_Options_PoisonImbue.DeadlyPoison) {
						return 'Deadly poison applied to slower (main hand) weapon.';
					}
					return '';
				},
			};
		},
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
	],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatSpellHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET_ASSASSINATION.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatAgility]: 1.86,
			[Stat.StatStrength]: 1.14,
			[Stat.StatAttackPower]: 1,
			[Stat.StatSpellCrit]: 0.28,
			[Stat.StatSpellHit]: 0.08,
			[Stat.StatMeleeHit]: 1.39,
			[Stat.StatMeleeCrit]: 1.32,
			[Stat.StatMeleeHaste]: 1.48,
			[Stat.StatArmorPenetration]: 0.84,
			[Stat.StatExpertise]: 0.98,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 2.94,
			[PseudoStat.PseudoStatOffHandDps]: 2.45,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.AssassinationTalents137.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			bloodlust: true,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
			icyTalons: true,
			leaderOfThePack: TristateEffect.TristateEffectImproved,
			abominationsMight: true,
			swiftRetribution: true,
			elementalOath: true,
			sanctifiedRetribution: true,
		}),
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			heartOfTheCrusader: true,
			mangle: true,
			sunderArmor: true,
			faerieFire: TristateEffect.TristateEffectImproved,
			shadowMastery: true,
			earthAndMoon: true,
			bloodFrenzy: true,
		}),
	},

	playerInputs: {
		inputs: [
			RogueInputs.ApplyPoisonsManually
		]
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		RogueInputs.MainHandImbue,
		RogueInputs.OffHandImbue,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellCritBuff,
		BuffDebuffInputs.SpellCritDebuff,
		BuffDebuffInputs.SpellHitDebuff,
		BuffDebuffInputs.SpellDamageDebuff
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			RogueInputs.StartingOverkillDuration,
			RogueInputs.VanishBreakTime,
			RogueInputs.AssumeBleedActive,
			RogueInputs.HonorOfThievesCritRate,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.AssassinationTalents137,
			Presets.AssassinationTalents182,
			Presets.AssassinationTalentsBF,
			Presets.CombatHackTalents,
			Presets.CombatCQCTalents,
			Presets.SubtletyTalents,
			Presets.HemoSubtletyTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_MUTILATE,
			Presets.ROTATION_PRESET_MUTILATE_EXPOSE,
			Presets.ROTATION_PRESET_RUPTURE_MUTILATE,
			Presets.ROTATION_PRESET_RUPTURE_MUTILATE_EXPOSE,
			Presets.ROTATION_PRESET_COMBAT,
			Presets.ROTATION_PRESET_COMBAT_EXPOSE,
			Presets.ROTATION_PRESET_COMBAT_CLEAVE_SND,
			Presets.ROTATION_PRESET_COMBAT_CLEAVE_SND_EXPOSE,
			Presets.ROTATION_PRESET_AOE,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET_ASSASSINATION,
			Presets.PRERAID_PRESET_COMBAT,
			Presets.P1_PRESET_ASSASSINATION,
			Presets.P1_PRESET_COMBAT,
			Presets.P1_PRESET_HEMO_SUB,
			Presets.P2_PRESET_ASSASSINATION,
			Presets.P2_PRESET_COMBAT,
			Presets.P3_PRESET_ASSASSINATION,
			Presets.P3_PRESET_COMBAT,
			Presets.P4_PRESET_ASSASSINATION,
			Presets.P4_PRESET_COMBAT,
			Presets.P5_PRESET_ASSASSINATION,
			Presets.P5_PRESET_COMBAT,
			Presets.P2_PRESET_HEMO_SUB,
			Presets.P3_PRESET_HEMO_SUB,
			Presets.P3_PRESET_DANCE_SUB,
		],
	},

	autoRotation: (player: Player<Spec.SpecRogue>): APLRotation => {
		const talentTree = player.getTalentTree();
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 5) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else if (talentTree == 0) {
			return Presets.ROTATION_PRESET_MUTILATE_EXPOSE.rotation.rotation!;
		} else if (talentTree == 1) {
			return Presets.ROTATION_PRESET_COMBAT_EXPOSE.rotation.rotation!;
		} else {
			// TODO: Need a sub rotation here
			return Presets.ROTATION_PRESET_MUTILATE_EXPOSE.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRogue,
			tooltip: 'Assassination Rogue',
			defaultName: 'Assassination',
			iconUrl: getSpecIcon(Class.ClassRogue, 0),

			talents: Presets.AssassinationTalents137.data,
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
					1: Presets.P1_PRESET_ASSASSINATION.gear,
					2: Presets.P2_PRESET_ASSASSINATION.gear,
					3: Presets.P3_PRESET_ASSASSINATION.gear,
					4: Presets.P4_PRESET_ASSASSINATION.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET_ASSASSINATION.gear,
					2: Presets.P2_PRESET_ASSASSINATION.gear,
					3: Presets.P3_PRESET_ASSASSINATION.gear,
					4: Presets.P4_PRESET_ASSASSINATION.gear,
				},
			},
		},
		{
			spec: Spec.SpecRogue,
			tooltip: 'Combat Rogue',
			defaultName: 'Combat',
			iconUrl: getSpecIcon(Class.ClassRogue, 1),

			talents: Presets.CombatCQCTalents.data,
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
					1: Presets.P1_PRESET_COMBAT.gear,
					2: Presets.P2_PRESET_COMBAT.gear,
					3: Presets.P3_PRESET_COMBAT.gear,
					4: Presets.P4_PRESET_COMBAT.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET_COMBAT.gear,
					2: Presets.P2_PRESET_COMBAT.gear,
					3: Presets.P3_PRESET_COMBAT.gear,
					4: Presets.P4_PRESET_COMBAT.gear,
				},
			},
		},
	],
});

export class RogueSimUI extends IndividualSimUI<Spec.SpecRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRogue>) {
		super(parentElem, player, SPEC_CONFIG);
		this.player.changeEmitter.on((c) => {
			const options = this.player.getSpecOptions()
			const encounter = this.sim.encounter
			if (!options.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return
				}
				if (encounter.targets.length > 3) {
					options.mhImbue = Rogue_Options_PoisonImbue.InstantPoison
					options.ohImbue = Rogue_Options_PoisonImbue.InstantPoison
				} else {
					if (mhWeaponSpeed <= ohWeaponSpeed) {
						options.mhImbue = Rogue_Options_PoisonImbue.DeadlyPoison
						options.ohImbue = Rogue_Options_PoisonImbue.InstantPoison
					} else {
						options.mhImbue = Rogue_Options_PoisonImbue.InstantPoison
						options.ohImbue = Rogue_Options_PoisonImbue.DeadlyPoison
					}
				}
			}
			this.player.setSpecOptions(c, options)
		});
		this.sim.encounter.changeEmitter.on((c) => {
			const options = this.player.getSpecOptions()
			const encounter = this.sim.encounter
			if (!options.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return
				}
				if (encounter.targets.length > 3) {
					options.mhImbue = Rogue_Options_PoisonImbue.InstantPoison
					options.ohImbue = Rogue_Options_PoisonImbue.InstantPoison
				} else {
					if (mhWeaponSpeed <= ohWeaponSpeed) {
						options.mhImbue = Rogue_Options_PoisonImbue.DeadlyPoison
						options.ohImbue = Rogue_Options_PoisonImbue.InstantPoison
					} else {
						options.mhImbue = Rogue_Options_PoisonImbue.InstantPoison
						options.ohImbue = Rogue_Options_PoisonImbue.DeadlyPoison
					}
				}
			}
			this.player.setSpecOptions(c, options)
		});
	}
}
