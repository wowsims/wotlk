import { ResourceType } from '../proto/api';
import { ResourceChangedLogGroup } from '../proto_utils/logs_parser';

export class ResourceMerger {
	protected readonly targets: ResourceType[];
	protected readonly aggregate_key: ResourceType;

	public constructor(targets: ResourceType[], aggregate_key: ResourceType) {
		this.targets = targets;
		this.aggregate_key = aggregate_key;
	}

	public mergeResources(resourceEvents: Record<ResourceType, ResourceChangedLogGroup[]>): Record<ResourceType, ResourceChangedLogGroup[]> {
		const mergedEvents: Record<ResourceType, ResourceChangedLogGroup[]> = { ...resourceEvents };
		let eventBuffer: ResourceChangedLogGroup[] = [];
		for (const resourceType of this.targets) {
			eventBuffer = eventBuffer.concat(mergedEvents[resourceType]);
			mergedEvents[resourceType] = [];
		}
		mergedEvents[this.aggregate_key] = eventBuffer;
		return mergedEvents;
	}
}
