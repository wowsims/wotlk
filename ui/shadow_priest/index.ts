import { Spec } from '/tbc/core/proto/common.js';
import { Sim } from '/tbc/core/sim.js';
import { Player } from '/tbc/core/player.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { ShadowPriestSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecShadowPriest>(Spec.SpecShadowPriest, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new ShadowPriestSimUI(document.body, player);
