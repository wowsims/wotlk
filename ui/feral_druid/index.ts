import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { Player } from '../core/player.js';
import { TypedEvent } from '../core/typed_event.js';

import { FeralDruidSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecFeralDruid>(Spec.SpecFeralDruid, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new FeralDruidSimUI(document.body, player);
