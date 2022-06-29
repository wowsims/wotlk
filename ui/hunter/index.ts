import { Spec } from '/wotlk/core/proto/common.js';
import { Sim } from '/wotlk/core/sim.js';
import { Player } from '/wotlk/core/player.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { HunterSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecHunter>(Spec.SpecHunter, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new HunterSimUI(document.body, player);
