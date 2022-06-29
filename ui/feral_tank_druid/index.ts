import { Spec } from '/wotlk/core/proto/common.js';
import { Sim } from '/wotlk/core/sim.js';
import { Player } from '/wotlk/core/player.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { FeralTankDruidSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecFeralTankDruid>(Spec.SpecFeralTankDruid, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new FeralTankDruidSimUI(document.body, player);
