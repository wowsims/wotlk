import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';

import {
	Mage,
	Mage_Rotation as MageRotation,
	Mage_Rotation_Type as RotationType,
	MageTalents as MageTalents,
	Mage_Options as MageOptions,
} from '../core/proto/mage.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as Tooltips from '../core/constants/tooltips.js';

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
				gear: Presets.ARCANE_P1_PRESET.gear,
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
				consumes: Presets.DefaultArcaneConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultArcaneRotation,
				// Default talents.
				talents: Presets.ArcaneTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultArcaneOptions,
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
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					curseOfElements: true,
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
					MageInputs.IgniteMunching,
					MageInputs.EvocationTicks,
					MageInputs.FocusMagicUptime,
					MageInputs.ReactionTime,
					OtherInputs.DistanceFromTarget,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: true,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.ArcaneTalents,
					Presets.FireTalents,
					Presets.FrostfireTalents,
					Presets.FrostTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.ARCANE_PRERAID_PRESET,
					Presets.ARCANE_P1_PRESET,
					Presets.ARCANE_P2_PRESET,
					Presets.FIRE_PRERAID_PRESET,
					Presets.FIRE_P1_PRESET,
					Presets.FIRE_P2_PRESET,
					Presets.FROST_P1_PRESET,
					Presets.FROST_P2_PRESET,
				],
			},
		});
	}
}
