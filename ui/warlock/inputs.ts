import {
	Warlock_Options as WarlockOptions,
	Warlock_Rotation_Type as RotationType,
	Warlock_Rotation_Preset as RotationPreset,
	Warlock_Rotation_PrimarySpell as PrimarySpell,
	Warlock_Rotation_SecondaryDot as SecondaryDot,
	Warlock_Rotation_SpecSpell as SpecSpell,
	Warlock_Rotation_Curse as Curse,
	Warlock_Options_WeaponImbue as WarlockWeaponImbue,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
} from '/wotlk/core/proto/warlock.js';
import { RaidTarget, Spec, Glyphs } from '/wotlk/core/proto/common.js';
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';

import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import { IconEnumPicker, IconEnumPickerConfig, IconEnumValueConfig } from '/wotlk/core/components/icon_enum_picker.js';
import * as Presets from './presets.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

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

export const WeaponImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, WarlockWeaponImbue>({
	fieldName: 'weaponImbue',
	values: [
		{ color: 'grey', value: WarlockWeaponImbue.NoWeaponImbue },
		{ actionId: ActionId.fromItemId(41174), value: WarlockWeaponImbue.GrandFirestone },
		{ actionId: ActionId.fromItemId(41196), value: WarlockWeaponImbue.GrandSpellstone },
	],
});

export const PetType = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, Summon>({
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
		newRotation.corruption = newValue;
		newRotation.primarySpell = PrimarySpell.ShadowBolt;
		newRotation.preset = RotationPreset.Manual;
		player.setRotation(eventID, newRotation);
	},
};


export const WarlockRotationConfig = {
	inputs: [
		{
			type: 'enum' as const,

			label: 'Spec',
			labelTooltip: 'Switches between spec rotation settings. Will also update talents to defaults for the selected spec.',
			values: [
				{ name: 'Affliction', value: RotationType.Affliction },
				{ name: 'Demonology', value: RotationType.Demonology },
				{ name: 'Destruction', value: RotationType.Destruction },
			],
			changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
			getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().type,
			setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
				var newRotation = player.getRotation();
				var newOptions = player.getSpecOptions();
				TypedEvent.freezeAllAndDo(() => {
					if (newValue == RotationType.Affliction) {
						player.setTalentsString(eventID, Presets.AfflictionTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.AfflictionTalents.data.glyphs || Glyphs.create());
						newRotation = Presets.AfflictionRotation
						newOptions = Presets.AfflictionOptions
					} else if (newValue == RotationType.Demonology) {
						player.setTalentsString(eventID, Presets.DemonologyTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.DemonologyTalents.data.glyphs || Glyphs.create());
						newRotation = Presets.DemonologyRotation
						newOptions = Presets.DemonologyOptions
					} else {
						player.setTalentsString(eventID, Presets.DestructionTalents.data.talentsString);
						player.setGlyphs(eventID, Presets.DestructionTalents.data.glyphs || Glyphs.create());
						newRotation = Presets.DestructionRotation
						newOptions = Presets.DestructionOptions
					}
					newRotation.type = newValue;
					newRotation.preset = RotationPreset.Automatic;
					player.setRotation(eventID, newRotation);
					player.setSpecOptions(eventID, newOptions);
				});
			},
		},
		InputHelpers.makeRotationEnumInput<Spec.SpecWarlock, RotationPreset>({
			fieldName: 'preset',
			label: 'Rotation Preset',
			labelTooltip: 'Automatic will select the spells for you if you have the last talent in a one of the trees. Otherwise you will have to manually select the spells you want to cast.',
			values: [
				{ name: "Manual", value: RotationPreset.Manual },
				{ name: "Automatic", value: RotationPreset.Automatic },
			],
		}),
		{
			type: 'enum' as const,
			label: 'Curse',
			labelTooltip: 'Manual curse selection. Choice ignored for an Automatic Rotation.',
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
