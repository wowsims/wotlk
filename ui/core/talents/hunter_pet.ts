import { Spec } from '../proto/common.js';
import { HunterPetTalents, Hunter_Options_PetType as PetType } from '../proto/hunter.js';
import { Player } from '../player.js';
import { Component } from '../components/component.js';
import { EnumPicker, EnumPickerConfig } from '../components/enum_picker.js';
import { SavedDataManager } from '../components/saved_data_manager.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { ActionId } from '../proto_utils/action_id.js';

import { TalentsConfig, TalentsPicker, newTalentsConfig } from './talents_picker.js';
import { protoToTalentString, talentStringToProto } from './factory.js';

import * as InputHelpers from '../components/input_helpers.js';

export function makePetTypeInputConfig(includeLabel: boolean): InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecHunter>, PetType> {
	return InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, PetType>({
		fieldName: 'petType',
		numColumns: 5,
		//label: includeLabel ? 'Pet' : '',
		values: [
			{ color: 'grey', value: PetType.PetNone },
			{ actionId: ActionId.fromPetName('Bat'), tooltip: 'Bat', value: PetType.Bat },
			{ actionId: ActionId.fromPetName('Bear'), tooltip: 'Bear', value: PetType.Bear },
			{ actionId: ActionId.fromPetName('Bird of Prey'), tooltip: 'Bird of Prey', value: PetType.BirdOfPrey },
			{ actionId: ActionId.fromPetName('Boar'), tooltip: 'Boar', value: PetType.Boar },
			{ actionId: ActionId.fromPetName('Carrion Bird'), tooltip: 'Carrion Bird', value: PetType.CarrionBird },
			{ actionId: ActionId.fromPetName('Cat'), tooltip: 'Cat', value: PetType.Cat },
			{ actionId: ActionId.fromPetName('Chimaera'), tooltip: 'Chimaera (Exotic)', value: PetType.Chimaera },
			{ actionId: ActionId.fromPetName('Core Hound'), tooltip: 'Core Hound (Exotic)', value: PetType.CoreHound },
			{ actionId: ActionId.fromPetName('Crab'), tooltip: 'Crab', value: PetType.Crab },
			{ actionId: ActionId.fromPetName('Crocolisk'), tooltip: 'Crocolisk', value: PetType.Crocolisk },
			{ actionId: ActionId.fromPetName('Devilsaur'), tooltip: 'Devilsaur (Exotic)', value: PetType.Devilsaur },
			{ actionId: ActionId.fromPetName('Dragonhawk'), tooltip: 'Dragonhawk', value: PetType.Dragonhawk },
			{ actionId: ActionId.fromPetName('Gorilla'), tooltip: 'Gorilla', value: PetType.Gorilla },
			{ actionId: ActionId.fromPetName('Hyena'), tooltip: 'Hyena', value: PetType.Hyena },
			{ actionId: ActionId.fromPetName('Moth'), tooltip: 'Moth', value: PetType.Moth },
			{ actionId: ActionId.fromPetName('Nether Ray'), tooltip: 'Nether Ray', value: PetType.NetherRay },
			{ actionId: ActionId.fromPetName('Raptor'), tooltip: 'Raptor', value: PetType.Raptor },
			{ actionId: ActionId.fromPetName('Ravager'), tooltip: 'Ravager', value: PetType.Ravager },
			{ actionId: ActionId.fromPetName('Rhino'), tooltip: 'Rhino', value: PetType.Rhino },
			{ actionId: ActionId.fromPetName('Scorpid'), tooltip: 'Scorpid', value: PetType.Scorpid },
			{ actionId: ActionId.fromPetName('Serpent'), tooltip: 'Serpent', value: PetType.Serpent },
			{ actionId: ActionId.fromPetName('Silithid'), tooltip: 'Silithid (Exotic)', value: PetType.Silithid },
			{ actionId: ActionId.fromPetName('Spider'), tooltip: 'Spider', value: PetType.Spider },
			{ actionId: ActionId.fromPetName('Spirit Beast'), tooltip: 'Spirit Beast (Exotic)', value: PetType.SpiritBeast },
			{ actionId: ActionId.fromPetName('Spore Bat'), tooltip: 'Spore Bat', value: PetType.SporeBat },
			{ actionId: ActionId.fromPetName('Tallstrider'), tooltip: 'Tallstrider', value: PetType.Tallstrider },
			{ actionId: ActionId.fromPetName('Turtle'), tooltip: 'Turtle', value: PetType.Turtle },
			{ actionId: ActionId.fromPetName('Warp Stalker'), tooltip: 'Warp Stalker', value: PetType.WarpStalker },
			{ actionId: ActionId.fromPetName('Wasp'), tooltip: 'Wasp', value: PetType.Wasp },
			{ actionId: ActionId.fromPetName('Wind Serpent'), tooltip: 'Wind Serpent', value: PetType.WindSerpent },
			{ actionId: ActionId.fromPetName('Wolf'), tooltip: 'Wolf', value: PetType.Wolf },
			{ actionId: ActionId.fromPetName('Worm'), tooltip: 'Worm (Exotic)', value: PetType.Worm },
		],
	});
}

