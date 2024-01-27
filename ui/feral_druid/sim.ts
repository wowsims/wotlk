import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import {
	APLAction,
	APLPrepullAction,
	APLListItem,
	APLRotation,
	APLRotation_Type as APLRotationType,
} from '../core/proto/apl.js';
import {
	Class,
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	Race,
	RaidBuffs,
	Spec,
	Stat,
  TristateEffect,
  WeaponImbue,
	SaygesFortune
} from '../core/proto/common.js';
import { FeralDruid_Rotation as DruidRotation } from '../core/proto/druid.js';
import { Gear } from '../core/proto_utils/gear.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { TypedEvent } from '../core/typed_event.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';
import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralDruid, {
	cssClass: 'feral-druid-sim-ui',
	cssScheme: 'druid',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],
	warnings: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMana,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatMP5,
	],
	epPseudoStats: [
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMana,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatMP5,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 2.20,
			[Stat.StatAgility]: 2.02,
			[Stat.StatAttackPower]: 1,
			[Stat.StatFeralAttackPower]: 1,
			[Stat.StatMeleeHit]: 8.21,
			[Stat.StatMeleeCrit]: 8.19,
			[Stat.StatMeleeHaste]: 4.17,
			[Stat.StatMana]: 0.04,
			[Stat.StatIntellect]: 0.67,
			[Stat.StatSpirit]: 0.08,
			[Stat.StatMP5]: 0.46,
		}, {
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultRotation,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			aspectOfTheLion: true,
			arcaneBrilliance: true,
			giftOfTheWild: TristateEffect.TristateEffectRegular,
			battleShout: TristateEffect.TristateEffectRegular,
		}),

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: IndividualBuffs.create({
			blessingOfMight: TristateEffect.TristateEffectImproved,
			blessingOfWisdom: TristateEffect.TristateEffectRegular,
			boonOfBlackfathom: true,
			ashenvalePvpBuff: true,
			saygesFortune: SaygesFortune.SaygesDamage,
		}),

		debuffs: Debuffs.create({
			judgementOfWisdom: false,
			giftOfArthas: false,
			exposeArmor: TristateEffect.TristateEffectMissing,
			faerieFire: false,
			sunderArmor: true,
			curseOfRecklessness: false,
			homunculi: 0,
			curseOfVulnerability: true,
			ancientCorrosivePoison: 30,
		}),

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.IntellectBuff,
		BuffDebuffInputs.BlessingOfWisdom,
		BuffDebuffInputs.ManaSpringTotem,
		BuffDebuffInputs.JudgementOfWisdom,
	],
  excludeBuffDebuffInputs: [
		WeaponImbue.ElementalSharpeningStone,
		WeaponImbue.DenseSharpeningStone,
		WeaponImbue.WildStrikes,
		BuffDebuffInputs.BleedDebuff,
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DruidInputs.LatencyMs,
			// DruidInputs.AssumeBleedActive,
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
		rotations: [
			Presets.SIMPLE_ROTATION_DEFAULT,
			Presets.APL_ROTATION_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.BlankPreset,
			Presets.DefaultGear,
		],
	},

	autoRotation: (_player: Player<Spec.SpecFeralDruid>): APLRotation => {
		return Presets.APL_ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecFeralDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const preroarDuration = Math.min(simple.preroarDuration, 33.0);
		const preRoar = APLPrepullAction.fromJsonString(`{"action":{"activateAura":{"auraId":{"spellId":407988}}},"doAtValue":{"const":{"val":"-${(34.0 - preroarDuration).toFixed(2)}s"}}}`);
		const preTF = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":5217,"rank":1}}},"doAtValue":{"const":{"val":"-3s"}}}`);
		const doRotation = APLAction.fromJsonString(`{"catOptimalRotationAction":{"maxWaitTime":${simple.maxWaitTime.toFixed(2)},"minCombosForRip":${simple.minCombosForRip.toFixed(0)},"maintainFaerieFire":${simple.maintainFaerieFire},"useShredTrick":${simple.useShredTrick}}}`);

		prepullActions.push(...[
			preroarDuration > 0 ? preRoar: null,
			simple.precastTigersFury ? preTF: null,
		].filter(a => a) as Array<APLPrepullAction>)

		actions.push(...[
			doRotation,
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
			spec: Spec.SpecFeralDruid,
			tooltip: specNames[Spec.SpecFeralDruid],
			defaultName: 'Cat',
			iconUrl: getSpecIcon(Class.ClassDruid, 3),

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
					1: Presets.DefaultGear.gear,
				},
				[Faction.Horde]: {
					1: Presets.DefaultGear.gear,
				},
			},
		},
	],
})

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, SPEC_CONFIG);
	}

	calcArpTarget(gear: Gear): number {
		let arpTarget = 1399;

		// First handle ArP proc trinkets
		if (gear.hasTrinket(45931)) {
			arpTarget -= 751;
		} else if (gear.hasTrinket(40256)) {
			arpTarget -= 612;
		}

		// Then check for Executioner enchant
		const weapon = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);

		if ((weapon != null) && (weapon!.enchant != null) && (weapon!.enchant!.effectId == 3225)) {
			arpTarget -= 120;
		}

		return arpTarget;
	}

	calcCritCap(gear: Gear): Stats {
		const baseCritCapPercentage = 77.8; // includes 3% Crit debuff
		let agiProcs = 0;

		if (gear.hasRelic(47668)) {
			agiProcs += 200;
		}

		if (gear.hasRelic(50456)) {
			agiProcs += 44*5;
		}

		if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
			agiProcs += 510;
		}

		if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
			agiProcs += 450;
		}

		if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
			agiProcs += 300;
		}

		return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - agiProcs*1.1*1.06*1.02/83.33) * 45.91);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
	}

	detectArpStackConfiguration(arpTarget: number): boolean {
		const currentArp = Stats.fromProto(this.player.getCurrentStats().finalStats).getStat(Stat.StatArmorPenetration);
		return (arpTarget > 1000) && (currentArp > 648) && (currentArp + 20 < arpTarget + 11);
	}
}
