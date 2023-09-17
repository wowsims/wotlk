import {
	Consumes,
	CustomRotation,
	CustomSpell,
	EquipmentSpec,
	Explosive,
	Flask,
	Food,
	Glyphs,
	PetFood,
	Potions,
	UnitReference,
	Spec
} from '../core/proto/common.js';
import { SavedRotation, SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	Deathknight_Options as DeathKnightOptions,
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Rotation_ArmyOfTheDead,
	Deathknight_Rotation_BloodRuneFiller,
	Deathknight_Rotation_CustomSpellOption as CustomSpellOption,
	Deathknight_Rotation_FrostRotationType,
	Deathknight_Rotation_Presence,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
	Deathknight_Rotation_DrwDiseases,
	Deathknight_Rotation_BloodSpell,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { APLRotation } from '../core/proto/apl.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost BL',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostUnholyTalents = {
	name: 'Frost UH',
	data: SavedTalents.create({
		talentsString: '01-32002350342203012300033101351-230200305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldTalents = {
	name: 'Unholy DW',
	data: SavedTalents.create({
		talentsString: '-320043500002-2300303050032152000150013133051',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldSSTalents = {
	name: 'Unholy DW SS',
	data: SavedTalents.create({
		talentsString: '-320033500002-2301303050032151000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const Unholy2HTalents = {
	name: 'Unholy 2H',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302003350032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDarkDeath,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyAoeTalents = {
	name: 'Unholy AOE',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302303050032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const BloodTalents = {
	name: 'Blood DPS',
	data: SavedTalents.create({
		talentsString: '2305120530003303231023001351--2302003050032',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDancingRuneWeapon,
			major2: DeathknightMajorGlyph.GlyphOfDeathStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultUnholyRotation = DeathKnightRotation.create({
	useDeathAndDecay: true,
	btGhoulFrenzy: true,
	refreshHornOfWinter: false,
	useGargoyle: true,
	useEmpowerRuneWeapon: true,
	holdErwArmy: false,
	preNerfedGargoyle: false,
	armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead.AsMajorCd,
	startingPresence: Deathknight_Rotation_Presence.Unholy,
	blPresence: Deathknight_Rotation_Presence.Blood,
	presence: Deathknight_Rotation_Presence.Blood,
	gargoylePresence: Deathknight_Rotation_Presence.Unholy,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
	bloodSpender: Deathknight_Rotation_BloodSpell.HS,
	useDancingRuneWeapon: true
});

export const DefaultUnholyOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastGhoulFrenzy: false,
	precastHornOfWinter: true,
	unholyFrenzyTarget: UnitReference.create(),
	diseaseDowntime: 2,
});

export const DefaultFrostRotation = DeathKnightRotation.create({
	useDeathAndDecay: false,
	btGhoulFrenzy: false,
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	preNerfedGargoyle: false,
	startingPresence: Deathknight_Rotation_Presence.Blood,
	presence: Deathknight_Rotation_Presence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
	avgAmsSuccessRate: 1.0,
	avgAmsHit: 10000.0,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
  	frostRotationType: Deathknight_Rotation_FrostRotationType.SingleTarget,
	armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead.PreCast,
  	frostCustomRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: CustomSpellOption.CustomDeathAndDecay }),
			CustomSpell.create({ spell: CustomSpellOption.CustomIcyTouch }),
			CustomSpell.create({ spell: CustomSpellOption.CustomPlagueStrike }),
			CustomSpell.create({ spell: CustomSpellOption.CustomPestilence }),
			CustomSpell.create({ spell: CustomSpellOption.CustomHowlingBlastRime }),
			CustomSpell.create({ spell: CustomSpellOption.CustomHowlingBlast }),
			CustomSpell.create({ spell: CustomSpellOption.CustomBloodBoil }),
			CustomSpell.create({ spell: CustomSpellOption.CustomObliterate }),
			CustomSpell.create({ spell: CustomSpellOption.CustomFrostStrike }),
		],
	}),
});

export const DefaultFrostOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: UnitReference.create(),
	diseaseDowntime: 0,
});

export const DefaultBloodRotation = DeathKnightRotation.create({
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	preNerfedGargoyle: false,
	startingPresence: Deathknight_Rotation_Presence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodStrike,
	armyOfTheDead: Deathknight_Rotation_ArmyOfTheDead.PreCast,
	holdErwArmy: false,
	useAms: false,
	drwDiseases: Deathknight_Rotation_DrwDiseases.Pestilence,
	bloodSpender: Deathknight_Rotation_BloodSpell.HS,
	useDancingRuneWeapon: true,
});

export const DefaultBloodOptions = DeathKnightOptions.create({
	drwPestiApply: true,
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: UnitReference.create(),
	diseaseDowntime: 0,
});

export const OtherDefaults = {
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.PotionOfSpeed,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const BLOOD_ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Blood Legacy',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultBloodRotation),
	}),
}

export const FROST_ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Frost Legacy',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultFrostRotation),
	}),
}

