import { ElementalShaman_Rotation_RotationType as RotationType, ShamanShield } from '/wotlk/core/proto/shaman.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput({
    fieldName: 'bloodlust',
    id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput({
    fieldName: 'shield',
    values: [
        { color: 'grey', value: ShamanShield.NoShield },
        { actionId: ActionId.fromSpellId(33736), value: ShamanShield.WaterShield },
        { actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
    ],
});
export const ElementalShamanRotationConfig = {
    inputs: [
        InputHelpers.makeRotationEnumInput({
            fieldName: 'type',
            label: 'Type',
            values: [
                {
                    name: 'Adaptive', value: RotationType.Adaptive,
                    tooltip: 'Dynamically adapts based on available mana to maximize CL casts without going OOM.',
                },
            ],
        }),
        InputHelpers.makeRotationBooleanInput({
            fieldName: 'inThunderstormRange',
            label: 'In Thunderstorm Range',
            labelTooltip: 'Thunderstorm will hit all targets when cast. Ignores knockback.',
            enableWhen: (player) => player.getTalents().thunderstorm,
        }),
    ],
};
