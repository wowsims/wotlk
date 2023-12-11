import {
	Consumes,
	CustomRotation,
	CustomSpell,
	Flask,
	Food,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Rotation_DemoShoutChoice as DemoShoutChoice,
	ProtectionWarrior_Rotation_ThunderClapChoice as ThunderClapChoice,
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	ProtectionWarrior_Rotation_SpellOption as SpellOption,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DefaultGear = PresetUtils.makePresetGear('Blank', BlankGear);

export const DefaultRotation = ProtectionWarriorRotation.create({
	customRotation: CustomRotation.create({
		spells: [
			CustomSpell.create({ spell: SpellOption.ShieldSlam }),
			CustomSpell.create({ spell: SpellOption.Revenge }),
			CustomSpell.create({ spell: SpellOption.Shout }),
			CustomSpell.create({ spell: SpellOption.ThunderClap }),
			CustomSpell.create({ spell: SpellOption.DemoralizingShout }),
			CustomSpell.create({ spell: SpellOption.MortalStrike }),
			CustomSpell.create({ spell: SpellOption.Devastate }),
			CustomSpell.create({ spell: SpellOption.SunderArmor }),
			CustomSpell.create({ spell: SpellOption.ConcussionBlow }),
			CustomSpell.create({ spell: SpellOption.Shockwave }),
		],
	}),
	demoShoutChoice: DemoShoutChoice.DemoShoutChoiceNone,
	thunderClapChoice: ThunderClapChoice.ThunderClapChoiceNone,
	hsRageThreshold: 30,
});

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '2500030023-302-053351225000012521030113321',
	}),
};

export const UATalents = {
	name: 'UA',
	data: SavedTalents.create({
		talentsString: '35023301230051002020120002-2-05035122500000252',
	}),
};

export const DefaultOptions = ProtectionWarriorOptions.create({
	shout: WarriorShout.WarriorShoutCommanding,
	useShatteringThrow: false,
	startingRage: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