export const UNHOLY_DW_ROTATION_PRESET_LEGACY_DEFAULT = {
	name: 'Unholy DW Legacy',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultUnholyRotation),
	}),
}

export const BLOOD_PESTI_ROTATION_PRESET_DEFAULT = {
	name: 'Blood Pesti APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultBloodRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-20s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":50689}}},"doAtValue":{"const":{"val":"-6s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"8.5s"}}}},"castSpell":{"spellId":{"spellId":64382,"tag":-1}}}},
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"resetSequence":{"sequenceName":"IT"}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"resetSequence":{"sequenceName":"PS"}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}}]}},"sequence":{"name":"IT","actions":[{"castSpell":{"spellId":{"spellId":59131}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}}]}},"sequence":{"name":"PS","actions":[{"castSpell":{"spellId":{"spellId":49921,"tag":1}}}]}}},
			  {"action":{"sequence":{"name":"Opener","actions":[{"castSpell":{"spellId":{"spellId":49016}}},{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"spellId":26297}}},{"castSpell":{"spellId":{"spellId":54758}}},{"castSpell":{"spellId":{"spellId":49924,"tag":1}}},{"castSpell":{"spellId":{"spellId":55262,"tag":1}}},{"castSpell":{"spellId":{"spellId":55262,"tag":1}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpEq","lhs":{"currentRuneCount":{"runeType":"RuneDeath"}},"rhs":{"const":{"val":"0"}}}}]}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":46584}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"3s"}}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"math":{"op":"OpAdd","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"4s"}}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"2s"}}}},{"auraIsActive":{"auraId":{"spellId":49028}}},{"sequenceIsComplete":{"sequenceName":"IT"}}]}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"spellId":50842}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":49016}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":26297}}},"rhs":{"const":{"val":"10"}}}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15"}}}},{"auraIsActive":{"auraId":{"spellId":33697}}}]}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"9"}}}}}}]}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":33697}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12"}}}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"55"}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"20"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"20"}}}}}},{"gcdIsReady":{}},{"spellIsReady":{"spellId":{"spellId":26297}}}]}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":49016}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":33697}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":46584}}},{"castSpell":{"spellId":{"spellId":49028}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"currentRuneActive":{"runeSlot":"SlotLeftBlood"}}}},{"sequenceIsComplete":{"sequenceName":"Opener"}},{"spellIsReady":{"spellId":{"spellId":45529}}},{"gcdIsReady":{}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"40"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}}]}}]}}]}},"castSpell":{"spellId":{"spellId":45529}}}},
			  {"action":{"condition":{"and":{"vals":[{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"8s"}}}}]}},{"not":{"val":{"or":{"vals":[{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneFrost"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}},{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneUnholy"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}}]}}}}]}},"castSpell":{"spellId":{"spellId":55262,"tag":1}}}},
			  {"action":{"condition":{"not":{"val":{}}},"castSpell":{"spellId":{"spellId":49924,"tag":1}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.3"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}}]}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const BLOOD_PESTI_DD_ROTATION_PRESET_DEFAULT = {
	name: 'Blood Pesti DD APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultBloodRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-20s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":50689}}},"doAtValue":{"const":{"val":"-6s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"8.5s"}}}},"castSpell":{"spellId":{"spellId":64382,"tag":-1}}}},
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"resetSequence":{"sequenceName":"IT"}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"resetSequence":{"sequenceName":"PS"}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}}]}},"sequence":{"name":"IT","actions":[{"castSpell":{"spellId":{"spellId":59131}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}}]}},"sequence":{"name":"PS","actions":[{"castSpell":{"spellId":{"spellId":49921,"tag":1}}}]}}},
			  {"action":{"sequence":{"name":"Opener","actions":[{"castSpell":{"spellId":{"spellId":49016}}},{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"spellId":26297}}},{"castSpell":{"spellId":{"spellId":54758}}},{"castSpell":{"spellId":{"spellId":49924,"tag":1}}},{"castSpell":{"spellId":{"spellId":55262,"tag":1}}},{"castSpell":{"spellId":{"spellId":55262,"tag":1}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpEq","lhs":{"currentRuneCount":{"runeType":"RuneDeath"}},"rhs":{"const":{"val":"0"}}}}]}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":46584}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"3s"}}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"math":{"op":"OpAdd","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"4s"}}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"2s"}}}},{"auraIsActive":{"auraId":{"spellId":49028}}},{"sequenceIsComplete":{"sequenceName":"IT"}}]}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"spellId":50842}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":49016}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":26297}}},"rhs":{"const":{"val":"10"}}}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15"}}}},{"auraIsActive":{"auraId":{"spellId":33697}}}]}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"9"}}}}}}]}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":33697}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12"}}}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"55"}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"20"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"20"}}}}}},{"gcdIsReady":{}},{"spellIsReady":{"spellId":{"spellId":26297}}}]}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":49016}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":33697}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":46584}}},{"castSpell":{"spellId":{"spellId":49028}}},{"resetSequence":{"sequenceName":"IT"}},{"resetSequence":{"sequenceName":"PS"}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"currentRuneActive":{"runeSlot":"SlotLeftBlood"}}}},{"sequenceIsComplete":{"sequenceName":"Opener"}},{"spellIsReady":{"spellId":{"spellId":45529}}},{"gcdIsReady":{}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"40"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}}]}}]}}]}},"castSpell":{"spellId":{"spellId":45529}}}},
			  {"action":{"condition":{"and":{"vals":[{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"8s"}}}}]}},{"not":{"val":{"or":{"vals":[{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneFrost"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}},{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneUnholy"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}}]}}}}]}},"castSpell":{"spellId":{"spellId":55262,"tag":1}}}},
			  {"action":{"condition":{"not":{"val":{}}},"castSpell":{"spellId":{"spellId":49924,"tag":1}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.3"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}}]}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const BLOOD_PESTI_AOE_ROTATION_PRESET_DEFAULT = {
	name: 'Blood Pesti AOE APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultBloodRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-20s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":50689}}},"doAtValue":{"const":{"val":"-6s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"8.5s"}}}},"castSpell":{"spellId":{"spellId":64382,"tag":-1}}}},
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"resetSequence":{"sequenceName":"IT"}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"resetSequence":{"sequenceName":"PS"}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}}]}},"sequence":{"name":"IT","actions":[{"castSpell":{"spellId":{"spellId":59131}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"cmp":{"op":"OpEq","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGe","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"1"}}}}]}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}}]}},"sequence":{"name":"PS","actions":[{"castSpell":{"spellId":{"spellId":49921,"tag":1}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49028}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},"sequence":{"name":"Pesti","actions":[{"castSpell":{"spellId":{"spellId":50842}}}]}}},
			  {"action":{"sequence":{"name":"Opener","actions":[{"castSpell":{"spellId":{"spellId":49016}}},{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"spellId":26297}}},{"castSpell":{"spellId":{"spellId":54758}}},{"castSpell":{"spellId":{"spellId":49924,"tag":1}}},{"castSpell":{"spellId":{"spellId":50842}}},{"castSpell":{"spellId":{"spellId":49941}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpEq","lhs":{"currentRuneCount":{"runeType":"RuneDeath"}},"rhs":{"const":{"val":"0"}}}}]}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}},{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":46584}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"3s"}}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"dotIsActive":{"spellId":{"spellId":55095}}},{"dotIsActive":{"spellId":{"spellId":55078}}}]}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"math":{"op":"OpAdd","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"4s"}}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"const":{"val":"2s"}}}},{"auraIsActive":{"auraId":{"spellId":49028}}},{"sequenceIsComplete":{"sequenceName":"IT"}}]}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"spellId":50842}}}]}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":49016}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":26297}}},"rhs":{"const":{"val":"10"}}}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15"}}}},{"auraIsActive":{"auraId":{"spellId":33697}}}]}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"15"}}}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"9"}}}}}}]}},{"gcdIsReady":{}}]}},{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":33697}}},{"spellIsReady":{"spellId":{"spellId":46584}}},{"spellIsReady":{"spellId":{"spellId":49028}}},{"gcdIsReady":{}}]}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12"}}}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"55"}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":33697}}},"rhs":{"const":{"val":"20"}}}}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":49016}}},"rhs":{"const":{"val":"20"}}}}}},{"gcdIsReady":{}},{"spellIsReady":{"spellId":{"spellId":26297}}}]}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":49016}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"castSpell":{"spellId":{"spellId":33697}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":54758}}},{"gcdIsReady":{}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":46584}}},{"castSpell":{"spellId":{"spellId":49028}}},{"resetSequence":{"sequenceName":"Pesti"}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"currentRuneActive":{"runeSlot":"SlotLeftBlood"}}}},{"sequenceIsComplete":{"sequenceName":"Opener"}},{"spellIsReady":{"spellId":{"spellId":45529}}},{"gcdIsReady":{}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"40"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49028}}}}},{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"auraId":{"spellId":49028}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}}]}}]}}]}},"castSpell":{"spellId":{"spellId":45529}}}},
			  {"action":{"condition":{"and":{"vals":[{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}}}},{"cmp":{"op":"OpGt","lhs":{"nextRuneCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"8s"}}}}]}},{"not":{"val":{"or":{"vals":[{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneFrost"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}},{"cmp":{"op":"OpGe","lhs":{"nextRuneCooldown":{"runeType":"RuneUnholy"}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}}}}]}}}}]}},"castSpell":{"spellId":{"spellId":49941}}}},
			  {"action":{"condition":{"not":{"val":{}}},"castSpell":{"spellId":{"spellId":49924,"tag":1}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.3"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49028}}},"rhs":{"const":{"val":"10"}}}},{"cmp":{"op":"OpEq","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"100"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":71227}}},"rhs":{"const":{"val":"1.5s"}}}}}}]}}]}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}},{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"runeCooldown":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"0.5"}}}},{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const FROST_BL_PESTI_ROTATION_PRESET_DEFAULT = {
	name: 'Frost BL Pesti APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultFrostRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-20s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":50689}}},"doAtValue":{"const":{"val":"-6s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"sequence":{"name":"Opener","actions":[{"castSpell":{"spellId":{"spellId":59131}}},{"castSpell":{"spellId":{"tag":1,"spellId":49921}}},{"castSpell":{"spellId":{"spellId":51271}}},{"castSpell":{"spellId":{"spellId":54758}}},{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":55268}}},{"castSpell":{"spellId":{"spellId":50842}}},{"castSpell":{"spellId":{"spellId":47568}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":55268}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"spellId":46584}}}]}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49921}}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"spellId":33697}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
			  {"action":{"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":51271}}},{"castSpell":{"spellId":{"spellId":45529}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"4s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":59052}}},"castSpell":{"spellId":{"spellId":51411}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":51425}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":46584}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":49930}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":55268}}}}
			]
		}`),
	}),
}

export const FROST_UH_PESTI_ROTATION_PRESET_DEFAULT = {
	name: 'Frost UH Pesti APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultFrostRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-20s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":50689}}},"doAtValue":{"const":{"val":"-6s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"sequence":{"name":"Opener","actions":[{"castSpell":{"spellId":{"spellId":59131}}},{"castSpell":{"spellId":{"tag":1,"spellId":49921}}},{"castSpell":{"spellId":{"spellId":51271}}},{"castSpell":{"spellId":{"spellId":54758}}},{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":55268}}},{"castSpell":{"spellId":{"tag":1,"spellId":49930}}},{"castSpell":{"spellId":{"spellId":47568}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":55268}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"tag":1,"spellId":51425}}},{"castSpell":{"spellId":{"spellId":46584}}}]}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49921}}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"1.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"spellId":33697}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":51271}}},{"spellCanCast":{"spellId":{"spellId":51271}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
			  {"action":{"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":51271}}},{"castSpell":{"spellId":{"spellId":45529}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"8.5s"}}}},{"dotIsActive":{"spellId":{"spellId":55095}}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":59052}}},"castSpell":{"spellId":{"spellId":51411}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":51425}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":46584}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":49930}}}},
			  {"action":{"castSpell":{"spellId":{"tag":1,"spellId":55268}}}}
			]
		}`),
	}),
}

