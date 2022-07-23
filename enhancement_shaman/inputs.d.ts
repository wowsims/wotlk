import { ShamanShield, ShamanImbue } from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
export declare const Bloodlust: InputHelpers.TypedIconPickerConfig<Player<Spec.SpecEnhancementShaman>, boolean>;
export declare const ShamanShieldInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecEnhancementShaman>, ShamanShield>;
export declare const ShamanImbueMH: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecEnhancementShaman>, ShamanImbue>;
export declare const ShamanImbueOH: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecEnhancementShaman>, ShamanImbue>;
export declare const DelayOffhandSwings: InputHelpers.TypedBooleanPickerConfig<Player<Spec.SpecEnhancementShaman>>;
export declare const EnhancementShamanRotationConfig: {
    inputs: never[];
};
