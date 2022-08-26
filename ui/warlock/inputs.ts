import {
	Warlock_Options as WarlockOptions,
	Warlock_Rotation_Type as RotationType,
	Warlock_Rotation_Preset as RotationPreset,
	Warlock_Rotation_PrimarySpell as PrimarySpell,
	Warlock_Rotation_SecondaryDot as SecondaryDot,
	Warlock_Rotation_SpecSpell as SpecSpell,
	Warlock_Rotation_Curse as Curse,
	Warlock_Options_WeaponImbue as WeaponImbue,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
} from '../core/proto/warlock.js';

import { RaidTarget, Spec, Glyphs, Debuffs, IndividualBuffs, RaidBuffs } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { SimUI, SimWarning } from '../core/sim_ui.js';

import { IconPickerConfig } from '../core/components/icon_picker.js';
import { IconEnumPicker, IconEnumPickerConfig, IconEnumValueConfig } from '../core/components/icon_enum_picker.js';
import * as Presets from './presets.js';
import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, Armor>({
	fieldName: 'armor',
	values: [
		{ color: 'grey', value: Armor.NoArmor },
		{ actionId: ActionId.fromSpellId(47893), value: Armor.FelArmor },
		{ actionId: ActionId.fromSpellId(47889), value: Armor.DemonArmor },
	],
});

export const WeaponImbueInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, WeaponImbue>({
	fieldName: 'weaponImbue',
	values: [
		{ color: 'grey', value: WeaponImbue.NoWeaponImbue },
		{ actionId: ActionId.fromItemId(41174), value: WeaponImbue.GrandFirestone },
		{ actionId: ActionId.fromItemId(41196), value: WeaponImbue.GrandSpellstone },
	],
});

export const PetInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, Summon>({
	fieldName: 'summon',
	values: [
		{ color: 'grey', value: Summon.NoSummon },
		{ actionId: ActionId.fromSpellId(688), value: Summon.Imp },
		{ actionId: ActionId.fromSpellId(712), value: Summon.Succubus },
		{ actionId: ActionId.fromSpellId(691), value: Summon.Felhunter },
		{ 
			actionId: ActionId.fromSpellId(30146), value: Summon.Felguard,
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().summonFelguard,
		},
	],
	changeEmitter: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
});

export const PrimarySpellInput = InputHelpers.makeRotationEnumIconInput<Spec.SpecWarlock, PrimarySpell>({
	fieldName: 'primarySpell',
	values: [
		{ actionId: ActionId.fromSpellId(47809), value: PrimarySpell.ShadowBolt },
		{ actionId: ActionId.fromSpellId(47838), value: PrimarySpell.Incinerate },
		{ actionId: ActionId.fromSpellId(47836), value: PrimarySpell.Seed },
	],
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
		const newRotation = player.getRotation();
		if (newValue == PrimarySpell.Seed && newRotation.corruption == true) {
			newRotation.corruption = false
		}
		newRotation.primarySpell = newValue
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
});

export const SecondaryDotInput = InputHelpers.makeRotationEnumIconInput<Spec.SpecWarlock, SecondaryDot>({
	fieldName: 'secondaryDot',
	values: [
		{ color: 'grey', value: SecondaryDot.NoSecondaryDot },
		{ actionId: ActionId.fromSpellId(47811), value: SecondaryDot.Immolate },
		{
			actionId: ActionId.fromSpellId(47843), value: SecondaryDot.UnstableAffliction,
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().unstableAffliction,
		},
	],
	changeEmitter: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
		const newRotation = player.getRotation();
		newRotation.secondaryDot = newValue;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
});

export const SpecSpellInput = InputHelpers.makeRotationEnumIconInput<Spec.SpecWarlock, SpecSpell>({
	fieldName: 'specSpell',
	values: [
		{ color: 'grey', value: SpecSpell.NoSpecSpell },
		{
			actionId: ActionId.fromSpellId(59164), value: SpecSpell.Haunt,
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().haunt,
		},
		{
			actionId: ActionId.fromSpellId(59172), value: SpecSpell.ChaosBolt,
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().chaosBolt,
		},
	],
	changeEmitter: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
		const newRotation = player.getRotation();
		newRotation.specSpell = newValue;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
});


export const CorruptionSpell = {
	type: 'icon' as const,
	id: ActionId.fromSpellId(47813),
	states: 2,
	extraCssClasses: [
		'Corruption-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().corruption,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		if (newRotation.primarySpell == PrimarySpell.Seed && newValue == true) {
			newRotation.primarySpell = PrimarySpell.ShadowBolt
		}
		newRotation.corruption = newValue;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};