export const UNHOLY_DW_ROTATION_PRESET_DEFAULT = {
	name: 'Unholy DW SS APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultUnholyRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":63560}}},"doAtValue":{"const":{"val":"-8s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55078}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":49921,"tag":1}}}},
			  {"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":66803}}}}},"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"const":{"val":"50s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"remainingTime":{}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":47568}}}}},{"cmp":{"op":"OpEq","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"1"}}}}]}},"castSpell":{"spellId":{"spellId":42650}}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":67383}}}}},{"spellIsReady":{"spellId":{"spellId":49206}}}]}},"castSpell":{"spellId":{"spellId":55271,"tag":1}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49938}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":42650}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49938}}},"rhs":{"const":{"val":"6s"}}}},"castSpell":{"spellId":{"spellId":55271,"tag":1}}}},
			  {"action":{"condition":{"and":{"vals":[{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49938}}},"rhs":{"const":{"val":"6s"}}}},{"spellIsReady":{"spellId":{"spellId":47568}}}]}},{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraRemainingTime":{"auraId":{"spellId":66803}}},"rhs":{"const":{"val":"10s"}}}},{"cmp":{"op":"OpLe","lhs":{"auraInternalCooldown":{"auraId":{"spellId":67117}}},"rhs":{"const":{"val":"0s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49938}}},"rhs":{"const":{"val":"6s"}}}},{"spellIsReady":{"spellId":{"spellId":47568}}}]}},"castSpell":{"spellId":{"spellId":49941}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49206}}}},
			  {"action":{"condition":{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"spellId":63560}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":48265}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49206}}}}}]}},"castSpell":{"spellId":{"spellId":50689}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const UNHOLY_2H_ROTATION_PRESET_DEFAULT = {
	name: 'Unholy 2H SS APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultUnholyRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":63560}}},"doAtValue":{"const":{"val":"-8s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55078}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":49921,"tag":1}}}},
			  {"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":66803}}}}},"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"const":{"val":"50s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"remainingTime":{}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":42650}}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"Pet","index":1,"owner":{"type":"Self"}},"auraId":{"spellId":63560}}}}},{"not":{"val":{"currentRuneActive":{"runeSlot":"SlotLeftBlood"}}}}]}},"castSpell":{"spellId":{"spellId":45529}}}},
			  {"action":{"condition":{"and":{"vals":[{"currentRuneDeath":{"runeSlot":"SlotLeftBlood"}},{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"Pet","index":1,"owner":{"type":"Self"}},"auraId":{"spellId":63560}}}}},{"cmp":{"op":"OpEq","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}}]}},"castSpell":{"spellId":{"spellId":63560}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":55271,"tag":1}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49206}}}},
			  {"action":{"condition":{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":48265}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49206}}}}}]}},"castSpell":{"spellId":{"spellId":50689}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT = {
	name: 'Unholy DND AOE APL',
	//enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	rotation: SavedRotation.create({
		specRotationOptionsJson: DeathKnightRotation.toJsonString(DefaultUnholyRotation),
		rotation: APLRotation.fromJsonString(`{
			"type": "TypeAPL",
			"prepullActions": [
			  {"action":{"castSpell":{"spellId":{"spellId":48265}}},"doAtValue":{"const":{"val":"-10s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":63560}}},"doAtValue":{"const":{"val":"-8s"}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}},"doAtValue":{"const":{"val":"-1.5s"}}},
			  {"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
			],
			"priorityList": [
			  {"action":{"autocastOtherCooldowns":{}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55095}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":59131}}}},
			  {"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55078}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":49921,"tag":1}}}},
			  {"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":66803}}}}},"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":26297}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"const":{"val":"50s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
			  {"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":49206}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49206}}},"rhs":{"remainingTime":{}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":47568}}}},
			  {"action":{"condition":{"auraIsActive":{"auraId":{"spellId":49206}}},"castSpell":{"spellId":{"spellId":42650}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49938}}}},
			  {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49938}}},"rhs":{"const":{"val":"6s"}}}},"castSpell":{"spellId":{"spellId":55271,"tag":1}}}},
			  {"action":{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"targetUnit":{"type":"Target","index":1},"spellId":{"spellId":55095}}}}}]}},"castSpell":{"spellId":{"spellId":50842}}}},
			  {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49938}}},"rhs":{"const":{"val":"6s"}}}},{"cmp":{"op":"OpLt","lhs":{"auraRemainingTime":{"auraId":{"spellId":66803}}},"rhs":{"const":{"val":"10s"}}}}]}},"castSpell":{"spellId":{"spellId":49930,"tag":1}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":49206}}}},
			  {"action":{"condition":{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},"castSpell":{"spellId":{"spellId":49895}}}},
			  {"action":{"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":45529}}},{"castSpell":{"spellId":{"spellId":63560}}}]}}},
			  {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":48265}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":49206}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49206}}}}}]}},"castSpell":{"spellId":{"spellId":50689}}}},
			  {"action":{"castSpell":{"spellId":{"spellId":57623}}}}
			]
		}`),
	}),
}

