import { Spec } from '/wotlk/core/proto/common.js';
import { Sim } from '/wotlk/core/sim.js';
import { Player } from '/wotlk/core/player.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { TankDeathknightSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecTankDeathknight>(Spec.SpecTankDeathknight, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new TankDeathknightSimUI(document.body, player);
