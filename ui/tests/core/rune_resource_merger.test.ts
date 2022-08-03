/// <reference types='jest' />
import { ResourceChangedLogGroup } from '../../core/proto_utils/logs_parser'
import { ResourceType } from '../../core/proto/api';
import { RuneResourceMerger } from '../../core/resources/rune/merger';
import { SimLogParams } from '../../core/proto_utils/sim_log_params';

const defaultSimLogParams: SimLogParams = {
	raw: '',
	logIndex: 0,
	timestamp: 0,
	source: null,
	target: null,
	actionId: null,
	threat: 0,
}
const emptyResourceChangedLogGroupArray: ResourceChangedLogGroup[] = [];
const health_event_a = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeHealth, 100, 90, []);
const health_event_b = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeHealth, 100, 90, []);
const bloodRuneEventA = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeBloodRune, 2, 1, []);
const bloodRuneEventB = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeBloodRune, 1, 2, []);
const frostRuneEventA = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeFrostRune, 2, 1, []);
const frostRuneEventB = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeFrostRune, 1, 0, []);
const unholyRuneEventA = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 2, 1, []);
const unholyRuneEventB = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, []);
const deathRuneEventA = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 0, 1, []);
const deathRuneEventB = new ResourceChangedLogGroup(defaultSimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, []);

describe('Rune resource merger', () => {
	it('Merges rune resources', () => {
		const input: Record<ResourceType, ResourceChangedLogGroup[]> = {
			[ResourceType.ResourceTypeNone]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeMana]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeEnergy]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRage]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeComboPoints]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeFocus]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRunicPower]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeHealth]: [health_event_a, health_event_b],
			[ResourceType.ResourceTypeBloodRune]: [bloodRuneEventA, bloodRuneEventB],
			[ResourceType.ResourceTypeFrostRune]: [frostRuneEventA, frostRuneEventB],
			[ResourceType.ResourceTypeUnholyRune]: [unholyRuneEventA, unholyRuneEventB],
			[ResourceType.ResourceTypeDeathRune]: [deathRuneEventA, deathRuneEventB],
		}
		const expected = input[ResourceType.ResourceTypeBloodRune]
			.concat(input[ResourceType.ResourceTypeFrostRune])
			.concat(input[ResourceType.ResourceTypeUnholyRune])
			.concat(input[ResourceType.ResourceTypeDeathRune]);
		const merger = new RuneResourceMerger();
		const output = merger.mergeRuneResources(input);
		expect(output[ResourceType.ResourceTypeBloodRune]).toStrictEqual(expected);
	})
})
