import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import { Hunter_Options_PetType as PetType } from '../proto/hunter.js';
import { ActionId } from '../proto_utils/action_id.js';


import * as InputHelpers from '../components/input_helpers.js';

export function makePetTypeInputConfig(_: boolean): InputHelpers.TypedIconEnumPickerConfig<Player<Spec.SpecHunter>, PetType> {
	return InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, PetType>({
		fieldName: 'petType',
		numColumns: 5,
		//label: includeLabel ? 'Pet' : '',
		values: [
			{ value: PetType.PetNone, tooltip: 'No Pet' },
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
