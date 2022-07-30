import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Encounter } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { EquipmentSpec } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { TypedEvent } from '/wotlk/core/typed_event.js';

import { ProtectionWarrior, ProtectionWarrior_Rotation as ProtectionWarriorRotation, WarriorTalents as WarriorTalents, ProtectionWarrior_Options as ProtectionWarriorOptions } from '/wotlk/core/proto/warrior.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as ProtectionWarriorInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionWarriorSimUI extends IndividualSimUI<Spec.SpecProtectionWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionWarrior>) {
		super(parentElem, player, {
			cssClass: 'protection-warrior-sim-ui',
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
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
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
				gear: Presets.P1_BALANCED_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.05,
					[Stat.StatStamina]: 1,
					[Stat.StatStrength]: 0.33,
					[Stat.StatAgility]: 0.6,
					[Stat.StatAttackPower]: 0.06,
					[Stat.StatExpertise]: 0.67,
					[Stat.StatMeleeHit]: 0.67,
					[Stat.StatMeleeCrit]: 0.28,
					[Stat.StatMeleeHaste]: 0.21,
					[Stat.StatArmorPenetration]: 0.19,
					[Stat.StatBlock]: 0.35,
					[Stat.StatBlockValue]: 0.59,
					[Stat.StatDodge]: 0.7,
					[Stat.StatParry]: 0.58,
					[Stat.StatDefense]: 0.8,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.ImpaleProtTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					powerWordFortitude: TristateEffect.TristateEffectRegular,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					windfuryTotem: TristateEffect.TristateEffectImproved,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					mangle: true,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
					faerieFire: TristateEffect.TristateEffectImproved,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				ProtectionWarriorInputs.ShoutPicker,
				ProtectionWarriorInputs.ShieldWall,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionWarriorInputs.ProtectionWarriorRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HpPercentForDefensives,
					ProtectionWarriorInputs.StartingRage,
					ProtectionWarriorInputs.PrecastShout,
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
					Presets.ImpDemoTalents,
					Presets.ImpaleProtTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_BALANCED_PRESET,
					Presets.P2_BALANCED_PRESET,
					Presets.P3_BALANCED_PRESET,
					Presets.P4_BALANCED_PRESET,
					Presets.P5_BALANCED_PRESET,
				],
			},
		});
	}
}
