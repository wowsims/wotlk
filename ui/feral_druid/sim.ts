import {
	Class,
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../core/proto/common.js';
import { APLRotation } from 'ui/core/proto/apl.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Gear } from '../core/proto_utils/gear.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { TypedEvent } from '../core/typed_event.js';
import { Player } from '../core/player.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
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
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
		Stat.StatMana,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 2.40,
			[Stat.StatAgility]: 2.39,
			[Stat.StatAttackPower]: 1,
			[Stat.StatMeleeHit]: 2.51,
			[Stat.StatMeleeCrit]: 2.23,
			[Stat.StatMeleeHaste]: 1.83,
			[Stat.StatArmorPenetration]: 2.08,
			[Stat.StatExpertise]: 2.44,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 16.5,
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
			arcaneBrilliance: true,
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			manaSpringTotem: TristateEffect.TristateEffectRegular,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
			battleShout: TristateEffect.TristateEffectImproved,
		}),
		partyBuffs: PartyBuffs.create({
			heroicPresence: true,
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			judgementOfWisdom: true,
			giftOfArthas: true,
			exposeArmor: true,
			faerieFire: TristateEffect.TristateEffectImproved,
			sunderArmor: true,
			curseOfWeakness: TristateEffect.TristateEffectRegular,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		IconInputs.IntellectBuff,
		IconInputs.MP5Buff,
		IconInputs.JudgementOfWisdom,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DruidInputs.LatencyMs,
			DruidInputs.AssumeBleedActive,
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
			Presets.ROTATION_PRESET_LEGACY_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.DefaultGear,
		],
	},

	autoRotation: (_player: Player<Spec.SpecFeralDruid>): APLRotation => {
		return Presets.ROTATION_PRESET_LEGACY_DEFAULT.rotation.rotation!;
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