export const P1_BLOOD_BIS_PRESET = {
	name: 'P1 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":44006,"enchant":3817,"gems":[41398,42702]},
		  {"id":44664,"gems":[39996]},
		  {"id":40557,"enchant":3808,"gems":[39996]},
		  {"id":40403,"enchant":3831},
		  {"id":40550,"enchant":3832,"gems":[42142,42142]},
		  {"id":40330,"enchant":3845,"gems":[42142,0]},
		  {"id":40552,"enchant":3604,"gems":[39996,0]},
		  {"id":40278,"gems":[39996,39996]},
		  {"id":40556,"enchant":3823,"gems":[39996,40037]},
		  {"id":40591,"enchant":3606},
		  {"id":40075},
		  {"id":39401},
		  {"id":40256},
		  {"id":42987},
		  {"id":40384,"enchant":3368},
		  {},
		  {"id":40207}
  ]}`),
};

export const P2_BLOOD_BIS_PRESET = {
	name: 'P2 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":46115,"enchant":3817,"gems":[41398,42702]},
		  {"id":45459,"gems":[39996]},
		  {"id":46117,"enchant":3808,"gems":[39996]},
		  {"id":46032,"enchant":3831,"gems":[39996,39996]},
		  {"id":46111,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[42142,0]},
		  {"id":46113,"enchant":3604,"gems":[39996,0]},
		  {"id":45241,"gems":[39996,45862,39996]},
		  {"id":45134,"enchant":3823,"gems":[39996,39996,39996]},
		  {"id":45599,"enchant":3606,"gems":[39996,39996]},
		  {"id":45534,"gems":[39996]},
		  {"id":46048,"gems":[39996]},
		  {"id":42987},
		  {"id":45931},
		  {"id":45516,"enchant":3368,"gems":[39996,39996]},
		  {},
		  {"id":45254}
  ]}`),
};

