import { StatMods } from '/wotlk/core/components/character_stats.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { EncounterPickerConfig } from '/wotlk/core/components/encounter_picker.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { EventID } from './typed_event.js';
import { Gear } from '/wotlk/core/proto_utils/gear.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { IndividualSimSettings } from '/wotlk/core/proto/ui.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { Player } from './player.js';
import { Profession } from '/wotlk/core/proto/common.js';
import { Race } from '/wotlk/core/proto/common.js';
import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { SavedDataConfig } from '/wotlk/core/components/saved_data_manager.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { SimUI, SimWarning } from './sim_ui.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { SpecOptions } from '/wotlk/core/proto_utils/utils.js';
import { SpecRotation } from '/wotlk/core/proto_utils/utils.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { StatWeightsResult } from '/wotlk/core/proto/api.js';
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
export declare type InputConfig<ModObject> = (InputHelpers.TypedBooleanPickerConfig<ModObject> | InputHelpers.TypedNumberPickerConfig<ModObject> | InputHelpers.TypedEnumPickerConfig<ModObject>);
export declare type IconInputConfig<ModObject, T> = (InputHelpers.TypedIconPickerConfig<ModObject, T> | InputHelpers.TypedIconEnumPickerConfig<ModObject, T>);
export interface InputSection {
    tooltip?: string;
    inputs: Array<InputConfig<Player<any>>>;
}
export interface IndividualSimUIConfig<SpecType extends Spec> {
    cssClass: string;
    knownIssues?: Array<string>;
    warnings?: Array<(simUI: IndividualSimUI<SpecType>) => SimWarning>;
    epStats: Array<Stat>;
    buffStats?: Array<Stat>;
    epReferenceStat: Stat;
    displayStats: Array<Stat>;
    modifyDisplayStats?: (player: Player<SpecType>) => StatMods;
    defaults: {
        gear: EquipmentSpec;
        epWeights: Stats;
        consumes: Consumes;
        rotation: SpecRotation<SpecType>;
        talents: SavedTalents;
        specOptions: SpecOptions<SpecType>;
        raidBuffs: RaidBuffs;
        partyBuffs: PartyBuffs;
        individualBuffs: IndividualBuffs;
        debuffs: Debuffs;
    };
    playerIconInputs: Array<IconInputConfig<Player<SpecType>, any>>;
    petConsumeInputs?: Array<IconInputConfig<Player<SpecType>, any>>;
    rotationInputs: InputSection;
    rotationIconInputs?: Array<IconInputConfig<Player<any>, any>>;
    otherInputs?: InputSection;
    additionalSections?: Record<string, InputSection>;
    additionalIconSections?: Record<string, Array<IconInputConfig<Player<any>, any>>>;
    customSections?: Array<(simUI: IndividualSimUI<SpecType>, parentElem: HTMLElement) => string>;
    encounterPicker: EncounterPickerConfig;
    presets: {
        gear: Array<PresetGear>;
        talents: Array<SavedDataConfig<Player<any>, SavedTalents>>;
        rotation?: Array<SavedDataConfig<Player<any>, string>>;
    };
}
export interface GearAndStats {
    gear: Gear;
    bonusStats?: Stats;
}
export interface PresetGear {
    name: string;
    gear: EquipmentSpec;
    tooltip?: string;
    enableWhen?: (obj: Player<any>) => boolean;
}
export interface Settings {
    raidBuffs: RaidBuffs;
    partyBuffs: PartyBuffs;
    individualBuffs: IndividualBuffs;
    consumes: Consumes;
    race: Race;
    professions?: Array<Profession>;
}
export declare abstract class IndividualSimUI<SpecType extends Spec> extends SimUI {
    readonly player: Player<SpecType>;
    readonly individualConfig: IndividualSimUIConfig<SpecType>;
    private raidSimResultsManager;
    private settingsMuuri;
    prevEpIterations: number;
    prevEpSimResult: StatWeightsResult | null;
    constructor(parentElem: HTMLElement, player: Player<SpecType>, config: IndividualSimUIConfig<SpecType>);
    private loadSettings;
    private addSidebarComponents;
    private addTopbarComponents;
    private addGearTab;
    private addSettingsTab;
    private addTalentsTab;
    private addDetailedResultsTab;
    private addLogTab;
    applyDefaults(eventID: EventID): void;
    getSavedGearStorageKey(): string;
    getSavedRotationStorageKey(): string;
    getSavedSettingsStorageKey(): string;
    getSavedTalentsStorageKey(): string;
    private recomputeSettingsLayout;
    getStorageKey(keyPart: string): string;
    toProto(): IndividualSimSettings;
    toLink(): string;
    fromProto(eventID: EventID, settings: IndividualSimSettings): void;
    splitRelevantOptions<T>(options: Array<StatOption<T>>): Array<T>;
}
export interface StatOption<T> {
    stats: Array<Stat>;
    item: T;
}
