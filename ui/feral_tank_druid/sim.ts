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
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { TypedEvent } from '../core/typed_event.js';

import {
	DruidTalents as DruidTalents,
	FeralTankDruid,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Options as DruidOptions
} from '../core/proto/druid.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

export class FeralTankDruidSimUI extends IndividualSimUI<Spec.SpecFeralTankDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralTankDruid>) {
		super(parentElem, player, {
			cssClass: 'feral-tank-druid-sim-ui',
			cssScheme: 'druid',
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
				Stat.StatDodge,
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
				Stat.StatDodge,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 3.5665,
					[Stat.StatBonusArmor]: 0.5187,
					[Stat.StatStamina]: 7.3021,
					[Stat.StatStrength]: 2.3786,
					[Stat.StatAgility]: 4.4974,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 2.6597,
					[Stat.StatMeleeHit]: 2.9282,
					[Stat.StatMeleeCrit]: 1.5143,
					[Stat.StatMeleeHaste]: 2.0983,
					[Stat.StatArmorPenetration]: 1.584,
					[Stat.StatDefense]: 1.8171,
					[Stat.StatDodge]: 2.0196,
					[Stat.StatHealth]: 0.4465,
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 0.0,
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
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					battleShout: TristateEffect.TristateEffectImproved,
					unleashedRage: true,
					windfuryTotem: TristateEffect.TristateEffectImproved,
					arcaneEmpowerment: true,
					moonkinAura: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					heroicPresence: true,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					renewedHope: true,
				}),
				debuffs: Debuffs.create({
					savageCombat: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					exposeArmor: true,
					frostFever: TristateEffect.TristateEffectImproved,
					masterPoisoner: true,
					ebonPlaguebringer: true,
					shadowMastery: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DruidInputs.FeralTankDruidRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
					IconInputs.HealthBuff,
					IconInputs.SpellCritBuff,
					IconInputs.SpellCritDebuff,
					IconInputs.SpellHitDebuff,
					IconInputs.SpellDamageDebuff,
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
					OtherInputs.InspirationUptime,
					OtherInputs.HpPercentForDefensives,
					DruidInputs.StartingRage,
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
					Presets.P1_PRESET, Presets.P2_PRESET
				],
			},
		});
	}
}