enum PetCategory {
    Cunning,
    Ferocity,
    Tenacity,
}

const petCategories: Record<PetType, PetCategory> = {
    [PetType.PetNone]: PetCategory.Ferocity,
    [PetType.Bat]: PetCategory.Cunning,
    [PetType.Bear]: PetCategory.Tenacity,
    [PetType.BirdOfPrey]: PetCategory.Cunning,
    [PetType.Boar]: PetCategory.Tenacity,
    [PetType.CarrionBird]: PetCategory.Ferocity,
    [PetType.Cat]: PetCategory.Ferocity,
    [PetType.Chimaera]: PetCategory.Cunning,
    [PetType.CoreHound]: PetCategory.Ferocity,
    [PetType.Crab]: PetCategory.Tenacity,
    [PetType.Crocolisk]: PetCategory.Tenacity,
    [PetType.Devilsaur]: PetCategory.Ferocity,
    [PetType.Dragonhawk]: PetCategory.Cunning,
    [PetType.Gorilla]: PetCategory.Tenacity,
    [PetType.Hyena]: PetCategory.Ferocity,
    [PetType.Moth]: PetCategory.Ferocity,
    [PetType.NetherRay]: PetCategory.Cunning,
    [PetType.Raptor]: PetCategory.Ferocity,
    [PetType.Ravager]: PetCategory.Cunning,
    [PetType.Rhino]: PetCategory.Tenacity,
    [PetType.Scorpid]: PetCategory.Tenacity,
    [PetType.Serpent]: PetCategory.Cunning,
    [PetType.Silithid]: PetCategory.Cunning,
    [PetType.Spider]: PetCategory.Cunning,
    [PetType.SpiritBeast]: PetCategory.Ferocity,
    [PetType.SporeBat]: PetCategory.Cunning,
    [PetType.Tallstrider]: PetCategory.Ferocity,
    [PetType.Turtle]: PetCategory.Tenacity,
    [PetType.WarpStalker]: PetCategory.Tenacity,
    [PetType.Wasp]: PetCategory.Ferocity,
    [PetType.WindSerpent]: PetCategory.Cunning,
    [PetType.Wolf]: PetCategory.Ferocity,
    [PetType.Worm]: PetCategory.Tenacity,
};

const categoryOrder = [PetCategory.Cunning, PetCategory.Ferocity, PetCategory.Tenacity];
const categoryClasses = ['cunning', 'ferocity', 'tenacity'];

export class HunterPetTalentsPicker extends Component {
    private readonly player: Player<Spec.SpecHunter>;
    private curCategory: PetCategory | null;
    private curTalents: HunterPetTalents;

