/// <reference types="jest" />
import { RuneResourceMerger } from "../detailed_results/timeline/rune_resource_merger"
import { ResourceChangedLogGroup, SimLogParams } from "../core/proto_utils/logs_parser"
import { ResourceType } from '../core/proto/api';

const defaultSimLogParams: SimLogParams = {
	raw: "",
	logIndex: 0,
	timestamp: 0,
	source: null,
	target: null,
	actionId: null,
	threat: 0,
}
const emptyResourceChangedLogGroupArray: ResourceChangedLogGroup[] = [];

describe("Rune resource merger", () => {
	it("Merges rune resources", () => {
		const input: Record<ResourceType, ResourceChangedLogGroup[]> = {
			[ResourceType.ResourceTypeNone]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeMana]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeEnergy]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRage]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeComboPoints]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeFocus]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRunicPower]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeHealth]: [
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeHealth, 100, 90, []),
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeHealth, 90, 80, [])
			],
			[ResourceType.ResourceTypeBloodRune]: [
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeBloodRune, 2, 1, []),
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeBloodRune, 1, 2, [])
			],
			[ResourceType.ResourceTypeFrostRune]: [
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeFrostRune, 2, 1, []),
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeFrostRune, 1, 0, [])
			],
			[ResourceType.ResourceTypeUnholyRune]: [
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 2, 1, []),
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, [])
			],
			[ResourceType.ResourceTypeDeathRune]: [
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 0, 1, []),
				new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, [])
			],
		}
		const merger = new RuneResourceMerger()
		merger.mergeRuneResources(input);
	})
})
