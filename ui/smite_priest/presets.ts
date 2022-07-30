import { Consumes } from '/wotlk/core/proto/common.js';
import { EquipmentSpec } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { Glyphs } from '/wotlk/core/proto/common.js';
import { ItemSpec } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { Faction } from '/wotlk/core/proto/common.js';
import { SavedTalents } from '/wotlk/core/proto/ui.js';
import { Player } from '/wotlk/core/player.js';

import { SmitePriest, SmitePriest_Rotation as Rotation, SmitePriest_Options as Options, SmitePriest_Rotation_RotationType } from '/wotlk/core/proto/priest.js';

import * as Enchants from '/wotlk/core/constants/enchants.js';
import * as Gems from '/wotlk/core/proto_utils/gems.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '5012300130505120501-005551002020052',
	}),
};

export const HolyTalents = {
	name: 'Holy',
	data: SavedTalents.create({
		talentsString: '50023011305-235050032002150520051',
	}),
};

export const DefaultRotation = Rotation.create({
	rotationType: SmitePriest_Rotation_RotationType.Basic,
});

export const DefaultOptions = Options.create({
	useShadowfiend: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfBlindingLight,
	food: Food.FoodBlackenedBasilisk,
	mainHandImbue: WeaponImbue.WeaponImbueSuperiorWizardOil,
	defaultPotion: Potions.SuperManaPotion,
});

export const P1_PRESET = {
	name: 'P1 Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	gear: EquipmentSpec.fromJsonString(`{
      "items": [
        {
          "id": 24266,
          "enchant": 29191,
          "gems": [
            28118,
            24030,
            24030
          ]
        },
        {
          "id": 28530
        },
        {
          "id": 29060,
          "enchant": 28886,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 28766,
          "enchant": 33150
        },
        {
          "id": 29056,
          "enchant": 24003,
          "gems": [
            24030,
            24030,
            24030
          ]
        },
        {
          "id": 24250,
          "enchant": 22534,
          "gems": [
            24030
          ]
        },
        {
          "id": 30725,
          "enchant": 28272,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 24256,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 30734,
          "enchant": 24274,
          "gems": [
            24030,
            24030,
            24030
          ]
        },
        {
          "id": 28517,
          "enchant": 35297,
          "gems": [
            24030,
            24030
          ]
        },
        {
          "id": 28793,
          "enchant": 22536
        },
        {
          "id": 29172,
          "enchant": 22536
        },
        {
          "id": 27683
        },
        {
          "id": 29370
        },
        {
          "id": 30723,
          "enchant": 22555,
          "gems": [
            30564,
            31867
          ]
        },
        {
          "id": 28734
        },
        {
          "id": 28673
        }
      ]
    }`),
};
