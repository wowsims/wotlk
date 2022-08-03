import { ActionId } from "./action_id";
import { Entity } from "./logs_parser";

export interface SimLogParams {
	raw: string,
	logIndex: number,
	timestamp: number,
	source: Entity | null,
	target: Entity | null,
	actionId: ActionId | null,
	threat: number,
}
