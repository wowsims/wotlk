import { ShamanShield, ShamanImbue } from '/wotlk/core/proto/shaman.js';
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
export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput({
    fieldName: 'imbueMH',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});
export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput({
    fieldName: 'imbueOH',
    values: [
        { color: 'grey', value: ShamanImbue.NoImbue },
        { actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
        { actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon },
        { actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
    ],
});
export const DelayOffhandSwings = InputHelpers.makeSpecOptionsBooleanInput({
    fieldName: 'delayOffhandSwings',
    label: 'Delay Offhand Swings',
    labelTooltip: 'Uses the startattack macro to delay OH swings, so they always follow within 0.5s of a MH swing.',
});
export const EnhancementShamanRotationConfig = {
    inputs: [],
};
