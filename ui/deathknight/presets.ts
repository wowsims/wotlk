import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { Player } from '/wotlk/core/player.js';

import {
	DeathKnightTalents as DeathKnightTalents,
	DeathKnight,
	DeathKnight_Rotation as DeathKnightRotation,
	DeathKnight_Options as DeathKnightOptions,
} from '/wotlk/core/proto/deathknight.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost Dps',
	data: '23050005-32005350352203012300033101351',
};

export const DefaultRotation = DeathKnightRotation.create({
	useScourgeStrike: false,
});

export const DefaultOptions = DeathKnightOptions.create({
	dualWhield: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfRelentlessAssault,
	food: Food.FoodRoastedClefthoof,
	defaultPotion: Potions.HastePotion,
	mainHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
	offHandImbue: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
});

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost BiS Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
            "name": "Spiked Titansteel Helm",
            "id": 41386,
            "enchant": {
                "name": "Arcanum of Torment",
                "id": 3817,
                "itemId": 50367
            },
            "gems": [
                {
                    "name": "Relentless Earthsiege Diamond",
                    "id": 41398
                },
                {
                    "name": "Nightmare Tear",
                    "id": 49110
                }
            ],
            "slot": "HEAD"
        },
        {
            "name": "Titanium Impact Choker",
            "id": 42645,
            "gems": [
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "NECK"
        },
        {
            "name": "Pauldrons of Berserking",
            "id": 34388,
            "enchant": {
                "name": "Greater Inscription of the Axe",
                "id": 3808,
                "itemId": 50335
            },
            "gems": [
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                },
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "SHOULDERS"
        },
        {
            "name": "Heroes' Scourgeborne Battleplate",
            "id": 39617,
            "enchant": {
                "name": "Enchant Chest - Powerful Stats",
                "id": 3832,
                "spellId": 60692
            },
            "gems": [
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                },
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "CHEST"
        },
        {
            "name": "Jorach's Crocolisk Skin Belt",
            "id": 40694,
            "gems": [
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                },
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "WAIST"
        },
        {
            "name": "Staggering Legplates",
            "id": 37193,
            "enchant": {
                "name": "Icescale Leg Armor",
                "id": 3823,
                "itemId": 38374
            },
            "gems": [
                {
                    "name": "Bold Dragon's Eye",
                    "id": 42142
                },
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "LEGS"
        },
        {
            "name": "Death-Inured Sabatons",
            "id": 44306,
            "enchant": {
                "name": "Nitro Boosts",
                "id": 3606,
                "itemId": 41118
            },
            "gems": [
                {
                    "name": "Bold Dragon's Eye",
                    "id": 42142
                },
                {
                    "name": "Bold Dragon's Eye",
                    "id": 42142
                }
            ],
            "slot": "FEET"
        },
        {
            "name": "Vengeance Bindings",
            "id": 41355,
            "enchant": {
                "name": "Enchant Bracers - Greater Assault",
                "id": 3845,
                "spellId": 44575
            },
            "slot": "WRISTS"
        },
        {
            "name": "Heroes' Scourgeborne Gauntlets",
            "id": 39618,
            "enchant": {
                "name": "Hyperspeed Accelerators",
                "id": 3604,
                "spellId": 54999
            },
            "gems": [
                {
                    "name": "Bold Scarlet Ruby",
                    "id": 39996
                }
            ],
            "slot": "HANDS"
        },
        {
            "name": "Hemorrhaging Circle",
            "id": 37642,
            "slot": "FINGER_1"
        },
        {
            "name": "Ring of the Kirin Tor",
            "id": 44935,
            "slot": "FINGER_2"
        },
        {
            "name": "Mirror of Truth",
            "id": 40684,
            "slot": "TRINKET_1"
        },
        {
            "name": "Darkmoon Card: Greatness",
            "id": 42987,
            "slot": "TRINKET_2"
        },
        {
            "name": "Cloak of Bloodied Waters",
            "id": 37647,
            "enchant": {
                "name": "Flexweave Underlay",
                "id": 3605,
                "itemId": 41111
            },
            "slot": "BACK"
        },
        {
            "name": "Titansteel Bonecrusher",
            "id": 41383,
            "enchant": {
                "name": "Rune of Razorice",
                "id": 3370,
                "spellId": 53343
            },
            "slot": "MAIN_HAND"
        },
        {
            "name": "Krol Cleaver",
            "id": 43611,
            "enchant": {
                "name": "Rune of the Fallen Crusader",
                "id": 3368,
                "spellId": 53344
            },
            "slot": "OFF_HAND"
        },
        {
            "name": "Sigil of Haunted Dreams",
            "id": 40715,
            "slot": "RANGED"
        }
	]}`),
};
