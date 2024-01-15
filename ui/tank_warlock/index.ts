import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { Player } from '../core/player.js';
import { TypedEvent } from '../core/typed_event.js';

import { TankWarlockSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecTankWarlock>(Spec.SpecTankWarlock, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

new TankWarlockSimUI(document.body, player);