    // Not saved to storage, just holds last-used values for this session.
    private savedSets: Array<HunterPetTalents>;

    constructor(parent: HTMLElement, player: Player<Spec.SpecHunter>) {
        super(parent, 'hunter-pet-talents-picker');
        this.player = player;

        this.rootElem.innerHTML = `
			<div class="pet-talents-container"></div>
		`;

        this.curCategory = this.getCategoryFromPlayer();
        this.curTalents = this.getPetTalentsFromPlayer();
        this.savedSets = defaultTalents.slice();
        this.savedSets[this.curCategory] = this.curTalents;
        this.rootElem.classList.add(categoryClasses[this.curCategory]);

        const talentsContainer = this.rootElem.getElementsByClassName('pet-talents-container')[0] as HTMLElement;

        const pickers = categoryOrder.map((category, i) => {
            const talentsConfig = petTalentsConfig[i];

            const pickerContainer = document.createElement('div');
            pickerContainer.classList.add('hunter-pet-talents-' + categoryClasses[i]);
            talentsContainer.appendChild(pickerContainer);

            const picker = new TalentsPicker(pickerContainer, player, {
                trees: talentsConfig,
                changedEvent: (player: Player<Spec.SpecHunter>) => player.specOptionsChangeEmitter,
                getValue: (player: Player<Spec.SpecHunter>) => protoToTalentString(this.getPetTalentsFromPlayer(), talentsConfig),
                setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: string) => {
                    const options = player.getSpecOptions();
                    options.petTalents = talentStringToProto(HunterPetTalents.create(), newValue, talentsConfig);
                    player.setSpecOptions(eventID, options);

                    this.savedSets[i] = options.petTalents;
                    this.curTalents = options.petTalents;
                },
                pointsPerRow: 3,
                maxPoints: 16,
            });

            const savedTalentsManager = new SavedDataManager<Player<Spec.SpecHunter>, string>(pickerContainer, this.player, {
                presetsOnly: true,
                label: 'Pet Talents',
                storageKey: '__NEVER_USED__',
                getData: (player: Player<Spec.SpecHunter>) => protoToTalentString(this.getPetTalentsFromPlayer(), talentsConfig),
                setData: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: string) => {
                    const options = player.getSpecOptions();
                    options.petTalents = talentStringToProto(HunterPetTalents.create(), newValue, talentsConfig);
                    player.setSpecOptions(eventID, options);

                    this.savedSets[i] = options.petTalents;
                    this.curTalents = options.petTalents;
                },
                changeEmitters: [this.player.specOptionsChangeEmitter],
                equals: (a: string, b: string) => a == b,
                toJson: (a: string) => a,
                fromJson: (obj: any) => '',
            });
            savedTalentsManager.addSavedData({
                name: 'Default',
                isPreset: true,
                data: protoToTalentString(defaultTalents[i], talentsConfig),
            });
            savedTalentsManager.addSavedData({
                name: 'BM',
                isPreset: true,
                data: protoToTalentString(defaultBMTalents[i], talentsConfig),
            });

            return picker;
        });

        player.specOptionsChangeEmitter.on(() => {
            const petCategory = this.getCategoryFromPlayer();
            const categoryIdx = categoryOrder.indexOf(petCategory);

            if (petCategory != this.curCategory) {
                this.curCategory = petCategory;
                this.rootElem.classList.remove(...categoryClasses);
                this.rootElem.classList.add(categoryClasses[categoryIdx]);

                const curTalents = this.getPetTalentsFromPlayer();
                if (!HunterPetTalents.equals(curTalents, this.curTalents)) {
                    // If the current talents have also changed, this was probably a load so we shouldn't switch sets.
                    this.curTalents = curTalents;
                    this.savedSets[this.curCategory] = this.curTalents;
                } else {
                    // Revert to the talents from last time the user was editing this category.
                    const options = this.player.getSpecOptions();
                    options.petTalents = this.savedSets[this.curCategory];
                    this.player.setSpecOptions(TypedEvent.nextEventID(), options);
                    this.curTalents = options.petTalents;
                }
            }
        });

        const updateIsBM = () => {
            const maxPoints = this.player.getTalents().beastMastery ? 20 : 16;
            pickers.forEach(picker => picker.setMaxPoints(maxPoints));
        };
        player.talentsChangeEmitter.on(updateIsBM);
        updateIsBM();
    }

    getPetTalentsFromPlayer(): HunterPetTalents {
        return this.player.getSpecOptions().petTalents || HunterPetTalents.create();
    }

    getCategoryFromPlayer(): PetCategory {
        const petType = this.player.getSpecOptions().petType;
        return petCategories[petType];
    }
}

