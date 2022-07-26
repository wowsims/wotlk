import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import { BalanceDruidSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecBalanceDruid>(Spec.SpecBalanceDruid, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new BalanceDruidSimUI(document.body, player);
