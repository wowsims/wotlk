import { BooleanPicker } from '../core/components/boolean_picker.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../core/components/icon_enum_picker.js';
import { IconPickerConfig } from '../core/components/icon_picker.js';
import { makeWeaponImbueInput } from '../core/components/icon_inputs.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Options as ShamanOptions,
	ShamanTotems,
	ShamanShield,
    ShamanImbue
} from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
import { WeaponImbue } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecEnhancementShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ color: 'grey', value: ShamanShield.NoShield },
		{ actionId: ActionId.fromSpellId(33736), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
    fieldName: 'imbueMH',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
    fieldName: 'imbueOH',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});

export const DelayOffhandSwings = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecEnhancementShaman>({
	fieldName: 'delayOffhandSwings',
	label: 'Delay Offhand Swings',
	labelTooltip: 'Uses the startattack macro to delay OH swings, so they always follow within 0.5s of a MH swing.',
});

export const EnhancementShamanRotationConfig = {
    inputs: 
        [

    ],
};
