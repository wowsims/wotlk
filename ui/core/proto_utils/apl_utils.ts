import {
    ActionID as ActionIdProto,
    Cooldowns,
} from '../proto/common.js';

import {
	APLAction,
	APLPrepullAction,
} from '../proto/apl.js';

export function prepullPotionAction(doAt?: string): APLPrepullAction {
	return APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"${doAt || '-1s'}"}}}`);
}

export function autocastCooldownsAction(startAt?: string): APLAction {
    if (startAt) {
      	return APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"${startAt}"}}}},"autocastOtherCooldowns":{}}`);
    } else {
      	return APLAction.fromJsonString(`{"autocastOtherCooldowns":{}}`);
    }
}

export function scheduledCooldownAction(schedule: string, actionId: ActionIdProto): APLAction {
    return APLAction.fromJsonString(`{"schedule":{"schedule":"${schedule}","innerAction":{"castSpell":{"spellId":${ActionIdProto.toJsonString(actionId)}}}}}`);
}

export function simpleCooldownActions(cooldowns: Cooldowns): Array<APLAction> {
    return cooldowns.cooldowns
    .filter(cd => cd.id)
    .map(cd => {
        const schedule = cd.timings.map(timing => timing.toFixed(1) + 's').join(', ');
        return scheduledCooldownAction(schedule, cd.id!);
    });
}

export function standardCooldownDefaults(cooldowns: Cooldowns, prepotAt?: string, startAutocastCDsAt?: string): [Array<APLPrepullAction>, Array<APLAction>] {
    return [
        [prepullPotionAction(prepotAt)],
        [
            autocastCooldownsAction(startAutocastCDsAt),
            simpleCooldownActions(cooldowns),
        ].flat(),
    ];
}
