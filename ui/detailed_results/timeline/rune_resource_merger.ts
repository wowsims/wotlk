import { ResourceType } from "/wotlk/core/proto/api";
import { ResourceChangedLogGroup } from "/wotlk/core/proto_utils/logs_parser";

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
        console.log(eventBuffer)
        return mergedEvents;
    }
}
