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
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { Stat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { TypedEvent } from '../core/typed_event.js';

import { FeralDruid, FeralDruid_Rotation as FeralDruidRotation, DruidTalents as DruidTalents, FeralDruid_Options as FeralDruidOptions } from '../core/proto/druid.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, {
			cssClass: 'feral-druid-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"Rotations need updating",
				"Lacerate doesn't snapshot correctly",
				"Form Stat changes may display/function incorrectly",
				"Bear / cat auto attacks don't change with form",
				"ilotp mana regen not implemented"
			],
			warnings: [
				// Warning that a frontal rotation is not implemented.
				(simUI: IndividualSimUI<Spec.SpecFeralDruid>) => {
					return {
						updateOn: simUI.player.inFrontOfTargetChangeEmitter,
						getContent: () => simUI.player.getInFrontOfTarget() ? 'Frontal rotation (without Shred) is not implemented.' : '',
					};
				},
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
				Stat.StatMana,
				Stat.StatMP5,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStrength]: 2.379,
					[Stat.StatAgility]: 2.15,
					[Stat.StatAttackPower]: 1,
					[Stat.StatMeleeHit]: 2.51,
					[Stat.StatMeleeCrit]: 2,
					[Stat.StatMeleeHaste]: 1.67,
					[Stat.StatArmorPenetration]: 1.66,
					[Stat.StatExpertise]: 2.28,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.StandardTalents.data,
				// Default spec-specific settings.
				specOptions: FeralDruidOptions.create({
					innervateTarget: RaidTarget.create({
						targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
					}),
					latencyMs: 100,
				}),
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					battleShout: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
				}),
				partyBuffs: PartyBuffs.create({
					braidedEterniumChain: true,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					bloodFrenzy: true,
					giftOfArthas: true,
					exposeArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					sunderArmor: true,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				DruidInputs.SelfInnervate,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DruidInputs.FeralDruidRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					DruidInputs.LatencyMs,
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
					Presets.StandardTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PreRaid_PRESET,
					Presets.P1_PRESET,
				],
			},
		});
	}
}