export const P3_BLOOD_BIS_PRESET = {
	name: 'P3 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		  {"id":48493,"enchant":3817,"gems":[41285,40142]},
		  {"id":47458,"gems":[40142]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47546,"enchant":3831,"gems":[42142]},
		  {"id":47449,"enchant":3832,"gems":[49110,42142,40142]},
		  {"id":48008,"enchant":3845,"gems":[40111,0]},
		  {"id":48492,"enchant":3604,"gems":[40142,0]},
		  {"id":47429,"gems":[40142,40142,40111]},
		  {"id":48494,"enchant":3823,"gems":[40142,40111]},
		  {"id":45599,"enchant":3606,"gems":[40111,40111]},
		  {"id":47993,"gems":[40111,45862]},
		  {"id":47413,"gems":[40142]},
		  {"id":45931},
		  {"id":47464},
		  {"id":47446,"enchant":3368,"gems":[42142,40141]},
		  {},
		  {"id":47673}
  ]}`),
};

export const P4_BLOOD_BIS_PRESET = {
	name: 'P4 Blood',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
        {"id":51312,"enchant":3817,"gems":[41398,40117]},
        {"id":50728,"gems":[40143]},
        {"id":51314,"enchant":3808,"gems":[40117]},
        {"id":50677,"enchant":3831,"gems":[42156]},
        {"id":51310,"enchant":3832,"gems":[40117,49110]},
        {"id":50659,"enchant":3845,"gems":[40162,0]},
        {"id":50675,"enchant":3604,"gems":[40143,40117,0]},
        {"id":50620,"gems":[40125,40117,40117]},
        {"id":51313,"enchant":3823,"gems":[40117,40117]},
        {"id":50639,"enchant":3606,"gems":[40125,40117]},
        {"id":50693,"gems":[40125]},
        {"id":52572,"gems":[40125]},
        {"id":50363},
        {"id":47464},
        {"id":49623,"enchant":3368,"gems":[40117,42153,42153]},
        {},
        {"id":47673}
  ]}`),
};

