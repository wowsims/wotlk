import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { Player } from '../core/player.js';
import { TypedEvent } from '../core/typed_event.js';

import { TankDeathknightSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecTankDeathknight>(Spec.SpecTankDeathknight, sim);
var hm = player.getHealingModel();
if (hm.cadenceSeconds == 0) {
    hm.cadenceSeconds = 2;
    player.setHealingModel(0, hm)
}
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new TankDeathknightSimUI(document.body, player);