export function getPetTalentsConfig(petType: PetType): TalentsConfig<HunterPetTalents> {
	const petCategory = petCategories[petType];
	const categoryIdx = categoryOrder.indexOf(petCategory);
	return petTalentsConfig[categoryIdx];
}

export const cunningDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    dive: true,
    boarsSpeed: true,
    mobility: 2,
    spikedCollar: 3,
    cornered: 2,
    feedingFrenzy: 2,
    wolverineBite: true,
    bullheaded: true,
    wildHunt: 1,
});
export const ferocityDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    dive: true,
    spikedCollar: 3,
    boarsSpeed: true,
    cullingTheHerd: 3,
    spidersBite: 3,
    rabid: true,
    callOfTheWild: true,
    wildHunt: 1,
});
export const tenacityDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    charge: true,
    greatStamina: 3,
    bloodOfTheRhino: 2,
    guardDog: 2,
    thunderstomp: true,
    graceOfTheMantis: 2,
    taunt: true,
    roarOfSacrifice: true,
    wildHunt: 1,
});
const defaultTalents = [cunningDefault, ferocityDefault, tenacityDefault];

export const cunningBMDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    dive: true,
    boarsSpeed: true,
    mobility: 2,
    spikedCollar: 3,
    cornered: 2,
    feedingFrenzy: 2,
    wolverineBite: true,
    bullheaded: true,
    graceOfTheMantis: 2,
    wildHunt: 2,
    roarOfSacrifice: true,
});
export const ferocityBMDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    dive: true,
    bloodthirsty: 1,
    spikedCollar: 3,
    boarsSpeed: true,
    cullingTheHerd: 3,
    spidersBite: 3,
    rabid: true,
    callOfTheWild: true,
    sharkAttack: 2,
    wildHunt: 2,
});
export const tenacityBMDefault: HunterPetTalents = HunterPetTalents.create({
    cobraReflexes: 2,
    charge: true,
    greatStamina: 3,
    spikedCollar: 3,
    bloodOfTheRhino: 2,
    guardDog: 2,
    thunderstomp: true,
    graceOfTheMantis: 2,
    taunt: true,
    roarOfSacrifice: true,
    wildHunt: 2,
});
const defaultBMTalents = [cunningBMDefault, ferocityBMDefault, tenacityBMDefault];

const cunningPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig([
    {
        name: 'Cunning',
        backgroundUrl: 'https://wotlkdb.com/static/images/wow/hunterpettalents/bg_3.jpg?85',
        talents: [
            {
                fieldName: 'cobraReflexes',
                location: {
                    rowIdx: 0,
                    colIdx: 0,
                },
                spellIds: [61682],
                maxPoints: 2,
            },
            {
                fieldName: 'dive',
                location: {
                    rowIdx: 0,
                    colIdx: 1,
                },
                spellIds: [23145],
                maxPoints: 1,
            },
            {
                fieldName: 'greatStamina',
                location: {
                    rowIdx: 0,
                    colIdx: 2,
                },
                spellIds: [61686],
                maxPoints: 3,
            },
            {
                fieldName: 'naturalArmor',
                location: {
                    rowIdx: 0,
                    colIdx: 3,
                },
                spellIds: [61689],
                maxPoints: 2,
            },
            {
                fieldName: 'boarsSpeed',
                location: {
                    rowIdx: 1,
                    colIdx: 0,
                },
                spellIds: [19596],
                maxPoints: 1,
            },
            {
                fieldName: 'mobility',
                location: {
                    rowIdx: 1,
                    colIdx: 1,
                },
                prereqLocation: {
                    rowIdx: 0,
                    colIdx: 1,
                },
                spellIds: [53554],
                maxPoints: 2,
            },
            {
                fieldName: 'owlsFocus',
                location: {
                    rowIdx: 1,
                    colIdx: 2,
                },
                spellIds: [53514, 53516],
                maxPoints: 2,
            },
            {
                fieldName: 'spikedCollar',
                location: {
                    rowIdx: 1,
                    colIdx: 3,
                },
                spellIds: [53182],
                maxPoints: 3,
            },
            {
                fieldName: 'cullingTheHerd',
                location: {
                    rowIdx: 2,
                    colIdx: 0,
                },
                spellIds: [61680, 61681, 52858],
                maxPoints: 3,
            },
            {
                fieldName: 'lionhearted',
                location: {
                    rowIdx: 2,
                    colIdx: 1,
                },
                spellIds: [53409, 53411],
                maxPoints: 2,
            },
            {
                fieldName: 'carrionFeeder',
                location: {
                    rowIdx: 2,
                    colIdx: 2,
                },
                spellIds: [54044],
                maxPoints: 1,
            },
            {
                fieldName: 'greatResistance',
                location: {
                    rowIdx: 3,
                    colIdx: 1,
                },
                spellIds: [53427, 53429],
                maxPoints: 3,
            },
            {
                fieldName: 'cornered',
                location: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [52234, 53497],
                maxPoints: 2,
            },
            {
                fieldName: 'feedingFrenzy',
                location: {
                    rowIdx: 3,
                    colIdx: 3,
                },
                prereqLocation: {
                    rowIdx: 1,
                    colIdx: 3,
                },
                spellIds: [53511],
                maxPoints: 2,
            },
            {
                fieldName: 'wolverineBite',
                location: {
                    rowIdx: 4,
                    colIdx: 0,
                },
                spellIds: [53508],
                maxPoints: 1,
            },
            {
                fieldName: 'roarOfRecovery',
                location: {
                    rowIdx: 4,
                    colIdx: 1,
                },
                spellIds: [53517],
                maxPoints: 1,
            },
            {
                fieldName: 'bullheaded',
                location: {
                    rowIdx: 4,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [53490],
                maxPoints: 1,
            },
            {
                fieldName: 'graceOfTheMantis',
                location: {
                    rowIdx: 4,
                    colIdx: 3,
                },
                spellIds: [53450],
                maxPoints: 2,
            },
            {
                fieldName: 'wildHunt',
                location: {
                    rowIdx: 5,
                    colIdx: 0,
                },
                prereqLocation: {
                    rowIdx: 4,
                    colIdx: 0,
                },
                spellIds: [62758, 62762],
                maxPoints: 2,
            },
            {
                fieldName: 'roarOfSacrifice',
                location: {
                    rowIdx: 5,
                    colIdx: 3,
                },
                prereqLocation: {
                    rowIdx: 4,
                    colIdx: 3,
                },
                spellIds: [53480],
                maxPoints: 1,
            },
        ],
    },
]);
const ferocityPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig([
    {
        name: 'Ferocity',
        backgroundUrl: 'https://wotlkdb.com/static/images/wow/hunterpettalents/bg_1.jpg?85',
        talents: [
            {
                fieldName: 'cobraReflexes',
                location: {
                    rowIdx: 0,
                    colIdx: 0,
                },
                spellIds: [61682],
                maxPoints: 2,
            },
            {
                fieldName: 'dive',
                location: {
                    rowIdx: 0,
                    colIdx: 1,
                },
                spellIds: [23145],
                maxPoints: 1,
            },
            {
                fieldName: 'greatStamina',
                location: {
                    rowIdx: 0,
                    colIdx: 2,
                },
                spellIds: [61686],
                maxPoints: 3,
            },
            {
                fieldName: 'naturalArmor',
                location: {
                    rowIdx: 0,
                    colIdx: 3,
                },
                spellIds: [61689],
                maxPoints: 2,
            },
            {
                fieldName: 'improvedCower',
                location: {
                    rowIdx: 1,
                    colIdx: 0,
                },
                spellIds: [53180],
                maxPoints: 2,
            },
            {
                fieldName: 'bloodthirsty',
                location: {
                    rowIdx: 1,
                    colIdx: 1,
                },
                spellIds: [53186],
                maxPoints: 2,
            },
            {
                fieldName: 'spikedCollar',
                location: {
                    rowIdx: 1,
                    colIdx: 2,
                },
                spellIds: [53182],
                maxPoints: 3,
            },
            {
                fieldName: 'boarsSpeed',
                location: {
                    rowIdx: 1,
                    colIdx: 3,
                },
                spellIds: [19596],
                maxPoints: 1,
            },
            {
                fieldName: 'cullingTheHerd',
                location: {
                    rowIdx: 2,
                    colIdx: 0,
                },
                spellIds: [61680, 61681, 52858],
                maxPoints: 3,
            },
            {
                fieldName: 'lionhearted',
                location: {
                    rowIdx: 2,
                    colIdx: 2,
                },
                spellIds: [53409, 53411],
                maxPoints: 2,
            },
            {
                fieldName: 'charge',
                location: {
                    rowIdx: 2,
                    colIdx: 3,
                },
                spellIds: [61685],
                maxPoints: 1,
            },
            {
                fieldName: 'heartOfThePheonix',
                location: {
                    rowIdx: 3,
                    colIdx: 1,
                },
                prereqLocation: {
                    rowIdx: 1,
                    colIdx: 1,
                },
                spellIds: [55709],
                maxPoints: 1,
            },
            {
                fieldName: 'spidersBite',
                location: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [53203],
                maxPoints: 3,
            },
            {
                fieldName: 'greatResistance',
                location: {
                    rowIdx: 3,
                    colIdx: 3,
                },
                spellIds: [53427, 53429],
                maxPoints: 3,
            },
            {
                fieldName: 'rabid',
                location: {
                    rowIdx: 4,
                    colIdx: 0,
                },
                spellIds: [53401],
                maxPoints: 1,
            },
            {
                fieldName: 'lickYourWounds',
                location: {
                    rowIdx: 4,
                    colIdx: 1,
                },
                prereqLocation: {
                    rowIdx: 3,
                    colIdx: 1,
                },
                spellIds: [53426],
                maxPoints: 1,
            },
            {
                fieldName: 'callOfTheWild',
                location: {
                    rowIdx: 4,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [53434],
                maxPoints: 1,
            },
            {
                fieldName: 'sharkAttack',
                location: {
                    rowIdx: 5,
                    colIdx: 0,
                },
                spellIds: [62759],
                maxPoints: 2,
            },
            {
                fieldName: 'wildHunt',
                location: {
                    rowIdx: 5,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 4,
                    colIdx: 2,
                },
                spellIds: [62758, 62762],
                maxPoints: 2,
            },
        ],
    },
]);
const tenacityPetTalentsConfig: TalentsConfig<HunterPetTalents> = newTalentsConfig([
    {
        name: 'Tenacity',
        backgroundUrl: 'https://wotlkdb.com/static/images/wow/hunterpettalents/bg_2.jpg?85',
        talents: [
            {
                fieldName: 'cobraReflexes',
                location: {
                    rowIdx: 0,
                    colIdx: 0,
                },
                spellIds: [61682],
                maxPoints: 2,
            },
            {
                fieldName: 'charge',
                location: {
                    rowIdx: 0,
                    colIdx: 1,
                },
                spellIds: [61685],
                maxPoints: 1,
            },
            {
                fieldName: 'greatStamina',
                location: {
                    rowIdx: 0,
                    colIdx: 2,
                },
                spellIds: [61686],
                maxPoints: 3,
            },
            {
                fieldName: 'naturalArmor',
                location: {
                    rowIdx: 0,
                    colIdx: 3,
                },
                spellIds: [61689],
                maxPoints: 2,
            },
            {
                fieldName: 'spikedCollar',
                location: {
                    rowIdx: 1,
                    colIdx: 0,
                },
                spellIds: [53182],
                maxPoints: 3,
            },
            {
                fieldName: 'boarsSpeed',
                location: {
                    rowIdx: 1,
                    colIdx: 1,
                },
                spellIds: [19596],
                maxPoints: 1,
            },
            {
                fieldName: 'bloodOfTheRhino',
                location: {
                    rowIdx: 1,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 0,
                    colIdx: 2,
                },
                spellIds: [53481],
                maxPoints: 2,
            },
            {
                fieldName: 'petBarding',
                location: {
                    rowIdx: 1,
                    colIdx: 3,
                },
                prereqLocation: {
                    rowIdx: 0,
                    colIdx: 3,
                },
                spellIds: [53175],
                maxPoints: 2,
            },
            {
                fieldName: 'cullingTheHerd',
                location: {
                    rowIdx: 2,
                    colIdx: 0,
                },
                spellIds: [61680, 61681, 52858],
                maxPoints: 3,
            },
            {
                fieldName: 'guardDog',
                location: {
                    rowIdx: 2,
                    colIdx: 1,
                },
                spellIds: [53178],
                maxPoints: 2,
            },
            {
                fieldName: 'lionhearted',
                location: {
                    rowIdx: 2,
                    colIdx: 2,
                },
                spellIds: [53409, 53411],
                maxPoints: 2,
            },
            {
                fieldName: 'thunderstomp',
                location: {
                    rowIdx: 2,
                    colIdx: 3,
                },
                spellIds: [63900],
                maxPoints: 1,
            },
            {
                fieldName: 'graceOfTheMantis',
                location: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [53450],
                maxPoints: 2,
            },
            {
                fieldName: 'greatResistance',
                location: {
                    rowIdx: 3,
                    colIdx: 3,
                },
                spellIds: [53427, 53429],
                maxPoints: 3,
            },
            {
                fieldName: 'lastStand',
                location: {
                    rowIdx: 4,
                    colIdx: 0,
                },
                spellIds: [53478],
                maxPoints: 1,
            },
            {
                fieldName: 'taunt',
                location: {
                    rowIdx: 4,
                    colIdx: 1,
                },
                prereqLocation: {
                    rowIdx: 2,
                    colIdx: 1,
                },
                spellIds: [53477],
                maxPoints: 1,
            },
            {
                fieldName: 'roarOfSacrifice',
                location: {
                    rowIdx: 4,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 3,
                    colIdx: 2,
                },
                spellIds: [53480],
                maxPoints: 1,
            },
            {
                fieldName: 'intervene',
                location: {
                    rowIdx: 4,
                    colIdx: 3,
                },
                spellIds: [53476],
                maxPoints: 1,
            },
            {
                fieldName: 'silverback',
                location: {
                    rowIdx: 5,
                    colIdx: 1,
                },
                spellIds: [62764],
                maxPoints: 2,
            },
            {
                fieldName: 'wildHunt',
                location: {
                    rowIdx: 5,
                    colIdx: 2,
                },
                prereqLocation: {
                    rowIdx: 4,
                    colIdx: 2,
                },
                spellIds: [62758, 62762],
                maxPoints: 2,
            },
        ],
    },
]);

const petTalentsConfig = [
    cunningPetTalentsConfig,
    ferocityPetTalentsConfig,
    tenacityPetTalentsConfig,
];
