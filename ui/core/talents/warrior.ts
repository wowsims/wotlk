import { WarriorTalents } from '../proto/warrior.js';


import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import WarriorTalentJson from './trees/warrior.json';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);