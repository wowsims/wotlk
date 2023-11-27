import { DruidTalents } from '../proto/druid.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import DruidTalentsJson from './trees/druid.json';

export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentsJson);

export const druidGlyphsConfig = {
	majorGlyphs: {
	},
	minorGlyphs: {
	},
};