export const P1_UNHOLY_2H_PRERAID_PRESET = {
	name: 'Pre-Raid 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":41386,"enchant":3817,"gems":[41400,49110]},
		{"id":37397},
		{"id":37627,"enchant":3808,"gems":[39996]},
		{"id":37647,"enchant":3831},
		{"id":39617,"enchant":3832,"gems":[42142,39996]},
		{"id":41355,"enchant":3845,"gems":[0]},
		{"id":39618,"enchant":3604,"gems":[39996,0]},
		{"id":40688,"gems":[39996,42142]},
		{"id":37193,"enchant":3823,"gems":[42142,39996]},
		{"id":44306,"enchant":3606,"gems":[39996,39996]},
		{"id":37642},
		{"id":44935},
		{"id":40684},
		{"id":42987},
		{"id":41257,"enchant":3368},
		{},
		{"id":40867}
  ]}`),
};

export const P1_UNHOLY_2H_BIS_PRESET = {
	name: 'P1 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
		{"id":44006,"enchant":3817,"gems":[41400,49110]},
		{"id":44664,"gems":[39996]},
		{"id":40557,"enchant":3808,"gems":[39996]},
		{"id":40403,"enchant":3831},
		{"id":40550,"enchant":3832,"gems":[42142,39996]},
		{"id":40330,"enchant":3845,"gems":[39996,0]},
		{"id":40552,"enchant":3604,"gems":[40038,0]},
		{"id":40278,"gems":[42142,42142]},
		{"id":40556,"enchant":3823,"gems":[39996,39996]},
		{"id":40591,"enchant":3606},
		{"id":39401},
		{"id":40075},
		{"id":40256},
		{"id":42987},
		{"id":40384,"enchant":3368},
		{},
		{"id":40207}
	  ]
    }`),
};

export const P4_UNHOLY_2H_BIS_PRESET = {
	name: 'P4 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {"id":51312,"enchant":3817,"gems":[41398,40111]},
        {"id":50647,"gems":[40111]},
        {"id":51314,"enchant":3808,"gems":[40111]},
        {"id":50677,"enchant":3831,"gems":[40146]},
        {"id":51310,"enchant":3832,"gems":[40111,40111]},
        {"id":50659,"enchant":3845,"gems":[40146,0]},
        {"id":51311,"enchant":3604,"gems":[40146,0]},
        {"id":50620,"gems":[40146,40111,40111]},
        {"id":50624,"enchant":3823,"gems":[40146,40111,49110]},
        {"id":50639,"enchant":3606,"gems":[40146,40111]},
        {"id":50693,"gems":[40146]},
        {"id":52572,"gems":[40146]},
        {"id":47464},
        {"id":50363},
        {"id":49623,"enchant":3368,"gems":[42142,42142,42142]},
        {},
        {"id":47673}
      ]
	}`),
};

export const P1_UNHOLY_DW_PRERAID_PRESET = {
	name: 'Pre-Raid DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":41386,"enchant":3817,"gems":[41400,49110]},
		{"id":37397},
		{"id":37627,"enchant":3808,"gems":[39996]},
		{"id":37647,"enchant":3831},
		{"id":39617,"enchant":3832,"gems":[42142,39996]},
		{"id":41355,"enchant":3845,"gems":[0]},
		{"id":39618,"enchant":3604,"gems":[39996,0]},
		{"id":40688,"gems":[39996,42142]},
		{"id":37193,"enchant":3823,"gems":[42142,39996]},
		{"id":44306,"enchant":3606,"gems":[39996,39996]},
		{"id":37642},
		{"id":44935},
		{"id":40684},
		{"id":42987},
		{"id":41383,"enchant":3368},
		{"id":40703,"enchant":3368},
		{"id":40867}
  ]}`),
};

