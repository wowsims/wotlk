import { Spec } from '../core/proto/common';
import { Sim } from '../core/sim';
import { Player } from '../core/player';
import { TypedEvent } from '../core/typed_event';
import { BalanceDruidSimUI } from './sim';

const sim = new Sim();
const player = new Player<Spec.SpecBalanceDruid>(Spec.SpecBalanceDruid, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);
export const simUI = new BalanceDruidSimUI(document.body, player);
