import { Warlock_Options as WarlockOptions,Warlock_Rotation_Type as RotationType, Warlock_Rotation_Preset as RotationPreset, Warlock_Rotation_PrimarySpell as PrimarySpell, Warlock_Rotation_SecondaryDot as SecondaryDot, Warlock_Rotation_SpecSpell as SpecSpell, Warlock_Rotation_Curse as Curse, Warlock_Options_Armor as Armor, Warlock_Options_Summon as Summon } from '/wotlk/core/proto/warlock.js';
import { RaidTarget, Spec, Glyphs } from '/wotlk/core/proto/common.js';
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';

import { IndividualSimIconPickerConfig } from '/wotlk/core/individual_sim_ui.js';
import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { IconEnumPicker, IconEnumPickerConfig, IconEnumValueConfig } from '/wotlk/core/components/icon_enum_picker.js';
import * as Presets from './presets.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const FelArmor = {
	id: ActionId.fromSpellId(47893),
	states: 2,
	extraCssClasses: [
		'fel-armor-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().armor == Armor.FelArmor,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.armor = newValue ? Armor.FelArmor : Armor.NoArmor;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const DemonArmor = {
	id: ActionId.fromSpellId(47889),
	states: 2,
	extraCssClasses: [
		'demon-armor-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().armor == Armor.DemonArmor,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.armor = newValue ? Armor.DemonArmor : Armor.NoArmor;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const SummonImp = {
	id: ActionId.fromSpellId(688),
	states: 2,
	extraCssClasses: [
		'SummonImp-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().summon == Summon.Imp,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.summon = newValue ? Summon.Imp : Summon.NoSummon;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const SummonSuccubus = {
	id: ActionId.fromSpellId(712),
	states: 2,
	extraCssClasses: [
		'SummonSuccubus-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().summon == Summon.Succubus,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.summon = newValue ? Summon.Succubus : Summon.NoSummon;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const SummonFelhunter = {
	id: ActionId.fromSpellId(691),
	states: 2,
	extraCssClasses: [
		'SummonFelhunter-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().summon == Summon.Felhunter,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.summon = newValue ? Summon.Felhunter : Summon.NoSummon;
		player.setSpecOptions(eventID, newOptions);
	},
};

export const SummonFelguard = {
	id: ActionId.fromSpellId(30146),
	states: 2,
	extraCssClasses: [
		'SummonFelguard-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions().summon == Summon.Felguard,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.summon = newValue ? Summon.Felguard : Summon.NoSummon;
		player.setSpecOptions(eventID, newOptions);
	},
	showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().summonFelguard,
};



export const PrimarySpellShadowbolt = {
	id: ActionId.fromSpellId(47809),
	states: 2,
	extraCssClasses: [
		'Shadowbolt-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Shadowbolt,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.primarySpell = newValue ? PrimarySpell.Shadowbolt : PrimarySpell.Shadowbolt;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};

export const PrimarySpellIncinerate = {
	id: ActionId.fromSpellId(47838),
	states: 2,
	extraCssClasses: [
		'Incinerate-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Incinerate,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.primarySpell = newValue ? PrimarySpell.Incinerate : PrimarySpell.Shadowbolt;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};

export const PrimarySpellSeed = {
	id: ActionId.fromSpellId(47836),
	states: 2,
	extraCssClasses: [
		'Seed-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Seed,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.primarySpell = newValue ? PrimarySpell.Seed : PrimarySpell.Shadowbolt;
		newRotation.preset = RotationPreset.Manual;
		newRotation.corruption = false;
		player.setRotation(eventID, newRotation);
	},
};

export const SecondaryDotImmolate = {
	id: ActionId.fromSpellId(47811),
	states: 2,
	extraCssClasses: [
		'Immolate-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().secondaryDot == SecondaryDot.Immolate,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.secondaryDot = newValue ? SecondaryDot.Immolate : SecondaryDot.NoSecondaryDot;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};

export const SecondaryDotUnstableAffliction= {
	id: ActionId.fromSpellId(47843),
	states: 2,
	extraCssClasses: [
		'UnstableAffliction-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().secondaryDot == SecondaryDot.UnstableAffliction,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.secondaryDot = newValue ? SecondaryDot.UnstableAffliction : SecondaryDot.NoSecondaryDot;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
	showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().unstableAffliction,
};

export const SpecSpellChaosBolt = {
	id: ActionId.fromSpellId(59172),
	states: 2,
	extraCssClasses: [
		'ChaosBolt-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().specSpell == SpecSpell.ChaosBolt,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.specSpell = newValue ? SpecSpell.ChaosBolt : SpecSpell.NoSpecSpell;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
	showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().chaosBolt,
};

export const SpecSpellHaunt = {
	id: ActionId.fromSpellId(59164),
	states: 2,
	extraCssClasses: [
		'Haunt-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().specSpell == SpecSpell.Haunt,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.specSpell = newValue ? SpecSpell.Haunt : SpecSpell.NoSpecSpell;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
	showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().haunt,
};

export const CorruptionSpell = {
	id: ActionId.fromSpellId(47813),
	states: 2,
	extraCssClasses: [
		'Corruption-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().corruption,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.corruption = newValue;
		newRotation.primarySpell = PrimarySpell.Shadowbolt;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};


export const WarlockRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'rotation-type-enum-picker',
				],
				label: 'Spec',
				labelTooltip: 'Switches between spec rotation settings. Will also update talents to defaults for the selected spec.',
				values: [
					{
						name: 'Affliction', value: RotationType.Affliction,
					},
					{
						name: 'Demonology', value: RotationType.Demonology,
					},
					{
						name: 'Destruction', value: RotationType.Destruction,
					},
				],
				changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().type,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
					const newRotation = player.getRotation();
					const newOptions = player.getSpecOptions();
					newRotation.type = newValue;
					newRotation.preset = RotationPreset.Automatic;
					TypedEvent.freezeAllAndDo(() => {
						if (newRotation.type == RotationType.Affliction) {
							player.setTalentsString(eventID, Presets.AfflictionTalents.data.talentsString);
							player.setGlyphs(eventID, Presets.AfflictionTalents.data.glyphs || Glyphs.create());
							newRotation.primarySpell = Presets.AfflictionRotation.primarySpell
							newRotation.secondaryDot = Presets.AfflictionRotation.secondaryDot
							newRotation.specSpell = Presets.AfflictionRotation.specSpell
							newRotation.curse = Presets.AfflictionRotation.curse
							newRotation.corruption = Presets.AfflictionRotation.corruption
							newRotation.detonateSeed = Presets.AfflictionRotation.detonateSeed
							newOptions.summon = Presets.AfflictionOptions.summon
							newOptions.armor = Presets.AfflictionOptions.armor
						} else if (newRotation.type == RotationType.Demonology) {
							player.setTalentsString(eventID, Presets.DemonologyTalents.data.talentsString);
							player.setGlyphs(eventID, Presets.DemonologyTalents.data.glyphs || Glyphs.create());
							newRotation.primarySpell = Presets.DemonologyRotation.primarySpell
							newRotation.secondaryDot = Presets.DemonologyRotation.secondaryDot
							newRotation.specSpell = Presets.DemonologyRotation.specSpell
							newRotation.curse = Presets.DemonologyRotation.curse
							newRotation.corruption = Presets.DemonologyRotation.corruption
							newRotation.detonateSeed = Presets.DemonologyRotation.detonateSeed
							newOptions.summon = Presets.DemonologyOptions.summon
							newOptions.armor = Presets.DemonologyOptions.armor
						} else {
							player.setTalentsString(eventID, Presets.DestructionTalents.data.talentsString);
							player.setGlyphs(eventID, Presets.DestructionTalents.data.glyphs || Glyphs.create());
							newRotation.primarySpell = Presets.DestructionRotation.primarySpell
							newRotation.secondaryDot = Presets.DestructionRotation.secondaryDot
							newRotation.specSpell = Presets.DestructionRotation.specSpell
							newRotation.curse = Presets.DestructionRotation.curse
							newRotation.corruption = Presets.DestructionRotation.corruption
							newRotation.detonateSeed = Presets.DestructionRotation.detonateSeed
							newOptions.summon = Presets.DestructionOptions.summon
							newOptions.armor = Presets.DestructionOptions.armor
						}
						player.setRotation(eventID, newRotation);
						player.setSpecOptions(eventID, newOptions);
					});
				},
			},
		},
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'rotation-preset-enum-picker',
				],
				label: 'Rotation Preset',
				labelTooltip: 'Automatic will select the spells for you if you have the last talent in a one of the trees. Otherwise you will have to manually select the spells you want to cast.',
				values: [
					{
						name: "Manual", value: RotationPreset.Manual,
					},
					{
						name: "Automatic", value: RotationPreset.Automatic,
					},
				],
				changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().preset,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
					const newRotation = player.getRotation();
					newRotation.preset = newValue;
					player.setRotation(eventID, newRotation);
				},
			},
		},
		{
			type: 'enum' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'curse-enum-picker',
				],
				label: 'Curse',
				labelTooltip: 'Manual curse selection. Choice ignored for an Automatic Rotation.',
				values: [
					{
						name: "None", value: Curse.NoCurse,
					},
					{
						name: "Elements", value: Curse.Elements,
					},
					{
						name: "Weakness", value: Curse.Weakness,
					},
					{
						name: "Doom", value: Curse.Doom,
					},
					{
						name: "Agony", value: Curse.Agony,
					},
					{
						name: "Tongues", value: Curse.Tongues,
					}
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
		},
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'detonate-seed-picker',
				],
				label: 'Detonate Seed on Cast',
				labelTooltip: 'Simulates raid doing damage to targets such that seed detonates immediately on cast.',
				changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
				getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().detonateSeed,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.detonateSeed = newValue;
					player.setRotation(eventID, newRotation);
				},
				showWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Seed,
			},
		},
	],
};