export const P1_UNHOLY_DW_BIS_PRESET = {
	name: 'P1 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":44006,"enchant":3817,"gems":[41398,42702]},
		{"id":39421},
		{"id":40557,"enchant":3808,"gems":[39996]},
		{"id":40403,"enchant":3831},
		{"id":40550,"enchant":3832,"gems":[42142,39996]},
		{"id":40330,"enchant":3845,"gems":[39996,0]},
		{"id":40347,"enchant":3604,"gems":[39996,0]},
		{"id":40278,"gems":[42142,42142]},
		{"id":40294,"enchant":3823},
		{"id":39706,"enchant":3606,"gems":[39996]},
		{"id":39401},
		{"id":40075},
		{"id":37390},
		{"id":42987},
		{"id":40402,"enchant":3368},
		{"id":40491,"enchant":3368},
		{"id":42620}
  ]}`),
};

export const P2_UNHOLY_DW_BIS_PRESET = {
	name: 'P2 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":45472,"enchant":3817,"gems":[41398,40041]},
		  {"id":46040,"gems":[39996]},
		  {"id":46117,"enchant":3808,"gems":[39996]},
		  {"id":45588,"enchant":3831,"gems":[39996]},
		  {"id":46111,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[39996,0]},
		  {"id":45481,"enchant":3604,"gems":[0]},
		  {"id":45241,"gems":[42142,45862,39996]},
		  {"id":45134,"enchant":3823,"gems":[40041,39996,40022]},
		  {"id":45599,"enchant":3606,"gems":[39996,39996]},
		  {"id":45534,"gems":[39996]},
		  {"id":45250},
		  {"id":45609},
		  {"id":42987},
		  {"id":46097,"enchant":3368,"gems":[39996]},
		  {"id":46036,"enchant":3368,"gems":[39996]},
		  {"id":45254}
  ]}`),
};

export const P3_UNHOLY_DW_BIS_PRESET = {
	name: 'P3 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
		{"id":48493,"enchant":3817,"gems":[41398,40146]},
		  {"id":47458,"gems":[40146]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47548,"enchant":3831,"gems":[40111]},
		  {"id":48491,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[40111,0]},
		  {"id":48492,"enchant":3604,"gems":[40146,0]},
		  {"id":47429,"gems":[40111,45862,40111]},
		  {"id":47465,"enchant":3823,"gems":[49110,40111,40146]},
		  {"id":45599,"enchant":3606,"gems":[40111,40111]},
		  {"id":47413,"gems":[40146]},
		  {"id":45534,"gems":[42142]},
		  {"id":47464},
		  {"id":45609},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":47673}
  ]}`),
};

export const P4_UNHOLY_DW_BIS_PRESET = {
	name: 'P4 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 2,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {"id":51312,"enchant":3817,"gems":[41398,40111]},
        {"id":50647,"gems":[40111]},
        {"id":51314,"enchant":3808,"gems":[40111]},
        {"id":50677,"enchant":3831,"gems":[40146]},
        {"id":51310,"enchant":3832,"gems":[42142,49110]},
        {"id":50659,"enchant":3845,"gems":[40146,0]},
        {"id":51311,"enchant":3604,"gems":[40146,0]},
        {"id":50620,"gems":[40146,40111,42142]},
        {"id":50624,"enchant":3823,"gems":[40111,42142,40111]},
        {"id":50639,"enchant":3606,"gems":[40146,40111]},
        {"id":52572,"gems":[40146]},
        {"id":51855,"gems":[40111]},
        {"id":47131},
        {"id":50363},
        {"id":50737,"enchant":3368,"gems":[40111]},
        {"id":50737,"enchant":3368,"gems":[40111]},
        {"id":47673}
	]}`),
};

