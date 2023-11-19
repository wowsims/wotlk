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

import * as TankDeathknightPresets from '../tank_deathknight/presets.js';
import * as DeathknightPresets from '../deathknight/presets.js';
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

import { TankDeathknightSimUI } from '../tank_deathknight/sim.js';
import { DeathknightSimUI } from '../deathknight/sim.js';
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
	[Spec.SpecTankDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new TankDeathknightSimUI(parentElem, player),
	[Spec.SpecDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new DeathknightSimUI(parentElem, player),
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
		spec: Spec.SpecTankDeathknight,
		rotation: TankDeathknightPresets.DefaultRotation,
		talents: TankDeathknightPresets.BloodTalents.data,
		specOptions: TankDeathknightPresets.DefaultOptions,
		consumes: TankDeathknightPresets.DefaultConsumes,
		defaultName: 'Blood Tank',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: TankDeathknightPresets.P1_BLOOD_PRESET.gear,
				2: TankDeathknightPresets.P2_BLOOD_PRESET.gear,
			},
			[Faction.Horde]: {
				1: TankDeathknightPresets.P1_BLOOD_PRESET.gear,
				2: TankDeathknightPresets.P2_BLOOD_PRESET.gear,
			},
		},
		tooltip: 'Blood Tank Death Knight',
		iconUrl: getSpecIcon(Class.ClassDeathknight, 0),
	},
	{
		spec: Spec.SpecDeathknight,
		rotation: DeathknightPresets.DefaultBloodRotation,
		talents: DeathknightPresets.BloodTalents.data,
		specOptions: DeathknightPresets.DefaultBloodOptions,
		consumes: DeathknightPresets.DefaultConsumes,
		defaultName: 'Blood DPS',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: DeathknightPresets.P1_BLOOD_PRESET.gear,
				2: DeathknightPresets.P2_BLOOD_PRESET.gear,
				3: DeathknightPresets.P3_BLOOD_PRESET.gear,
			},
			[Faction.Horde]: {
				1: DeathknightPresets.P1_BLOOD_PRESET.gear,
				2: DeathknightPresets.P2_BLOOD_PRESET.gear,
				3: DeathknightPresets.P3_BLOOD_PRESET.gear,
			},
		},
		tooltip: 'Blood DPS Death Knight',
		iconUrl: getSpecIcon(Class.ClassDeathknight, 3),
	},
	{
		spec: Spec.SpecDeathknight,
		rotation: DeathknightPresets.DefaultFrostRotation,
		talents: DeathknightPresets.FrostTalents.data,
		specOptions: DeathknightPresets.DefaultFrostOptions,
		consumes: DeathknightPresets.DefaultConsumes,
		defaultName: 'Frost',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: DeathknightPresets.P1_FROST_PRESET.gear,
				2: DeathknightPresets.P2_FROST_PRESET.gear,
				3: DeathknightPresets.P3_FROST_PRESET.gear,
			},
			[Faction.Horde]: {
				1: DeathknightPresets.P1_FROST_PRESET.gear,
				2: DeathknightPresets.P2_FROST_PRESET.gear,
				3: DeathknightPresets.P3_FROST_PRESET.gear,
			},
		},
		otherDefaults: DeathknightPresets.OtherDefaults,
		tooltip: 'Frost Death Knight',
		iconUrl: getSpecIcon(Class.ClassDeathknight, 1),
	},
	{
		spec: Spec.SpecDeathknight,
		rotation: DeathknightPresets.DefaultUnholyRotation,
		talents: DeathknightPresets.UnholyDualWieldTalents.data,
		specOptions: DeathknightPresets.DefaultUnholyOptions,
		consumes: DeathknightPresets.DefaultConsumes,
		defaultName: 'Unholy',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: DeathknightPresets.P1_UNHOLY_DW_PRESET.gear,
				2: DeathknightPresets.P2_UNHOLY_DW_PRESET.gear,
				3: DeathknightPresets.P3_UNHOLY_DW_PRESET.gear,
			},
			[Faction.Horde]: {
				1: DeathknightPresets.P1_UNHOLY_DW_PRESET.gear,
				2: DeathknightPresets.P2_UNHOLY_DW_PRESET.gear,
				3: DeathknightPresets.P3_UNHOLY_DW_PRESET.gear,
			},
		},
		otherDefaults: DeathknightPresets.OtherDefaults,
		tooltip: 'Dual-Wield Unholy DK',
		iconUrl: getSpecIcon(Class.ClassDeathknight, 2),
	},
	{
		spec: Spec.SpecBalanceDruid,
		rotation: BalanceDruidPresets.DefaultRotation,
		talents: BalanceDruidPresets.Phase2Talents.data,
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
				1: BalanceDruidPresets.P1_PRESET.gear,
				2: BalanceDruidPresets.P2_PRESET.gear,
				3: BalanceDruidPresets.P3_PRESET_ALLI.gear,
			},
			[Faction.Horde]: {
				1: BalanceDruidPresets.P1_PRESET.gear,
				2: BalanceDruidPresets.P2_PRESET.gear,
				3: BalanceDruidPresets.P3_PRESET_HORDE.gear,
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
				1: FeralDruidPresets.P1_PRESET.gear,
				2: FeralDruidPresets.P2_PRESET.gear,
				3: FeralDruidPresets.P3_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralDruidPresets.P1_PRESET.gear,
				2: FeralDruidPresets.P2_PRESET.gear,
				3: FeralDruidPresets.P3_PRESET.gear,
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
				1: FeralTankDruidPresets.P1_PRESET.gear,
				2: FeralTankDruidPresets.P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralTankDruidPresets.P1_PRESET.gear,
				2: FeralTankDruidPresets.P2_PRESET.gear,
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
				1: RestorationDruidPresets.P1_PRESET.gear,
				2: RestorationDruidPresets.P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RestorationDruidPresets.P1_PRESET.gear,
				2: RestorationDruidPresets.P2_PRESET.gear,
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
				1: HunterPresets.MM_P1_PRESET.gear,
				2: HunterPresets.MM_P2_PRESET.gear,
				3: HunterPresets.MM_P3_PRESET.gear,
				4: HunterPresets.MM_P4_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.MM_P1_PRESET.gear,
				2: HunterPresets.MM_P2_PRESET.gear,
				3: HunterPresets.MM_P3_PRESET.gear,
				4: HunterPresets.MM_P4_PRESET.gear,
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
				1: HunterPresets.MM_P1_PRESET.gear,
				2: HunterPresets.MM_P2_PRESET.gear,
				3: HunterPresets.MM_P3_PRESET.gear,
				4: HunterPresets.MM_P4_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.MM_P1_PRESET.gear,
				2: HunterPresets.MM_P2_PRESET.gear,
				3: HunterPresets.MM_P3_PRESET.gear,
				4: HunterPresets.MM_P4_PRESET.gear,
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
				1: HunterPresets.SV_P1_PRESET.gear,
				2: HunterPresets.SV_P2_PRESET.gear,
				3: HunterPresets.SV_P3_PRESET.gear,
				4: HunterPresets.SV_P4_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.SV_P1_PRESET.gear,
				2: HunterPresets.SV_P2_PRESET.gear,
				3: HunterPresets.SV_P3_PRESET.gear,
				4: HunterPresets.SV_P4_PRESET.gear,
			},
		},
		tooltip: 'Survival Hunter',
		iconUrl: getSpecIcon(Class.ClassHunter, 2),
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultSimpleRotation,
		talents: MagePresets.ArcaneTalents.data,
		specOptions: MagePresets.DefaultArcaneOptions,
		consumes: MagePresets.DefaultArcaneConsumes,
		otherDefaults: MagePresets.OtherDefaults,
		defaultName: 'Arcane',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.ARCANE_P1_PRESET.gear,
				2: MagePresets.ARCANE_P2_PRESET.gear,
				3: MagePresets.ARCANE_P3_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.ARCANE_P1_PRESET.gear,
				2: MagePresets.ARCANE_P2_PRESET.gear,
				3: MagePresets.ARCANE_P3_PRESET_HORDE.gear,
			},
		},
		tooltip: 'Arcane Mage',
		iconUrl: getSpecIcon(Class.ClassMage, 0),
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultSimpleRotation,
		talents: MagePresets.FireTalents.data,
		specOptions: MagePresets.DefaultFireOptions,
		consumes: MagePresets.DefaultFireConsumes,
		otherDefaults: MagePresets.OtherDefaults,
		defaultName: 'TTW Fire',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
				2: MagePresets.FIRE_P2_PRESET.gear,
				3: MagePresets.FIRE_P3_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
				2: MagePresets.FIRE_P2_PRESET.gear,
				3: MagePresets.FIRE_P3_PRESET_HORDE.gear,
			},
		},
		tooltip: 'TTW Fire Mage',
		iconUrl: getSpecIcon(Class.ClassMage, 1),
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultSimpleRotation,
		talents: MagePresets.FrostfireTalents.data,
		specOptions: MagePresets.DefaultFFBOptions,
		consumes: MagePresets.DefaultFireConsumes,
		otherDefaults: MagePresets.OtherDefaults,
		defaultName: 'FFB Fire',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
				2: MagePresets.FFB_P2_PRESET.gear,
				3: MagePresets.FFB_P3_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
				2: MagePresets.FFB_P2_PRESET.gear,
				3: MagePresets.FFB_P3_PRESET_HORDE.gear,
			},
		},
		tooltip: 'FFB Fire Mage',
		iconUrl: "https://wow.zamimg.com/images/wow/icons/medium/ability_mage_frostfirebolt.jpg",
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
				1: RoguePresets.P1_PRESET_ASSASSINATION.gear,
				2: RoguePresets.P2_PRESET_ASSASSINATION.gear,
				3: RoguePresets.P3_PRESET_ASSASSINATION.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET_ASSASSINATION.gear,
				2: RoguePresets.P2_PRESET_ASSASSINATION.gear,
				3: RoguePresets.P3_PRESET_ASSASSINATION.gear,
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
				1: RoguePresets.P1_PRESET_COMBAT.gear,
				2: RoguePresets.P2_PRESET_COMBAT.gear,
				3: RoguePresets.P3_PRESET_COMBAT.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET_COMBAT.gear,
				2: RoguePresets.P2_PRESET_COMBAT.gear,
				3: RoguePresets.P3_PRESET_COMBAT.gear,
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
				1: ElementalShamanPresets.P1_PRESET.gear,
				2: ElementalShamanPresets.P2_PRESET.gear,
				3: ElementalShamanPresets.P3_PRESET_ALLI.gear,
				4: ElementalShamanPresets.P4_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ElementalShamanPresets.P1_PRESET.gear,
				2: ElementalShamanPresets.P2_PRESET.gear,
				3: ElementalShamanPresets.P3_PRESET_HORDE.gear,
				4: ElementalShamanPresets.P4_PRESET.gear,
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
				1: EnhancementShamanPresets.P1_PRESET.gear,
				2: EnhancementShamanPresets.P2_PRESET_FT.gear,
				3: EnhancementShamanPresets.P3_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: EnhancementShamanPresets.P1_PRESET.gear,
				2: EnhancementShamanPresets.P2_PRESET_FT.gear,
				3: EnhancementShamanPresets.P3_PRESET_HORDE.gear,
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
				1: RestorationShamanPresets.P1_PRESET.gear,
				2: RestorationShamanPresets.P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RestorationShamanPresets.P1_PRESET.gear,
				2: RestorationShamanPresets.P2_PRESET.gear,
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
				1: HealingPriestPresets.DISC_P1_PRESET.gear,
				2: HealingPriestPresets.DISC_P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.DISC_P1_PRESET.gear,
				2: HealingPriestPresets.DISC_P2_PRESET.gear,
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
				1: HealingPriestPresets.HOLY_P1_PRESET.gear,
				2: HealingPriestPresets.HOLY_P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.HOLY_P1_PRESET.gear,
				2: HealingPriestPresets.HOLY_P2_PRESET.gear,
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
		rotation: WarriorPresets.ArmsRotation,
		talents: WarriorPresets.ArmsTalents.data,
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
				1: WarriorPresets.P1_ARMS_PRESET.gear,
				2: WarriorPresets.P2_ARMS_PRESET.gear,
				3: WarriorPresets.P3_ARMS_4P_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_ARMS_PRESET.gear,
				2: WarriorPresets.P2_ARMS_PRESET.gear,
				3: WarriorPresets.P3_ARMS_4P_PRESET_HORDE.gear,
			},
		},
		tooltip: 'Arms Warrior',
		iconUrl: getSpecIcon(Class.ClassWarrior, 0),
	},
	{
		spec: Spec.SpecWarrior,
		rotation: WarriorPresets.DefaultRotation,
		talents: WarriorPresets.FuryTalents.data,
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
				1: WarriorPresets.P1_FURY_PRESET.gear,
				2: WarriorPresets.P2_FURY_PRESET.gear,
				3: WarriorPresets.P3_FURY_PRESET_ALLIANCE.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_FURY_PRESET.gear,
				2: WarriorPresets.P2_FURY_PRESET.gear,
				3: WarriorPresets.P3_FURY_PRESET_HORDE.gear,
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
				1: ProtectionWarriorPresets.P1_BALANCED_PRESET.gear,
				2: ProtectionWarriorPresets.P2_SURVIVAL_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionWarriorPresets.P1_BALANCED_PRESET.gear,
				2: ProtectionWarriorPresets.P2_SURVIVAL_PRESET.gear,
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
				1: HolyPaladinPresets.P1_PRESET.gear,
				2: HolyPaladinPresets.P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HolyPaladinPresets.P1_PRESET.gear,
				2: HolyPaladinPresets.P2_PRESET.gear,
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
				1: ProtectionPaladinPresets.P1_PRESET.gear,
				2: ProtectionPaladinPresets.P2_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionPaladinPresets.P1_PRESET.gear,
				2: ProtectionPaladinPresets.P2_PRESET.gear,
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
				1: RetributionPaladinPresets.P1_PRESET.gear,
				2: RetributionPaladinPresets.P2_PRESET.gear,
				3: RetributionPaladinPresets.P3_PRESET.gear,
				4: RetributionPaladinPresets.P4_PRESET.gear,
				5: RetributionPaladinPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RetributionPaladinPresets.P1_PRESET.gear,
				2: RetributionPaladinPresets.P2_PRESET.gear,
				3: RetributionPaladinPresets.P3_PRESET.gear,
				4: RetributionPaladinPresets.P4_PRESET.gear,
				5: RetributionPaladinPresets.P5_PRESET.gear,
			},
		},
		tooltip: 'Retribution Paladin',
		iconUrl: getSpecIcon(Class.ClassPaladin, 2),
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.AfflictionRotation,
		talents: WarlockPresets.AfflictionTalents.data,
		specOptions: WarlockPresets.AfflictionOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Affliction',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_AFFLICTION_PRESET.gear,
				2: WarlockPresets.P2_AFFLICTION_PRESET.gear,
				3: WarlockPresets.P3_AFFLICTION_ALLIANCE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_AFFLICTION_PRESET.gear,
				2: WarlockPresets.P2_AFFLICTION_PRESET.gear,
				3: WarlockPresets.P3_AFFLICTION_HORDE_PRESET.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Affliction Warlock',
		iconUrl: getSpecIcon(Class.ClassWarlock, 0),
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DemonologyRotation,
		talents: WarlockPresets.DemonologyTalents.data,
		specOptions: WarlockPresets.DemonologyOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Demonology',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_DEMODESTRO_PRESET.gear,
				2: WarlockPresets.P2_DEMODESTRO_PRESET.gear,
				3: WarlockPresets.P3_DEMO_ALLIANCE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_DEMODESTRO_PRESET.gear,
				2: WarlockPresets.P2_DEMODESTRO_PRESET.gear,
				3: WarlockPresets.P3_DEMO_HORDE_PRESET.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Demonology Warlock',
		iconUrl: getSpecIcon(Class.ClassWarlock, 1),
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DestructionRotation,
		talents: WarlockPresets.DestructionTalents.data,
		specOptions: WarlockPresets.DestructionOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Destruction',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_DEMODESTRO_PRESET.gear,
				2: WarlockPresets.P2_DEMODESTRO_PRESET.gear,
				3: WarlockPresets.P3_DESTRO_ALLIANCE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_DEMODESTRO_PRESET.gear,
				2: WarlockPresets.P2_DEMODESTRO_PRESET.gear,
				3: WarlockPresets.P3_DESTRO_HORDE_PRESET.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Destruction Warlock',
		iconUrl: getSpecIcon(Class.ClassWarlock, 2),
	},
];

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];
