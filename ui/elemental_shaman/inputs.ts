import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { ElementalShaman_Options_ThunderstormRange, ElementalShaman_Rotation_BloodlustUse, ElementalShaman_Rotation_RotationType as RotationType, ShamanShield } from '../core/proto/shaman.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import { EventID } from 'ui/core/typed_event.js';
import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const InThunderstormRange = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecElementalShaman>({
	fieldName: 'thunderstormRange',
	// id: ActionId.fromSpellId(59159),
	label: "Thunderstorm In Range",
	labelTooltip: "When set to true, thunderstorm casts will cause damage.",
	getValue: (player: Player<Spec.SpecElementalShaman>) => player.getSpecOptions().thunderstormRange == ElementalShaman_Options_ThunderstormRange.TSInRange,
	setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		if (newValue) {
			newOptions.thunderstormRange = ElementalShaman_Options_ThunderstormRange.TSInRange;
		} else {
			newOptions.thunderstormRange = ElementalShaman_Options_ThunderstormRange.TSOutofRange;
		}
		player.setSpecOptions(eventID, newOptions);
	},
});

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecElementalShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ElementalShamanRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecElementalShaman, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{
					name: 'Adaptive', value: RotationType.Adaptive,
					tooltip: 'Dynamically adapts based on available mana to maximize CL casts without going OOM.',
				},
				{
					name: 'Manual', value: RotationType.Manual,
					tooltip: 'Allows custom selection of which spells to use and to modify cast conditions.',
				},
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'bloodlust',
			label: 'Use Bloodlust',
			labelTooltip: 'Player will cast bloodlust',
			getValue: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().bloodlust == ElementalShaman_Rotation_BloodlustUse.UseBloodlust,
			setValue: (eventID: EventID, player: Player<Spec.SpecElementalShaman>, newValue: boolean) => {
				const newRotation = player.getRotation();
				if (newValue) {
					newRotation.bloodlust = ElementalShaman_Rotation_BloodlustUse.UseBloodlust;
				} else {
					newRotation.bloodlust = ElementalShaman_Rotation_BloodlustUse.NoBloodlust;
				}
				player.setRotation(eventID, newRotation);
			},
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecElementalShaman>({
			fieldName: 'lvbFsWaitMs',
			label: 'Max wait for LvB/FS (ms)',
			labelTooltip: 'Amount of time the sim will wait if FS is about to fall off or LvB CD is about to come up. Setting to 0 will default to 175ms',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'useChainLightning',
			label: 'Use Chain Lightning in Rotation',
			labelTooltip: 'Use Chain Lightning in rotation',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'useClOnlyGap',
			label: 'Use CL only as gap filler',
			labelTooltip: 'Use CL to fill short gaps in LvB CD instead of on CD.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual && player.getRotation().useChainLightning,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecElementalShaman>({
			fieldName: 'clMinManaPer',
			label: 'Min mana percent to use Chain Lightning',
			labelTooltip: 'Customize minimum mana level to cast Chain Lightning. 0 will spam until OOM.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual && player.getRotation().useChainLightning,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'useFireNova',
			label: 'Use Fire Nova in Rotation',
			labelTooltip: 'Fire Nova will hit all targets when cast.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecElementalShaman>({
			fieldName: 'fnMinManaPer',
			label: 'Min mana percent to use FireNova',
			labelTooltip: 'Customize minimum mana level to cast Fire Nova. 0 will spam until OOM.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual && player.getRotation().useFireNova,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'overwriteFlameshock',
			label: 'Allow Flameshock to be overwritten',
			labelTooltip: 'Will use flameshock at the end of the duration even if its still ticking if there isn\'t enough time to cast lavaburst before expiring.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecElementalShaman>({
			fieldName: 'alwaysCritLvb',
			label: 'Only cast Lavaburst with FS',
			labelTooltip: 'Will only cast Lavaburst if Flameshock will be active when the cast finishes.',
			enableWhen: (player: Player<Spec.SpecElementalShaman>) => player.getRotation().type == RotationType.Manual,
		}),
	],
};
