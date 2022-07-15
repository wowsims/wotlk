import { Warlock_Options as WarlockOptions, Warlock_Rotation_PrimarySpell as PrimarySpell, Warlock_Rotation_SecondaryDot as SecondaryDot, Warlock_Rotation_SpecSpell as SpecSpell, Warlock_Rotation_Curse as Curse, Warlock_Options_Armor as Armor, Warlock_Options_Summon as Summon } from '/wotlk/core/proto/warlock.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
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
// import { IndividualSimIconPickerConfig } from '/wotlk/core/components/icon_inputs.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

// function makeBooleanPersonalBuffInput(id: ActionId, optionsFieldName: keyof Armor): IconPickerConfig<Player<any>, boolean> {
// 	return {
// 		id: id,
// 		states: 2,
// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.specOptionsChangeEmitter,
// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getSpecOptions()[optionsFieldName] as boolean,
// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
// 			const newOptions = player.getSpecOptions();
// 			(newOptions[optionsFieldName] as boolean) = newValue;
// 			player.setSpecOptions(eventID, newOptions);
// 		},
// 	}
// }

// export const FelArmor = makeBooleanPersonalBuffInput(ActionId.fromSpellId(47893), 'FelArmor')

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
		player.setRotation(eventID, newRotation);
	},
};

export const SecondaryDotUnstableAffliction= {
	id: ActionId.fromSpellId(47843),
	states: 2,
	extraCssClasses: [
		'UnstableAffliction-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter || player.specOptionsChangeEmitter || player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().secondaryDot == SecondaryDot.UnstableAffliction,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.secondaryDot = newValue ? SecondaryDot.UnstableAffliction : SecondaryDot.NoSecondaryDot;
		player.setRotation(eventID, newRotation);
	},
	enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().unstableAffliction
};

export const SpecSpellChaosBolt = {
	id: ActionId.fromSpellId(59172),
	states: 2,
	extraCssClasses: [
		'ChaosBolt-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().specSpell == SpecSpell.ChaosBolt,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.specSpell = newValue ? SpecSpell.ChaosBolt : SpecSpell.NoSpecSpell;
		player.setRotation(eventID, newRotation);
	},
	enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().chaosBolt
};

export const SpecSpellHaunt = {
	id: ActionId.fromSpellId(59164),
	states: 2,
	extraCssClasses: [
		'Haunt-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().specSpell == SpecSpell.Haunt,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.specSpell = newValue ? SpecSpell.Haunt : SpecSpell.NoSpecSpell;
		player.setRotation(eventID, newRotation);
	},
	enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().haunt
};

export const CorruptionSpell = {
	id: ActionId.fromSpellId(47813),
	states: 2,
	extraCssClasses: [
		'Corruption-picker',
	],
	changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
	getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().corruption,
	setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		const newRotation = player.getRotation();
		newRotation.corruption = newValue;
		player.setRotation(eventID, newRotation);
	},
	enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell != PrimarySpell.Seed,
};


export const WarlockRotationConfig = {
	inputs: [
		// {
		// 	type: 'enum' as const,
		// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		// 	config: {
		// 		extraCssClasses: [
		// 			'primary-spell-enum-picker',
		// 		],
		// 		label: 'Primary Spell',
		// 		labelTooltip: 'Choose primary spell to cast',
		// 		values: [
		// 			{
		// 				name: 'Shadowbolt', value: PrimarySpell.Shadowbolt,
		// 			},
		// 			{
		// 				name: 'Incinerate', value: PrimarySpell.Incinerate,
		// 			},
		// 			{
		// 				name: 'Seed of Corruption', value: PrimarySpell.Seed,
		// 			},
		// 		],
		// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
		// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell,
		// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => {
		// 			const newRotation = player.getRotation();
		// 			newRotation.primarySpell = newValue;
		// 			player.setRotation(eventID, newRotation);
		// 		},
		// 	},
		// },
		// {
		// 	type: 'boolean' as const,
		// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		// 	config: {
		// 		id: ActionId.fromSpellId(47811),
		// 		states: 2,
		// 		extraCssClasses: [
		// 			'immolate-picker',
		// 		],
		// 		label: 'Use Immolate',
		// 		labelTooltip: 'Use Immolate as the next cast after the dot expires.',
		// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
		// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().immolate,
		// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		// 			const newRotation = player.getRotation();
		// 			newRotation.immolate = newValue;
		// 			player.setRotation(eventID, newRotation);
		// 		},
		// 	},
		// },
		// {
		// 	type: 'boolean' as const,
		// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		// 	config: {
		// 		extraCssClasses: [
		// 			'corruption-picker',
		// 		],
		// 		label: 'Use Corruption',
		// 		labelTooltip: 'Use Corruption as the next cast after the dot expires.',
		// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.rotationChangeEmitter,
		// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().corruption,
		// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		// 			const newRotation = player.getRotation();
		// 			newRotation.corruption = newValue;
		// 			player.setRotation(eventID, newRotation);
		// 		},
		// 	},
		// },
		{
			type: 'boolean' as const,
			getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
			config: {
				extraCssClasses: [
					'unstableaffliction-picker',
				],
				label: 'Use Unstable Affliction',
				labelTooltip: 'Use Unstable Affliction as the next cast after the dot expires.',
				changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter,
				getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().secondaryDot == SecondaryDot.UnstableAffliction,
				setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
					const newRotation = player.getRotation();
					newRotation.secondaryDot = newValue ? SecondaryDot.UnstableAffliction : SecondaryDot.NoSecondaryDot;
					player.setRotation(eventID, newRotation);
				},
				enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().unstableAffliction,
			},
		},
		// {
		// 	type: 'boolean' as const,
		// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		// 	config: {
		// 		extraCssClasses: [
		// 			'haunt-picker',
		// 		],
		// 		label: 'Use Haunt',
		// 		labelTooltip: 'Use Haunt as the next cast after the buff expires.',
		// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter,
		// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().haunt,
		// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		// 			const newRotation = player.getRotation();
		// 			newRotation.haunt = newValue;
		// 			player.setRotation(eventID, newRotation);
		// 		},
		// 		enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().haunt,
		// 	},
		// },
		// {
		// 	type: 'boolean' as const,
		// 	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
		// 	config: {
		// 		extraCssClasses: [
		// 			'ChaosBolt-picker',
		// 		],
		// 		label: 'Use Chaos Bolt',
		// 		labelTooltip: 'Use Chaos Bolt as the next cast when CD is up.',
		// 		changedEvent: (player: Player<Spec.SpecWarlock>) => player.talentsChangeEmitter,
		// 		getValue: (player: Player<Spec.SpecWarlock>) => player.getRotation().chaosBolt,
		// 		setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => {
		// 			const newRotation = player.getRotation();
		// 			newRotation.chaosBolt = newValue;
		// 			player.setRotation(eventID, newRotation);
		// 		},
		// 		enableWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().chaosBolt,
		// 	},
		// },
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
				enableWhen: (player: Player<Spec.SpecWarlock>) => player.getRotation().primarySpell == PrimarySpell.Seed,
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
				labelTooltip: 'No tooltip yet',
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
					player.setRotation(eventID, newRotation);
				},
			},
		},
	],
};
