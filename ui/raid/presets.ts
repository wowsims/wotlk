import { IndividualSimUI, OtherDefaults } from '../core/individual_sim_ui.js';

import {
	Class,
	Consumes,
	EquipmentSpec,
	Faction,
	Race,
	Spec
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	getSpecIcon,
	specNames,
	SpecOptions,
	SpecRotation,
} from '../core/proto_utils/utils.js';

import { Player } from '../core/player.js';

import * as BalanceDruidPresets from '../balance_druid/presets.js';
import * as FeralDruidPresets from '../feral_druid/presets.js';
import * as FeralTankDruidPresets from '../feral_tank_druid/presets.js';
import * as RestorationDruidPresets from '../restoration_druid/presets.js';
import * as ElementalShamanPresets from '../elemental_shaman/presets.js';
import * as EnhancementShamanPresets from '../enhancement_shaman/presets.js';
import * as RestorationShamanPresets from '../restoration_shaman/presets.js';
import * as HunterPresets from '../hunter/presets.js';
import * as MagePresets from '../mage/presets.js';
import * as RoguePresets from '../rogue/presets.js';
import * as HolyPaladinPresets from '../holy_paladin/presets.js';
import * as ProtectionPaladinPresets from '../protection_paladin/presets.js';
import * as RetributionPaladinPresets from '../retribution_paladin/presets.js';
import * as HealingPriestPresets from '../healing_priest/presets.js';
import * as ShadowPriestPresets from '../shadow_priest/presets.js';
import * as WarriorPresets from '../warrior/presets.js';
import * as ProtectionWarriorPresets from '../protection_warrior/presets.js';
import * as WarlockPresets from '../warlock/presets.js';

import { BalanceDruidSimUI } from '../balance_druid/sim.js';
import { FeralDruidSimUI } from '../feral_druid/sim.js';
import { FeralTankDruidSimUI } from '../feral_tank_druid/sim.js';
import { RestorationDruidSimUI } from '../restoration_druid/sim.js';
import { ElementalShamanSimUI } from '../elemental_shaman/sim.js';
import { EnhancementShamanSimUI } from '../enhancement_shaman/sim.js';
import { RestorationShamanSimUI } from '../restoration_shaman/sim.js';
import { HunterSimUI } from '../hunter/sim.js';
import { MageSimUI } from '../mage/sim.js';
import { RogueSimUI } from '../rogue/sim.js';
import { HolyPaladinSimUI } from '../holy_paladin/sim.js';
import { ProtectionPaladinSimUI } from '../protection_paladin/sim.js';
import { RetributionPaladinSimUI } from '../retribution_paladin/sim.js';
import { HealingPriestSimUI } from '../healing_priest/sim.js';
import { ShadowPriestSimUI } from '../shadow_priest/sim.js';
import { WarriorSimUI } from '../warrior/sim.js';
import { ProtectionWarriorSimUI } from '../protection_warrior/sim.js';
import { WarlockSimUI } from '../warlock/sim.js';

export const specSimFactories: Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>> = {
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecFeralTankDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralTankDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	[Spec.SpecHunter]: (parentElem: HTMLElement, player: Player<any>) => new HunterSimUI(parentElem, player),
	[Spec.SpecMage]: (parentElem: HTMLElement, player: Player<any>) => new MageSimUI(parentElem, player),
	[Spec.SpecRogue]: (parentElem: HTMLElement, player: Player<any>) => new RogueSimUI(parentElem, player),
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	[Spec.SpecHealingPriest]: (parentElem: HTMLElement, player: Player<any>) => new HealingPriestSimUI(parentElem, player),
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	[Spec.SpecWarrior]: (parentElem: HTMLElement, player: Player<any>) => new WarriorSimUI(parentElem, player),
	[Spec.SpecProtectionWarrior]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionWarriorSimUI(parentElem, player),
	[Spec.SpecWarlock]: (parentElem: HTMLElement, player: Player<any>) => new WarlockSimUI(parentElem, player),
};

