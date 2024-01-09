import {Class, Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, Race, RaidBuffs, Spec, Stat, TristateEffect} from '../core/proto/common.js';
import {APLAction, APLListItem, APLPrepullAction, APLRotation} from '../core/proto/apl.js';
import {Stats} from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import {Player} from '../core/player.js';
import {IndividualSimUI, registerSpecConfig} from '../core/individual_sim_ui.js';
import {
	Mage_Rotation as MageRotation,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
} from '../core/proto/mage.js';

import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';

import * as MageInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecMage, {
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

		if (player.getTalentTree() === 0) {
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
			demonicPactSp: 500,
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
			MageInputs.WaterElementalDisobeyChance,
			OtherInputs.ReactionTime,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.nibelungAverageCasts,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_SIMPLE,
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
			if (player.getTalents().iceShards > 0) {
				return Presets.FROSTFIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
			}
			return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		} else {
			return Presets.FROST_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecMage>, simple: MageRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const prepullMirrorImage = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":55342}}},"doAtValue":{"const":{"val":"-2s"}}}`);

		const berserking = APLAction.fromJsonString(`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":26297}}}`);
		const hyperspeedAcceleration = APLAction.fromJsonString(`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":54758}}}`);
		const combatPot = APLAction.fromJsonString(`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}`);
		const evocation = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpLe","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"25%"}}}},"castSpell":{"spellId":{"spellId":12051}}}`);

		const arcaneBlastBelowStacks = APLAction.fromJsonString(`{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"auraId":{"spellId":36032}}},"rhs":{"const":{"val":"4"}}}},{"and":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"auraId":{"spellId":36032}}},"rhs":{"const":{"val":"3"}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(simple.only3ArcaneBlastStacksBelowManaPercent * 100).toFixed(0)}%"}}}}]}}]}},"castSpell":{"spellId":{"spellId":42897}}}`);
		const arcaneMissilesWithMissileBarrageBelowMana = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44401}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(simple.missileBarrageBelowManaPercent * 100).toFixed(0)}%"}}}}]}},"castSpell":{"spellId":{"spellId":42846}}}`);
		const arcaneMisslesWithMissileBarrage = APLAction.fromJsonString(`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44401}}},"castSpell":{"spellId":{"spellId":42846}}}`);
		const arcaneBlastAboveMana = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(simple.blastWithoutMissileBarrageAboveManaPercent * 100).toFixed(0)}%"}}}},"castSpell":{"spellId":{"spellId":42897}}}`);
		const arcaneMissiles = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42846}}}`);
		const arcaneBarrage = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":44781}}}`);

		const maintainImpScorch = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":12873},"maxOverlap":{"const":{"val":"4s"}}}},"castSpell":{"spellId":{"spellId":42859}}}`);
		const pyroWithHotStreak = APLAction.fromJsonString(`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44448}}},"castSpell":{"spellId":{"spellId":42891}}}`);
		const livingBomb = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12s"}}}}]}},"multidot":{"spellId":{"spellId":55360},"maxDots":10,"maxOverlap":{"const":{"val":"0ms"}}}}`);
		const cheekyFireBlastFinisher = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"spellCastTime":{"spellId":{"spellId":42859}}}}},"castSpell":{"spellId":{"spellId":42873}}}`);
		const scorchFinisher = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"4s"}}}},"castSpell":{"spellId":{"spellId":42859}}}`);
		const fireball = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42833}}}`);
		const frostfireBolt = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":47610}}}`);
		const scorch = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42859}}}`);

		const deepFreeze = APLAction.fromJsonString(`{"condition":{"auraIsActive":{"auraId":{"spellId":44545}}},"castSpell":{"spellId":{"spellId":44572}}}`);
		const frostfireBoltWithBrainFreeze = APLAction.fromJsonString(`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44549}}},"castSpell":{"spellId":{"spellId":47610}}}`);
		const frostbolt = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42842}}}`);
		const iceLance = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"spellId":44545}}},"rhs":{"const":{"val":"1"}}}},"castSpell":{"spellId":{"spellId":42914}}}`);

		prepullActions.push(prepullMirrorImage);
		if (player.getTalents().improvedScorch > 0 && simple.maintainImprovedScorch) {
			actions.push(maintainImpScorch);
		}

		const talentTree = player.getTalentTree();
		if (talentTree == 0) { // Arcane
			actions.push(...[
				berserking,
				hyperspeedAcceleration,
				combatPot,
				simple.missileBarrageBelowManaPercent > 0 ? arcaneMissilesWithMissileBarrageBelowMana : null,
				arcaneBlastBelowStacks,
				arcaneMisslesWithMissileBarrage,
				evocation,
				arcaneBlastAboveMana,
				simple.useArcaneBarrage ? arcaneBarrage : null,
				arcaneMissiles,
			].filter(a => a) as Array<APLAction>)
		} else if (talentTree == 1) { // Fire
			actions.push(...[
				pyroWithHotStreak,
				livingBomb,
				cheekyFireBlastFinisher,
				scorchFinisher,
				simple.primaryFireSpell == PrimaryFireSpell.Fireball
					? fireball
					: (simple.primaryFireSpell == PrimaryFireSpell.FrostfireBolt
						? frostfireBolt : scorch),
			].filter(a => a) as Array<APLAction>)
		} else if (talentTree == 2) { // Frost
			actions.push(...[
				berserking,
				hyperspeedAcceleration,
				evocation,
				deepFreeze,
				frostfireBoltWithBrainFreeze,
				simple.useIceLance ? iceLance : null,
				frostbolt,
			].filter(a => a) as Array<APLAction>)
		}

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action => APLListItem.create({
				action: action,
			}))
		});
	},

	raidSimPresets: [
		{
			spec: Spec.SpecMage,
			tooltip: 'Arcane Mage',
			defaultName: 'Arcane',
			iconUrl: getSpecIcon(Class.ClassMage, 0),

			talents: Presets.ArcaneTalents.data,
			specOptions: Presets.DefaultArcaneOptions,
			consumes: Presets.DefaultArcaneConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.ARCANE_P1_PRESET.gear,
					2: Presets.ARCANE_P2_PRESET.gear,
					3: Presets.ARCANE_P3_PRESET_ALLIANCE.gear,
					4: Presets.ARCANE_P4_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.ARCANE_P1_PRESET.gear,
					2: Presets.ARCANE_P2_PRESET.gear,
					3: Presets.ARCANE_P3_PRESET_HORDE.gear,
					4: Presets.ARCANE_P4_PRESET_HORDE.gear,
				},
			},
		},
		{
			spec: Spec.SpecMage,
			tooltip: 'TTW Fire Mage',
			defaultName: 'TTW Fire',
			iconUrl: getSpecIcon(Class.ClassMage, 1),

			talents: Presets.FireTalents.data,
			specOptions: Presets.DefaultFireOptions,
			consumes: Presets.DefaultFireConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.FIRE_P1_PRESET.gear,
					2: Presets.FIRE_P2_PRESET.gear,
					3: Presets.FIRE_P3_PRESET_ALLIANCE.gear,
					4: Presets.FIRE_P4_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.FIRE_P1_PRESET.gear,
					2: Presets.FIRE_P2_PRESET.gear,
					3: Presets.FIRE_P3_PRESET_HORDE.gear,
					4: Presets.FIRE_P4_PRESET_HORDE.gear,
				},
			},
		},
		{
			spec: Spec.SpecMage,
			tooltip: 'FFB Fire Mage',
			defaultName: 'FFB Fire',
			iconUrl: "https://wow.zamimg.com/images/wow/icons/medium/ability_mage_frostfirebolt.jpg",

			talents: Presets.FrostfireTalents.data,
			specOptions: Presets.DefaultFFBOptions,
			consumes: Presets.DefaultFireConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.FIRE_P1_PRESET.gear,
					2: Presets.FFB_P2_PRESET.gear,
					3: Presets.FFB_P3_PRESET_ALLIANCE.gear,
					4: Presets.FFB_P4_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.FIRE_P1_PRESET.gear,
					2: Presets.FFB_P2_PRESET.gear,
					3: Presets.FFB_P3_PRESET_HORDE.gear,
					4: Presets.FFB_P4_PRESET_HORDE.gear,
				},
			},
		},
	],
});

export class MageSimUI extends IndividualSimUI<Spec.SpecMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
