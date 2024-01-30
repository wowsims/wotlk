import {
	Class,
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { Player } from '../core/player.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { TankGemOptimizer } from '../core/components/suggest_gems_action.js';

import {
	FeralTankDruid_Rotation as DruidRotation,
} from '../core/proto/druid.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralTankDruid, {
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
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
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
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
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
		BuffDebuffInputs.HealthBuff,
		BuffDebuffInputs.SpellCritBuff,
		BuffDebuffInputs.SpellCritDebuff,
		BuffDebuffInputs.SpellHitDebuff,
		BuffDebuffInputs.SpellDamageDebuff,
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
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_SIMPLE,
			Presets.ROTATION_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_PRESET,
			Presets.P2_PRESET,
			Presets.P3_PRESET,
			Presets.P4_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecFeralTankDruid>): APLRotation => {
		return Presets.ROTATION_PRESET_SIMPLE.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecFeralTankDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const emergencyLacerate = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}`);
		const demoRoar = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":48560},"maxOverlap":{"const":{"val":"1.5s"}}}},"castSpell":{"spellId":{"spellId":48560}}}`);
		const mangle = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":48564}}}`);
		const delayFaerieFireForMangle = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48564}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":48564}}},"rhs":{"const":{"val":"1.0s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48564}}}}}`);
		const faerieFire = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":16857}}}`);
		const delayFillersForMangle = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48564}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":48564}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48564}}}}}`);
		const lacerate = APLAction.fromJsonString(`{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"${simple.lacerateTime.toFixed(1)}s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}`);
		const swipe = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"${(simple.maulRageThreshold + 15).toFixed(0)}"}}}},"castSpell":{"spellId":{"spellId":48562}}}`);
		const queueMaul = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"${simple.maulRageThreshold.toFixed(0)}"}}}},"castSpell":{"spellId":{"spellId":48480,"tag":1}}}`);
		const waitForFaerieFire = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":16857}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":16857}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":16857}}}}}`);

		actions.push(...[
			emergencyLacerate,
			simple.maintainDemoralizingRoar ? demoRoar : null,
			mangle,
			delayFaerieFireForMangle,
			faerieFire,
			delayFillersForMangle,
			lacerate,
			swipe,
			queueMaul,
			waitForFaerieFire,
		].filter(a => a) as Array<APLAction>)

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action => APLListItem.create({
				action: action,
			}))
		});
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFeralTankDruid,
			tooltip: specNames[Spec.SpecFeralTankDruid],
			defaultName: 'Bear',
			iconUrl: getSpecIcon(Class.ClassDruid, 1),

			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
			},
		},
	],
});

export class FeralTankDruidSimUI extends IndividualSimUI<Spec.SpecFeralTankDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralTankDruid>) {
		super(parentElem, player, SPEC_CONFIG);
		const _gemOptimizer = new TankGemOptimizer(this);
	}
}
