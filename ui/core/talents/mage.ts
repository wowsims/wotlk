import { MageTalents } from '../proto/mage.js';

import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import MageTalentJson from './trees/mage.json';

export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);