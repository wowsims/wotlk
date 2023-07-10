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
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { EquipmentSpec } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { TypedEvent } from '../core/typed_event.js';

import { ProtectionWarrior, ProtectionWarrior_Rotation as ProtectionWarriorRotation, WarriorTalents as WarriorTalents, ProtectionWarrior_Options as ProtectionWarriorOptions } from '../core/proto/warrior.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as ProtectionWarriorInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionWarriorSimUI extends IndividualSimUI<Spec.SpecProtectionWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionWarrior>) {
		super(parentElem, player, {
			cssClass: 'protection-warrior-sim-ui',
			cssScheme: 'warrior',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmor,
				Stat.StatBonusArmor,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			epPseudoStats: [
				PseudoStat.PseudoStatMainHandDps,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatBonusArmor,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P2_SURVIVAL_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.174,
					[Stat.StatBonusArmor]: 0.155,
					[Stat.StatStamina]: 2.336,
					[Stat.StatStrength]: 1.555,
					[Stat.StatAgility]: 2.771,
					[Stat.StatAttackPower]: 0.32,
					[Stat.StatExpertise]: 1.44,
					[Stat.StatMeleeHit]: 1.432,
					[Stat.StatMeleeCrit]: 0.925,
					[Stat.StatMeleeHaste]: 0.431,
					[Stat.StatArmorPenetration]: 1.055,
					[Stat.StatBlock]: 1.320,
					[Stat.StatBlockValue]: 1.373,
					[Stat.StatDodge]: 2.606,
					[Stat.StatParry]: 2.649,
					[Stat.StatDefense]: 3.305,
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 6.081,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.StandardTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					abominationsMight: true,
					swiftRetribution: true,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
					sanctifiedRetribution: true,
					devotionAura: TristateEffect.TristateEffectImproved,
					stoneskinTotem: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					retributionAura: true,
					thorns: TristateEffect.TristateEffectImproved,
					shadowProtection: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					blessingOfSanctuary: true,
				}),
				debuffs: Debuffs.create({
					sunderArmor: true,
					mangle: true,
					vindication: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					insectSwarm: true,
					bloodFrenzy: true,
					judgementOfLight: true,
					heartOfTheCrusader: true,
					frostFever: TristateEffect.TristateEffectImproved,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				ProtectionWarriorInputs.ShoutPicker,
				ProtectionWarriorInputs.ShatteringThrow,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionWarriorInputs.ProtectionWarriorRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.HealthBuff,
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HealingCadenceVariation,
					OtherInputs.BurstWindow,
					OtherInputs.HpPercentForDefensives,
					OtherInputs.InspirationUptime,
					ProtectionWarriorInputs.StartingRage,
					ProtectionWarriorInputs.Munch,
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
					Presets.UATalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_PRERAID_BALANCED_PRESET,
					Presets.P1_BALANCED_PRESET,
					Presets.P2_SURVIVAL_PRESET,
				],
			},
		});
	}
}
