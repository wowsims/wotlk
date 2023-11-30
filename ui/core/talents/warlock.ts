import { WarlockTalents } from '../proto/warlock.js';


import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import WarlockTalentJson from './trees/warlock.json';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);

export const warlockGlyphsConfig = {
	majorGlyphs: {
	},
	minorGlyphs: {
	},
};
