import { HunterTalents } from '../proto/hunter.js';


import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import HunterTalentJson from './trees/hunter.json';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);