export const WarlockRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,

			label: 'SIM PRESETS',
			labelTooltip: 'Quick switch between sim spec presets. Will UPDATE TALENTS and SPELLS to defaults.',
			values: [
				{ name: 'Affliction', value: RotationType.Affliction },
				{ name: 'Demonology', value: RotationType.Demonology },
				{ name: 'Destruction', value: RotationType.Destruction },
			],
			changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().type,
			setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
				var newRotation = player.getRotation();
				var newOptions: WarlockOptions;
				var newGlyphs: Glyphs;
				var newTalents: string;
				// var newIndividualBuffs = player.getBuffs();
				// const raid = player.getRaid();
				// var newDebuffs = raid?.getDebuffs();
				// var newRaidBuffs = raid?.getBuffs();
				if (newValue == RotationType.Affliction) {
					newTalents = Presets.AfflictionTalents.data.talentsString
					newGlyphs = Presets.AfflictionTalents.data.glyphs || Glyphs.create();
					newRotation = Presets.AfflictionRotation
					newOptions = Presets.AfflictionOptions
					// if (newDebuffs != undefined) {
					// 	newDebuffs.shadowMastery = false
					// }
				} else if (newValue == RotationType.Demonology) {
					newTalents = Presets.DemonologyTalents.data.talentsString
					newGlyphs = Presets.DemonologyTalents.data.glyphs || Glyphs.create();
					newRotation = Presets.DemonologyRotation
					newOptions = Presets.DemonologyOptions
					// if (newDebuffs != undefined) {
					// 	newDebuffs.shadowMastery = false
					// }
					// if (newRaidBuffs != undefined) {
					// 	newRaidBuffs.demonicPact = 0
					// }
				} else if (newValue == RotationType.Destruction) {
					newTalents = Presets.DestructionTalents.data.talentsString
					newGlyphs = Presets.DestructionTalents.data.glyphs || Glyphs.create();
					newRotation = Presets.DestructionRotation
					newOptions = Presets.DestructionOptions
					// newIndividualBuffs.improvedSoulLeech = false
					// if (newDebuffs != undefined) {
					// 	newDebuffs.shadowMastery = true
					// }
				}
				newRotation.type = newValue;
				newRotation.preset = RotationPreset.Automatic;
				TypedEvent.freezeAllAndDo(() => {
					player.setTalentsString(eventID, newTalents);
					player.setSpecOptions(eventID, newOptions);
					player.setGlyphs(eventID, newGlyphs);
					player.setRotation(eventID, newRotation);
					// player.setBuffs(eventID, newIndividualBuffs);
					// raid?.setDebuffs(eventID, newDebuffs || Debuffs.create());
					// raid?.setBuffs(eventID, newRaidBuffs || RaidBuffs.create());
				});
			},
		},

		{
			type: 'enum' as const,
			label: 'Spell & Talent',
			labelTooltip: 'Putting it on Automatic will UPDATE talents and spells to defaults.',
			values: [
				{name: "Manual", value: RotationPreset.Manual},
				{name: "Automatic", value: RotationPreset.Automatic},
			],
			changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().preset,
			setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
				var newRotation = player.getRotation();
				if (newValue == RotationPreset.Automatic) {
					var newOptions: WarlockOptions;
					var newGlyphs: Glyphs;
					var newTalents: string;
					// var newIndividualBuffs = player.getBuffs();
					// const raid = player.getRaid();
					// var newDebuffs = raid?.getDebuffs();
					if (newRotation.type == RotationType.Affliction) {
						newTalents = Presets.AfflictionTalents.data.talentsString
						newGlyphs = Presets.AfflictionTalents.data.glyphs || Glyphs.create()
						newRotation = Presets.AfflictionRotation
						newOptions = Presets.AfflictionOptions
						// if (newDebuffs != undefined) {
						// 	newDebuffs.shadowMastery = false
						// }
					} else if (newRotation.type == RotationType.Demonology) {
						newTalents = Presets.DemonologyTalents.data.talentsString
						newGlyphs = Presets.DemonologyTalents.data.glyphs || Glyphs.create()
						newRotation = Presets.DemonologyRotation
						newOptions = Presets.DemonologyOptions
						// if (newDebuffs != undefined) {
						// 	newDebuffs.shadowMastery = false
						// }
					} else if (newRotation.type == RotationType.Destruction) {
						newTalents = Presets.DestructionTalents.data.talentsString
						newGlyphs = Presets.DestructionTalents.data.glyphs || Glyphs.create()
						newRotation = Presets.DestructionRotation
						newOptions = Presets.DestructionOptions
						// newIndividualBuffs.improvedSoulLeech = false
						// if (newDebuffs != undefined) {
						// 	newDebuffs.shadowMastery = true
						// }
					}
				}
				newRotation.preset = newValue;
				const raid = player.getRaid();
				TypedEvent.freezeAllAndDo(() => {
					if (newValue == RotationPreset.Automatic) {
						player.setTalentsString(eventID, newTalents);
						player.setSpecOptions(eventID, newOptions);
						player.setGlyphs(eventID, newGlyphs);
						// player.setBuffs(eventID, newIndividualBuffs);
						// raid?.setDebuffs(eventID, newDebuffs || Debuffs.create());
					}
					player.setRotation(eventID, newRotation);
				});
			},
		},
		{
			type: 'enum' as const,
			label: 'Curse',
			labelTooltip: 'Manual curse selection.',
			values: [
				{ name: "None", value: Curse.NoCurse },
				{ name: "Elements", value: Curse.Elements },
				{ name: "Weakness", value: Curse.Weakness },
				{ name: "Doom", value: Curse.Doom },
				{ name: "Agony", value: Curse.Agony },
				{ name: "Tongues", value: Curse.Tongues }
			],
			changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().curse,
			setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
				const newRotation = player.getRotation();
				newRotation.curse = newValue;
				newRotation.preset = RotationPreset.Manual;
				player.setRotation(eventID, newRotation);
			},
		},
		InputHelpers.makeRotationBooleanInput<Spec.SpecWarlock>({
			fieldName: 'detonateSeed',
			label: 'Detonate Seed on Cast',
			labelTooltip: 'Simulates raid doing damage to targets such that seed detonates immediately on cast.',
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Seed,
		}),
	],
};