// Configuration necessary for creating new players.
export interface PresetSpecSettings<SpecType extends Spec> {
	spec: Spec,
	rotation: SpecRotation<SpecType>,
	talents: SavedTalents,
	specOptions: SpecOptions<SpecType>,
	consumes: Consumes,

	defaultName: string,
	defaultFactionRaces: Record<Faction, Race>,
	defaultGear: Record<Faction, Record<number, EquipmentSpec>>,
	otherDefaults?: OtherDefaults,

	tooltip: string,
	iconUrl: string,
}

export const playerPresets: Array<PresetSpecSettings<any>> = [
	{
		spec: Spec.SpecBalanceDruid,
		rotation: BalanceDruidPresets.DefaultRotation,
		talents: BalanceDruidPresets.BalanceTalents.data,
		specOptions: BalanceDruidPresets.DefaultOptions,
		consumes: BalanceDruidPresets.DefaultConsumes,
		otherDefaults: BalanceDruidPresets.OtherDefaults,
		defaultName: 'Balance',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: BalanceDruidPresets.DEFAULT_PRESET.gear,
			},
			[Faction.Horde]: {
				1: BalanceDruidPresets.DEFAULT_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecBalanceDruid],
		iconUrl: getSpecIcon(Class.ClassDruid, 0),
	},
	{
		spec: Spec.SpecFeralDruid,
		rotation: FeralDruidPresets.DefaultRotation,
		talents: FeralDruidPresets.StandardTalents.data,
		specOptions: FeralDruidPresets.DefaultOptions,
		consumes: FeralDruidPresets.DefaultConsumes,
		defaultName: 'Cat',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: FeralDruidPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: FeralDruidPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecFeralDruid],
		iconUrl: getSpecIcon(Class.ClassDruid, 3),
	},
	{
		spec: Spec.SpecFeralTankDruid,
		rotation: FeralTankDruidPresets.DefaultRotation,
		talents: FeralTankDruidPresets.StandardTalents.data,
		specOptions: FeralTankDruidPresets.DefaultOptions,
		consumes: FeralTankDruidPresets.DefaultConsumes,
		defaultName: 'Bear',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: FeralTankDruidPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: FeralTankDruidPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecFeralTankDruid],
		iconUrl: getSpecIcon(Class.ClassDruid, 1),
	},
	{
		spec: Spec.SpecRestorationDruid,
		rotation: RestorationDruidPresets.DefaultRotation,
		talents: RestorationDruidPresets.CelestialFocusTalents.data,
		specOptions: RestorationDruidPresets.DefaultOptions,
		consumes: RestorationDruidPresets.DefaultConsumes,
		defaultName: 'Restoration',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RestorationDruidPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: RestorationDruidPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecRestorationDruid],
		iconUrl: getSpecIcon(Class.ClassDruid, 2),
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.BeastMasteryTalents.data,
		specOptions: HunterPresets.BMDefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'Beast Mastery',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.BeastMasteryDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.BeastMasteryDefaultGear.gear,
			},
		},
		tooltip: 'Beast Mastery Hunter',
		iconUrl: getSpecIcon(Class.ClassHunter, 0),
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.MarksmanTalents.data,
		specOptions: HunterPresets.DefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'Marksmanship',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.MarksmanDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.MarksmanDefaultGear.gear,
			},
		},
		tooltip: 'Marksmanship Hunter',
		iconUrl: getSpecIcon(Class.ClassHunter, 1),
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.SurvivalTalents.data,
		specOptions: HunterPresets.DefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'Survival',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.SurvivalDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.SurvivalDefaultGear.gear,
			},
		},
		tooltip: 'Survival Hunter',
		iconUrl: getSpecIcon(Class.ClassHunter, 2),
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.ROTATION_PRESET_DEFAULT.rotation,
		talents: MagePresets.DefaultTalents.data,
		specOptions: MagePresets.DefaultOptions,
		consumes: MagePresets.DefaultConsumes,
		otherDefaults: MagePresets.OtherDefaults,
		defaultName: 'Fire Mage',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.DEFAULT_GEAR.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.DEFAULT_GEAR.gear,
			},
		},
		tooltip: 'Fire Mage',
		iconUrl: getSpecIcon(Class.ClassMage, 0),
	},
	{
		spec: Spec.SpecRogue,
		rotation: RoguePresets.DefaultRotation,
		talents: RoguePresets.AssassinationTalents137.data,
		specOptions: RoguePresets.DefaultOptions,
		consumes: RoguePresets.DefaultConsumes,
		defaultName: 'Assassination',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RoguePresets.AssassinationDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.AssassinationDefaultGear.gear,
			},
		},
		tooltip: 'Assassination Rogue',
		iconUrl: getSpecIcon(Class.ClassRogue, 0),
	},
	{
		spec: Spec.SpecRogue,
		rotation: RoguePresets.DefaultRotation,
		talents: RoguePresets.CombatCQCTalents.data,
		specOptions: RoguePresets.DefaultOptions,
		consumes: RoguePresets.DefaultConsumes,
		defaultName: 'Combat',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RoguePresets.CombatDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.CombatDefaultGear.gear,
			},
		},
		tooltip: 'Combat Rogue',
		iconUrl: getSpecIcon(Class.ClassRogue, 1),
	},
	{
		spec: Spec.SpecElementalShaman,
		rotation: ElementalShamanPresets.DefaultRotation,
		talents: ElementalShamanPresets.StandardTalents.data,
		specOptions: ElementalShamanPresets.DefaultOptions,
		consumes: ElementalShamanPresets.DefaultConsumes,
		defaultName: 'Elemental',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ElementalShamanPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: ElementalShamanPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecElementalShaman],
		iconUrl: getSpecIcon(Class.ClassShaman, 0),
	},
	{
		spec: Spec.SpecEnhancementShaman,
		rotation: EnhancementShamanPresets.DefaultRotation,
		talents: EnhancementShamanPresets.StandardTalents.data,
		specOptions: EnhancementShamanPresets.DefaultOptions,
		consumes: EnhancementShamanPresets.DefaultConsumes,
		defaultName: 'Enhancement',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: EnhancementShamanPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: EnhancementShamanPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecEnhancementShaman],
		iconUrl: getSpecIcon(Class.ClassShaman, 1),
	},
	{
		spec: Spec.SpecRestorationShaman,
		rotation: RestorationShamanPresets.DefaultRotation,
		talents: RestorationShamanPresets.RaidHealingTalents.data,
		specOptions: RestorationShamanPresets.DefaultOptions,
		consumes: RestorationShamanPresets.DefaultConsumes,
		defaultName: 'Restoration',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RestorationShamanPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: RestorationShamanPresets.DefaultGear.gear,
			},
		},
		tooltip: specNames[Spec.SpecRestorationShaman],
		iconUrl: getSpecIcon(Class.ClassShaman, 2),
	},
	{
		spec: Spec.SpecHealingPriest,
		rotation: HealingPriestPresets.DiscDefaultRotation,
		talents: HealingPriestPresets.DiscTalents.data,
		specOptions: HealingPriestPresets.DefaultOptions,
		consumes: HealingPriestPresets.DefaultConsumes,
		defaultName: 'Discipline',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HealingPriestPresets.DiscDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.DiscDefaultGear.gear,
			},
		},
		tooltip: 'Discipline Priest',
		iconUrl: getSpecIcon(Class.ClassPriest, 0),
	},
	{
		spec: Spec.SpecHealingPriest,
		rotation: HealingPriestPresets.HolyDefaultRotation,
		talents: HealingPriestPresets.HolyTalents.data,
		specOptions: HealingPriestPresets.DefaultOptions,
		consumes: HealingPriestPresets.DefaultConsumes,
		defaultName: 'Holy',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HealingPriestPresets.HolyDefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.HolyDefaultGear.gear,
			},
		},
		tooltip: 'Holy Priest',
		iconUrl: getSpecIcon(Class.ClassPriest, 1),
	},
	{
		spec: Spec.SpecShadowPriest,
		rotation: ShadowPriestPresets.DefaultRotation,
		talents: ShadowPriestPresets.StandardTalents.data,
		specOptions: ShadowPriestPresets.DefaultOptions,
		consumes: ShadowPriestPresets.DefaultConsumes,
		defaultName: 'Shadow',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ShadowPriestPresets.BLANK_GEAR_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ShadowPriestPresets.BLANK_GEAR_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecShadowPriest],
		iconUrl: getSpecIcon(Class.ClassPriest, 2),
	},
	{
		spec: Spec.SpecWarrior,
		rotation: [],
		talents: WarriorPresets.Talent25.data,
		specOptions: WarriorPresets.DefaultOptions,
		consumes: WarriorPresets.DefaultConsumes,
		defaultName: 'Arms',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {

			},
			[Faction.Horde]: {

			},
		},
		tooltip: 'Arms Warrior',
		iconUrl: getSpecIcon(Class.ClassWarrior, 0),
	},
	{
		spec: Spec.SpecWarrior,
		rotation: [],
		talents: WarriorPresets.Talent25.data,
		specOptions: WarriorPresets.DefaultOptions,
		consumes: WarriorPresets.DefaultConsumes,
		defaultName: 'Fury',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
			},
			[Faction.Horde]: {
			},
		},
		tooltip: 'Fury Warrior',
		iconUrl: getSpecIcon(Class.ClassWarrior, 1),
	},
	{
		spec: Spec.SpecProtectionWarrior,
		rotation: ProtectionWarriorPresets.DefaultRotation,
		talents: ProtectionWarriorPresets.StandardTalents.data,
		specOptions: ProtectionWarriorPresets.DefaultOptions,
		consumes: ProtectionWarriorPresets.DefaultConsumes,
		defaultName: 'Protection',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ProtectionWarriorPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: ProtectionWarriorPresets.DefaultGear.gear,
			},
		},
		tooltip: 'Protection Warrior',
		iconUrl: getSpecIcon(Class.ClassWarrior, 2),
	},
	{
		spec: Spec.SpecHolyPaladin,
		rotation: HolyPaladinPresets.DefaultRotation,
		talents: HolyPaladinPresets.StandardTalents.data,
		specOptions: HolyPaladinPresets.DefaultOptions,
		consumes: HolyPaladinPresets.DefaultConsumes,
		defaultName: 'Holy',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HolyPaladinPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: HolyPaladinPresets.DefaultGear.gear,
			},
		},
		tooltip: 'Holy Paladin',
		iconUrl: getSpecIcon(Class.ClassPaladin, 0),
	},
	{
		spec: Spec.SpecProtectionPaladin,
		rotation: ProtectionPaladinPresets.DefaultRotation,
		talents: ProtectionPaladinPresets.GenericAoeTalents.data,
		specOptions: ProtectionPaladinPresets.DefaultOptions,
		consumes: ProtectionPaladinPresets.DefaultConsumes,
		defaultName: 'Protection',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ProtectionPaladinPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: ProtectionPaladinPresets.DefaultGear.gear,
			},
		},
		tooltip: 'Protection Paladin',
		iconUrl: getSpecIcon(Class.ClassPaladin, 1),
	},
	{
		spec: Spec.SpecRetributionPaladin,
		rotation: RetributionPaladinPresets.DefaultRotation,
		talents: RetributionPaladinPresets.AuraMasteryTalents.data,
		specOptions: RetributionPaladinPresets.DefaultOptions,
		consumes: RetributionPaladinPresets.DefaultConsumes,
		defaultName: 'Retribution',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RetributionPaladinPresets.DefaultGear.gear,
			},
			[Faction.Horde]: {
				1: RetributionPaladinPresets.DefaultGear.gear,
			},
		},
		tooltip: 'Retribution Paladin',
		iconUrl: getSpecIcon(Class.ClassPaladin, 2),
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.ROTATION_PRESET_DEFAULT,
		talents: WarlockPresets.DefaultTalents.data,
		specOptions: WarlockPresets.DefaultOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Destruction Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.DEFAULT_GEAR.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.DEFAULT_GEAR.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Destruction Warlock',
		iconUrl: getSpecIcon(Class.ClassWarlock, 2),
	},
];

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];
