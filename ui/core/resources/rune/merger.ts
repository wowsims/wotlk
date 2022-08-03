import { ResourceType } from "../../proto/api";
import { ResourceChangedLogGroup } from "../../proto_utils/logs_parser";

const RUNE_TYPES = [
	ResourceType.ResourceTypeBloodRune,
	ResourceType.ResourceTypeFrostRune,
	ResourceType.ResourceTypeUnholyRune,
	ResourceType.ResourceTypeDeathRune
]

export class RuneResourceMerger {
	public mergeRuneResources(resourceEvents: Record<ResourceType, ResourceChangedLogGroup[]>): Record<ResourceType, ResourceChangedLogGroup[]> {
		const mergedEvents: Record<ResourceType, ResourceChangedLogGroup[]> = { ...resourceEvents };
		let eventBuffer: ResourceChangedLogGroup[] = [];
		for (const runeType of RUNE_TYPES) {
			eventBuffer = eventBuffer.concat(mergedEvents[runeType]);
			delete mergedEvents[runeType];
		}
		mergedEvents[ResourceType.ResourceTypeBloodRune] = eventBuffer; // We do not have a generic rune type so blood rune is the placeholde.
		return mergedEvents;
	}
}
