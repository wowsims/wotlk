import { RogueTalents } from '../proto/rogue.js';

import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import RogueTalentJson from './trees/rogue.json';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);
