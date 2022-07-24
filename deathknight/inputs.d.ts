import { Spec } from '/wotlk/core/proto/common.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
export declare const StartingRunicPower: InputHelpers.TypedNumberPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const PetUptime: InputHelpers.TypedNumberPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const PrecastGhoulFrenzy: InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const PrecastHornOfWinter: InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const RefreshHornOfWinter: InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const DiseaseRefreshDuration: InputHelpers.TypedNumberPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const UseDeathAndDecay: InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const UseArmyOfTheDead: InputHelpers.TypedEnumPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const UnholyPresenceOpener: InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>;
export declare const DeathKnightRotationConfig: {
    inputs: (InputHelpers.TypedNumberPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>> | InputHelpers.TypedBooleanPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>> | InputHelpers.TypedEnumPickerConfig<import("/wotlk/core/player").Player<Spec.SpecDeathKnight>>)[];
};
