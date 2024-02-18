import {
	Class,
	Cooldowns, 
	Debuffs, 
	Faction,
	IndividualBuffs, 
	PartyBuffs, 
	Race,
	RaidBuffs, 
	Spec, 
	Stat, 
	PseudoStat, 
	TristateEffect
} from '../core/proto/common.js';

import {
	APLAction,
	APLPrepullAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import { ProtectionWarrior_Rotation as ProtectionWarriorRotation } from '../core/proto/warrior.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';

import * as ProtectionWarriorInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionWarrior, {
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
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatResilience,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P3_PRESET.gear,
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
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.HealthBuff,
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
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_DEFAULT,
			Presets.ROTATION_PRESET_SIMPLE,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_BALANCED_PRESET,
			Presets.P4_PRERAID_PRESET,
			Presets.P1_BALANCED_PRESET,
			Presets.P2_SURVIVAL_PRESET,
			Presets.P3_PRESET,
			Presets.P4_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecProtectionWarrior>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},
	
	simpleRotation: (player: Player<Spec.SpecProtectionWarrior>, simple: ProtectionWarriorRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const preShout = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":47440}}},"doAtValue":{"const":{"val":"-10s"}}}`);

		const heroicStrike = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"30"}}}},"castSpell":{"spellId":{"tag":1,"spellId":47450}}}`);
		const shieldSlam = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":47488}}}`);
		const revenge = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":57823}}}`);
		const refreshShout = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"sourceUnit":{"type":"Self"},"auraId":{"spellId":47440},"maxOverlap":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":47440}}}`);
		const refreshTclap = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":47502},"maxOverlap":{"const":{"val":"2s"}}}},"castSpell":{"spellId":{"spellId":47502}}}`);
		const refreshDemo = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":47437},"maxOverlap":{"const":{"val":"2s"}}}},"castSpell":{"spellId":{"spellId":25203}}}`);
		const devastate = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":47498}}}`);

		prepullActions.push(preShout);

		actions.push(...[
			heroicStrike,
			shieldSlam,
			revenge,
			refreshShout,
			refreshTclap,
			refreshDemo,
			devastate,
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
			spec: Spec.SpecProtectionWarrior,
			tooltip: 'Protection Warrior',
			defaultName: 'Protection',
			iconUrl: getSpecIcon(Class.ClassWarrior, 2),

			talents: Presets.StandardTalents.data,
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
					1: Presets.P1_BALANCED_PRESET.gear,
					2: Presets.P2_SURVIVAL_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_BALANCED_PRESET.gear,
					2: Presets.P2_SURVIVAL_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
			},
		},
	],
});

export class ProtectionWarriorSimUI extends IndividualSimUI<Spec.SpecProtectionWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
