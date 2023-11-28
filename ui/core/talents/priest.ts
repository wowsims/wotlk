import { PriestTalents } from '../proto/priest.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import PriestTalentJson from './trees/priest.json';

export const priestTalentsConfig: TalentsConfig<PriestTalents> = newTalentsConfig(PriestTalentJson);

export const priestGlyphsConfig = {
	majorGlyphs: {
	},
	minorGlyphs: {
	},
};