export const P1_FROST_PRE_BIS_PRESET = {
	name: 'Pre-Raid Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":41386,"enchant":3817,"gems":[41398,49110]},
		{"id":42645,"gems":[42142]},
		{"id":34388,"enchant":3808,"gems":[39996,39996]},
		{"id":37647,"enchant":3831},
		{"id":39617,"enchant":3832,"gems":[42142,39996]},
		{"id":41355,"enchant":3845,"gems":[0]},
		{"id":39618,"enchant":3604,"gems":[39996,0]},
		{"id":37171,"gems":[39996,39996]},
		{"id":37193,"enchant":3823,"gems":[42142,39996]},
		{"id":44306,"enchant":3606,"gems":[39996,39996]},
		{"id":42642,"gems":[39996]},
		{"id":44935},
		{"id":40684},
		{"id":42987},
		{"id":41383,"enchant":3370},
		{"id":43611,"enchant":3368},
		{"id":40715}
  ]}`),
};

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":44006,"enchant":3817,"gems":[41398,42702]},
		{"id":44664,"gems":[39996]},
		{"id":40557,"enchant":3808,"gems":[39996]},
		{"id":40403,"enchant":3831},
		{"id":40550,"enchant":3832,"gems":[42142,39996]},
		{"id":40330,"enchant":3845,"gems":[39996,0]},
		{"id":40552,"enchant":3604,"gems":[39996,0]},
		{"id":40278,"gems":[39996,42142]},
		{"id":40556,"enchant":3823,"gems":[42142,39996]},
		{"id":40591,"enchant":3606},
		{"id":39401},
		{"id":40075},
		{"id":40256},
		{"id":42987},
		{"id":40189,"enchant":3370},
		{"id":40189,"enchant":3368},
		{"id":40207}
  ]}`),
};

export const P2_FROST_BIS_PRESET = {
	name: 'P2 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":46115,"enchant":3817,"gems":[41398,42702]},
		  {"id":45459,"gems":[39996]},
		  {"id":46117,"enchant":3808,"gems":[39996]},
		  {"id":46032,"enchant":3831,"gems":[39996,39996]},
		  {"id":46111,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[39996,0]},
		  {"id":46113,"enchant":3604,"gems":[39996,0]},
		  {"id":45241,"gems":[42142,45862,39996]},
		  {"id":45134,"enchant":3823,"gems":[39996,39996,39996]},
		  {"id":45599,"enchant":3606,"gems":[39996,39996]},
		  {"id":45608,"gems":[39996]},
		  {"id":45534,"gems":[39996]},
		  {"id":45931},
		  {"id":42987},
		  {"id":46097,"enchant":3370,"gems":[39996]},
		  {"id":46097,"enchant":3368,"gems":[39996]},
		  {"id":40207}
  ]}`),
};

export const P3_FROST_BIS_PRESET = {
	name: 'P3 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":48493,"enchant":3817,"gems":[41398,40142]},
		  {"id":45459,"gems":[40111]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47548,"enchant":3831,"gems":[40111]},
		  {"id":48491,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[40111,0]},
		  {"id":47492,"enchant":3604,"gems":[49110,40111,0]},
		  {"id":45241,"gems":[40111,42142,40111]},
		  {"id":48494,"enchant":3823,"gems":[40142,40111]},
		  {"id":47473,"enchant":3606,"gems":[40142,40111]},
		  {"id":46966,"gems":[40111]},
		  {"id":45534,"gems":[40111]},
		  {"id":47464},
		  {"id":45931},
		  {"id":47528,"enchant":3370,"gems":[40111]},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":40207}
  ]}`),
};

export const P4_FROST_BIS_PRESET = {
	name: 'P4 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
        {"id":51312,"enchant":3817,"gems":[41398,42153]},
        {"id":50728,"gems":[40117]},
        {"id":51314,"enchant":3808,"gems":[42153]},
        {"id":47548,"enchant":3831,"gems":[40117]},
        {"id":51310,"enchant":3832,"gems":[42153,40117]},
        {"id":50659,"enchant":3845,"gems":[40117,0]},
        {"id":51311,"enchant":3604,"gems":[40117,0]},
        {"id":50620,"enchant":3601,"gems":[40143,40117,40117]},
        {"id":51817,"enchant":3823,"gems":[49110,40117,40143]},
        {"id":50639,"enchant":3606,"gems":[40143,40117]},
        {"id":52572,"gems":[40117]},
        {"id":50693,"gems":[40117]},
        {"id":50363},
        {"id":47464},
        {"id":50737,"enchant":3370,"gems":[40117]},
        {"id":50737,"enchant":3368,"gems":[40117]},
        {"id":40207}
  ]}`),
};

export const P1_FROSTSUBUNH_BIS_PRESET = {
	name: 'P1 Frost Sub Unh',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 1,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
		{"id":44006,"enchant":3817,"gems":[41398,42702]},
		{"id":44664,"gems":[40003]},
		{"id":40557,"enchant":3808,"gems":[40003]},
		{"id":40403,"enchant":3831},
		{"id":40550,"enchant":3832,"gems":[42142,40003]},
		{"id":40330,"enchant":3845,"gems":[39996,0]},
		{"id":40552,"enchant":3604,"gems":[40058,0]},
		{"id":40278,"gems":[39996,42142]},
		{"id":40556,"enchant":3823,"gems":[42142,39996]},
		{"id":40591,"enchant":3606},
		{"id":39401},
		{"id":40075},
		{"id":40256},
		{"id":42987},
		{"id":40189,"enchant":3370},
		{"id":40189,"enchant":3368},
		{"id":40207}
  ]}`),
};
