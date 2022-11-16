import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { Player } from '../core/player.js';
import { TypedEvent } from '../core/typed_event.js';

import { ProtectionWarriorSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecProtectionWarrior>(Spec.SpecProtectionWarrior, sim);
player.enableHealing();

sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new ProtectionWarriorSimUI(document.body, player);
