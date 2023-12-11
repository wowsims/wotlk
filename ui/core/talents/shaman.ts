import { ShamanTalents } from '../proto/shaman.js';


import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import ShamanTalentJson from './trees/shaman.json';

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig(ShamanTalentJson);
