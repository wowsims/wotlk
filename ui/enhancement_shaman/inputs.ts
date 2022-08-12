import { BooleanPicker } from '../core/components/boolean_picker.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../core/components/icon_enum_picker.js';
import { IconPickerConfig } from '../core/components/icon_picker.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Options as ShamanOptions,
	ShamanTotems,
	ShamanShield,
    ShamanImbue,
    ShamanSyncType
} from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
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
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
    fieldName: 'imbueMh',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
    fieldName: 'imbueOh',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman, ShamanSyncType>({
	fieldName: 'syncType',
	label: 'Sync/Stagger Setting',
	labelTooltip: 'Choose your sync or stagger option, Perfect Sync makes your weapons always attack at the same time, which is ideal for mixed imbues. Delayed Offhand is similar but additionally adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window, ideal for matched imbues.',
    values: [
        { name: 'None', value: ShamanSyncType.NoSync },
        { name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
        { name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
    ],
});

export const EnhancementShamanRotationConfig = {
    inputs:
        [
        InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
            fieldName: 'lavaburstWeave',
            label: 'Enable Weaving Lava Burst',
            labelTooltip: 'Not particularily useful for dual wield, mostly a 2h option',
            enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().maelstromWeapon > 0,
        }),
        InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
            fieldName: 'lightningboltWeave',
            label: 'Enable Weaving Lightning Bolt',
            labelTooltip: 'Will provide a DPS increase, but is harder to execute',
            enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().maelstromWeapon > 0,
        }),
        InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, number>({
            fieldName: 'maelstromweaponMinStack',
            label: 'Min number of MW for weaving',
            labelTooltip: 'Minimum number of Maelstrom Weapon stacks to consider weaving',
            values: [
                { name: '1', value: 1 },
                { name: '2', value: 2 },
                { name: '3', value: 3 },
                { name: '4', value: 4 },
                { name: '5', value: 5 },
            ],
            enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().maelstromWeapon > 0,
         }),
    ],
};
