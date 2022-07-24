import { Warlock_Rotation_Type as RotationType, Warlock_Rotation_Preset as RotationPreset, Warlock_Rotation_PrimarySpell as PrimarySpell, Warlock_Rotation_SecondaryDot as SecondaryDot, Warlock_Rotation_SpecSpell as SpecSpell, Warlock_Rotation_Curse as Curse, Warlock_Options_WeaponImbue as WeaponImbue, Warlock_Options_Armor as Armor, Warlock_Options_Summon as Summon } from '/wotlk/core/proto/warlock.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
export declare const ArmorInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, Armor>;
export declare const WeaponImbueInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, WeaponImbue>;
export declare const PetInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, Summon>;
export declare const PrimarySpellInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, PrimarySpell>;
export declare const SecondaryDotInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, SecondaryDot>;
export declare const SpecSpellInput: InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecWarlock>, SpecSpell>;
export declare const CorruptionSpell: {
    type: "icon";
    id: ActionId;
    states: number;
    extraCssClasses: string[];
    changedEvent: (player: Player<Spec.SpecWarlock>) => TypedEvent<void>;
    getValue: (player: Player<Spec.SpecWarlock>) => boolean;
    setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: boolean) => void;
};
export declare const WarlockRotationConfig: {
    inputs: (InputHelpers.TypedBooleanPickerConfig<Player<Spec.SpecWarlock>> | {
        type: "enum";
        label: string;
        labelTooltip: string;
        values: {
            name: string;
            value: RotationType;
        }[];
        changedEvent: (player: Player<Spec.SpecWarlock>) => TypedEvent<void>;
        getValue: (player: Player<Spec.SpecWarlock>) => RotationType;
        setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => void;
    } | {
        type: "enum";
        label: string;
        labelTooltip: string;
        values: {
            name: string;
            value: RotationPreset;
        }[];
        changedEvent: (player: Player<Spec.SpecWarlock>) => TypedEvent<void>;
        getValue: (player: Player<Spec.SpecWarlock>) => RotationPreset;
        setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => void;
    } | {
        type: "enum";
        label: string;
        labelTooltip: string;
        values: {
            name: string;
            value: Curse;
        }[];
        changedEvent: (player: Player<Spec.SpecWarlock>) => TypedEvent<void>;
        getValue: (player: Player<Spec.SpecWarlock>) => Curse;
        setValue: (eventID: EventID, player: Player<Spec.SpecWarlock>, newValue: number) => void;
    })[];
};
