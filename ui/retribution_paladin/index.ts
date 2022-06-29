import { Spec } from '/wotlk/core/proto/common.js';
import { Sim } from '/wotlk/core/sim.js';
import { Player } from '/wotlk/core/player.js';
import { TypedEvent } from '/wotlk/core/typed_event.js';

import { RetributionPaladinSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecRetributionPaladin>(Spec.SpecRetributionPaladin, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new RetributionPaladinSimUI(document.body, player);