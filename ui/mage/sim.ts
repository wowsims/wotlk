import {
	Debuffs,
	IndividualBuffs,
	PartyBuffs,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect
} from '../core/proto/common.js';
import {
	APLAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';
import {Stats} from '../core/proto_utils/stats.js';
import {Player} from '../core/player.js';
import {IndividualSimUI} from '../core/individual_sim_ui.js';

import {Mage_Rotation_Type as RotationType,} from '../core/proto/mage.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

import * as MageInputs from './inputs.js';
import * as Presets from './presets.js';

export class MageSimUI extends IndividualSimUI<Spec.SpecMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parentElem, player, {
			cssClass: 'mage-sim-ui',
			cssScheme: 'mage',
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
			modifyDisplayStats: (player: Player<Spec.SpecMage>) => {
				let stats = new Stats();
				if (player.getRotation().type == RotationType.Arcane) {
					stats = stats.addStat(Stat.StatSpellHit, player.getTalents().arcaneFocus * 1 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
				}

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.FIRE_P3_PRESET_HORDE.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.48,
					[Stat.StatSpirit]: 0.42,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellHit]: 0.38,
					[Stat.StatSpellCrit]: 0.58,
					[Stat.StatSpellHaste]: 0.94,
					[Stat.StatMP5]: 0.09,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultFireConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultFireRotation,
				// Default talents.
				talents: Presets.Phase3FireTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultFireOptions,
				other: Presets.OtherDefaults,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectImproved,
					wrathOfAirTotem: true,
					divineSpirit: true,
					swiftRetribution: true,
					sanctifiedRetribution: true,
					demonicPact: 500,
					moonkinAura: TristateEffect.TristateEffectImproved,
					arcaneBrilliance: true,
				}),
				partyBuffs: PartyBuffs.create({
					manaTideTotems: 1,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					innervates: 0,
					vampiricTouch: true,
					focusMagic: true,
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					ebonPlaguebringer: true,
					shadowMastery: true,
					heartOfTheCrusader: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				MageInputs.Armor,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: MageInputs.MageRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				//Should add hymn of hope, revitalize, and 
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					MageInputs.FocusMagicUptime,
					OtherInputs.ReactionTime,
					OtherInputs.DistanceFromTarget,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: true,
			},

			presets: {
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.ARCANE_ROTATION_PRESET_DEFAULT,
					Presets.FIRE_ROTATION_PRESET_DEFAULT,
					Presets.FROSTFIRE_ROTATION_PRESET_DEFAULT,
					Presets.FROST_ROTATION_PRESET_DEFAULT,
					Presets.ARCANE_ROTATION_PRESET_AOE,
					Presets.FIRE_ROTATION_PRESET_AOE,
					Presets.FROST_ROTATION_PRESET_AOE,
				],
				// Preset talents that the user can quickly select.
				talents: [
					Presets.ArcaneTalents,
					Presets.FireTalents,
					Presets.FrostfireTalents,
					Presets.FrostTalents,
					Presets.Phase3FireTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.ARCANE_PRERAID_PRESET,
					Presets.FIRE_PRERAID_PRESET,
					Presets.ARCANE_P1_PRESET,
					Presets.FIRE_P1_PRESET,
					Presets.FROST_P1_PRESET,
					Presets.ARCANE_P2_PRESET,
					Presets.FIRE_P2_PRESET,
					Presets.FROST_P2_PRESET,
					Presets.FFB_P2_PRESET,
					Presets.ARCANE_P3_PRESET_ALLIANCE,
					Presets.ARCANE_P3_PRESET_HORDE,
					Presets.FROST_P3_PRESET_ALLIANCE,
					Presets.FROST_P3_PRESET_HORDE,
					Presets.FIRE_P3_PRESET_ALLIANCE,
					Presets.FIRE_P3_PRESET_HORDE,
					Presets.FFB_P3_PRESET_ALLIANCE,
					Presets.FFB_P3_PRESET_HORDE,
					Presets.FIRE_P4_PRESET_HORDE,
					Presets.FIRE_P4_PRESET_ALLIANCE,
					Presets.FFB_P4_PRESET_HORDE,
					Presets.FFB_P4_PRESET_ALLIANCE,
					Presets.ARCANE_P4_PRESET_HORDE,
					Presets.ARCANE_P4_PRESET_ALLIANCE,
				],
			},

			autoRotation: (player: Player<Spec.SpecMage>): APLRotation => {
				const talentTree = player.getTalentTree();
				const numTargets = player.sim.encounter.targets.length;
				if (numTargets > 3) {
					if (talentTree == 0) {
						return Presets.ARCANE_ROTATION_PRESET_AOE.rotation.rotation!;
					} else if (talentTree == 1) {
						return Presets.FIRE_ROTATION_PRESET_AOE.rotation.rotation!;
					} else {
						return Presets.FROST_ROTATION_PRESET_AOE.rotation.rotation!;
					}
				} else if (talentTree == 0) {
					return Presets.ARCANE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				} else if (talentTree == 1) {
					return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				} else {
					return Presets.FROST_ROTATION_PRESET_DEFAULT.rotation.rotation!;
				}
			},
		});
	}
}
